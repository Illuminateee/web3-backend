package handlers

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateWalletRequest struct {
    Password string `json:"password" binding:"required"`
    StoreCredentials bool  `json:"store_credentials"`
}

// ImportWalletRequest represents the request for wallet import
type ImportWalletRequest struct {
    Mnemonic         string `json:"mnemonic" binding:"required"`
    Password         string `json:"password" binding:"required"`
    Index            uint32 `json:"index,omitempty"`
    StoreCredentials bool   `json:"store_credentials"`
}

// SwapTokensRequest represents the request for token swapping
type SwapTokensRequest struct {
    Mnemonic       string  `json:"mnemonic" binding:"required"`
    FromToken      string  `json:"from_token" binding:"required"`
    ToToken        string  `json:"to_token" binding:"required"`
    Amount         string  `json:"amount" binding:"required"`
    SlippageTolerance float64 `json:"slippage_tolerance,omitempty"`
}

// CreateWalletHandler creates a new wallet for the authenticated user
func (h *Handler) CreateWalletHandler(c *gin.Context) {
    // Get user ID from auth middleware
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }

    var req CreateWalletRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
	uid, err := uuid.Parse(userID.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    // Create wallet
    wallet, mnemonic, err := h.WalletService.CreateUserWallet(uid, req.StoreCredentials)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create wallet: %v", err)})
        return
    }

    // Return wallet info and mnemonic
    c.JSON(http.StatusOK, gin.H{
        "message": "Wallet created successfully",
        "wallet_address": wallet.WalletAddress,
        "mnemonic": mnemonic,
        "important_notice": "SAVE THIS MNEMONIC SECURELY! It will not be stored on our servers and cannot be recovered if lost.",
    })
}

func (h *Handler) ImportWalletHandler(c *gin.Context) {
    // Get user ID from auth middleware
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }

    var req ImportWalletRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
	uid, err := uuid.Parse(userID.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Default index to 0 if not provided
    index := req.Index
    
    // Import wallet
    wallet, err := h.WalletService.ImportWalletFromMnemonic(uid, req.Mnemonic, index, req.StoreCredentials)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to import wallet: %v", err)})
        return
    }

    // Log the wallet import
    h.ActivityLoggerService.LogFromRequest(c, "import_wallet", 
        "User imported a wallet", 
        "wallet", wallet.WalletAddress, 
        "success", "")

    // Return wallet info
    response := gin.H{
        "message": "Wallet imported successfully",
        "wallet_address": wallet.WalletAddress,
    }

    if req.StoreCredentials {
        response["backup_status"] = "Wallet credentials securely stored"
    } else {
        response["backup_status"] = "Not backed up. Keep your mnemonic safe."
    }

    c.JSON(http.StatusOK, response)
}



// GetUserWalletBalanceHandler gets the balance of the authenticated user's wallet
func (h *Handler) GetUserWalletBalanceHandler(c *gin.Context) {
    // Get user ID from the context (set by auth middleware)
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }

    // Get user from database
    var user models.User
    if result := h.DB.Where("uuid = ?", userID).First(&user); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Ensure user has a wallet address
    if user.WalletAddress == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No wallet address associated with this account"})
        return
    }

    // Validate the address
    if !common.IsHexAddress(user.WalletAddress) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid wallet address format"})
        return
    }

    ethAddress := common.HexToAddress(user.WalletAddress)

    // Connect to Ethereum node
    client, err := ethclient.Dial(h.GetEthereumNodeURL())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Ethereum node"})
        return
    }
    defer client.Close()

    // Get ETH balance
    ethBalance, err := client.BalanceAt(context.Background(), ethAddress, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Get ETH price in USD and IDR for reference
    ethPriceUSD, ethPriceIDR, err := h.GetEthPrices(c.Request.Context())
    if err != nil {
        log.Printf("Warning: Failed to get ETH prices: %v", err)
        // Continue anyway as this is just reference info
    }

    // Calculate ETH value
    ethAmount := new(big.Float).SetInt(ethBalance)
    ethAmount = ethAmount.Quo(ethAmount, big.NewFloat(1e18)) // Convert Wei to ETH
    ethFloat, _ := ethAmount.Float64()

    // Calculate USD and IDR values
    usdValue := ethFloat * ethPriceUSD
    idrValue := ethFloat * ethPriceIDR

    // Get token balance
    var tokenBalance *big.Int
    
    // Using PaymentGateway to get token balance rather than TokenClient
    if h.BlockchainService != nil && h.BlockchainService.PaymentGateway != nil {

        tokenBalance, err = h.BlockchainService.PaymentGateway.GetTokenBalance(
            context.Background(), 
            ethAddress,
        )
        
        if err != nil {
            log.Printf("Warning: Failed to get CIFO token balance: %v", err)
        }
    }

    h.ActivityLoggerService.LogFromRequest(c, "check_balance", 
    "User checked wallet balance", 
    "wallet", user.WalletAddress, 
    "success", "")   
    // Prepare response
    response := gin.H{
        "wallet_address":  user.WalletAddress,
        "eth_balance_wei": ethBalance.String(),
        "eth_balance":     fmt.Sprintf("%.6f ETH", ethFloat),
        "eth_value_usd":   fmt.Sprintf("$%.2f USD", usdValue),
        "eth_value_idr":   fmt.Sprintf("Rp %.2f IDR", idrValue),
    }

    // Add token balance if available
    if tokenBalance != nil {
        tokenAmount := new(big.Float).SetInt(tokenBalance)
        tokenAmount = tokenAmount.Quo(tokenAmount, big.NewFloat(1e18)) // Assuming 18 decimals
        tokenFloat, _ := tokenAmount.Float64()

        response["cifo_balance"] = tokenBalance.String()
        response["cifo_formatted"] = fmt.Sprintf("%.6f CIFO", tokenFloat)
    }

    c.JSON(http.StatusOK, response)
}

