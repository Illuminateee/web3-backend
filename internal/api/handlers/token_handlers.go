package handlers

import (
    "net/http"
    "strconv"

    "git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// SwapTokensHandler handles token swap requests
func (h *Handler) SwapTokensHandler(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "User not authenticated",
        })
        return
    }

    // Parse request body
    var req struct {
        Mnemonic          string  `json:"mnemonic" binding:"required"`
        FromToken         string  `json:"from_token" binding:"required"`
        ToToken           string  `json:"to_token" binding:"required"`
        Amount            string  `json:"amount" binding:"required"`
        SlippageTolerance float64 `json:"slippage_tolerance"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid request",
            "errors":  err.Error(),
        })
        return
    }

    // Validate token types
    if (req.FromToken != "ETH" && req.FromToken != "CIFO") ||
        (req.ToToken != "ETH" && req.ToToken != "CIFO") {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid token types. Supported tokens: ETH, CIFO",
        })
        return
    }

    // Use a default slippage tolerance if not provided
    if req.SlippageTolerance <= 0 {
        req.SlippageTolerance = 0.5 // Default 0.5% slippage
    }

    // Parse user ID
    uid, err := uuid.Parse(userID.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid user ID",
        })
        return
    }

    // Execute the swap
    txHash, err := h.SwapService.SwapTokens(
        uid,
        req.Mnemonic,
        req.FromToken,
        req.ToToken,
        req.Amount,
        req.SlippageTolerance,
    )

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "Failed to swap tokens",
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Token swap initiated",
        "tx_hash": txHash,
    })
}

// GetWalletTransactionsHandler retrieves wallet transactions for the authenticated user
func (h *Handler) GetWalletTransactionsHandler(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "User not authenticated",
        })
        return
    }

    // Parse pagination parameters
    page := 1
    pageSize := 10

    if p, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && p > 0 {
        page = p
    }

    if ps, err := strconv.Atoi(c.DefaultQuery("page_size", "10")); err == nil && ps > 0 && ps <= 100 {
        pageSize = ps
    }

    // Parse wallet address from query if provided
    walletAddress := c.Query("wallet_address")

    // Parse user ID
    uid, err := uuid.Parse(userID.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid user ID",
        })
        return
    }

    // Build a query to get transactions for this user
    var transactions []models.WalletTransaction
    query := h.DB.Where("user_id = ?", uid)

    // Add wallet address filter if provided
    if walletAddress != "" {
        query = query.Where("wallet_address = ?", walletAddress)
    }

    // Apply pagination
    offset := (page - 1) * pageSize
    
    // Count total records for pagination
    var total int64
    if err := query.Model(&models.WalletTransaction{}).Count(&total).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "Failed to count transactions",
            "error":   err.Error(),
        })
        return
    }

    // Get paginated results
    if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&transactions).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "Failed to fetch wallet transactions",
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":      "success",
        "transactions": transactions,
        "pagination": gin.H{
            "total":      total,
            "page":       page,
            "page_size":  pageSize,
            "total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
        },
    })
}