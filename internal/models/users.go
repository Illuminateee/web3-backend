package models

import (
	"time"

	"github.com/google/uuid"
)


type User struct {
    UUID          uuid.UUID `gorm:"primary_key;type:uuid" json:"uuid"`
    Username      string    `gorm:"uniqueIndex;not null" json:"username"`
    Password      string    `gorm:"not null" json:"-"` 
    WalletAddress string    `gorm:"uniqueIndex;not null" json:"wallet_address"`
    Email         string    `gorm:"uniqueIndex" json:"email"`
    RecoveryHash  string    `gorm:"-" json:"-"` // For key recovery
    TwoFactorSecret string   `gorm:"column:two_factor_secret" json:"-"`
    TwoFactorEnabled bool    `gorm:"column:two_factor_enabled;default:false" json:"two_factor_enabled"`

    Phone          string    `gorm:"column:phone" json:"phone"`
    FullName       string    `gorm:"column:full_name" json:"full_name"`

    RecoveryCode     string     `json:"-" gorm:"column:recovery_code"`
    RecoveryExpires  *time.Time `json:"-" gorm:"column:recovery_expires"`
    LastActivity     time.Time  `json:"last_activity" gorm:"column:last_activity"`
    
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type Recovery struct {
    UUID         uuid.UUID  `gorm:"primary_key;type:uuid"`
    UserID       uuid.UUID  `gorm:"index;not null"`
    Token        string     `gorm:"not null"`
    KeystoreJSON []byte     `gorm:"not null"`
    ExpiresAt    time.Time  `gorm:"not null"`
    Used         bool       `gorm:"default:false"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}