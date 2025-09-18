package handlers

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"

    "git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
    "git.winteraccess.id/walanja/web3-tokensale-be/internal/services"
    "github.com/ethereum/go-ethereum/common"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

)

// FiatToCryptoRequest represents a request to convert fiat to crypto through Transak
type FiatToCryptoRequest struct {
    FiatAmount        float64 `json:"fiat_amount" binding:"required"`
    FiatCurrency      string  `json:"fiat_currency" binding:"required"`
    WalletAddress     string  `json:"wallet_address" binding:"required"`
    SwapToCifo        bool    `json:"swap_to_cifo"`
    CifoAmount        float64 `json:"cifo_amount"`
    Email             string  `json:"email"`
    Name              string  `json:"name"`
    Phone             string  `json:"phone"`
    SuccessURL        string  `json:"success_url"`
    CancelURL         string  `json:"cancel_url"`
}

// FiatToCIFORequest represents a request to buy CIFO tokens with fiat through Transak+Uniswap
type FiatToCIFORequest struct {
    CifoAmount      float64 `json:"cifo_amount" binding:"required"`
    FiatCurrency    string  `json:"fiat_currency" binding:"required"`
    WalletAddress   string  `json:"wallet_address" binding:"required"`
    Email           string  `json:"email"`
    Name            string  `json:"name"`
    Phone           string  `json:"phone"`
    SuccessURL      string  `json:"success_url"`
    CancelURL       string  `json:"cancel_url"`
}

// CreateTransakOrderHandler creates a new order for buying ETH with fiat through Transak
func (h *Handler) CreateTransakOrderHandler(c *gin.Context) {
    // Get authenticated user
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }
    
    // Parse request
    var req FiatToCryptoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }
    
    // Validate request
    if req.FiatAmount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Fiat amount must be greater than 0"})
        return
    }
    
    if !common.IsHexAddress(req.WalletAddress) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet address"})
        return
    }
    
    // Normalize currency
    req.FiatCurrency = strings.ToUpper(req.FiatCurrency)
    if req.FiatCurrency != "IDR" && req.FiatCurrency != "USD" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Currency must be either IDR or USD"})
        return
    }
    
    // Set default URLs if not provided
    if req.SuccessURL == "" {
        req.SuccessURL = "https://yourwebsite.com/payment/success"
    }
    if req.CancelURL == "" {
        req.CancelURL = "https://yourwebsite.com/payment/cancel"
    }
    
    // Calculate expected ETH amount based on current rates
    ethPriceUSD, ethPriceIDR, err := h.GetEthPrices(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ETH prices"})
        return
    }
    
    // Calculate expected ETH amount
    var ethAmount float64
    if req.FiatCurrency == "USD" {
        ethAmount = req.FiatAmount / ethPriceUSD
    } else {
        ethAmount = req.FiatAmount / ethPriceIDR
    }
    
    // If planning to swap to CIFO, calculate CIFO amount if not provided
    var cifoAmount float64
    if req.SwapToCifo {
        if req.CifoAmount > 0 {
            // Use provided CIFO amount
            cifoAmount = req.CifoAmount
        } else {
            // Calculate CIFO amount based on ETH
            cifoAmountStr, err := h.getCifoAmount(ethAmount)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate CIFO amount"})
                return
            }
            
            cifoAmount, err = strconv.ParseFloat(cifoAmountStr, 64)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CIFO amount"})
                return
            }
        }
    }
    
    // Parse user ID
    uid, err := uuid.Parse(userID.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    
    // Create transaction record
    txUUID := uuid.New()
    orderID := fmt.Sprintf("TRANSAK-%s", txUUID.String()[:8])
    
transaction := models.Transaction{
    UUID:               txUUID,
    UserID:             uid,
    PaymentID:          orderID,
    WalletAddress:      req.WalletAddress,
    FiatCurrency:       req.FiatCurrency,
    FiatAmount:         req.FiatAmount,
    EthAmount:          ethAmount,
    TokenAmount:        cifoAmount,
    TokenSymbol:        func() string { if req.SwapToCifo { return "CIFO" } else { return "ETH" } }(),
    Status:             models.TransactionStatusPending,
    PaymentMethod:      "transak",
    EthPriceAtPurchase: func() float64 { if req.FiatCurrency == "USD" { return ethPriceUSD } else { return ethPriceIDR } }(),
    TransactionType:    func() string { if req.SwapToCifo { return "fiat_to_cifo" } else { return "fiat_to_eth" } }(),
    SwapType:           func() string { if req.SwapToCifo { return "uniswap" } else { return "" } }(),
    MinTokenAmount:     cifoAmount * 0.97, // 3% slippage if swapping to CIFO
    CreatedAt:          time.Now(),
    UpdatedAt:          time.Now(),
}
    
    // Save transaction to database
    if err := h.DB.Create(&transaction).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction record"})
        log.Printf("Error creating transaction: %v", err)
        return
    }
    
    // Create Transak order
    transakResp, err := h.TransakService.CreateOrder(c.Request.Context(), &transaction)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Transak order: " + err.Error()})
        
        // Update transaction with error
        transaction.Status = "transak_failed"
        transaction.ErrorMessage = err.Error()
        h.DB.Save(&transaction)
        
        return
    }
    
    // Return response with Transak checkout link
    c.JSON(http.StatusOK, gin.H{
        "transaction": gin.H{
            "id":              transaction.UUID.String(),
            "payment_id":      transaction.PaymentID,
            "created_at":      transaction.CreatedAt,
            "eth_amount":      ethAmount,
            "cifo_amount":     cifoAmount,
            "swap_to_cifo":    req.SwapToCifo,
            "status":          transaction.Status,
        },
        "transak": gin.H{
            "order_id":      transakResp.Data.ID,
            "checkout_link": transakResp.Data.CheckoutLink,
            "status":        transakResp.Data.Status,
        },
        "destination_wallet": req.WalletAddress,
        "costs": gin.H{
            "fiat_amount":      req.FiatAmount,
            "fiat_currency":    req.FiatCurrency,
            "eth_price":        func() float64 { if req.FiatCurrency == "USD" { return ethPriceUSD } else { return ethPriceIDR } }(),
            "expected_eth":     ethAmount,
        },
        "next_steps": func() []string {
            if req.SwapToCifo {
                return []string{
                    "Complete payment through Transak checkout link",
                    "After ETH arrives in your wallet, swap to CIFO using Uniswap",
                }
            } else {
                return []string{
                    "Complete payment through Transak checkout link",
                    "",
                }
            }
        }(),
    })
}

