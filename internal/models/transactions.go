package models

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

// Transaction status constants
const (
    TransactionStatusPending   = "pending"
    TransactionStatusCompleted = "completed"
    TransactionStatusFailed    = "failed"
    TransactionStatusRefunded  = "refunded"
    TransactionStatusProcessing = "processing"
    
)

// Transaction represents a token purchase transaction
type Transaction struct {
    UUID               uuid.UUID  `gorm:"primary_key;type:uuid" json:"uuid"`
    UserID             uuid.UUID  `gorm:"type:uuid;index" json:"user_id"` // Explicitly specify UUID type
    PaymentID          string     `gorm:"uniqueIndex;not null" json:"payment_id"`
    WalletAddress      string     `gorm:"not null" json:"wallet_address"`
    FiatCurrency       string     `gorm:"not null" json:"fiat_currency"`
    FiatAmount         float64    `gorm:"not null" json:"fiat_amount"`
    EthAmount          float64    `gorm:"not null" json:"eth_amount"`
    TokenAmount        float64    `gorm:"not null" json:"token_amount"`
    TokenSymbol        string     `gorm:"not null" json:"token_symbol"`
    Status             string     `gorm:"not null" json:"status"`
    PaymentMethod      string     `gorm:"not null" json:"payment_method"`
    PaymentReference   string     `json:"payment_reference,omitempty"`
    BlockchainTxHash   string     `json:"blockchain_tx_hash,omitempty"`
    SwapType           string     `gorm:"default:uniswap"` // Type of swap (e.g., "uniswap")
    SwapTxHash         string     `gorm:"default:null"` // Transaction hash for the swap
    MinTokenAmount     float64    `gorm:"default:0"` 
    TransakStatus     string     `gorm:"default:null"` // Status from Transak
    ErrorMessage       string     `json:"error_message,omitempty"`
    EthPriceAtPurchase float64    `gorm:"not null" json:"eth_price_at_purchase"`
    CreatedAt          time.Time  `json:"created_at"`
    UpdatedAt          time.Time  `json:"updated_at"`
    CompletedAt        *time.Time `json:"completed_at,omitempty"`

    // Transaction type ('purchase' or 'send')
    TransactionType    string     `gorm:"index"`

    // Gas fee information
    GasFee             float64    // Gas fee in ETH
    GasFeeFiat         float64    // Gas fee in fiat currency

    BlockchainRegistered bool  `gorm:"default:false"` // Set to true when created in blockchain
    BlockchainCompleted  bool  `gorm:"default:false"`
}

// TokenAmountInWei converts token amount to wei (with 18 decimals)
func (t *Transaction) TokenAmountInWei() *big.Int {
    // Convert token amount (e.g. 1.5) to wei (e.g. 1500000000000000000)
    tokenAmountBig := new(big.Float).Mul(
        big.NewFloat(t.TokenAmount),
        new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
    )
    tokenAmountInt, _ := tokenAmountBig.Int(nil)
    return tokenAmountInt
}

// FiatAmountInSmallestUnit converts fiat amount to smallest unit (with 2 decimals)
func (t *Transaction) FiatAmountInSmallestUnit() *big.Int {
    // Convert fiat amount (e.g. 100.50) to smallest unit (e.g. 10050)
    fiatAmountBig := new(big.Float).Mul(
        big.NewFloat(t.FiatAmount),
        new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(2), nil)),
    )
    fiatAmountInt, _ := fiatAmountBig.Int(nil)
    return fiatAmountInt
}