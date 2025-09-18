package middleware

import (
    "net/http"
    "strings"

    "git.winteraccess.id/walanja/web3-tokensale-be/internal/api/auth"
    "github.com/gin-gonic/gin"
)

// AuthMiddleware creates middleware for JWT authentication
func AuthMiddleware(tokenService *auth.TokenService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        
        // Check if Authorization header exists
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // Check if it's a Bearer token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        
        // Validate token
        claims, err := tokenService.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }
        
        // Set user info for handlers to use
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("address", claims.Address)
        c.Set("user_email", claims.Email)
        c.Set("wallet_address", claims.WalletAddress)
        
        c.Next()
    }
}