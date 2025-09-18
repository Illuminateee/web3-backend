package models

import (
    "time"

    "github.com/google/uuid"
)

// TransakPayment represents a payment made through Transak
type TransakPayment struct {
    UUID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"uuid"`
    TransactionID     uuid.UUID  `gorm:"type:uuid;index" json:"transaction_id"`
    TransakOrderID    string     `gorm:"index" json:"transak_order_id"`
    TransakStatus     string     `json:"transak_status"`
    WalletAddress     string     `json:"wallet_address"`
    CryptoAmount      float64    `json:"crypto_amount"`
    FiatAmount        float64    `json:"fiat_amount"`
    CryptoCurrency    string     `json:"crypto_currency"`
    FiatCurrency      string     `json:"fiat_currency"`
    TransactionHash   string     `json:"transaction_hash"`
    CheckoutLink      string     `json:"checkout_link"`
    CreatedAt         time.Time  `json:"created_at"`
    UpdatedAt         time.Time  `json:"updated_at"`
}