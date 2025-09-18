package routes

import (
	"net/http"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/api/handlers"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/api/middleware"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handler *handlers.Handler, authMiddleware gin.HandlerFunc, secretKey string, activityLogger *services.ActivityLoggerService) {
    // API v1 group
    v1 := router.Group("/api/v1")
    {
        // Token quote endpoints
        v1.GET("/quote/:token", handler.GetCifoQuoteHandler)

        // Token conversion endpoints
        v1.GET("/convert/token-to-eth", handler.CifoToEthHandler)
        v1.GET("/convert/eth-to-token", handler.EthToCifoHandler)

        // On-ramp endpoint
        v1.POST("/onramp/session", handler.CreateOnRampSessionHandler)

        // Midtrans payment endpoints
        v1.POST("/payment/midtrans", handler.CreateMidtransPaymentHandler)
        v1.POST("/payment/midtrans/webhook", handler.ProcessMidtransWebhookHandler)

        v1.POST("/ethereum/account", handler.CreateAccountHandler)
        v1.POST("/ethereum/import", handler.ImportAccountHandler)
        v1.GET("/ethereum/balance/:address", handler.GetAccountBalanceHandler)

        

        transactionGroup := v1.Group("/transactions")
        {
            // Get transaction details by blockchain transaction hash
            v1.GET("/tx/:hash", handler.GetTransactionDetailsHandler)
        
            // Get details of a specific transaction by payment_id or UUID
            transactionGroup.GET("/:id", handler.GetTransactionStatusHandler)
        
            // Get a list of transactions for the authenticated user (with pagination)
            transactionGroup.GET("", authMiddleware, handler.GetUserTransactionsHandler)
        }

        transakGroup := v1.Group("/transak")
        {
            transakGroup.Use(authMiddleware)
            transakGroup.POST("/order", handler.CreateTransakOrderHandler)
            transakGroup.POST("/cifo", handler.CreateFiatToCIFOHandler)
            transakGroup.GET("/status/:order_id", handler.GetTransakOrderStatusHandler)
        }

        // Transak webhook endpoint (no authentication)
        v1.POST("/payment/transak/webhook", handler.ProcessTransakWebhookHandler)

        // Midtrans-to-Transak flow
        v1.POST("/payment/midtrans-to-transak", authMiddleware, handler.CreateMidtransToTransakHandler)
        v1.POST("/payment/midtrans-to-transak/webhook", handler.ProcessMidtransTransakWebhookHandler)
        
        sendGroup := v1.Group("/send")
        {
            sendGroup.Use(authMiddleware)
            sendGroup.POST("/create-payment", handler.CreateFiatToTokenPaymentHandler)
        }

        purchaseGroup := v1.Group("/purchase")
        {
            purchaseGroup.Use(authMiddleware)
            // from wallet ethereum cifo one way
            purchaseGroup.POST("/cifo", handler.PurchaseCifoHandler)
            purchaseGroup.POST("/cifo/auto-swap", authMiddleware, handler.AutoSwapCifoPurchaseHandler)
            purchaseGroup.POST("/payment/autoswap-webhook", handler.ProcessAutoSwapWebhookHandler)
        }

        router.Use(middleware.SessionTimeout(handler.DB, secretKey))
        authGroup := v1.Group("/auth")
        {
            // Existing endpoints
            authGroup.Use(middleware.ActivityLoggingMiddleware(activityLogger))
            authGroup.POST("/login", handler.LoginHandler)
            authGroup.POST("/login/2fa", handler.LoginWith2FAHandler)

            authGroup.POST("/register", handler.RegisterUserHandler)
            
            // New endpoints
            authGroup.GET("/refresh", handler.RefreshTokenHandler)
            authGroup.POST("/logout", handler.LogoutHandler)

            // Recovery endpoints
            authGroup.POST("/recovery/request", handler.RequestWalletRecoveryHandler)
            authGroup.POST("/recovery/verify", handler.VerifyWalletRecoveryHandler)

            // 2FA endpoints
            authGroup.POST("/2fa/setup", handler.Setup2FAHandler)
            authGroup.POST("/2fa/verify", handler.Verify2FAHandler)
            authGroup.POST("/2fa/recover", handler.Request2FARecovery)  // Add this line
            authGroup.POST("/2fa/verify-recovery", handler.Verify2FARecoveryHandler)  // Add this line

            // Protected routes (require authentication)
            protected := authGroup.Group("/")
            protected.Use(authMiddleware)
            {
                protected.POST("/2fa/disable", handler.Disable2FAHandler)
            }
        }

        walletGroup := v1.Group("/wallet")
        {
            walletGroup.Use(authMiddleware)
                walletGroup.GET("/balance", handler.GetUserWalletBalanceHandler)
                walletGroup.POST("/create", handler.CreateWalletHandler)
                walletGroup.POST("/import", handler.ImportWalletHandler)
                walletGroup.POST("/swap", handler.SwapTokensHandler)
                walletGroup.GET("/transactions", handler.GetWalletTransactionsHandler)
                // Add new wallet backup and recovery routes
                walletGroup.POST("/enable-backup", handler.EnableWalletBackupHandler)
                walletGroup.POST("/recover", handler.RecoverWalletHandler)
        }

        // CIFO token specific endpoints for convenience
        cifoGroup := v1.Group("/cifo")
        {
            cifoGroup.GET("/quote", handler.GetCifoQuoteHandler)
            cifoGroup.GET("/convert-to-eth", handler.CifoToEthHandler)
            cifoGroup.GET("/convert-from-eth", handler.EthToCifoHandler)
            cifoGroup.GET("/convert-from-fiat", handler.FiatToTokenHandler)

            // Health check endpoint
            cifoGroup.GET("/health", func(c *gin.Context) {
                ctx := c.Request.Context()

                // Get ETH prices as health check
                ethPriceUSD, ethPriceIDR, err := handler.GetEthPrices(ctx)

                if err != nil {
                    c.JSON(http.StatusServiceUnavailable, gin.H{
                        "status":  "error",
                        "message": "Failed to fetch ETH prices",
                        "error":   err.Error(),
                    })
                    return
                }

                c.JSON(http.StatusOK, gin.H{
                    "status": "healthy",
                    "eth_price": gin.H{
                        "usd": ethPriceUSD,
                        "idr": ethPriceIDR,
                    },
                })
            })
        }
    }
}