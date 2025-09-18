package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest structure for user login
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// LoginResponse structure
type LoginResponse struct {
    UUID          string    `json:"uuid"`
    Username      string    `json:"username"`
    WalletAddress string    `json:"wallet_address"`
    AccessToken   string    `json:"access_token"`
    RefreshToken  string    `json:"refresh_token,omitempty"`
    ExpiresAt     time.Time `json:"expires_at"`
}

// RefreshTokenRequest represents a request to refresh an access token
type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest represents a request to logout (invalidate tokens)
type LogoutRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

// LoginHandler authenticates a user and returns their wallet info
func (h *Handler) LoginHandler(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Username and password required."})
        return
    }

    // Find user by username
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

    // Check if 2FA is enabled for this user
    if user.TwoFactorEnabled {
        // Generate a temporary token for the 2FA verification step
        tempToken, err := h.TokenService.GenerateTempToken(user.UUID, user.Username)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate temporary token"})
            return
        }
        
        // Return a response indicating 2FA is required
        c.JSON(http.StatusOK, gin.H{
            "message": "2FA required",
            "requires_2fa": true,
            "temp_token": tempToken,
            "username": user.Username,
        })
        return
    }

    // Update last login time
    h.DB.Model(&user).Updates(map[string]interface{}{
        "updated_at":    time.Now(),
        "last_activity": time.Now(),
    })

    // Generate a JWT token
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
    
    h.ActivityLoggerService.LogFromRequest(c, "login", "User successfully logged in", "user", user.Username, "success", "")

    // Return user data with tokens
    c.JSON(http.StatusOK, LoginResponse{
        UUID:          user.UUID.String(),
        Username:      user.Username,
        WalletAddress: user.WalletAddress,
        AccessToken:   token,
        RefreshToken:  refreshToken,
        ExpiresAt:     expiresAt,
    })
}

func (h *Handler) GetAccountBalanceHandler(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address is required"})
		return
	}

	// Validate the address
	if !common.IsHexAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Ethereum address"})
		return
	}

	ethAddress := common.HexToAddress(address)

	// Connect to Ethereum node
	client, err := ethclient.Dial(h.GetEthereumNodeURL())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Ethereum node"})
		return
	}
	defer client.Close()

	// Get balance
	balance, err := client.BalanceAt(context.Background(), ethAddress, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account balance"})
		return
	}

	// Get ETH price in USD and IDR for reference
	ethPriceUSD, ethPriceIDR, err := h.GetEthPrices(c.Request.Context())
	if err != nil {
		log.Printf("Warning: Failed to get ETH prices: %v", err)
		// Continue anyway as this is just reference info
	}

	// Return balance information
	c.JSON(http.StatusOK, gin.H{
		"address":       address,
		"balance_wei":   balance.String(),
		"balance_eth":   fmt.Sprintf("%.18f", float64(balance.Int64())/1e18),
		"eth_price_usd": ethPriceUSD,
		"eth_price_idr": ethPriceIDR,
	})
}

// RefreshTokenHandler handles token refresh requests
func (h *Handler) RefreshTokenHandler(c *gin.Context) {
    // Get refresh token from header
    refreshToken := c.GetHeader("X-Refresh-Token")
    if refreshToken == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
        return
    }

    // Validate refresh token
    claims, err := h.TokenService.ValidateToken(refreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
        return
    }

    // Ensure it's a refresh token
    if claims.TokenType != "refresh" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
        return
    }
    
    // Get user
    userID, err := uuid.Parse(claims.UserID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in token"})
        return
    }

    var user models.User
    if result := h.DB.Where("uuid = ?", userID).First(&user); result.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    // Update last activity time
    h.DB.Model(&user).Update("last_activity", time.Now())

    // Generate new access token
    newToken, err := h.TokenService.GenerateToken(user.UUID, user.Username, user.WalletAddress)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
        return
    }

    // Calculate expiry time
    expiresAt := time.Now().Add(h.TokenService.GetConfig().TokenDuration)

    c.JSON(http.StatusOK, gin.H{
        "access_token": newToken,
        "expires_at":   expiresAt,
    })
}

// LogoutHandler handles user logout
func (h *Handler) LogoutHandler(c *gin.Context) {
    // Get refresh token from header
    refreshToken := c.GetHeader("X-Refresh-Token")
    if refreshToken == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
        return
    }

    // Get authorization header for access token
    authHeader := c.GetHeader("Authorization")
    if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
        accessToken := authHeader[7:]
        // Blacklist access token
        h.TokenService.BlacklistToken(accessToken)
    }

    // Blacklist refresh token
    h.TokenService.BlacklistToken(refreshToken)

    c.JSON(http.StatusOK, gin.H{
        "message": "Successfully logged out",
    })
}
func (h *Handler) GetEthereumNodeURL() string {
	// You can make this configurable via environment variables later
	return "https://geth-geth.ede2390e1937cf50.dyndns.dappnode.io"	
}

