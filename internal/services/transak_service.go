package services

import (
    "bytes"
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "time"

    "git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type TransakService struct {
    DB           *gorm.DB
    APIKey       string
    APISecret    string
    Environment  string
    WebhookURL   string
    RedirectURL  string
    PartnerID    string
    isInitialized bool
}

// TransakOrderResponse represents the response from Transak order creation API
type TransakOrderResponse struct {
    Status  string                 `json:"status"`
    Data    TransakOrderData       `json:"data"`
    Error   bool                   `json:"error"`
    Message string                 `json:"message"`
    Meta    map[string]interface{} `json:"meta,omitempty"`
}

// TransakOrderData contains the order details
type TransakOrderData struct {
    ID                 string  `json:"id"`
    WalletAddress      string  `json:"walletAddress"`
    CryptoAmount       float64 `json:"cryptoAmount"`
    FiatAmount         float64 `json:"fiatAmount"`
    CryptoCurrency     string  `json:"cryptoCurrency"`
    FiatCurrency       string  `json:"fiatCurrency"`
    Network            string  `json:"network"`
    PaymentMethod      string  `json:"paymentMethod"`
    Status             string  `json:"status"`
    WebhookURL         string  `json:"webhookUrl,omitempty"`
    RedirectURL        string  `json:"redirectUrl,omitempty"`
    TransactionHash    string  `json:"transactionHash,omitempty"`
    TransactionLink    string  `json:"transactionLink,omitempty"`
    CheckoutLink       string  `json:"checkoutLink,omitempty"`
    PartnerOrderID     string  `json:"partnerOrderId,omitempty"`
    CreatedAt          string  `json:"createdAt"`
    UpdatedAt          string  `json:"updatedAt"`
    CompletedAt        string  `json:"completedAt,omitempty"`
}

// TransakWebhookPayload represents the webhook payload sent by Transak
type TransakWebhookPayload struct {
    ID                 string  `json:"id"`
    OrderID            string  `json:"orderId"`
    Status             string  `json:"status"`
    FiatCurrency       string  `json:"fiatCurrency"`
    CryptoCurrency     string  `json:"cryptoCurrency"`
    FiatAmount         float64 `json:"fiatAmount"`
    CryptoAmount       float64 `json:"cryptoAmount"`
    WalletAddress      string  `json:"walletAddress"`
    TransactionHash    string  `json:"transactionHash"`
    TransactionLink    string  `json:"transactionLink"`
    Network            string  `json:"network"`
    PartnerOrderID     string  `json:"partnerOrderId"`
    Signature          string  `json:"signature"`
}

// NewTransakService creates a new Transak service
func NewTransakService(db *gorm.DB) *TransakService {
    apiKey := os.Getenv("TRANSAK_API_KEY")
    apiSecret := os.Getenv("TRANSAK_API_SECRET")
    env := os.Getenv("APP_ENV")
    transakEnv := "STAGING"

    if env == "production" {
        transakEnv = "PRODUCTION"
    }

    baseUrl := os.Getenv("API_BASE_URL")
    if baseUrl == "" {
        baseUrl = "http://localhost:8080"
    }

    service := &TransakService{
        DB:           db,
        APIKey:       apiKey,
        APISecret:    apiSecret,
        Environment:  transakEnv,
        WebhookURL:   fmt.Sprintf("%s/api/v1/payment/transak/webhook", baseUrl),
        RedirectURL:  fmt.Sprintf("%s/payment/success", os.Getenv("APP_URL")),
        PartnerID:    "cifo",
        isInitialized: apiKey != "" && apiSecret != "",
    }
    
    return service
}

// IsInitialized returns whether the service is properly initialized
func (s *TransakService) IsInitialized() bool {
    return s.isInitialized
}

