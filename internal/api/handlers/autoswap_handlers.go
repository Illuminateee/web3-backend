package handlers

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AutoSwapCifoPurchaseHandler handles purchase of CIFO tokens with automatic ETH swap
func (h *Handler) AutoSwapCifoPurchaseHandler(c *gin.Context) {
    // Get the authenticated user
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }
    
    // Parse request
    var req struct {
        CifoAmount       float64 `json:"cifo_amount" binding:"required"`      // Amount of CIFO tokens to buy
        FiatCurrency     string  `json:"fiat_currency" binding:"required"`    // IDR or USD
        WalletAddress    string  `json:"wallet_address" binding:"required"`   // User's wallet address
        Email            string  `json:"email"`
        Name             string  `json:"name"`
        Phone            string  `json:"phone"`
        SuccessURL       string  `json:"success_url"`
        CancelURL        string  `json:"cancel_url"`
    }
    
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
    req.FiatCurrency = strings.ToLower(req.FiatCurrency)
    if req.FiatCurrency != "idr" && req.FiatCurrency != "usd" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Currency must be either 'idr' or 'usd'"})
        return
    }
    
    // Set default URLs if not provided
    if req.SuccessURL == "" {
        req.SuccessURL = "https://yourwebsite.com/payment/success"
    }
    if req.CancelURL == "" {
        req.CancelURL = "https://yourwebsite.com/payment/cancel"
    }
    
    // Step 1: Calculate how much ETH is needed to get the requested CIFO amount
    // First, get CIFO amount for 1 ETH to determine the exchange rate
    cifoPerEth, err := h.getCifoAmount(1.0)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get CIFO amount: " + err.Error()})
        return
    }
    
    cifoPerEthFloat, err := strconv.ParseFloat(cifoPerEth, 64)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CIFO amount: " + err.Error()})
        return
    }
    
    // Calculate ETH required based on CIFO amount requested
    ethRequired := req.CifoAmount / cifoPerEthFloat
    
    // Step 2: Get ETH prices in USD and IDR
    ethPriceUSD, ethPriceIDR, err := h.GetEthPrices(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ETH prices"})
        return
    }
    
    // Step 3: Calculate the fiat amount required
    var fiatAmount float64
    if req.FiatCurrency == "usd" {
        fiatAmount = ethRequired * ethPriceUSD
    } else {
        fiatAmount = ethRequired * ethPriceIDR
    }
    
    // Add slippage buffer (3%) to ensure successful swap
    fiatAmount = fiatAmount * 1.03
    
    // Step 4: Get gas deposit requirement from contract
    gasDeposit, err := h.BlockchainService.PaymentGateway.GetRequiredGasDeposit(context.Background())
    if err != nil {
        log.Printf("Failed to get required gas deposit: %v", err)
        // Fallback to default value
        gasDeposit = big.NewInt(5000000000000000) // 0.005 ETH
    }
    
    // Convert gas deposit to ETH for display
    gasDepositEth := new(big.Float).Quo(
        new(big.Float).SetInt(gasDeposit), 
        new(big.Float).SetInt(big.NewInt(1e18)),
    )
    
    gasDepositFloat, _ := gasDepositEth.Float64()
    
    // Calculate gas fee in fiat currency
    var gasFeeFiat float64
    if req.FiatCurrency == "usd" {
        gasFeeFiat = gasDepositFloat * ethPriceUSD
    } else {
        gasFeeFiat = gasDepositFloat * ethPriceIDR
    }
    
    // Step 5: Create transaction record
    uid, err := uuid.Parse(userID.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    
    txUUID := uuid.New()
    orderID := fmt.Sprintf("CIFO-%s", txUUID.String()[:8])
    
    transaction := models.Transaction{
        UUID:               txUUID,
        UserID:             uid,
        PaymentID:          orderID,
        WalletAddress:      req.WalletAddress,
        FiatCurrency:       strings.ToUpper(req.FiatCurrency),
        FiatAmount:         fiatAmount,
        EthAmount:          ethRequired,
        TokenAmount:        req.CifoAmount,
        TokenSymbol:        "CIFO",
        Status:             models.TransactionStatusPending,
        PaymentMethod:      "midtrans",
        EthPriceAtPurchase: ethPriceUSD,
        TransactionType:    "auto_swap",
        SwapType:           "uniswap",
        MinTokenAmount:     req.CifoAmount * 0.97, // 3% slippage allowance
        GasFee:             gasDepositFloat,
        GasFeeFiat:         gasFeeFiat,
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }
    
    // Save transaction to database
    if err := h.DB.Create(&transaction).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction record"})
        log.Printf("Error creating transaction: %v", err)
        return
    }
    
    // Step 6: Create Midtrans payment
    midtransRequest := FiatToTokenRequest{
        FiatAmount:        fiatAmount,
        FiatCurrency:      req.FiatCurrency,
        DestinationWallet: req.WalletAddress,
        PaymentMethod:     "midtrans",
        Email:             req.Email,
        Name:              req.Name,
        Phone:             req.Phone,
        SuccessURL:        req.SuccessURL,
        CancelURL:         req.CancelURL,
    }
    
    paymentResponse, err := h.createMidtransSnap(transaction, midtransRequest)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment: " + err.Error()})
        return
    }
    
    // Step 7: Return response with all details
    c.JSON(http.StatusOK, gin.H{
        "transaction": gin.H{
            "id":              transaction.UUID.String(),
            "payment_id":      transaction.PaymentID,
            "created_at":      transaction.CreatedAt,
            "cifo_amount":     req.CifoAmount,
            "eth_required":    ethRequired,
            "status":          transaction.Status,
        },
        "payment": paymentResponse,
        "destination_wallet": req.WalletAddress,
        "costs": gin.H{
            "eth_required":     ethRequired,
            "fiat_amount":      fiatAmount,
            "fiat_currency":    strings.ToUpper(req.FiatCurrency),
            "gas_fee_eth":      gasDepositFloat,
            "gas_fee_fiat":     formatCurrencyAmount(gasFeeFiat, req.FiatCurrency),
            "total_fiat":       formatCurrencyAmount(fiatAmount + gasFeeFiat, req.FiatCurrency),
            "eth_price":        func() float64 { if req.FiatCurrency == "usd" { return ethPriceUSD } else { return ethPriceIDR } }(),
            "slippage_buffer":  "3%",
        },
        "exchange_rate": gin.H{
            "cifo_per_eth": cifoPerEthFloat,
            "eth_per_cifo": 1/cifoPerEthFloat,
        },
    })
}

