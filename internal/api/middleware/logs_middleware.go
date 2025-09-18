package middleware

import (
    "git.winteraccess.id/walanja/web3-tokensale-be/internal/services"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// ActivityLoggingMiddleware creates middleware for logging API requests
func ActivityLoggingMiddleware(activityLogger *services.ActivityLoggerService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Process request
        c.Next()

        // After request is processed
        userID, exists := c.Get("user_id")
        if !exists {
            return // Skip logging if not authenticated
        }

        // Only log for authenticated routes
        status := c.Writer.Status()
        path := c.Request.URL.Path
        method := c.Request.Method
        username, _ := c.Get("username")
        usernameStr, _ := username.(string)

        uid, err := uuid.Parse(userID.(string))
        if err != nil {
            return
        }

        // Determine status text
        statusText := "success"
        errorMsg := ""
        if status >= 400 {
            statusText = "failure"
            if err, exists := c.Get("error"); exists {
                errorMsg = err.(string)
            }
        }

        // Log the API request
        activityLogger.LogUserActivity(
            c.Request.Context(),
            uid,
            usernameStr,
            "api_request",
            method+" "+path,
            "route",
            path,
            statusText,
            errorMsg,
        )
    }
}