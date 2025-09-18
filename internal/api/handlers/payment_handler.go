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

// FiatToTokenRequest represents a request to convert fiat to tokens
type FiatToTokenRequest struct {
	FiatAmount        float64 `json:"fiat_amount"`
	FiatCurrency      string  `json:"fiat_currency"`
	DestinationWallet string  `json:"destination_wallet"`
	PaymentMethod     string  `json:"payment_method"`
	Email             string  `json:"email"`
	Name              string  `json:"name"`
	Phone             string  `json:"phone"`
	SuccessURL        string  `json:"success_url"`
	CancelURL         string  `json:"cancel_url"`
}

func (h *Handler) CreateFiatToTokenPaymentHandler(c *gin.Context) {
    var req FiatToTokenRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Amount and currency required."})
        return
    }

    // Get authenticated user from context
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }

    // Convert userID to UUID
    userUUID, err := uuid.Parse(userID.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    if req.FiatAmount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Fiat amount must be greater than 0"})
        return
    }

    if req.FiatCurrency == "" {
        req.FiatCurrency = "idr" // Default to IDR
    }

    req.FiatCurrency = strings.ToLower(req.FiatCurrency)
    if req.FiatCurrency != "idr" && req.FiatCurrency != "usd" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Currency must be either 'idr' or 'usd'"})
        return
    }

    if req.DestinationWallet == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Destination wallet address is required"})
        return
    }

    if req.PaymentMethod == "" {
        req.PaymentMethod = "midtrans" // Default to Midtrans
    }

    req.PaymentMethod = strings.ToLower(req.PaymentMethod)
    if req.PaymentMethod != "midtrans" && req.PaymentMethod != "stripe" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Payment method must be either 'midtrans' or 'stripe'"})
        return
    }

    // Set default success and cancel URLs if not provided
    if req.SuccessURL == "" {
        req.SuccessURL = "https://yourwebsite.com/payment/success"
    }
    if req.CancelURL == "" {
        req.CancelURL = "https://yourwebsite.com/payment/cancel"
    }

    // Get ETH prices
    ethPriceUSD, ethPriceIDR, err := h.GetEthPrices(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ETH prices"})
        return
    }

    // Calculate ETH amount
    ethPrice := ethPriceUSD
    if req.FiatCurrency == "idr" {
        ethPrice = ethPriceIDR
    }
    ethAmount := req.FiatAmount / ethPrice

    // Calculate token amount
    tokenAmount, err := h.getTokenAmount(ethAmount)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate token amount"})
        return
    }

    // Parse token amount as float
    tokenAmountFloat, err := strconv.ParseFloat(tokenAmount, 64)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token amount"})
        return
    }

    gasDeposit, err := h.BlockchainService.PaymentGateway.GetRequiredGasDeposit(context.Background())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gas deposit"})
         // Fallback to default value
        gasDeposit = big.NewInt(5000000000000000) // 0.005 ETH
    }

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

    // Generate unique transaction ID
    txUUID := uuid.New()
    orderID := fmt.Sprintf("TOKEN-%s", txUUID.String()[:8])

    // Create transaction record in database
    transaction := models.Transaction{
        UUID:               txUUID,
        UserID:             userUUID,
        PaymentID:          orderID,
        WalletAddress:      req.DestinationWallet,
        FiatCurrency:       strings.ToUpper(req.FiatCurrency),
        FiatAmount:         req.FiatAmount,
        EthAmount:          ethAmount,
        TokenAmount:        tokenAmountFloat,
        TokenSymbol:        "CIFO", // Using CIFO token
        Status:             models.TransactionStatusPending,
        PaymentMethod:      req.PaymentMethod,
        EthPriceAtPurchase: ethPrice,
        TransactionType:    "send", // Explicitly set transaction type to 'send'
        SwapType:           "",     // No swap type for sending
        GasFee:             gasDepositFloat, // Store gas fee in ETH
        GasFeeFiat:         gasFeeFiat,      // Store gas fee in fiat
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }

    // Save transaction to database
    if err := h.DB.Create(&transaction).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction record"})
        log.Printf("Error creating transaction: %v", err)
        return
    }

    go func() {
        if err := h.createBlockchainPayment(&transaction); err != nil {
            log.Printf("Error creating blockchain payment for %s: %v", transaction.PaymentID, err)
            
            // Update transaction with error information
            transaction.ErrorMessage = fmt.Sprintf("Blockchain registration failed: %v", err)
            h.DB.Save(&transaction)
        } else {
            // Update transaction with blockchain registration status
            transaction.BlockchainRegistered = true
            h.DB.Save(&transaction)
        }
	}()

    // Process payment based on payment method
    if req.PaymentMethod == "midtrans" {
        paymentResponse, err := h.createMidtransSnap(transaction, req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment: " + err.Error()})
            return
        }
        c.JSON(http.StatusOK, paymentResponse)
    } else if req.PaymentMethod == "stripe" {
        paymentResponse, err := h.createStripeCheckout(transaction, req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment: " + err.Error()})
            return
        }
        c.JSON(http.StatusOK, paymentResponse)
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported payment method"})
        return
    }
}

