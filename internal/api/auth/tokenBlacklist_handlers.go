package auth

import (
    "sync"
    "time"
)

// TokenBlacklist stores blacklisted tokens
type TokenBlacklist struct {
    blacklist map[string]time.Time
    mu        sync.RWMutex
}

// NewTokenBlacklist creates a new token blacklist
func NewTokenBlacklist() *TokenBlacklist {
    bl := &TokenBlacklist{
        blacklist: make(map[string]time.Time),
    }
    
    // Start a goroutine to clean up expired tokens
    go bl.cleanupRoutine()
    
    return bl
}

// Add adds a token to the blacklist
func (bl *TokenBlacklist) Add(token string, expiry time.Time) {
    bl.mu.Lock()
    defer bl.mu.Unlock()
    bl.blacklist[token] = expiry
}

// IsBlacklisted checks if a token is blacklisted
func (bl *TokenBlacklist) IsBlacklisted(token string) bool {
    bl.mu.RLock()
    defer bl.mu.RUnlock()
    _, exists := bl.blacklist[token]
    return exists
}

// cleanupRoutine periodically removes expired tokens
func (bl *TokenBlacklist) cleanupRoutine() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        bl.cleanup()
    }
}

// cleanup removes expired tokens from the blacklist
func (bl *TokenBlacklist) cleanup() {
    bl.mu.Lock()
    defer bl.mu.Unlock()
    
    now := time.Now()
    for token, expiry := range bl.blacklist {
        if now.After(expiry) {
            delete(bl.blacklist, token)
        }
    }
}