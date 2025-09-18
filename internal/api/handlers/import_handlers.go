package handlers

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ImportAccountRequest struct {
	Username   string `json:"username" binding:"required"`
	PrivateKey string `json:"private_key" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// ImportAccountHandler imports an existing Ethereum account and associates it with a username
func (h *Handler) ImportAccountHandler(c *gin.Context) {
	var req ImportAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Username, private key and password required."})
		return
	}

	// Validate username
	if len(req.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 3 characters"})
		return
	}

	// Check if username already exists
	var count int64
	if result := h.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check username availability"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	}

	// Password strength validation
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters"})
		return
	}

	// Clean up the private key format if needed
	cleanKey := req.PrivateKey
	if len(cleanKey) >= 2 && cleanKey[:2] == "0x" {
		cleanKey = cleanKey[2:]
	}

	// Parse the private key
	privateKeyBytes, err := hex.DecodeString(cleanKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key format"})
		return
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key"})
		return
	}

	// Create temp directory for keystore
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("eth-keystore-%s", randomString(8)))
	if err := os.MkdirAll(tempDir, 0700); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create keystore directory"})
		return
	}
	defer os.RemoveAll(tempDir)

	// Create keystore
	ks := keystore.NewKeyStore(tempDir, keystore.StandardScryptN, keystore.StandardScryptP)
	if ks == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create keystore"})
		return
	}

	// Import the private key to create an account
	account, err := ks.ImportECDSA(privateKey, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import private key"})
		return
	}

	// Get the JSON representation
	jsonData, err := ks.Export(account, req.Password, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export keystore JSON"})
		return
	}

	// Hash the password for database storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to secure password"})
		return
	}

	// Create user record in database
	now := time.Now()
	userID := uuid.New()
	user := models.User{
		UUID:          userID,
		Username:      req.Username,
		Password:      string(hashedPassword),
		WalletAddress: account.Address.Hex(),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Save to database
	if result := h.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user record"})
		return
	}

	// Return the account information
	c.JSON(http.StatusOK, AccountResponse{
		UUID:         userID.String(),
		Username:     req.Username,
		Address:      account.Address.Hex(),
		KeystoreJSON: jsonData,
		CreatedAt:    now,
	})
}