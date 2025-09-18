// token to fiat
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/api/auth"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/config"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

// ExchangeRateResponse represents the response from the Exchange Rate API
type ExchangeRateResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

// Handler struct contains services needed for the handlers
type Handler struct {
	DB *gorm.DB
	PriceService *services.PriceService
	BlockchainService *services.BlockchainService
	Config *config.Config
	TokenService *auth.TokenService// Token service for JWT handling
	RecoveryService *services.RecoveryService // Recovery service for wallet recovery
	TOTPService *auth.TOTPService
	TransakService *services.TransakService

	ActivityLoggerService *services.ActivityLoggerService // Activity logger service
	WalletStorageService *services.WalletStorageService // Wallet storage service
	EncryptionService *services.EncryptionService // Encryption service for secure storage
	// Add cache for exchange rates
	exchangeRateCache     map[string]float64
	exchangeRateLastFetch time.Time
	exchangeRateCacheTTL  time.Duration

	WalletService    *services.WalletService  
    SwapService      *services.SwapService    
}

// NewHandler creates a new Handler instance
func NewHandler(    db *gorm.DB, 
    priceService *services.PriceService, 
    blockchainService *services.BlockchainService, 
    cfg *config.Config, 
    tokenService *auth.TokenService, 
    totpService *auth.TOTPService, 
    recoveryService *services.RecoveryService,
    walletService *services.WalletService,
	transakService *services.TransakService,activityLoggerService *services.ActivityLoggerService, walletStorageService *services.WalletStorageService,
	encryptionService *services.EncryptionService,
    swapService *services.SwapService) *Handler {
    return &Handler{
		DB: 			   db,
        PriceService:          priceService,
        BlockchainService:     blockchainService,
		Config: cfg,
		TokenService: tokenService,
		RecoveryService: recoveryService,
		TOTPService: totpService,
		WalletService: walletService,
		TransakService: transakService,
		ActivityLoggerService: activityLoggerService,
		WalletStorageService: walletStorageService,
        EncryptionService:    encryptionService,
        exchangeRateCache:     make(map[string]float64),
        exchangeRateCacheTTL:  15 * time.Minute, // Cache rates for 15 minutes
        exchangeRateLastFetch: time.Time{},
    }
}

// GetCifoQuoteHandler returns CIFO token quote information
func (h *Handler) GetCifoQuoteHandler(c *gin.Context) {
	// Get ETH amount from query if provided
	ethAmountStr := c.Query("amount")
	var ethAmount float64 = 1.0 // Default to 1 ETH
	
	if ethAmountStr != "" {
		var err error
		ethAmount, err = strconv.ParseFloat(ethAmountStr, 64)
		if err != nil || ethAmount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ETH amount"})
			return
		}
	}

	cifoAmount, err := h.getCifoAmount(ethAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch CIFO data"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"ethAmount": ethAmount,
		"cifoAmount": cifoAmount,
	})
}

// EthToCifoHandler converts ETH to CIFO tokens
func (h *Handler) EthToCifoHandler(c *gin.Context) {
    // Get ETH amount from query
    ethAmountStr := c.Query("amount")
    if ethAmountStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ETH amount is required"})
        return
    }
    
    ethAmount, err := strconv.ParseFloat(ethAmountStr, 64)
    if err != nil || ethAmount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ETH amount"})
        return
    }
    
    // Get the CIFO amount for the ETH
    cifoAmount, err := h.getCifoAmount(ethAmount)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch CIFO rate"})
        return
    }
    
    // Parse CIFO amount for calculations
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CIFO amount"})
        return
    }
    
    // Get ETH price in USD and IDR for reference
    ethPriceUSD, ethPriceIDR, err := h.GetEthPrices(c.Request.Context())
    if err != nil {
        log.Printf("Warning: Failed to get ETH prices: %v", err)
        // Continue anyway as this is just reference info
    }
    
    // Calculate equivalent value in USD and IDR
    totalUSD := ethAmount * ethPriceUSD
    totalIDR := ethAmount * ethPriceIDR
    
    c.JSON(http.StatusOK, gin.H{
        "eth_amount": ethAmount,
        "cifo_amount": cifoAmount,
        "eth_price_usd": ethPriceUSD,
        "eth_price_idr": ethPriceIDR,
        "total_usd": totalUSD,
        "total_idr": totalIDR,
    })
}