// CreateFiatToCIFOHandler creates an order specifically for buying CIFO with fiat
func (h *Handler) CreateFiatToCIFOHandler(c *gin.Context) {
    // Get authenticated user
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }
    
    // Parse request
    var req FiatToCIFORequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }
    
    // Validate request
    if req.CifoAmount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "CIFO amount must be greater than 0"})
        return
    }
    
    if !common.IsHexAddress(req.WalletAddress) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet address"})
        return
    }
    
    // Normalize currency
    req.FiatCurrency = strings.ToUpper(req.FiatCurrency)
    if req.FiatCurrency != "IDR" && req.FiatCurrency != "USD" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Currency must be either IDR or USD"})
        return
    }
    
    // Set default URLs if not provided
    if req.SuccessURL == "" {
        req.SuccessURL = "https://yourwebsite.com/payment/success"
    }
    if req.CancelURL == "" {
        req.CancelURL = "https://yourwebsite.com/payment/cancel"
    }
    
    // Step 1: Calculate ETH required for the CIFO amount
    cifoPerEth, err := h.getCifoAmount(1.0)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get CIFO rate"})
        return
    }
    
    cifoPerEthFloat, err := strconv.ParseFloat(cifoPerEth, 64)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CIFO rate"})
        return
    }
    
    // Calculate ETH required
    ethRequired := req.CifoAmount / cifoPerEthFloat
    
    // Step 2: Calculate fiat amount required based on current rates
    ethPriceUSD, ethPriceIDR, err := h.GetEthPrices(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ETH prices"})
        return
    }
    
    // Calculate fiat amount
    var fiatAmount float64
    if req.FiatCurrency == "USD" {
        fiatAmount = ethRequired * ethPriceUSD
    } else {
        fiatAmount = ethRequired * ethPriceIDR
    }
    
    // Add 3% slippage buffer
    fiatAmount = fiatAmount * 1.03
    
    // Parse user ID
    uid, err := uuid.Parse(userID.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    
    // Create transaction record
    txUUID := uuid.New()
    orderID := fmt.Sprintf("CIFO-%s", txUUID.String()[:8])
    
    transaction := models.Transaction{
        UUID:               txUUID,
        UserID:             uid,
        PaymentID:          orderID,
        WalletAddress:      req.WalletAddress,
        FiatCurrency:       req.FiatCurrency,
        FiatAmount:         fiatAmount,
        EthAmount:          ethRequired,
        TokenAmount:        req.CifoAmount,
        TokenSymbol:        "CIFO",
        Status:             models.TransactionStatusPending,
        PaymentMethod:      "transak",
        EthPriceAtPurchase: func() float64 { if req.FiatCurrency == "USD" { return ethPriceUSD } else { return ethPriceIDR } }(),
        TransactionType:    "fiat_to_cifo",
        SwapType:           "uniswap",
        MinTokenAmount:     req.CifoAmount * 0.97, // 3% slippage
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }
    
    // Save transaction to database
    if err := h.DB.Create(&transaction).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction record"})
        log.Printf("Error creating transaction: %v", err)
        return
    }
    
    // Create Transak order
    transakResp, err := h.TransakService.CreateOrder(c.Request.Context(), &transaction)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Transak order: " + err.Error()})
        
        // Update transaction with error
        transaction.Status = "transak_failed"
        transaction.ErrorMessage = err.Error()
        h.DB.Save(&transaction)
        
        return
    }
    
    // Generate Uniswap swap URL
    uniswapURL := generateUniswapURL(req.WalletAddress, ethRequired, req.CifoAmount)
    
    // Return response with Transak checkout link and next steps
    c.JSON(http.StatusOK, gin.H{
        "transaction": gin.H{
            "id":              transaction.UUID.String(),
            "payment_id":      transaction.PaymentID,
            "created_at":      transaction.CreatedAt,
            "cifo_amount":     req.CifoAmount,
            "eth_required":    ethRequired,
            "status":          transaction.Status,
        },
        "transak": gin.H{
            "order_id":      transakResp.Data.ID,
            "checkout_link": transakResp.Data.CheckoutLink,
            "status":        transakResp.Data.Status,
        },
        "uniswap": gin.H{
            "swap_url": uniswapURL,
            "note":     "Use this URL after ETH arrives in your wallet to swap for CIFO",
        },
        "destination_wallet": req.WalletAddress,
        "costs": gin.H{
            "eth_required":    ethRequired,
            "fiat_amount":     fiatAmount,
            "fiat_currency":   req.FiatCurrency,
            "eth_price":       func() float64 { if req.FiatCurrency == "USD" { return ethPriceUSD } else { return ethPriceIDR } }(),
            "slippage_buffer": "3%",
        },
        "exchange_rate": gin.H{
            "cifo_per_eth": cifoPerEthFloat,
            "eth_per_cifo": 1/cifoPerEthFloat,
        },
        "next_steps": []string{
            "Complete payment through Transak checkout link",
            "After ETH arrives in your wallet, use the Uniswap link to swap to CIFO",
        },
    })
}

