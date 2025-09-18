package middleware

import (
    "net/http"
    "strings"
    "time"

    "git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
    "gorm.io/gorm"
)

// SessionTimeout middleware checks if the user session has expired
func SessionTimeout(db *gorm.DB, secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get token from Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.Next()
            return
        }

        // Check if bearer token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.Next()
            return
        }

        tokenString := parts[1]

        // Parse token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secretKey), nil
        })

        if err != nil || !token.Valid {
            c.Next()
            return
        }

        // Extract user ID from claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.Next()
            return
        }

        userID, ok := claims["user_id"].(string)
        if !ok {
            c.Next()
            return
        }

        // Get user from database
        var user models.User
        if result := db.Where("uuid = ?", userID).First(&user); result.Error != nil {
            c.Next()
            return
        }

        // Check if last activity is older than 30 minutes
        if time.Since(user.LastActivity) > 30*time.Minute {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Session expired, please login again",
            })
            return
        }

        // Update last activity
        db.Model(&user).Update("last_activity", time.Now())

        // Continue
        c.Next()
    }
}