// CifoToEthHandler converts CIFO tokens to ETH
func (h *Handler) CifoToEthHandler(c *gin.Context) {
    // Get CIFO amount from query
    cifoAmountStr := c.Query("amount")
    if cifoAmountStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "CIFO amount is required"})
        return
    }
    
    cifoAmount, err := strconv.ParseFloat(cifoAmountStr, 64)
    if err != nil || cifoAmount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CIFO amount"})
        return
    }
    
    // Get current exchange rate by first getting the amount of CIFO for 1 ETH
    cifoPerEth, err := h.getCifoAmount(1.0)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch CIFO rate"})
        return
    }
    
    // Convert CIFO amount to ETH
    cifoPerEthFloat, err := strconv.ParseFloat(cifoPerEth, 64)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid CIFO rate"})
        return
    }
    
    ethAmount := cifoAmount / cifoPerEthFloat
    
    // Get ETH price in USD and IDR for reference
    ethPriceUSD, ethPriceIDR, err := h.GetEthPrices(c.Request.Context())
    if err != nil {
        log.Printf("Warning: Failed to get ETH prices: %v", err)
        // Continue anyway as this is just reference info
    }
    
    // Calculate total values in USD and IDR based on ETH amount
    totalUSD := ethAmount * ethPriceUSD
    totalIDR := ethAmount * ethPriceIDR
    
    c.JSON(http.StatusOK, gin.H{
        "cifo_amount": cifoAmount,
        "eth_amount": ethAmount,
        "eth_price_usd": ethPriceUSD,
        "eth_price_idr": ethPriceIDR,
        "total_usd": totalUSD,
        "total_idr": totalIDR,
    })
}

