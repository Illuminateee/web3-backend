package handlers

import (
    "fmt"
    "net/http"
    "time"

    "git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

// EnableWalletBackupRequest represents a request to enable wallet backup
type EnableWalletBackupRequest struct {
    Password string `json:"password" binding:"required"`
}

// RecoverWalletRequest represents a request to recover wallet credentials
type RecoverWalletRequest struct {
    Password string `json:"password" binding:"required"`
}

// EnableWalletBackupHandler enables wallet credential backup for a user
func (h *Handler) EnableWalletBackupHandler(c *gin.Context) {
    // Get authenticated user
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }

    var req EnableWalletBackupRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Get user from database
    var user models.User
    if result := h.DB.Where("uuid = ?", userID).First(&user); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Verify password
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
        h.ActivityLoggerService.LogFromRequest(c, "enable_wallet_backup", 
            "User attempted to enable wallet backup", 
            "wallet", user.WalletAddress, 
            "failure", "Invalid password")
        return
    }

    // Ensure user has a wallet address
    if user.WalletAddress == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No wallet associated with this account"})
        return
    }

    // Check if credentials are already stored
    uid, _ := uuid.Parse(userID.(string))
    hasCredentials := h.WalletService.HasStoredCredentials(uid, user.WalletAddress)
    if hasCredentials {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet backup already enabled"})
        return
    }

    // Create a new wallet with backup (to replace existing wallet)
    wallet, mnemonic, err := h.WalletService.CreateUserWallet(uid, true)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to backup wallet: %v", err)})
        return
    }

    // Update user with new wallet address if needed
    if wallet.WalletAddress != user.WalletAddress {
        h.DB.Model(&user).Updates(map[string]interface{}{
            "wallet_address": wallet.WalletAddress,
            "updated_at":     time.Now(),
        })
    }

    h.ActivityLoggerService.LogFromRequest(c, "enable_wallet_backup", 
        "User enabled wallet backup", 
        "wallet", wallet.WalletAddress, 
        "success", "")

    // Return success with mnemonic for the user to record as well
    c.JSON(http.StatusOK, gin.H{
        "message": "Wallet backup enabled successfully",
        "wallet_address": wallet.WalletAddress,
        "mnemonic": mnemonic,
        "important_notice": "SAVE THIS MNEMONIC SECURELY as an additional backup!",
    })
}

// RecoverWalletHandler recovers wallet credentials for a user
func (h *Handler) RecoverWalletHandler(c *gin.Context) {
    // Get authenticated user
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }

    var req RecoverWalletRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Get user from database
    var user models.User
    if result := h.DB.Where("uuid = ?", userID).First(&user); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Verify password
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
        h.ActivityLoggerService.LogFromRequest(c, "recover_wallet", 
            "User attempted to recover wallet credentials", 
            "wallet", user.WalletAddress, 
            "failure", "Invalid password")
        return
    }

    // Ensure user has a wallet address
    if user.WalletAddress == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No wallet associated with this account"})
        return
    }

    // Retrieve wallet credentials
    uid, _ := uuid.Parse(userID.(string))
    credentials, err := h.WalletService.RecoverWalletFromStorage(uid, user.WalletAddress)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "No wallet backup found"})
        h.ActivityLoggerService.LogFromRequest(c, "recover_wallet", 
            "User attempted to recover wallet credentials", 
            "wallet", user.WalletAddress, 
            "failure", "No backup found")
        return
    }

    h.ActivityLoggerService.LogFromRequest(c, "recover_wallet", 
        "User recovered wallet credentials", 
        "wallet", user.WalletAddress, 
        "success", "")

    // Return the credentials
    c.JSON(http.StatusOK, gin.H{
        "message": "Wallet credentials recovered successfully",
        "wallet_address": user.WalletAddress,
        "credentials": credentials,
        "important_notice": "SAVE THESE CREDENTIALS SECURELY!",
    })
}