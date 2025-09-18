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

type PurchaseCifoHandlerRequest struct {
	TokenAmount       float64 `json:"token_amount"`
	DestinationWallet string  `json:"destination_wallet,omitempty"`
	FiatCurrency      string  `json:"fiat_currency"`
	Name			  string  `json:"name"`
	Email             string  `json:"email"`
	Phone             string  `json:"phone"`
	SuccessURL        string  `json:"success_url"`
	CancelURL         string  `json:"cancel_url"`
}

func (h *Handler) PurchaseCifoHandler(c *gin.Context) {
	var req PurchaseCifoHandlerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }

    if req.DestinationWallet == "" {
        var user models.User
        if err := h.DB.Where("uuid = ?", userID).First(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user wallet"})
            return
        }
        
        if user.WalletAddress == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "No wallet address associated with your account"})
            return
        }
        
        req.DestinationWallet = user.WalletAddress
    }

	    // Validate wallet address
    if !common.IsHexAddress(req.DestinationWallet) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid destination wallet address"})
        return
    }


	// validate the request

	if req.TokenAmount <= 0 {
		c.JSON(400, gin.H{"error": "Token amount must be greater than 0"})
		return
	}

	if req.DestinationWallet == "" {
		c.JSON(400, gin.H{"error": "Destination wallet is required"})
		return
	}

	if !common.IsHexAddress(req.DestinationWallet) {
		c.JSON(400, gin.H{"error": "Invalid destination wallet address"})
		return
	}

	if req.FiatCurrency == "" {
		req.FiatCurrency = "idr" // default to idr if not provided
	}

	if req.SuccessURL == "" {
        req.SuccessURL = "https://yourwebsite.com/payment/success"
    }

    if req.CancelURL == "" {
        req.CancelURL = "https://yourwebsite.com/payment/cancel"
    }

	cifoPerEth, err := h.getCifoAmount(1.0)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get CIFO amount: " + err.Error()})
		return
	}

	cifoPerEthFloat, err := strconv.ParseFloat(cifoPerEth, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse CIFO amount: " + err.Error()})
		return
	}

	ethRequired := req.TokenAmount / cifoPerEthFloat

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

	 // Add 3% slippage buffer to ensure successful swap
	fiatAmount = fiatAmount * 1.03
	
    // Get the required gas deposit from contract (in wei)
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

    txUUID := uuid.New()
    orderID := fmt.Sprintf("CIFO-%s", txUUID.String()[:8])

    transaction := models.Transaction{
        UUID:               txUUID,
        PaymentID:          orderID,
        WalletAddress:      req.DestinationWallet,
        FiatCurrency:       strings.ToUpper(req.FiatCurrency),
        FiatAmount:         fiatAmount,
        EthAmount:          ethRequired,
        TokenAmount:        req.TokenAmount,
        TokenSymbol:        "CIFO",
        Status:             models.TransactionStatusPending,
        PaymentMethod:      "midtrans",
        EthPriceAtPurchase: ethPriceUSD,
        TransactionType:    "purchase", // Explicitly set transaction type to 'purchase'
        SwapType:           "uniswap",  // Use uniswap for purchases
        MinTokenAmount:     req.TokenAmount * 0.97, // 3% slippage allowance
        GasFee:             gasDepositFloat,        // Store gas fee in ETH
        GasFeeFiat:         gasFeeFiat,             // Store gas fee in fiat
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }

	    // Save transaction to database
    if err := h.DB.Create(&transaction).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction record"})
        log.Printf("Error creating transaction: %v", err)
        return
    }

    midtransRequest := FiatToTokenRequest{
        FiatAmount:        fiatAmount,
        FiatCurrency:      req.FiatCurrency,
        DestinationWallet: req.DestinationWallet,
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

    c.JSON(http.StatusOK, gin.H{
        "transaction": gin.H{
            "id":           transaction.UUID.String(),
            "payment_id":   transaction.PaymentID,
            "created_at":   transaction.CreatedAt,
            "token_amount": req.TokenAmount,
            "token_symbol": "CIFO", 
        },
        "payment": paymentResponse,
        "destination_wallet": req.DestinationWallet,
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
    })
}

// Helper function to format currency display
func formatDisplayAmount(amount float64, currency string) string {
    if strings.ToLower(currency) == "usd" {
        return fmt.Sprintf("$%.2f", amount)
    }
    return fmt.Sprintf("Rp %s", formatIDRPrice(amount, false))
}

func formatCurrencyAmount(amount float64, currency string) string {
    if strings.ToLower(currency) == "usd" {
        return fmt.Sprintf("$%.2f", amount)
    }
    return fmt.Sprintf("Rp %s", formatIDRPrice(amount, false))
}