// CreateOrder creates a new order for buying ETH with fiat
func (s *TransakService) CreateOrder(ctx context.Context, transaction *models.Transaction) (*TransakOrderResponse, error) {
    if !s.IsInitialized() {
        return nil, fmt.Errorf("transak service not initialized, check API keys")
    }

    apiURL := "https://staging-api.transak.com/api/v2/order/create"
    if s.Environment == "PRODUCTION" {
        apiURL = "https://api.transak.com/api/v2/order/create"
    }

    // Prepare payload for Transak
    payload := map[string]interface{}{
        "partnerApiKey":   s.APIKey,
        "partnerOrderId":  transaction.PaymentID,
        "walletAddress":   transaction.WalletAddress,
        "fiatCurrency":    transaction.FiatCurrency,
        "cryptoCurrency":  "ETH", // Always buying ETH first
        "fiatAmount":      transaction.FiatAmount,
        "network":         "ethereum",
        "paymentMethod":   "credit_debit_card", // Default to card, can be customized
        "defaultPaymentMethod": "credit_debit_card",
        "disablePaymentMethods": "bank_transfer,gbp_bank_transfer",
        "webhookUrl":      s.WebhookURL,
        "redirectUrl":     s.RedirectURL,
        "partnerCustomerId": transaction.UserID.String(),
        "exchangeScreenTitle": "Buy ETH with " + transaction.FiatCurrency,
    }

    // Convert payload to JSON
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal payload: %w", err)
    }

    // Create HTTP request
    req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Add headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    
    // Add API key header
    req.Header.Set("x-api-key", s.APIKey)

    // Execute request
    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    // Read response body
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    // Parse response
    var transakResp TransakOrderResponse
    if err := json.Unmarshal(respBody, &transakResp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w, body: %s", err, string(respBody))
    }

    // Check if order creation was successful
    if transakResp.Error || resp.StatusCode >= 400 {
        return nil, fmt.Errorf("failed to create transak order: %s, status: %d", transakResp.Message, resp.StatusCode)
    }

    // Update transaction with Transak reference
    transaction.PaymentReference = transakResp.Data.ID
    transaction.TransakStatus = transakResp.Data.Status
    transaction.Status = "transak_initiated"
    transaction.UpdatedAt = time.Now()

    // Create Transak payment record
    transakPayment := models.TransakPayment{
        UUID:              uuid.New(),
        TransactionID:     transaction.UUID,
        TransakOrderID:    transakResp.Data.ID,
        TransakStatus:     transakResp.Data.Status,
        WalletAddress:     transaction.WalletAddress,
        CryptoAmount:      transakResp.Data.CryptoAmount,
        FiatAmount:        transakResp.Data.FiatAmount,
        CryptoCurrency:    "ETH",
        FiatCurrency:      transaction.FiatCurrency,
        TransactionHash:   "",
        CheckoutLink:      transakResp.Data.CheckoutLink,
        CreatedAt:         time.Now(),
        UpdatedAt:         time.Now(),
    }

    // Save both records in a transaction
    err = s.DB.Transaction(func(tx *gorm.DB) error {
        if err := tx.Save(transaction).Error; err != nil {
            return fmt.Errorf("failed to update transaction: %w", err)
        }
        if err := tx.Create(&transakPayment).Error; err != nil {
            return fmt.Errorf("failed to create transak payment record: %w", err)
        }
        return nil
    })

    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }

    return &transakResp, nil
}

// VerifyWebhookSignature verifies the signature of a Transak webhook payload
func (s *TransakService) VerifyWebhookSignature(payload *TransakWebhookPayload) bool {
    if !s.IsInitialized() {
        return false
    }

    providedSignature := payload.Signature
    
    // Remove signature from the payload for verification
    payload.Signature = ""
    
    // Convert payload to JSON
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return false
    }
    
    // Create HMAC
    h := hmac.New(sha256.New, []byte(s.APISecret))
    h.Write(payloadBytes)
    calculatedSignature := hex.EncodeToString(h.Sum(nil))
    
    // Restore signature to payload
    payload.Signature = providedSignature
    
    // Compare signatures
    return hmac.Equal([]byte(providedSignature), []byte(calculatedSignature))
}

