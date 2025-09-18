package models

import (
    "time"

    "github.com/google/uuid"
)

// ActivityLog represents a record of user activity in the system
type ActivityLog struct {
    UUID        uuid.UUID `gorm:"primary_key;type:uuid"`
    UserID      uuid.UUID `gorm:"index;not null"`
    Username    string    `gorm:"index"`
    Action      string    `gorm:"not null"` // login, check_balance, create_transaction, etc.
    Description string    `gorm:"type:text"`
    IPAddress   string    `gorm:"size:45"` // Supports IPv6 addresses
    UserAgent   string    `gorm:"type:text"`
    Resource    string    // Resource identifier (e.g., transaction hash)
    ResourceID  string    // ID of the resource being accessed/modified
    Status      string    // success, failure, error
    ErrorMsg    string    `gorm:"type:text"`
    CreatedAt   time.Time `gorm:"index"`

    // Relationships
    User User `gorm:"foreignKey:UserID"`
}