package api

import (
	"context"
	"fmt"
	"log"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/api/auth"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/api/handlers"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/api/middleware"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/api/routes"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/blockchain"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/config"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/database"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/services"
	"git.winteraccess.id/walanja/web3-tokensale-be/pkg/ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
    router  *gin.Engine
    config  *config.Config
    handler *handlers.Handler
    db      *gorm.DB
    walletDB *gorm.DB
}

func NewServer(cfg *config.Config) (*Server, error) {
    // Initialize Ethereum client
    log.Printf("Connecting to database at %s:%s...", cfg.DBHost, cfg.DBPort)
    db, err := database.Connect(cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }

    log.Printf("Connecting to Ethereum RPC: %s", cfg.EthereumRPC)
    log.Printf("Using payment gateway address: %s", cfg.PaymentGatewayAddress.Hex())
    log.Printf("Using token address: %s", cfg.TokenAddress.Hex())
    
    ethClient, err := ethereum.NewClient(cfg.EthereumRPC)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
    }
    

    networkID, err := ethClient.Client.NetworkID(context.Background())
    if err != nil {
        return nil, fmt.Errorf("failed to get network ID: %v", err)
    }
    log.Printf("Connected to network with ID: %s", networkID.String())

    code, err := ethClient.Client.CodeAt(context.Background(), cfg.PaymentGatewayAddress, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to check contract code: %v", err)
    }
    if len(code) == 0 {
        return nil, fmt.Errorf("no contract code found at address %s - verify you're connected to the correct network", cfg.PaymentGatewayAddress.Hex())
    }

    // Initialize Uniswap client
    uniswapClient := ethereum.NewUniswapClient(ethClient.RPCClient, cfg)

    // Initialize payment gateway client
    // Pass the Ethereum client directly - your constructor should accept a client directly
	paymentGateway, err := blockchain.NewPaymentGatewayClient(
        cfg.EthereumRPC,
        cfg.PaymentGatewayAddress.Hex(), 
        cfg.PrivateKey,
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to initialize payment gateway client: %v", err)
    }

    mainDB, err := database.Connect(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    walletDB, err := database.ConnectWalletDB(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to wallet database: %v", err)
    }

    encryptionService, err := services.NewEncryptionService(cfg.WalletDB.EncryptKey)
    if err != nil {
        log.Fatalf("Failed to initialize encryption service: %v", err)
    }

    // Initialize services
    priceService := services.NewPriceService(ethClient, uniswapClient)
    blockchainService, err := services.NewBlockchainService(cfg, paymentGateway)
    // Initialize activity logger
    activityLogger := services.NewActivityLoggerService(mainDB)

    // jwt service
    jwtConfig := auth.JWTConfig{
        SecretKey:     cfg.JWTSecret,
        TokenDuration: cfg.JWTExpiration,
    }
    tokenService := auth.NewTokenService(jwtConfig)
    
    // Initialize TOTP service
    totpService := auth.NewTOTPService("Web3Tokensale")
    
    // Initialize recovery service
    recoveryService := services.NewRecoveryService(
        db,
        cfg.SendGridAPIKey,
        cfg.AppURL,
        cfg.FromEmail,
        cfg.FromName,
    )

    if err != nil {
        return nil, fmt.Errorf("failed to initialize blockchain service: %v", err)
    }

    walletStorageService := services.NewWalletStorageService(
        walletDB,
        mainDB,
        encryptionService,
        activityLogger,
    )
    // Initialize wallet service
    walletService := services.NewWalletService(mainDB, walletStorageService)
    transakService := services.NewTransakService(db)
    activityLogger = services.NewActivityLoggerService(db)

    // Initialize swap service
    swapService := services.NewSwapService(
        db, 
        ethClient.Client, 
        walletService,
        common.HexToAddress(cfg.UniswapRouterAddress),
        common.HexToAddress(cfg.WrappedEthAddress),
        cfg.TokenAddress, // Now passing the address directly
    )
    // Initialize handlers
    handler := handlers.NewHandler(db, priceService, blockchainService, cfg, tokenService, totpService, recoveryService, walletService, transakService, activityLogger,walletStorageService,encryptionService,swapService)

    // Initialize router
    router := gin.Default()

    authMiddleware := middleware.AuthMiddleware(tokenService)

    // Setup routes
    router.Use(middleware.CorsMiddleware())
    routes.SetupRoutes(router, handler, authMiddleware,jwtConfig.SecretKey, activityLogger)

    return &Server{
        router:  router,
        config:  cfg,
        handler: handler,
        db:      db,
        walletDB: walletDB,
    }, nil
}

func (s *Server) Run() error {
    addr := fmt.Sprintf(":%s", s.config.Port)
    return s.router.Run(addr)
}