// ProcessTransakWebhookHandler processes webhooks from Transak
func (h *Handler) ProcessTransakWebhookHandler(c *gin.Context) {
    // Read the request body
    body, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
        return
    }
    
    // Parse webhook payload
    var payload services.TransakWebhookPayload
    if err := json.Unmarshal(body, &payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
        return
    }
    
    // Process the webhook
    if err := h.TransakService.ProcessWebhook(&payload); err != nil {
        log.Printf("Error processing Transak webhook: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
        return
    }
    
    // If status is COMPLETED and we need to auto-swap, trigger swap process
    if payload.Status == "COMPLETED" {
        // Find transaction
        var transaction models.Transaction
        if err := h.DB.Where("payment_reference = ?", payload.OrderID).
            Or("payment_id = ?", payload.PartnerOrderID).
            First(&transaction).Error; err != nil {
            log.Printf("Error finding transaction for auto-swap: %v", err)
        } else {
            // Check if auto-swap is needed
            if transaction.SwapType == "uniswap" && transaction.Status == "ready_for_swap" {
                // Update status to indicate we're preparing to swap
                transaction.Status = "preparing_swap"
                h.DB.Save(&transaction)
                
                // Here you would typically notify the user to complete the swap
                // For a fully automated solution, you would need to handle the swap
                // using the private key, which is not recommended
                
                log.Printf("Transaction %s ready for ETH to CIFO swap", transaction.UUID.String())
            }
        }
    }
    
    c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetTransakOrderStatusHandler gets the status of a Transak order
func (h *Handler) GetTransakOrderStatusHandler(c *gin.Context) {
    // Get order ID from path parameter
    orderID := c.Param("order_id")
    if orderID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
        return
    }
    
    // Get order status from Transak
    transakResp, err := h.TransakService.GetOrderStatus(c.Request.Context(), orderID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order status: " + err.Error()})
        return
    }
    
    // Find our transaction record
    var transaction models.Transaction
    if err := h.DB.Where("payment_reference = ?", orderID).First(&transaction).Error; err != nil {
        // Try to find by payment ID
        if err := h.DB.Where("payment_id = ?", orderID).First(&transaction).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
            return
        }
    }
    
    // Return the status
    c.JSON(http.StatusOK, gin.H{
        "transaction": gin.H{
            "id":             transaction.UUID.String(),
            "payment_id":     transaction.PaymentID,
            "status":         transaction.Status,
            "transak_status": transaction.TransakStatus,
            "created_at":     transaction.CreatedAt,
            "updated_at":     transaction.UpdatedAt,
            "completed_at":   transaction.CompletedAt,
        },
        "transak": gin.H{
            "order_id":         transakResp.Data.ID,
            "status":           transakResp.Data.Status,
            "crypto_amount":    transakResp.Data.CryptoAmount,
            "fiat_amount":      transakResp.Data.FiatAmount,
            "transaction_hash": transakResp.Data.TransactionHash,
            "wallet_address":   transakResp.Data.WalletAddress,
        },
        "next_steps": getNextStepsForTransak(transaction.Status, transaction.SwapType),
    })
}

