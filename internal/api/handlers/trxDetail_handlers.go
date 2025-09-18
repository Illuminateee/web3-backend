package handlers

import (
	"net/http"
	"strconv"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTransactionDetailsHandler(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction hash is required"})
		return
	}

	ctx := c.Request.Context()
	txHash := common.HexToHash(hash)

	// Get transaction details
	tx, pending, err := h.BlockchainService.PaymentGateway.GetEthClient().TransactionByHash(ctx, txHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transaction: " + err.Error()})
		return
	}

	// Get transaction receipt
	receipt, err := h.BlockchainService.PaymentGateway.GetEthClient().TransactionReceipt(ctx, txHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transaction receipt: " + err.Error()})
		return
	}

	// Format transaction data
	
	var from string
    
    // Get the sender using the correct approach
    signer := types.LatestSignerForChainID(tx.ChainId())
    sender, err := types.Sender(signer, tx)
    if err == nil {
        from = sender.Hex()
    }

	result := gin.H{
		"hash":      tx.Hash().Hex(),
		"from":      from,
		"to":        tx.To().String(),
		"value":     tx.Value().String(),
		"gas_price": tx.GasPrice().String(),
		"gas":       tx.Gas(),
		"nonce":     tx.Nonce(),
		"pending":   pending,
		"receipt": map[string]interface{}{
			"status":            receipt.Status,
			"gas_used":          receipt.GasUsed,
			"block_number":      receipt.BlockNumber.Uint64(),
			"block_hash":        receipt.BlockHash.Hex(),
			"transaction_index": receipt.TransactionIndex,
		},
	}

	// Check if related to a payment
	var transaction models.Transaction
	if err := h.DB.Where("blockchain_tx_hash = ?", hash).First(&transaction).Error; err == nil {
		// Found related transaction
		result["payment"] = map[string]interface{}{
			"payment_id":         transaction.PaymentID,
			"destination_wallet": transaction.WalletAddress,
			"fiat_amount":        transaction.FiatAmount,
			"fiat_currency":      transaction.FiatCurrency,
			"token_amount":       transaction.TokenAmount,
			"status":             transaction.Status,
			"created_at":         transaction.CreatedAt,
			"completed_at":       transaction.CompletedAt,
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetTransactionStatusHandler(c *gin.Context) {
    // Get transaction ID from path
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
        return
    }
    
    // Find transaction in database
    var transaction models.Transaction
    if err := h.DB.Where("payment_id = ? OR uuid = ?", id, id).First(&transaction).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
        return
    }
    
    // Format response
    response := gin.H{
        "id":                transaction.PaymentID,
        "uuid":              transaction.UUID,
        "status":            transaction.Status,
        "created_at":        transaction.CreatedAt,
        "updated_at":        transaction.UpdatedAt,
        "transaction_type":  transaction.TransactionType, // Include transaction type
        "token_details": gin.H{
            "amount":  transaction.TokenAmount,
            "symbol":  transaction.TokenSymbol,
            "wallet":  transaction.WalletAddress,
        },
        "payment_details": gin.H{
            "fiat_amount":   transaction.FiatAmount,
            "fiat_currency": transaction.FiatCurrency,
            "method":        transaction.PaymentMethod,
            "gas_fee_eth":   transaction.GasFee,
            "gas_fee_fiat":  transaction.GasFeeFiat,
            "total_fiat":    transaction.FiatAmount + transaction.GasFeeFiat,
        },
    }
    
    // Add blockchain details if available
    if transaction.BlockchainCompleted {
        response["blockchain_details"] = gin.H{
            "tx_hash":      transaction.BlockchainTxHash,
            "completed_at": transaction.CompletedAt,
        }
    }
    
    if transaction.ErrorMessage != "" {
        response["error"] = transaction.ErrorMessage
    }
    
    c.JSON(http.StatusOK, response)
}

func (h *Handler) GetUserTransactionsHandler(c *gin.Context) {
    // Get user ID from JWT token
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    
    // Parse query parameters
    limit := 10 // Default limit
    offset := 0 // Default offset
    
    if limitStr := c.Query("limit"); limitStr != "" {
        if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
            limit = parsedLimit
        }
    }
    
    if offsetStr := c.Query("offset"); offsetStr != "" {
        if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
            offset = parsedOffset
        }
    }
    
    // Status filter
    status := c.Query("status")
    
    // Prepare the database query
    query := h.DB.Where("user_id = ?", userID)
    
    // Apply status filter if provided
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    // Get total count
    var total int64
    if err := query.Model(&models.Transaction{}).Count(&total).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count transactions"})
        return
    }
    
    // Get transactions with pagination
    var transactions []models.Transaction
    if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&transactions).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
        return
    }
    
    // Format response
    formattedTransactions := make([]gin.H, 0, len(transactions))
    for _, tx := range transactions {
        formattedTx := gin.H{
            "id":               tx.PaymentID,
            "uuid":             tx.UUID,
            "status":           tx.Status,
            "created_at":       tx.CreatedAt,
            "updated_at":       tx.UpdatedAt,
            "transaction_type": tx.TransactionType, // Include transaction type
            "token_details": gin.H{
                "amount": tx.TokenAmount,
                "symbol": tx.TokenSymbol,
                "wallet": tx.WalletAddress,
            },
            "payment_details": gin.H{
                "fiat_amount":   tx.FiatAmount,
                "fiat_currency": tx.FiatCurrency,
                "method":        tx.PaymentMethod,
                "gas_fee_eth":   tx.GasFee,
                "gas_fee_fiat":  tx.GasFeeFiat,
                "total_fiat":    tx.FiatAmount + tx.GasFeeFiat,
            },
        }
        
        // Add blockchain details if available
        if tx.BlockchainCompleted {
            formattedTx["blockchain_details"] = gin.H{
                "tx_hash":      tx.BlockchainTxHash,
                "completed_at": tx.CompletedAt,
            }
        }
        
        // Add swap details only for purchases with swap type
        if tx.TransactionType == "purchase" && tx.SwapType != "" {
            formattedTx["swap_details"] = gin.H{
                "swap_type":        tx.SwapType,
                "swap_tx_hash":     tx.SwapTxHash,
                "min_token_amount": tx.MinTokenAmount,
            }
        }
        
        formattedTransactions = append(formattedTransactions, formattedTx)
    }
    
    c.JSON(http.StatusOK, gin.H{
        "total": total,
        "limit": limit,
        "offset": offset,
        "transactions": formattedTransactions,
    })
}