// ProcessAutoSwapWebhookHandler processes Midtrans webhook for auto-swap transactions
func (h *Handler) ProcessAutoSwapWebhookHandler(c *gin.Context) {
    // This function should be called by Midtrans when a payment is completed

    // Parse the webhook notification
    var notification map[string]interface{}
    if err := c.ShouldBindJSON(&notification); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook data"})
        return
    }
    
    // Extract transaction info
    orderID, ok := notification["order_id"].(string)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing order_id in notification"})
        return
    }
    
    transactionStatus, ok := notification["transaction_status"].(string)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing transaction_status in notification"})
        return
    }
    
    // Find the transaction in our database
    var transaction models.Transaction
    if err := h.DB.Where("payment_id = ?", orderID).First(&transaction).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
        return
    }
    
    // Verify this is an auto_swap transaction
    if transaction.TransactionType != "auto_swap" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Not an auto_swap transaction"})
        return
    }
    
    // Process the payment status
    if transactionStatus == "settlement" || transactionStatus == "capture" {
        // Payment is successful, proceed with the swap
        
        // Update transaction status to processing
        transaction.Status = models.TransactionStatusProcessing
        transaction.UpdatedAt = time.Now()
        h.DB.Save(&transaction)
        
        // Execute the swap asynchronously
        go h.executeAutoSwap(&transaction)
        
        c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Payment successful, swap initiated"})
        return
    } else if transactionStatus == "cancel" || transactionStatus == "deny" || transactionStatus == "expire" {
        // Payment failed
        transaction.Status = models.TransactionStatusFailed
        transaction.ErrorMessage = "Payment " + transactionStatus
        transaction.UpdatedAt = time.Now()
        h.DB.Save(&transaction)
        
        c.JSON(http.StatusOK, gin.H{"status": "failed", "message": "Payment failed: " + transactionStatus})
        return
    } else {
        // Still pending or other status
        c.JSON(http.StatusOK, gin.H{"status": "pending", "message": "Payment status: " + transactionStatus})
        return
    }
}