func (h *Handler) processTokenPurchase(transaction *models.Transaction) error {
    ctx := context.Background()
    
    // Convert token amount to blockchain format with correct decimals (assuming 18 decimals)
    tokenAmountBig := new(big.Float).Mul(
        big.NewFloat(transaction.TokenAmount),
        new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
    )
    tokenAmountInt, _ := tokenAmountBig.Int(nil)
    
    // Convert fiat amount to blockchain format
    fiatAmountBig := new(big.Float).Mul(
        big.NewFloat(transaction.FiatAmount),
        new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
    )
    fiatAmountInt, _ := fiatAmountBig.Int(nil)
    
    // Create payment on blockchain
    paymentID := transaction.PaymentID
    gateway := strings.ToLower(transaction.PaymentMethod)
    gasDeposit := big.NewInt(5000000000000000) // 0.005 ETH from your contract
    
    // Check if payment exists
	
    exists, err := h.BlockchainService.PaymentGateway.CheckPaymentExists(ctx, paymentID)
    if err != nil {
        return fmt.Errorf("error checking payment existence: %v", err)
    }
    
    if !exists {
		// Create payment in contract
		_, err = h.BlockchainService.PaymentGateway.CreatePayment(
			ctx, 
			paymentID, 
			tokenAmountInt, 
			fiatAmountInt, 
			gateway,
			common.HexToAddress(transaction.WalletAddress), // Add destination wallet
			gasDeposit,
		)
		if err != nil {
			return fmt.Errorf("failed to create payment in contract: %v", err)
		}
	}
    
    // Process payment in contract
    txHash, err := h.BlockchainService.PaymentGateway.ProcessPaymentCallback(ctx, paymentID, 1, nil) // 1 = completed status
    if err != nil {
        return fmt.Errorf("failed to process payment callback: %v", err)
    }
    
    // Get transaction hash
    // In a real implementation, you'd get this from the blockchain transaction receipt
    // For now, we'll just use a placeholder
    transaction.BlockchainTxHash = txHash
    
    return nil
}

func (h *Handler) createBlockchainPayment(transaction *models.Transaction) error {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()
    
    // Convert token amount to blockchain format with 18 decimals
    tokenAmountBig := new(big.Float).Mul(
        big.NewFloat(transaction.TokenAmount),
        new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
    )
    tokenAmountInt, _ := tokenAmountBig.Int(nil)
    
    // Convert fiat amount to blockchain format - use 18 decimals for consistency
    fiatAmountBig := new(big.Float).Mul(
        big.NewFloat(transaction.FiatAmount),
        new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
    )
    fiatAmountInt, _ := fiatAmountBig.Int(nil)
    
    paymentID := transaction.PaymentID
    gateway := strings.ToLower(transaction.PaymentMethod)
    destinationWallet := common.HexToAddress(transaction.WalletAddress)
    
    // Get required gas deposit from contract
    gasDeposit, err := h.BlockchainService.PaymentGateway.GetRequiredGasDeposit(ctx)
    if err != nil {
        log.Printf("Failed to get required gas deposit: %v", err)
        // Fallback to default gas deposit value
        gasDeposit = big.NewInt(5000000000000000) // 0.005 ETH
    }
    
    // Check if payment already exists first
    exists, err := h.BlockchainService.PaymentGateway.CheckPaymentExists(ctx, paymentID)
    if err != nil {
        return fmt.Errorf("error checking payment existence: %v", err)
    }
    
    if exists {
        log.Printf("Payment %s already exists in blockchain, skipping creation", paymentID)
        return nil
    }
    
    // Create payment in contract (this reserves the tokens but doesn't transfer them)
    log.Printf("Creating blockchain payment: id=%s, tokens=%s, wallet=%s", 
        paymentID, tokenAmountInt.String(), destinationWallet.Hex())
    
    txHash, err := h.BlockchainService.PaymentGateway.CreatePayment(
        ctx,
        paymentID,
        tokenAmountInt,
        fiatAmountInt,
        gateway,
        destinationWallet,
        gasDeposit,
    )
    
    if err != nil {
        return fmt.Errorf("failed to create payment in blockchain: %v", err)
    }
    
    log.Printf("Successfully created payment %s in blockchain", paymentID)
    transaction.BlockchainTxHash = txHash
    return nil
}