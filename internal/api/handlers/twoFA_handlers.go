package handlers

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Setup2FARequest represents the 2FA setup request
type Setup2FARequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// Verify2FARequest represents the 2FA verification request
type Verify2FARequest struct {
    Username string `json:"username"`
    Code     string `json:"code"`
}

// Login2FARequest represents the 2FA login request
type Login2FARequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Code     string `json:"code"`
}

type RecoveryFARequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
}

type VerifyRecoveryFARequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Code     string `json:"code"`
}



// Setup2FAHandler initiates 2FA setup
func (h *Handler) Setup2FAHandler(c *gin.Context) {
    var req Setup2FARequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Authenticate user first
    var user models.User
    if result := h.DB.Where("username = ?", req.Username).First(&user); result.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    // Verify password
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    // Generate TOTP secret
    key, err := h.TOTPService.GenerateSecret(req.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate 2FA secret"})
        return
    }

    // Store secret temporarily (not activated yet)
    h.DB.Model(&user).Updates(map[string]interface{}{
        "two_factor_secret": key.Secret(),
        "updated_at":        time.Now(),
    })

    // Return data needed for QR code generation
    c.JSON(http.StatusOK, gin.H{
        "message":      "2FA setup initiated. Scan the QR code with your authenticator app.",
        "secret":       key.Secret(),
        "qr_code_url":  key.URL(),
        "provisioning_uri": h.TOTPService.GetTOTPProvisioningURI(user.Username, key.Secret()),
    })
}

// Verify2FAHandler verifies and activates 2FA
func (h *Handler) Verify2FAHandler(c *gin.Context) {
    var req Verify2FARequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Get user
    var user models.User
    if result := h.DB.Where("username = ?", req.Username).First(&user); result.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
        return
    }

    // Ensure 2FA is not already enabled
    if user.TwoFactorEnabled {
        c.JSON(http.StatusBadRequest, gin.H{"error": "2FA is already enabled"})
        return
    }

    // Validate TOTP code
    if !h.TOTPService.ValidateCode(user.TwoFactorSecret, req.Code) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
            h.ActivityLoggerService.LogFromRequest(c, "enabled_2fa", 
        "User activated 2 fa for their account", 
        "account", user.Email, 
        "failed", "")  
        return
    }

    // Enable 2FA
    h.DB.Model(&user).Updates(map[string]interface{}{
        "two_factor_enabled": true,
        "updated_at":        time.Now(),
    })

    c.JSON(http.StatusOK, gin.H{
        "message": "Two-factor authentication enabled successfully",
    })

    h.ActivityLoggerService.LogFromRequest(c, "enabled_2fa", 
    "User activated 2 fa for their account", 
    "account", user.Email, 
    "success", "")  
}

// LoginWith2FAHandler handles login with 2FA
func (h *Handler) LoginWith2FAHandler(c *gin.Context) {
    var req Login2FARequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Get user
    var user models.User
    if result := h.DB.Where("username = ?", req.Username).First(&user); result.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    // Verify password
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    // If 2FA is enabled, verify code
    if user.TwoFactorEnabled {
        if !h.TOTPService.ValidateCode(user.TwoFactorSecret, req.Code) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
            return
        }
    }

    // Update last login time
    h.DB.Model(&user).Updates(map[string]interface{}{
        "updated_at":    time.Now(),
        "last_activity": time.Now(),
    })

    // Generate JWT token
    token, err := h.TokenService.GenerateToken(user.UUID, user.Username, user.WalletAddress)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Generate refresh token
    refreshToken, err := h.TokenService.GenerateRefreshToken(user.UUID, user.Username, user.WalletAddress)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
        return
    }
    // Calculate expiry time
    expiresAt := time.Now().Add(h.TokenService.GetConfig().TokenDuration)

    c.JSON(http.StatusOK, LoginResponse{
        UUID:          user.UUID.String(),
        Username:      user.Username,
        WalletAddress: user.WalletAddress,
        AccessToken:   token,
        RefreshToken:  refreshToken,
        ExpiresAt:     expiresAt,
    })
}

