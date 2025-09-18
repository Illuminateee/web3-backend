package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/api"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize server
	server, err := api.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Print startup message
	fmt.Printf("Starting web3-tokensale-be server on port %s\n", cfg.Port)
	fmt.Printf("Connected to Ethereum RPC: %s\n", cfg.EthereumRPC)

	// Run the server in a separate goroutine
	go func() {
		server.Run()
		if err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
		fmt.Println("Server is running...")
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down...")
}