// executeAutoSwap handles the actual Uniswap swap process
func (h *Handler) executeAutoSwap(transaction *models.Transaction) {
    ctx := context.Background()
    
    // Step 1: Log that we're starting the swap
    log.Printf("Starting auto-swap for transaction %s, CIFO amount: %.8f", 
               transaction.UUID.String(), transaction.TokenAmount)
    
    // Step 2: Parse destination wallet address
    destWallet := common.HexToAddress(transaction.WalletAddress)
    
    // Step 3: Prepare ETH amount for the swap
    ethAmountFloat := new(big.Float).SetFloat64(transaction.EthAmount)
    ethAmountInt, _ := ethAmountFloat.Mul(ethAmountFloat, big.NewFloat(1e18)).Int(nil) // Convert to wei
    
    // Step 4: Get token address - this is the CIFO token
    cifoTokenAddress := common.HexToAddress("0x7c74e84955891dfbdaf465be3d809f9605f93436")
    
    // Step 5: Get Uniswap client
    uniswap := h.PriceService.GetUniswapClient()
    if uniswap == nil {
        transaction.Status = models.TransactionStatusFailed
        transaction.ErrorMessage = "Uniswap client not initialized"
        h.DB.Save(transaction)
        log.Printf("Error: Uniswap client not initialized for transaction %s", transaction.UUID.String())
        return
    }

    // Step 6: Parse min tokens amount (with slippage protection)
    minTokensFloat := new(big.Float).SetFloat64(transaction.MinTokenAmount)
    minTokensInt, _ := minTokensFloat.Mul(minTokensFloat, big.NewFloat(1e18)).Int(nil) // Convert to wei
    
    // Step 7: Prepare swap path (ETH -> CIFO)
    path := []common.Address{
        uniswap.Config.WethAddress, // WETH address
        cifoTokenAddress,           // CIFO token address
    }
    
    log.Printf("Executing Uniswap swap: %.8f ETH -> min %.8f CIFO tokens to %s",
        transaction.EthAmount, transaction.MinTokenAmount, destWallet.Hex())
    
    // Step 8: Execute the swap
    txHash, err := uniswap.SwapExactETHForTokens(
        ctx,
        ethAmountInt,
        minTokensInt,
        path,
        destWallet,
        time.Now().Add(20*time.Minute).Unix(), // 20 min deadline
    )

    if err != nil {
        transaction.Status = models.TransactionStatusFailed
        transaction.ErrorMessage = fmt.Sprintf("Uniswap swap failed: %v", err)
        h.DB.Save(transaction)
        log.Printf("Error: Failed to execute Uniswap swap for transaction %s: %v", 
                  transaction.UUID.String(), err)
        return
    }
    
    // Step 9: Update transaction with swap details
    now := time.Now()
    transaction.Status = models.TransactionStatusCompleted
    transaction.BlockchainTxHash = txHash
    transaction.SwapTxHash = txHash
    transaction.BlockchainCompleted = true
    transaction.CompletedAt = &now
    transaction.UpdatedAt = now
    
    if err := h.DB.Save(transaction).Error; err != nil {
        log.Printf("Warning: Failed to update transaction after successful swap: %v", err)
    }
    
    log.Printf("Successfully completed auto-swap of %.8f CIFO tokens via Uniswap, tx: %s",
        transaction.TokenAmount, txHash)
}