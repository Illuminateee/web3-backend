package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

type Config struct {
    Port                  string
    Environment           string
    EthereumRPC           string
    UniswapV3URL          string
    StablecoinAddress     common.Address
    TokenAddress          common.Address           // CIFO token address
    WethAddress           common.Address           // WETH address
    PaymentGatewayAddress common.Address           // Payment gateway contract address
    PrivateKey            string    
    WalletPrivateKey    string                   

    // jwt configuration
    JWTSecret     string
    JWTExpiration time.Duration

    // Email/Recovery configuration
    SendGridAPIKey string
    AppURL         string
    FromEmail      string
    FromName       string

    TransakAPIKey    string
    TransakSecretKey string
    TransakBaseURL  string

    // Database configuration
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string

    WalletDB WalletDBConfig `mapstructure:",squash"`

    // Add these fields for Uniswap integration
    UniswapRouterAddress string
    WrappedEthAddress    string
}

type WalletDBConfig struct {
    Host        string
    User        string
    Password    string
    DBName      string
    Port        string
    DBSSLMode  string
    EncryptKey  string
}

func LoadConfig() (*Config, error) {
    err := godotenv.Load()
    if err != nil {
        // Try to load from parent directory
        err = godotenv.Load("../.env")
        if err != nil {
            // Try one more directory up
            err = godotenv.Load("../../.env")
            if err != nil {
                log.Println("Warning: .env file not found, using default values")
            }
        }
    }

    // Default addresses
    defaultCifo := "0x5FbDB2315678afecb367f032d93F642f64180aa3"        // TestToken
    defaultWeth := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"        // WETH on Ethereum mainnet
    defaultPaymentGateway := "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512" // Replace with your actual contract address

    config := &Config{
        Port:                  getEnv("PORT", "8080"),
        Environment:           getEnv("APP_ENV", "development"),
        EthereumRPC:           getEnv("ETHEREUM_RPC", "https://geth-geth.ede2390e1937cf50.dyndns.dappnode.io"),
        UniswapV3URL:          getEnv("UNISWAP_V3_URL", "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3"),
        StablecoinAddress:     common.HexToAddress(getEnv("STABLECOIN_ADDRESS", defaultCifo)),
        TokenAddress:          common.HexToAddress(getEnv("TOKEN_ADDRESS", defaultCifo)),
        WethAddress:           common.HexToAddress(getEnv("WETH_ADDRESS", defaultWeth)),
        PaymentGatewayAddress: common.HexToAddress(getEnv("PAYMENT_GATEWAY_ADDRESS", defaultPaymentGateway)),
        PrivateKey:            getEnv("PRIVATE_KEY", "1e88f382bfed1d0597d717d10063c1ea7149d106b36078edbe5913b5d8f0327e"), // Never include default private keys in code

        JWTSecret:    getEnv("JWT_SECRET", "your_jwt_secret"),
        JWTExpiration: time.Duration(getEnvAsInt("JWT_EXPIRATION_HOURS", 24)) * time.Hour,

        SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),
        AppURL:         getEnv("APP_URL", "http://localhost:3000"),
        FromEmail:      getEnv("FROM_EMAIL", "no-reply@example.com"),
        FromName:       getEnv("FROM_NAME", "Web3 Tokensale"),

        
        TransakAPIKey:    getEnv("TRANSAK_API_KEY", ""),
        TransakSecretKey: getEnv("TRANSAK_SECRET_KEY", ""),
        TransakBaseURL:  getEnv("TRANSAK_BASE_URL", "https://global.transak.com"),

        
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "web3_tokensale"),


        // Wallet database config
        WalletDB: WalletDBConfig{
            Host:       getEnv("WALLET_DB_HOST", "localhost"),
            Port:       getEnv("WALLET_DB_PORT", "5432"),
            User:       getEnv("WALLET_DB_USER", "postgres"),
            Password:   getEnv("WALLET_DB_PASSWORD", ""),
            DBName:     getEnv("WALLET_DB_NAME", "wallet_credentials"),
            DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
            EncryptKey: getEnv("WALLET_DB_ENCRYPT_KEY", ""),
        },
    }

        // Validate encryption key if provided
    if config.WalletDB.EncryptKey != "" && len(config.WalletDB.EncryptKey) != 64 {
        log.Println("Warning: WALLET_DB_ENCRYPT_KEY should be 64 hex characters (32 bytes)")
    }

    return config, nil
}

func getEnvAsInt(key string, defaultVal int) int {
    valueStr := getEnv(key, "")
    if valueStr == "" {
        return defaultVal
    }
    
    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultVal
    }
    
    return value
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}