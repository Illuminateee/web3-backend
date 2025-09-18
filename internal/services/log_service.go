package services

import (
    "context"
    "time"

    "git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// ActivityLoggerService handles user activity logging
type ActivityLoggerService struct {
    DB *gorm.DB
}

// NewActivityLoggerService creates a new activity logger service
func NewActivityLoggerService(db *gorm.DB) *ActivityLoggerService {
    return &ActivityLoggerService{
        DB: db,
    }
}

// LogUserActivity logs a user activity
func (s *ActivityLoggerService) LogUserActivity(ctx context.Context, userID uuid.UUID, username string, action, description, resource, resourceID string, status string, errorMsg string) error {
    // Extract IP and user agent from context if available
    var ipAddress, userAgent string
    if ginCtx, ok := ctx.(*gin.Context); ok {
        ipAddress = getClientIP(ginCtx)
        userAgent = ginCtx.Request.UserAgent()
    }

    log := models.ActivityLog{
        UUID:        uuid.New(),
        UserID:      userID,
        Username:    username,
        Action:      action,
        Description: description,
        IPAddress:   ipAddress,
        UserAgent:   userAgent,
        Resource:    resource,
        ResourceID:  resourceID,
        Status:      status,
        ErrorMsg:    errorMsg,
        CreatedAt:   time.Now(),
    }

    return s.DB.Create(&log).Error
}

// LogFromRequest logs activity from a request
func (s *ActivityLoggerService) LogFromRequest(c *gin.Context, action, description, resource, resourceID string, status string, errorMsg string) {
    // Get user info from context
    userID, exists := c.Get("user_id")
    if !exists {
        return // Don't log if not authenticated
    }

    username, _ := c.Get("username")
    usernameStr, _ := username.(string)

    uid, err := uuid.Parse(userID.(string))
    if err != nil {
        return
    }

    // Log the activity
    log := models.ActivityLog{
        UUID:        uuid.New(),
        UserID:      uid,
        Username:    usernameStr,
        Action:      action,
        Description: description,
        IPAddress:   getClientIP(c),
        UserAgent:   c.Request.UserAgent(),
        Resource:    resource,
        ResourceID:  resourceID,
        Status:      status,
        ErrorMsg:    errorMsg,
        CreatedAt:   time.Now(),
    }

    s.DB.Create(&log)
}

// getClientIP extracts the client IP from the request
func getClientIP(c *gin.Context) string {
    // Check for X-Forwarded-For header
    if xForwardedFor := c.GetHeader("X-Forwarded-For"); xForwardedFor != "" {
        return xForwardedFor
    }
    // Check for X-Real-IP header
    if xRealIP := c.GetHeader("X-Real-IP"); xRealIP != "" {
        return xRealIP
    }
    // Use RemoteAddr as fallback
    return c.ClientIP()
}