// Helper function to generate Uniswap URL
func generateUniswapURL(walletAddress string, ethAmount float64, cifoAmount float64) string {
    cifoTokenAddress := "0x1234567890123456789012345678901234567890" // Replace with actual CIFO address
    
    // Format amounts
    ethStr := fmt.Sprintf("%.6f", ethAmount)
    
    // Create Uniswap URL for ETH->CIFO swap
    return fmt.Sprintf("https://app.uniswap.org/#/swap?inputCurrency=ETH&outputCurrency=%s&exactAmount=%s&exactField=input", 
        cifoTokenAddress, ethStr)
}

// Helper function to get next steps based on status
func getNextStepsForTransak(status string, swapType string) []string {
    switch status {
    case models.TransactionStatusPending:
        steps := []string{
            "Complete payment through Transak checkout",
            "Wait for ETH to arrive in your wallet",
        }
        if swapType == "uniswap" {
            steps = append(steps, "Swap ETH for CIFO using Uniswap")
        }
        return steps
    case "transak_initiated":
        steps := []string{
            "Complete payment through Transak checkout",
            "Wait for ETH to arrive in your wallet",
        }
        if swapType == "uniswap" {
            steps = append(steps, "Swap ETH for CIFO using Uniswap")
        }
        return steps
    case "transak_completed":
        if swapType == "uniswap" {
            return []string{
                "ETH has arrived in your wallet",
                "Swap ETH for CIFO using Uniswap",
            }
        }
        return []string{
            "Transaction complete! ETH has arrived in your wallet",
        }
    case "ready_for_swap":
        return []string{
            "ETH has arrived in your wallet",
            "Swap ETH for CIFO using Uniswap",
        }
    case "transak_failed":
        return []string{
            "Transaction failed. Please contact support.",
        }
    default:
        if strings.HasPrefix(status, "transak_") {
            return []string{
                "Transaction is being processed by Transak",
                "Wait for completion",
            }
        }
        return []string{
            "Wait for transaction to complete",
        }
    }
}