// ProcessWebhook processes a webhook notification from Transak
func (s *TransakService) ProcessWebhook(payload *TransakWebhookPayload) error {
    if !s.VerifyWebhookSignature(payload) {
        return fmt.Errorf("invalid webhook signature")
    }

    // Find the transaction by Transak order ID
    var transakPayment models.TransakPayment
    if err := s.DB.Where("transak_order_id = ?", payload.OrderID).First(&transakPayment).Error; err != nil {
        // Try finding by partner order ID if not found by Transak order ID
        if err := s.DB.Where("transak_order_id = ?", payload.PartnerOrderID).First(&transakPayment).Error; err != nil {
            return fmt.Errorf("transak payment not found: %w", err)
        }
    }

    // Find the main transaction
    var transaction models.Transaction
    if err := s.DB.Where("uuid = ?", transakPayment.TransactionID).First(&transaction).Error; err != nil {
        return fmt.Errorf("transaction not found: %w", err)
    }

    // Update Transak payment record
    transakPayment.TransakStatus = payload.Status
    transakPayment.CryptoAmount = payload.CryptoAmount
    transakPayment.TransactionHash = payload.TransactionHash
    transakPayment.UpdatedAt = time.Now()

    // Update transaction status based on Transak status
    var txStatus string
    switch payload.Status {
    case "COMPLETED":
        txStatus = "transak_completed"
        // ETH has been sent to user's wallet, now we can proceed with swap if needed
        if transaction.SwapType == "uniswap" {
            txStatus = "ready_for_swap"
        }
    case "FAILED":
        txStatus = "transak_failed"
    case "CANCELLED":
        txStatus = "transak_cancelled"
    case "REFUNDED":
        txStatus = "transak_refunded"
    default:
        txStatus = "transak_" + payload.Status
    }

    // Update transaction
    transaction.Status = txStatus
    transaction.TransakStatus = payload.Status
    transaction.EthAmount = payload.CryptoAmount // Update with actual ETH amount
    transaction.BlockchainTxHash = payload.TransactionHash
    transaction.UpdatedAt = time.Now()

    if payload.Status == "COMPLETED" {
        transaction.BlockchainCompleted = true
        now := time.Now()
        transaction.CompletedAt = &now
    }

    // Save both records in a transaction
    return s.DB.Transaction(func(tx *gorm.DB) error {
        if err := tx.Save(&transakPayment).Error; err != nil {
            return fmt.Errorf("failed to update transak payment: %w", err)
        }
        if err := tx.Save(&transaction).Error; err != nil {
            return fmt.Errorf("failed to update transaction: %w", err)
        }
        return nil
    })
}

// GetOrderStatus gets the status of an order from Transak
func (s *TransakService) GetOrderStatus(ctx context.Context, orderID string) (*TransakOrderResponse, error) {
    if !s.IsInitialized() {
        return nil, fmt.Errorf("transak service not initialized, check API keys")
    }

    baseURL := "https://staging-api.transak.com/api/v2/order/"
    if s.Environment == "PRODUCTION" {
        baseURL = "https://api.transak.com/api/v2/order/"
    }
    
    apiURL := baseURL + url.PathEscape(orderID)

    // Create HTTP request
    req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Add headers
    req.Header.Set("Accept", "application/json")
    
    // Add API key header
    req.Header.Set("x-api-key", s.APIKey)

    // Execute request
    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    // Read response body
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    // Parse response
    var transakResp TransakOrderResponse
    if err := json.Unmarshal(respBody, &transakResp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w, body: %s", err, string(respBody))
    }

    // Check if response is successful
    if transakResp.Error || resp.StatusCode >= 400 {
        return nil, fmt.Errorf("failed to get transak order: %s, status: %d", transakResp.Message, resp.StatusCode)
    }

    return &transakResp, nil
}