// Disable2FAHandler disables 2FA for a user
func (h *Handler) Disable2FAHandler(c *gin.Context) {
    // Get user ID from authenticated session
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get user
    var user models.User
    if result := h.DB.Where("uuid = ?", userID).First(&user); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Disable 2FA
    h.DB.Model(&user).Updates(map[string]interface{}{
        "two_factor_enabled": false,
        "two_factor_secret":  "",
        "updated_at":         time.Now(),
    })

    c.JSON(http.StatusOK, gin.H{
        "message": "Two-factor authentication disabled successfully",
    })

    h.ActivityLoggerService.LogFromRequest(c, "disable_2fa", 
    "User disable 2 fa for their account", 
    "account", user.Email, 
    "success", "")  
}

func (h *Handler) Request2FARecovery(c *gin.Context){
    var req RecoveryFARequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    // Get user
    var user models.User
    if result := h.DB.Where("username = ? AND email = ?", req.Username, req.Email).First(&user); result.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or email"})
        return
    }

    // ensure 2FA is enabled
    if !user.TwoFactorEnabled {
        c.JSON(http.StatusBadRequest, gin.H{"error": "2FA is not enabled"})
        return
    }
    // Generate recovery code
    recoveryCode := generateRandomCode(6)
    // hash the recovery code
    hashedCode, err := bcrypt.GenerateFromPassword([]byte(recoveryCode), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recovery code"})
        return
    }

    // Store recovery code with expiration (30 minutes)
    h.DB.Model(&user).Updates(map[string]interface{}{
        "recovery_code":      string(hashedCode),
        "recovery_expires":   time.Now().Add(30 * time.Minute),
        "updated_at":         time.Now(),
    })

    // Send recovery code via email
    err = h.RecoveryService.SendRecoveryFAEmail(user.Email, user.Username, recoveryCode)
    if err != nil {
        log.Printf("Failed to send recovery email: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send recovery email"})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "Recovery code has been sent to your email. Code valid for 30 minutes.",
    })
}

func (h *Handler) Verify2FARecoveryHandler(c *gin.Context) {
    var req VerifyRecoveryFARequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Get user
    var user models.User
    if result := h.DB.Where("username = ? AND email = ?", req.Username, req.Email).First(&user); result.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or email"})
        return
    }

    // Check if recovery code exists and is not expired
    if user.RecoveryExpires.Before(time.Now()) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Recovery code expired. Please request a new one."})
        return
    }

    // Verify recovery code
    err := bcrypt.CompareHashAndPassword([]byte(user.RecoveryCode), []byte(req.Code))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid recovery code"})
        return
    }

    // Generate new 2FA secret
    key, err := h.TOTPService.GenerateSecret(req.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new 2FA secret"})
        return
    }

    // Update user with new 2FA secret and clear recovery code
    h.DB.Model(&user).Updates(map[string]interface{}{
        "two_factor_secret":  key.Secret(),
        "recovery_code":      "",
        "recovery_expires":   nil,
        "updated_at":         time.Now(),
    })

    // Return data needed for QR code generation
    c.JSON(http.StatusOK, gin.H{
        "message":           "2FA has been reset. Scan the new QR code with your authenticator app.",
        "secret":            key.Secret(),
        "qr_code_url":       key.URL(),
        "provisioning_uri":  h.TOTPService.GetTOTPProvisioningURI(user.Username, key.Secret()),
    })

    h.ActivityLoggerService.LogFromRequest(c, "recovery_2fa", 
    "User recover 2 fa code for their account", 
    "account", user.Email, 
    "success", "")  
}
// Helper function to generate random codes
func generateRandomCode(length int) string {
    const digits = "0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = digits[rand.Intn(len(digits))]
    }
    return string(b)
}