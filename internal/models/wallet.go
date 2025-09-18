package models

import (
    "time"

    "github.com/google/uuid"
)

// Wallet represents a user's blockchain wallet
type Wallet struct {
    UUID          uuid.UUID `gorm:"primary_key;type:uuid"`
    UserID        uuid.UUID `gorm:"index;not null"`
    WalletAddress string    `gorm:"not null"`
    CreatedAt     time.Time
    UpdatedAt     time.Time

    // Relationships
    User User `gorm:"foreignKey:UserID"`
}

// WalletTransaction represents a transaction initiated by the user
type WalletTransaction struct {
    UUID          uuid.UUID `gorm:"primary_key;type:uuid"`
    UserID        uuid.UUID `gorm:"index;not null"`
    WalletAddress string    `gorm:"not null"`
    TxHash        string    `gorm:"index"`
    TxType        string    `gorm:"not null"` // swap, transfer, etc.
    Amount        string    `gorm:"not null"`
    TokenSymbol   string    `gorm:"not null"`
    ToAddress     string    `gorm:"index"`
    Status        string    `gorm:"not null"` // pending, confirmed, failed
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type EncryptedWallet struct {
    UUID           uuid.UUID `gorm:"primary_key;type:uuid"`
    UserUUID       uuid.UUID `gorm:"index;not null"` // Reference to user in main DB
    WalletAddress  string    `gorm:"uniqueIndex;not null"`
    EncPrivateKey  []byte    `gorm:"type:bytea"` // Encrypted private key
    EncKeystoreJSON []byte   `gorm:"type:bytea"` // Encrypted keystore JSON
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

// WalletMnemonic stores encrypted mnemonic phrases
type WalletMnemonic struct {
    UUID           uuid.UUID `gorm:"primary_key;type:uuid"`
    UserUUID       uuid.UUID `gorm:"index;not null"`
    WalletAddress  string    `gorm:"uniqueIndex;not null"`
    EncMnemonic    []byte    `gorm:"type:bytea"` // Encrypted mnemonic phrase
    PathIndex      uint32    `gorm:"not null"`   // HD path index used
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

// WalletBackup represents a wallet backup record
type WalletBackup struct {
    UUID           uuid.UUID `gorm:"primary_key;type:uuid"`
    UserUUID       uuid.UUID `gorm:"index;not null"`
    WalletAddress  string    `gorm:"index;not null"`
    BackupType     string    `gorm:"not null"` // mnemonic, private_key, keystore
    BackupData     []byte    `gorm:"type:bytea"`
    CreatedAt      time.Time
}