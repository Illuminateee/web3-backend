package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"

	"time"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MidtransNotification represents the notification structure from Midtrans
type MidtransNotification struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

// MidtransResponse is the response structure for Midtrans callbacks
type MidtransResponse struct {
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

// MidtransToTransakRequest represents a request to initiate a Midtrans payment that will fund a Transak purchase
type MidtransToTransakRequest struct {
    FiatToCryptoRequest
    // Additional Midtrans-specific fields can be added here
}


// ProcessMidtransWebhookHandler handles payment notifications from Midtrans
func (h *Handler) ProcessMidtransWebhookHandler(c *gin.Context) {
    // Read request body
    bodyBytes, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
        return
    }

    // Parse notification JSON
    var notification MidtransNotification
    if err := json.Unmarshal(bodyBytes, &notification); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification format"})
        return
    }

    log.Printf("Received Midtrans notification: order_id=%s status=%s", 
        notification.OrderID, notification.TransactionStatus)

    // Create context with timeout
    _, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
    defer cancel()

    // Find transaction by payment ID
    var transaction models.Transaction
    if err := h.DB.Where("payment_id = ?", notification.OrderID).First(&transaction).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        }
        return
    }

    switch notification.TransactionStatus {
    case "settlement", "capture":
        // Payment completed successfully
        log.Printf("Payment for %s completed successfully", notification.OrderID)
        
        // Update transaction status
        transaction.Status = models.TransactionStatusProcessing
        now := time.Now()
        transaction.UpdatedAt = now
        h.DB.Save(&transaction)
        
        // Process the transaction based on swap type
        if transaction.SwapType == "uniswap" {
            // Process through Uniswap in background
            go func() {
                err := h.processUniswapPurchase(&transaction)
                if err != nil {
                    log.Printf("Error processing Uniswap purchase: %v", err)
                }
            }()
        } else {
            // Process regular token purchase through blockchain service
            go func() {
                // Prepare transaction data for blockchain service
                transactionData := map[string]interface{}{
                    "destination_wallet": transaction.WalletAddress,
                    "token_amount": transaction.TokenAmountInWei(),
                    "fiat_amount": transaction.FiatAmountInSmallestUnit(),
                }

                err := h.BlockchainService.ProcessPayment(
                    context.Background(), 
                    notification.OrderID, 
                    true,
                    "midtrans",
                    transactionData,
                )
                
                if err != nil {
                    log.Printf("Error processing blockchain payment: %v", err)
                    transaction.Status = models.TransactionStatusFailed
                    transaction.ErrorMessage = fmt.Sprintf("Blockchain error: %v", err)
                    h.DB.Save(&transaction)
                }
            }()
        }
        
    case "pending":
        // Payment is pending - just update status
        transaction.Status = models.TransactionStatusPending
        transaction.UpdatedAt = time.Now()
        h.DB.Save(&transaction)
        
    case "deny", "cancel", "expire", "failure":
        // Payment failed
        transaction.Status = models.TransactionStatusFailed
        transaction.UpdatedAt = time.Now()
        transaction.ErrorMessage = fmt.Sprintf("Payment %s: %s", 
            notification.TransactionStatus, notification.StatusMessage)
        h.DB.Save(&transaction)
    }
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CreateMidtransPaymentHandler creates a new Midtrans payment session
func (h *Handler) CreateMidtransPaymentHandler(c *gin.Context) {
	// Get parameters from request
	var req struct {
		DestinationAddress string  `json:"destination_address"`
		Email              string  `json:"email"`
		Amount             float64 `json:"amount"` // Amount in ETH
		Name               string  `json:"name"`   // Customer name
		Phone              string  `json:"phone"`  // Customer phone
		CallbackURL        string  `json:"callback_url"`
		RedirectURL        string  `json:"redirect_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Set defaults if not provided
	if req.DestinationAddress == "" {
		req.DestinationAddress = h.PriceService.GetTokenAddress().Hex()
	}

	if req.Amount <= 0 {
		req.Amount = 1.0 // Default to 1 ETH
	}

	// Get current ETH prices
	_, ethPriceIDR, err := h.GetEthPrices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch ETH price",
			"details": err.Error(),
		})
		return
	}

	// Calculate IDR amount
	amountInIDR := int64(req.Amount * ethPriceIDR)

	// Get CIFO amount for the ETH
	cifoAmount, err := h.getCifoAmount(req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch CIFO conversion rate",
			"details": err.Error(),
		})
		return
	}

	// Generate a unique order ID
	orderID := fmt.Sprintf("CIFO-%d", time.Now().UnixNano())

	// Create Midtrans payment request
	midtransURL := "https://api.midtrans.com/v2/charge"
	if os.Getenv("APP_ENV") != "production" {
		midtransURL = "https://api.sandbox.midtrans.com/v2/charge"
	}

	// Prepare the payment request payload
	paymentPayload := map[string]interface{}{
		"payment_type": "bank_transfer",
		"transaction_details": map[string]interface{}{
			"order_id":     orderID,
			"gross_amount": amountInIDR,
		},
		"customer_details": map[string]interface{}{
			"email": req.Email,
			"name":  req.Name,
			"phone": req.Phone,
		},
		"item_details": []map[string]interface{}{
			{
				"id":       "CIFO-TOKEN",
				"price":    amountInIDR,
				"quantity": 1,
				"name":     fmt.Sprintf("CIFO Token Purchase (%s tokens)", cifoAmount),
			},
		},
		"bank_transfer": map[string]interface{}{
			"bank": "bca",
		},
		"metadata": map[string]interface{}{
			"crypto_destination_address": req.DestinationAddress,
			"crypto_currency":            "cifo",
			"crypto_network":             "ethereum",
			"eth_amount":                 fmt.Sprintf("%.4f", req.Amount),
			"cifo_amount":                cifoAmount,
			"eth_price_idr":              fmt.Sprintf("%.0f", ethPriceIDR),
		},
	}

	// Get Midtrans server key
	midtransServerKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if midtransServerKey == "" {
		log.Println("Warning: Using test Midtrans server key")
		midtransServerKey = "SB-Mid-server-Replace // Replace with your test key
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(paymentPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create payment request",
		})
		return
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", midtransURL, strings.NewReader(string(payloadBytes)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create payment request",
		})
		return
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.SetBasicAuth(midtransServerKey, "")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to connect to payment gateway",
			"details": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read payment response",
		})
		return
	}

	// Parse response
	var midtransResp map[string]interface{}
	if err := json.Unmarshal(respBody, &midtransResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse payment response",
		})
		return
	}

	// Store payment information in database (not implemented here)
	// This would typically include order ID, destination address, amount, etc.

	// Return payment details to client
	c.JSON(http.StatusOK, gin.H{
		"order_id": orderID,
		"payment":  midtransResp,
		"destination": gin.H{
			"address":  req.DestinationAddress,
			"currency": "cifo",
			"network":  "ethereum",
		},
		"pricing": gin.H{
			"eth_amount":    req.Amount,
			"cifo_amount":   cifoAmount,
			"eth_price_idr": ethPriceIDR,
			"idr_amount":    formatIDRPrice(float64(amountInIDR), false),
		},
	})
}

// Helper function to verify Midtrans signature
func verifyMidtransSignature(notification MidtransNotification, serverKey string) bool {
	// In a real implementation, you would verify the signature here
	// This is a simplified example - refer to Midtrans documentation for actual implementation

	// Example verification:
	// dataToSign := notification.OrderID + notification.StatusCode + notification.GrossAmount + serverKey
	// expectedSignature := sha512.Sum512([]byte(dataToSign))
	// return notification.SignatureKey == hex.EncodeToString(expectedSignature[:])

	// For this example, we'll assume it's valid
	return true
}

// processSuccessfulPayment handles successful payment notifications
func (h *Handler) processSuccessfulPayment(notification MidtransNotification) error {
	log.Printf("Processing successful payment for order %s", notification.OrderID)

	// Create context with timeout for blockchain operations
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Get the transaction details from the database based on order ID
	var transaction models.Transaction
	result := h.DB.Where("payment_id = ?", notification.OrderID).First(&transaction)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Transaction not found for order ID: %s", notification.OrderID)
			// Transaction not found in database, could be a test notification or error
			// Look for metadata in the payment to identify the type of transaction

			// Determine the transaction type from the order ID prefix
			transactionType := ""
			if strings.HasPrefix(notification.OrderID, "TOKEN-") || strings.HasPrefix(notification.OrderID, "FIAT-TO-TOKEN") {
				transactionType = "fiat_to_token"
			} else if strings.HasPrefix(notification.OrderID, "CIFO-") {
				transactionType = "fiat_to_token" // Legacy format - assume fiat-to-token
			} else {
				return fmt.Errorf("unknown transaction type for order ID: %s", notification.OrderID)
			}

			// Process using order ID directly since we don't have DB record
			log.Printf("Processing %s payment without database record", transactionType)
			return h.processPaymentWithoutRecord(ctx, notification, transactionType)
		}

		return fmt.Errorf("database error when fetching transaction: %v", result.Error)
	}

	log.Printf("Found transaction record: ID=%s, Type=%s, Amount=%v, Wallet=%s",
		transaction.UUID, transaction.TokenSymbol, transaction.TokenAmount, transaction.WalletAddress)

	// Update transaction status to completed
	now := time.Now()
	transaction.Status = models.TransactionStatusCompleted
	transaction.CompletedAt = &now
	transaction.UpdatedAt = now

	// Save the updated status first
	if err := h.DB.Save(&transaction).Error; err != nil {
		log.Printf("Error updating transaction status: %v", err)
		// Continue processing anyway as we want to process the blockchain transaction
	}

	// Process the blockchain transaction
	tokenAmountBig := new(big.Float).Mul(
		big.NewFloat(transaction.TokenAmount),
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
	)
	tokenAmountInt, _ := tokenAmountBig.Int(nil)

	// Convert fiat amount to blockchain format for the payment gateway
	fiatAmountBig := new(big.Float).Mul(
		big.NewFloat(transaction.FiatAmount),
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(2), nil)), // 2 decimal places for fiat
	)
	fiatAmountInt, _ := fiatAmountBig.Int(nil)

	paymentID := transaction.PaymentID
	destinationWallet := transaction.WalletAddress
	gateway := "midtrans"

	// Check if payment already exists in the contract
	exists, err := h.BlockchainService.PaymentGateway.CheckPaymentExists(ctx, paymentID)
	if err != nil {
		log.Printf("Error checking if payment exists in contract: %v", err)
		// Continue and try to create/process anyway
	}

	// Get required gas deposit from contract
	gasDeposit, err := h.BlockchainService.PaymentGateway.GetRequiredGasDeposit(ctx)
	if err != nil {
		log.Printf("Failed to get required gas deposit: %v", err)
		// Fallback to default gas deposit value
		gasDeposit = big.NewInt(5000000000000000) // 0.005 ETH from your contract
	}

	if !exists {
		log.Printf("Creating new payment in contract: id=%s, tokens=%s, fiat=%s, wallet=%s",
			paymentID, tokenAmountInt.String(), fiatAmountInt.String(), destinationWallet)
	
		// Create the payment in the contract with destination wallet
		_,err = h.BlockchainService.PaymentGateway.CreatePayment(
			ctx, 
			paymentID, 
			tokenAmountInt, 
			fiatAmountInt, 
			gateway,
			common.HexToAddress(destinationWallet), // Convert to common.Address
			gasDeposit,
		)
	
		if err != nil {
			blockchainError := fmt.Sprintf("Failed to create payment in contract: %v", err)
			log.Println(blockchainError)
	
			// Update transaction with error
			transaction.ErrorMessage = blockchainError
			h.DB.Save(&transaction)
	
			return fmt.Errorf(blockchainError)
		}
	
		log.Printf("Successfully created payment %s in contract", paymentID)
	} else {
		log.Printf("Payment %s already exists in contract", paymentID)
	}

	// Process the payment callback to mark it as successful in the contract
	log.Printf("Processing payment callback for: %s", paymentID)
	txHash,err := h.BlockchainService.PaymentGateway.ProcessPaymentCallback(ctx, paymentID, 1, nil)
	if err != nil {
		blockchainError := fmt.Sprintf("Failed to process payment callback: %v", err)
		log.Println(blockchainError)

		// Update transaction with error
		transaction.ErrorMessage = blockchainError
		h.DB.Save(&transaction)

		return fmt.Errorf(blockchainError)
	}

	// Get transaction hash after successful callback
	// This could be from an event emitted by the smart contract
	// For now, we'll use a placeholder approach
	paymentDetails, err := h.BlockchainService.PaymentGateway.GetPaymentDetails(ctx, paymentID)
	if err == nil && paymentDetails != nil {
		// If we can get payment details, update with blockchain info
		transaction.BlockchainTxHash = "0x" + paymentDetails.Timestamp.String() // Placeholder
	}

	// Save final transaction state
	transaction.ErrorMessage = ""
	transaction.BlockchainCompleted = true
	transaction.BlockchainTxHash = txHash
	h.DB.Save(&transaction)

	log.Printf("Payment successful for order %s, amount: %s",
		notification.OrderID, notification.GrossAmount)

	return nil
}

func (h *Handler) createMidtransSnap(transaction models.Transaction, req FiatToTokenRequest) (gin.H, error) {
	// Calculate amount for payment (in lowest currency unit)
	var amountForPayment int64
	if transaction.FiatCurrency == "IDR" {
		amountForPayment = int64(transaction.FiatAmount)
	} else {
		// Convert USD to IDR for Midtrans (or use the USD amount in cents)
		_, ethPriceIDR, err := h.GetEthPrices(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to fetch ETH prices: %v", err)
		}
		amountForPayment = int64(transaction.FiatAmount * ethPriceIDR / transaction.EthPriceAtPurchase)
	}

	// Create Midtrans Snap payment request
	midtransURL := "https://app.sandbox.midtrans.com/snap/v1/transactions"
	if os.Getenv("APP_ENV") == "production" {
		midtransURL = "https://app.midtrans.com/snap/v1/transactions"
	}

	var itemDetails []map[string]interface{}

    // Main token purchase
    tokenPurchaseItem := map[string]interface{}{
        "id":       "TOKEN-PURCHASE",
        "price":    int64(transaction.FiatAmount),
        "quantity": 1,
        "name":     fmt.Sprintf("Purchase of %.4f CIFO tokens", transaction.TokenAmount),
    }

	    // Gas fee as separate item
    gasFeeItem := map[string]interface{}{
        "id":       "GAS-FEE",
        "price":    int64(transaction.GasFeeFiat),
        "quantity": 1,
        "name":     "Network gas fee for token transaction",
    }

	itemDetails = append(itemDetails, tokenPurchaseItem, gasFeeItem)

    snapPayload := map[string]interface{}{
        "transaction_details": map[string]interface{}{
            "order_id":     transaction.PaymentID,
            "gross_amount": amountForPayment,
        },
        "customer_details": map[string]interface{}{
            "email":      req.Email,
            "first_name": req.Name,
            "phone":      req.Phone,
        },
        "item_details": itemDetails,
        "callbacks": map[string]interface{}{
            "finish": req.SuccessURL,
        },
        "metadata": map[string]interface{}{
            "transaction_id":   transaction.UUID.String(),
            "transaction_type": "fiat_to_token",
            "wallet_address":   transaction.WalletAddress,
            "token_amount":     fmt.Sprintf("%.4f", transaction.TokenAmount),
            "eth_amount":       fmt.Sprintf("%.6f", transaction.EthAmount),
            "eth_price":        fmt.Sprintf("%.2f", transaction.EthPriceAtPurchase),
            "gas_fee_eth":      fmt.Sprintf("%.6f", transaction.GasFee),
            "gas_fee_fiat":     fmt.Sprintf("%.2f", transaction.GasFeeFiat),
        },
    }
	// Get Midtrans server key
	midtransServerKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if midtransServerKey == "" {
		log.Println("Warning: Using test Midtrans server key")
		midtransServerKey = "SB-Mid-server-REPLACE_WITH_YOUR_SERVER_KEY" // Replace with your test key
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(snapPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment request: %v", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", midtransURL, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create payment request: %v", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.SetBasicAuth(midtransServerKey, "")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment gateway: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read payment response: %v", err)
	}

	// Parse response
	var snapResp map[string]interface{}
	if err := json.Unmarshal(respBody, &snapResp); err != nil {
		return nil, fmt.Errorf("failed to parse payment response: %v", err)
	}

	// Update transaction with payment reference
	if token, ok := snapResp["token"].(string); ok {
		h.DB.Model(&transaction).Update("PaymentReference", token)
	}

    // Return payment response
    return gin.H{
        "transaction_id": transaction.UUID.String(),
        "order_id":       transaction.PaymentID,
        "snap_token":     snapResp["token"],
        "snap_url":       snapResp["redirect_url"],
        "destination": gin.H{
            "wallet_address": transaction.WalletAddress,
            "token_symbol":   transaction.TokenSymbol,
        },
        "pricing": gin.H{
            "fiat_amount":   transaction.FiatAmount,
            "fiat_currency": transaction.FiatCurrency,
            "eth_amount":    transaction.EthAmount,
            "token_amount":  transaction.TokenAmount,
            "eth_price":     transaction.EthPriceAtPurchase,
            "gas_fee_eth":   transaction.GasFee,
            "gas_fee_fiat":  transaction.GasFeeFiat,
            "total_fiat":    transaction.FiatAmount + transaction.GasFeeFiat,
        },
    }, nil
}

func (h *Handler) processPaymentWithoutRecord(ctx context.Context, notification MidtransNotification, transactionType string) error {
	// This is a fallback handler for cases where we don't have a database record
	log.Printf("Processing %s payment without database record for order %s", transactionType, notification.OrderID)

	// Extract token amount and wallet address from metadata or order ID
	// For real implementation, you should have this data in your database
	// Here we're using default values for demonstration

	paymentID := notification.OrderID
	fiatAmount, _ := strconv.ParseFloat(notification.GrossAmount, 64)

	// Convert fiat amount to ETH using current rates
	ethPriceUSD, ethPriceIDR, err := h.GetEthPrices(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch ETH prices: %v", err)
	}

	var ethAmount float64
	if notification.Currency == "USD" {
		ethAmount = fiatAmount / ethPriceUSD
	} else {
		// Assume IDR by default
		ethAmount = fiatAmount / ethPriceIDR
	}

	// Get token amount for the ETH amount
	tokenAmount, err := h.getTokenAmount(ethAmount)
	if err != nil {
		return fmt.Errorf("failed to calculate token amount: %v", err)
	}

	tokenAmountFloat, err := strconv.ParseFloat(tokenAmount, 64)
	if err != nil {
		return fmt.Errorf("invalid token amount format: %v", err)
	}

	// Convert to token amount with decimals for blockchain
	tokenAmountBig := new(big.Float).Mul(
		big.NewFloat(tokenAmountFloat),
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
	)
	tokenAmountInt, _ := tokenAmountBig.Int(nil)

	// Convert fiat amount to blockchain format
	fiatAmountBig := new(big.Float).Mul(
		big.NewFloat(fiatAmount),
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(2), nil)),
	)
	fiatAmountInt, _ := fiatAmountBig.Int(nil)

	// Use default wallet address as destination if metadata doesn't have it
	// In production, you should extract this from metadata or order ID
	destinationWallet := h.PriceService.GetTokenAddress().Hex()
	gateway := "midtrans"

	// Get required gas deposit
	gasDeposit, err := h.BlockchainService.PaymentGateway.GetRequiredGasDeposit(ctx)
	if err != nil {
		log.Printf("Failed to get required gas deposit: %v", err)
		// Fallback to default value
		gasDeposit = big.NewInt(5000000000000000) // 0.005 ETH
	}

	log.Printf("Using destination wallet: %s", destinationWallet)

	// Check if payment exists in contract
	exists, err := h.BlockchainService.PaymentGateway.CheckPaymentExists(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("error checking payment existence: %v", err)
	}

	if !exists {
		// Create payment in contract
		log.Printf("Creating payment in contract: id=%s, tokens=%s",
			paymentID, tokenAmountInt.String())
	
		_, err = h.BlockchainService.PaymentGateway.CreatePayment(
			ctx, 
			paymentID, 
			tokenAmountInt, 
			fiatAmountInt, 
			gateway,
			common.HexToAddress(destinationWallet), // Convert string to common.Address
			gasDeposit,
		)
	
		if err != nil {
			return fmt.Errorf("failed to create payment in contract: %v", err)
		}
	}

	// Process payment callback
	log.Printf("Processing payment callback for: %s", paymentID)
	_, err = h.BlockchainService.PaymentGateway.ProcessPaymentCallback(ctx, paymentID, 1, nil)
	if err != nil {
		return fmt.Errorf("failed to process payment callback: %v", err)
	}

	log.Printf("Successfully processed payment %s without database record", paymentID)
	return nil
}

// CreateMidtransToTransakHandler creates a new flow that uses Midtrans to fund Transak
func (h *Handler) CreateMidtransToTransakHandler(c *gin.Context) {
    // Get authenticated user
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }
    
    // Parse request
    var req MidtransToTransakRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }
    
    // Validate request (similar to other handlers)
    if req.FiatAmount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Fiat amount must be greater than 0"})
        return
    }
    
    // Parse user ID
    uid, err := uuid.Parse(userID.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    
    // Create transaction record
    txUUID := uuid.New()
    orderID := fmt.Sprintf("MIDTRANS-TRANSAK-%s", txUUID.String()[:8])
    
	// Determine token symbol and transaction type based on SwapToCifo flag
	tokenSymbol := "ETH"
	transactionType := "fiat_to_eth"
	swapType := ""
	if req.SwapToCifo {
		tokenSymbol = "CIFO"
		transactionType = "fiat_to_cifo"
		swapType = "uniswap"
	}

	transaction := models.Transaction{
		UUID:               txUUID,
		UserID:             uid,
		PaymentID:          orderID,
		WalletAddress:      req.WalletAddress,
		FiatCurrency:       req.FiatCurrency,
		FiatAmount:         req.FiatAmount,
		EthAmount:          0, // Will be updated after Transak purchase
		TokenAmount:        req.CifoAmount,
		TokenSymbol:        tokenSymbol,
		Status:             models.TransactionStatusPending,
		PaymentMethod:      "midtrans_transak",
		TransactionType:    transactionType,
		SwapType:           swapType,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
    
    // Save transaction to database
    if err := h.DB.Create(&transaction).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction record"})
        log.Printf("Error creating transaction: %v", err)
        return
    }
    
    // Create Midtrans payment
    midtransRequest := FiatToTokenRequest{
        FiatAmount:        req.FiatAmount,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Midtrans payment: " + err.Error()})
		return
	}
	
	// Add flow information to the payment response
	paymentResponse["flow"] = gin.H{
		"description": "This flow will use Midtrans for IDR payment, then Transak for ETH purchase",
		"steps": func() []string {
			steps := []string{
				"1. Complete the Midtrans payment",
				"2. Our system will initiate a Transak order for ETH purchase",
				"3. ETH will be delivered to your wallet address",
			}
			if req.SwapToCifo {
				steps = append(steps, "4. Swap ETH for CIFO using Uniswap")
			}
			return steps
		}(),
	}
	
	paymentResponse["destination_wallet"] = req.WalletAddress
	
	c.JSON(http.StatusOK, paymentResponse)
}

// ProcessMidtransTransakWebhookHandler handles successful Midtrans payments for the Midtrans-to-Transak flow
func (h *Handler) ProcessMidtransTransakWebhookHandler(c *gin.Context) {
    // Parse webhook notification
    var notification MidtransNotification
    if err := c.ShouldBindJSON(&notification); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook data"})
        return
    }
    
    // Extract order ID
    orderID := notification.OrderID
    
    // Check transaction status
    if notification.TransactionStatus != "settlement" && 
       notification.TransactionStatus != "capture" {
        // Not a completed payment, no need to process further
        c.JSON(http.StatusOK, gin.H{"status": "received"})
        return
    }
    
    // Find the transaction
    var transaction models.Transaction
    if result := h.DB.Where("payment_id = ?", orderID).First(&transaction); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
        return
    }
    
    // Check if this is a Midtrans-to-Transak flow
    if transaction.PaymentMethod != "midtrans_transak" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Not a Midtrans-to-Transak transaction"})
        return
    }
    
    // Update transaction status
    transaction.Status = "midtrans_completed"
    transaction.UpdatedAt = time.Now()
    if err := h.DB.Save(&transaction).Error; err != nil {
        log.Printf("Error updating transaction status: %v", err)
    }
    
    // Now initiate Transak order - this runs asynchronously
    go func() {
        // Create a background context with timeout
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
        defer cancel()
        
        // Initiate Transak order
        transaction.Status = "initiating_transak"
        h.DB.Save(&transaction)
        
        transakResp, err := h.TransakService.CreateOrder(ctx, &transaction)
        if err != nil {
            log.Printf("Failed to create Transak order for transaction %s: %v", transaction.UUID, err)
            
            transaction.Status = "transak_failed"
            transaction.ErrorMessage = fmt.Sprintf("Failed to create Transak order: %v", err)
            h.DB.Save(&transaction)
            return
        }
        
        // Update transaction with Transak information
        transaction.PaymentReference = transakResp.Data.ID
        transaction.TransakStatus = transakResp.Data.Status
        transaction.Status = "transak_initiated"
        transaction.UpdatedAt = time.Now()
        h.DB.Save(&transaction)
        
        log.Printf("Successfully initiated Transak order %s for transaction %s", 
            transakResp.Data.ID, transaction.UUID)
    }()
    
    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "message": "Payment completed, Transak order will be initiated",
    })
}