// CreateOnRampSessionHandler creates a Stripe checkout session for purchasing CIFO tokens
func (h *Handler) CreateOnRampSessionHandler(c *gin.Context) {
	// Set stripe key
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		stripeKey = "sk_test_51RCAH1DAlaHHowO8xmYoaEeQucLbJMxJX8tqm4lZtRUSJN4QkVbQq4nB7jUvCoDJgky5VZprnC52NgVPuGSfJYfi00vdgOep3g"
		log.Println("Warning: Using default Stripe test key. Set STRIPE_SECRET_KEY environment variable for production.")
	}
	stripe.Key = stripeKey

	// Get parameters from request
	var req struct {
		DestinationAddress  string  `json:"destination_address"`
		Email               string  `json:"email"`
		Amount              float64 `json:"amount"`   // Amount in ETH
		Currency            string  `json:"currency"` // Payment currency (USD or IDR)
		SuccessURL          string  `json:"success_url"`
		CancelURL           string  `json:"cancel_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// If no JSON provided, use defaults
		req.DestinationAddress = h.PriceService.GetTokenAddress().Hex()
		req.Amount = 1.0  // Default to 1 ETH
		req.Currency = "usd"
	}

	// Set defaults if not provided
	if req.DestinationAddress == "" {
		req.DestinationAddress = h.PriceService.GetTokenAddress().Hex()
	}
	if req.Currency == "" {
		req.Currency = "usd"
	}
	req.Currency = strings.ToLower(req.Currency)
	if req.Currency != "usd" && req.Currency != "idr" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported currency. Use 'usd' or 'idr'.",
		})
		return
	}

	// Set default URLs if not provided
	if req.SuccessURL == "" {
		req.SuccessURL = "https://yourwebsite.com/success"
	}
	if req.CancelURL == "" {
		req.CancelURL = "https://yourwebsite.com/cancel"
	}

	// Get current ETH prices
	ethPriceUSD, ethPriceIDR, err := h.GetEthPrices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch ETH price",
			"details": err.Error(),
		})
		return
	}

	// Calculate amounts
	var amountInCurrency int64
	var ethAmountToBuy float64 = req.Amount

	if ethAmountToBuy <= 0 {
		ethAmountToBuy = 1.0
	}

	if req.Currency == "usd" {
		amountInCurrency = int64(ethAmountToBuy * ethPriceUSD * 100) // Convert to cents
	} else {
		amountInCurrency = int64(math.Round(ethAmountToBuy * ethPriceIDR * 100))
	}

	// Get CIFO amount for the ETH
	cifoAmount, err := h.getCifoAmount(ethAmountToBuy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch CIFO conversion rate",
			"details": err.Error(),
		})
		return
	}

	// Create product name and description
	productName := "CIFO Token Purchase"
	description := fmt.Sprintf("Buy approximately %s CIFO tokens", cifoAmount)

	// Create line items for the session
	var lineItems = []*stripe.CheckoutSessionLineItemParams{
		{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String(req.Currency),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(productName),
					Description: stripe.String(description),
				},
				UnitAmount: stripe.Int64(amountInCurrency),
			},
			Quantity: stripe.Int64(1),
		},
	}

	// Set up metadata
	isZeroDecimal := "false"
	if req.Currency == "idr" {
		isZeroDecimal = "true"
	}

	metadata := map[string]string{
		"crypto_destination_address": req.DestinationAddress,
		"crypto_currency":            "cifo",
		"crypto_network":             "ethereum",
		"eth_amount":                 fmt.Sprintf("%.4f", ethAmountToBuy),
		"cifo_amount":                cifoAmount,
		"eth_price_usd":              fmt.Sprintf("%.2f", ethPriceUSD),
		"eth_price_idr":              fmt.Sprintf("%.0f", ethPriceIDR),
		"payment_currency":           req.Currency,
		"is_zero_decimal_currency":   isZeroDecimal,
	}

	// Create session parameters
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems:  lineItems,
		Metadata:   metadata,
		SuccessURL: stripe.String(req.SuccessURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(req.CancelURL),
	}

	// Add customer email if provided
	if req.Email != "" {
		params.CustomerEmail = stripe.String(req.Email)
	}

	// Create the session
	s, err := session.New(params)
	if err != nil {
		log.Printf("Error creating checkout session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create checkout session",
			"details": err.Error(),
		})
		return
	}

	// Format display amount
	var displayAmount string
	if req.Currency == "usd" {
		displayAmount = fmt.Sprintf("%.2f", float64(amountInCurrency)/100.0)
	} else {
		displayAmount = formatIDRPrice(float64(amountInCurrency)/100.0, true)
	}

	// Return session info
	c.JSON(http.StatusOK, gin.H{
		"sessionId":    s.ID,
		"url":          s.URL,
		"clientSecret": s.ClientSecret,
		"successUrl":   req.SuccessURL,
		"cancelUrl":    req.CancelURL,
		"destination": gin.H{
			"address":  req.DestinationAddress,
			"currency": "cifo",
			"network":  "ethereum",
		},
		"pricing": gin.H{
			"eth_amount":       ethAmountToBuy,
			"cifo_amount":      cifoAmount,
			"eth_price_usd":    ethPriceUSD,
			"eth_price_idr":    ethPriceIDR,
			"payment_currency": req.Currency,
			"payment_amount":   displayAmount,
		},
	})
}

// Helper function to extract integer part from decimal string
func extractIntegerPart(decimalStr string) string {
	parts := strings.Split(decimalStr, ".")
	return parts[0]
}

// getCifoAmount gets CIFO amount for a specific ETH amount
func (h *Handler) getCifoAmount(ethAmount float64) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/token/0x7C74E84955891dfbdaf465bE3d809F9605f93436?amount=%f", ethAmount))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	
	quote, ok := data["quote"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}
	
	rawQuoteStr, ok := quote["rawQuote"].(string)
	if !ok {
		return "", fmt.Errorf("invalid quote format")
	}
	
	return extractIntegerPart(rawQuoteStr), nil
}

// FetchExchangeRates fetches the latest exchange rates from the API
func (h *Handler) FetchExchangeRates() error {
	// Check if cache is still valid
	if !h.exchangeRateLastFetch.IsZero() && time.Since(h.exchangeRateLastFetch) < h.exchangeRateCacheTTL {
		return nil // Use cached rates
	}

	// Fetch new rates
	resp, err := http.Get("https://api.exchangerate-api.com/v4/latest/USD")
	if err != nil {
		return fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	var rateResp ExchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&rateResp); err != nil {
		return fmt.Errorf("failed to decode exchange rates: %w", err)
	}

	// Update cache
	h.exchangeRateCache = rateResp.Rates
	h.exchangeRateLastFetch = time.Now()

	return nil
}

// GetExchangeRate returns the exchange rate for a given currency against USD
func (h *Handler) GetExchangeRate(currency string) (float64, error) {
	// Ensure we have up-to-date rates
	if err := h.FetchExchangeRates(); err != nil {
		return 0, err
	}

	// Get rate for requested currency
	rate, ok := h.exchangeRateCache[strings.ToUpper(currency)]
	if !ok {
		return 0, fmt.Errorf("exchange rate not available for currency: %s", currency)
	}

	return rate, nil
}


// getEthPrices gets ETH prices in USD and IDR from real-time exchange rates
func (h *Handler) GetEthPrices(ctx context.Context) (usdPrice float64, idrPrice float64, err error) {
	// Call external exchange API for real-time ETH/USD price
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	if err != nil {
		log.Printf("Error fetching ETH/USD rate: %v", err)
		// Fallback to hardcoded value
		return 1810.75, 1810.75 * 15800.0, nil
	}
	defer resp.Body.Close()
	
	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Error decoding ETH/USD rate data: %v", err)
		// Fallback to hardcoded value
		return 1810.75, 1810.75 * 15800.0, nil
	}
	
	// Extract the ETH/USD price
	ethData, exists := data["ethereum"]
	if !exists {
		log.Println("Warning: Could not find ethereum data in API response")
		// Fallback to hardcoded value
		return 1810.75, 1810.75 * 15800.0, nil
	}
	
	usdPrice, exists = ethData["usd"]
	if !exists {
		log.Println("Warning: Could not find USD price in API response")
		// Fallback to hardcoded value
		return 1810.75, 1810.75 * 15800.0, nil
	}
	
	// Get IDR/USD exchange rate from API
	idrRate, err := h.GetExchangeRate("IDR")
	if err != nil {
		log.Printf("Warning: Could not fetch IDR exchange rate: %v", err)
		// Fallback to approximate rate if API call fails
		idrRate = 15800.0
	}
	
	// Convert USD price to IDR using exchange rate
	idrPrice = usdPrice * idrRate
	
	return usdPrice, idrPrice, nil
}
// GetCurrencyExchangeRate gets the exchange rate between two currencies
func (h *Handler) GetCurrencyExchangeRate(ctx context.Context, fromCurrency, toCurrency string) (float64, error) {
	// We'll use USD as the base currency since the API provides rates against USD
	// First, normalize currency codes to uppercase
	fromCurrency = strings.ToUpper(fromCurrency)
	toCurrency = strings.ToUpper(toCurrency)
	
	// If one of the currencies is USD, we can get the rate directly
	if fromCurrency == "USD" {
		return h.GetExchangeRate(toCurrency)
	}
	
	if toCurrency == "USD" {
		rate, err := h.GetExchangeRate(fromCurrency)
		if err != nil {
			return 0, err
		}
		// Invert the rate for USD to fromCurrency
		return 1.0 / rate, nil
	}
	
	// For non-USD pairs, we need to convert via USD
	fromRate, err := h.GetExchangeRate(fromCurrency)
	if err != nil {
		return 0, fmt.Errorf("failed to get exchange rate for %s: %w", fromCurrency, err)
	}
	
	toRate, err := h.GetExchangeRate(toCurrency)
	if err != nil {
		return 0, fmt.Errorf("failed to get exchange rate for %s: %w", toCurrency, err)
	}
	
	// Calculate the cross rate
	return toRate / fromRate, nil
}

// formatIDRPrice formats a number for IDR display
func formatIDRPrice(num float64, showDecimals bool) string {
	// Format with or without decimals
	var formattedNum string
	if showDecimals {
		formattedNum = fmt.Sprintf("%.2f", num)
	} else {
		formattedNum = fmt.Sprintf("%.0f", num)
	}
	
	// Split integer and decimal parts
	decimalPos := strings.Index(formattedNum, ".")
	integerPart := formattedNum
	decimalPart := ""
	
	if decimalPos >= 0 {
		integerPart = formattedNum[:decimalPos]
		decimalPart = formattedNum[decimalPos:]
	}
	
	// Add thousand separators
	var result strings.Builder
	for i, digit := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}
	
	result.WriteString(decimalPart)
	return result.String()
}

// testing purpose only
func (h *Handler) getTokenAmount(ethAmount float64) (string, error) {
	// For test token, conversion rate is 0.01 ETH per token
	// So 1 ETH = 100 tokens
	tokenAmount := ethAmount / 0.01
	return fmt.Sprintf("%.8f", tokenAmount), nil
}

func (h *Handler) createStripeCheckout(transaction models.Transaction, req FiatToTokenRequest) (gin.H, error) {
    // Your Stripe checkout implementation here
    // This is a placeholder for the Stripe payment flow
    
    // Return placeholder response
    return gin.H{
        "transaction_id": transaction.UUID.String(),
        "message": "Stripe checkout not implemented yet",
    }, nil
}