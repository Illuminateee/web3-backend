package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "runtime"
)

func main() {
    // Get project root directory (assuming script is run from project root or the script itself is in internal/scripts)
    workingDir, err := os.Getwd()
    if err != nil {
        log.Fatalf("Failed to get working directory: %v", err)
    }

    // Check if we're in the scripts directory or project root
    if strings.HasSuffix(workingDir, filepath.Join("internal", "scripts")) {
        // We're running from the scripts directory, go up one level
        workingDir = filepath.Dir(filepath.Dir(workingDir))
    }

    // Now define paths relative to project root
    contractsDir := filepath.Join(workingDir, "internal", "bindings", "contracts")
    buildDir := filepath.Join(workingDir, "internal", "bindings", "build")
    genDir := filepath.Join(workingDir, "internal", "bindings", "generated")
    nodeModulesDir := filepath.Join(workingDir, "node_modules")

    fmt.Printf("Working directory: %s\n", workingDir)
    fmt.Printf("Contracts directory: %s\n", contractsDir)
    fmt.Printf("Build directory: %s\n", buildDir)
    fmt.Printf("Generated directory: %s\n", genDir)
    fmt.Printf("Node modules directory: %s\n", nodeModulesDir)

    // Ensure directories exist
    ensureDirExists(contractsDir)
    ensureDirExists(buildDir)
    ensureDirExists(genDir)

    // Check if node_modules exists
    if _, err := os.Stat(nodeModulesDir); os.IsNotExist(err) {
        log.Println("Warning: node_modules directory not found. Installing @openzeppelin/contracts...")
        installCmd := exec.Command("npm", "install", "@openzeppelin/contracts")
        installCmd.Dir = workingDir
        if output, err := installCmd.CombinedOutput(); err != nil {
            log.Fatalf("Failed to install dependencies: %v\n%s", err, output)
        }
    }

    // Define contracts to compile
    contracts := map[string]string{
        "TestToken": "TestToken",
        "FiatToTokenPaymentGateway": "FiatToTokenPaymentGateway",
    }

    // Check if Docker is installed
    if _, err := exec.LookPath("docker"); err != nil {
        log.Fatal("Docker not found. Please install Docker to compile Solidity contracts.")
    }

    // Check if abigen is installed
    if _, err := exec.LookPath("abigen"); err != nil {
        log.Fatal("abigen not found in PATH. Please install with: go install github.com/ethereum/go-ethereum/cmd/abigen@latest")
    }
    
    // Clean build directory first
    if err := os.RemoveAll(buildDir); err != nil {
        log.Fatalf("Failed to clean build directory: %v", err)
    }
    ensureDirExists(buildDir)

    // Compile contracts using Docker
    fmt.Println("Compiling contracts using Docker solc image...")

    for filename, contractName := range contracts {
        contractPath := filepath.Join(contractsDir, filename+".sol")
        
        // Check if contract file exists
        if _, err := os.Stat(contractPath); os.IsNotExist(err) {
            log.Fatalf("Contract file not found: %s", contractPath)
        }

        fmt.Printf("Compiling %s...\n", filename)
        
        // Convert paths for Docker mounting
        // For Windows, convert backslashes to forward slashes and adjust path format for Docker
        sourcesMount := contractsDir
        outputMount := buildDir
        nodeModulesMount := nodeModulesDir
        
        if runtime.GOOS == "windows" {
            // Convert Windows paths to Docker format
            sourcesMount = convertWindowsPathToDocker(contractsDir)
            outputMount = convertWindowsPathToDocker(buildDir)
            nodeModulesMount = convertWindowsPathToDocker(nodeModulesDir)
        }
        
        args := []string{
			"run", "--rm",
			"--volume", fmt.Sprintf("%s:/sources", sourcesMount),
			"--volume", fmt.Sprintf("%s:/output", outputMount),
			"--volume", fmt.Sprintf("%s:/node_modules", nodeModulesMount),
			"ethereum/solc:0.8.20",  // Specify solc version explicitly
			"/sources/" + filename + ".sol",
			"--abi", "--bin",
			"--optimize", "--optimize-runs", "200",
			"--base-path", "/sources",
			"--include-path", "/node_modules",
			"--output-dir", "/output",
			"--overwrite",  // Add this flag to allow overwriting existing files
		}
        
        cmd := exec.Command("docker", args...)
        output, err := cmd.CombinedOutput()
        
        if err != nil {
            fmt.Printf("Error compiling %s: %v\n", filename, err)
            fmt.Println("Output:", string(output))
            continue
        }
        
        // Check if the output files were created with the expected names
        abiFile := filepath.Join(buildDir, contractName+".abi")
        binFile := filepath.Join(buildDir, contractName+".bin")
        
        if _, err := os.Stat(abiFile); os.IsNotExist(err) {
            // Try to find any .abi file that might have been created
            files, err := ioutil.ReadDir(buildDir)
            if err != nil {
                log.Fatalf("Failed to read build directory: %v", err)
            }
            
            abiFound := false
            for _, file := range files {
                if strings.HasSuffix(file.Name(), ".abi") {
                    log.Printf("Found ABI file with different name: %s", file.Name())
                    // Rename the file to what we expect
                    oldPath := filepath.Join(buildDir, file.Name())
                    if err := os.Rename(oldPath, abiFile); err != nil {
                        log.Fatalf("Failed to rename ABI file: %v", err)
                    }
                    abiFound = true
                    break
                }
            }
            
            if !abiFound {
                log.Printf("No ABI file found for %s", contractName)
                continue
            }
        }
        
        if _, err := os.Stat(binFile); os.IsNotExist(err) {
            // Try to find any .bin file that might have been created
            files, err := ioutil.ReadDir(buildDir)
            if err != nil {
                log.Fatalf("Failed to read build directory: %v", err)
            }
            
            binFound := false
            for _, file := range files {
                if strings.HasSuffix(file.Name(), ".bin") {
                    log.Printf("Found BIN file with different name: %s", file.Name())
                    // Rename the file to what we expect
                    oldPath := filepath.Join(buildDir, file.Name())
                    if err := os.Rename(oldPath, binFile); err != nil {
                        log.Fatalf("Failed to rename BIN file: %v", err)
                    }
                    binFound = true
                    break
                }
            }
            
            if !binFound {
                log.Printf("No BIN file found for %s", contractName)
                continue
            }
        }
        
        fmt.Printf("Successfully compiled %s\n", filename)
    }
    
    // List files in build directory for debugging
    files, err := ioutil.ReadDir(buildDir)
    if err != nil {
        log.Printf("Failed to read build directory: %v", err)
    } else {
        fmt.Println("Files in build directory:")
        for _, file := range files {
            fmt.Printf("  %s\n", file.Name())
        }
    }

    // Generate Go bindings
    fmt.Println("Generating Go bindings...")

    for _, contractName := range contracts {
        abiFile := filepath.Join(buildDir, contractName+".abi")
        binFile := filepath.Join(buildDir, contractName+".bin")
        
        // Check if ABI file exists
        if _, err := os.Stat(abiFile); os.IsNotExist(err) {
            log.Printf("ABI file not found: %s", abiFile)
            continue
        }
        
        // Check if BIN file exists
        if _, err := os.Stat(binFile); os.IsNotExist(err) {
            log.Printf("BIN file not found: %s", binFile)
            continue
        }
        
        // Create package directory with lowercase name
        pkgName := strings.ToLower(contractName)
        pkgDir := filepath.Join(genDir, pkgName)
        ensureDirExists(pkgDir)
        
        outFile := filepath.Join(pkgDir, pkgName+".go")
        
        fmt.Printf("Generating binding for %s...\n", contractName)
        
        args := []string{
            "--abi", abiFile,
            "--bin", binFile,
            "--pkg", pkgName,
            "--out", outFile,
            "--type", contractName,
        }
        
        cmd := exec.Command("abigen", args...)
        output, err := cmd.CombinedOutput()
        
        if err != nil {
            fmt.Printf("Error generating binding for %s: %v\n", contractName, err)
            fmt.Println("Output:", string(output))
            continue
        }
        
        fmt.Printf("Successfully generated binding for %s at %s\n", contractName, outFile)
    }

    fmt.Println("Binding generation complete!")
}

// ensureDirExists creates a directory if it doesn't exist
func ensureDirExists(path string) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        err = os.MkdirAll(path, 0755)
        if err != nil {
            log.Fatalf("Failed to create directory %s: %v", path, err)
        }
        fmt.Printf("Created directory: %s\n", path)
    }
}

// convertWindowsPathToDocker converts a Windows path to a format usable by Docker
func convertWindowsPathToDocker(windowsPath string) string {
    if runtime.GOOS != "windows" {
        return windowsPath
    }
    
    // Convert C:\path\to\dir to /c/path/to/dir format for Docker
    if len(windowsPath) > 2 && windowsPath[1] == ':' {
        drive := strings.ToLower(string(windowsPath[0]))
        path := strings.ReplaceAll(windowsPath[2:], "\\", "/")
        return "/" + drive + path
    }
    
    // Just convert backslashes to forward slashes if not a drive path
    return strings.ReplaceAll(windowsPath, "\\", "/")
}