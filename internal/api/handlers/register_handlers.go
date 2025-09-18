package handlers

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
)

// CreateAccountRequest structure
type CreateAccountRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

// AccountResponse response structure
type AccountResponse struct {
    UUID        string `json:"uuid"`
    Username    string `json:"username"`
    Email    string `json:"email"`
    Address     string `json:"address"`
    KeystoreJSON []byte `json:"keystore_json,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}

type RegisterUserRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
    FullName string `json:"full_name"`
    Phone    string `json:"phone"`
}

// RegisterUserResponse represents the response after successful registration
type RegisterUserResponse struct {
    UUID      string    `json:"uuid"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    FullName  string    `json:"full_name,omitempty"`
    Phone     string    `json:"phone,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}

// CreateAccountHandler creates a new Ethereum account and associates it with a username/password
func (h *Handler) CreateAccountHandler(c *gin.Context) {
    var req CreateAccountRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Username, Password, or Email required."})
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

    // check email availability

    var emailCount int64
    if result := h.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&emailCount); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email availability"})
        return
    }

    if emailCount > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
        return
    }

    // Password strength validation
    if len(req.Password) < 8 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters"})
        return
    }

    // Create a temporary directory for the keystore
    tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("eth-keystore-%s", randomString(8)))
    if err := os.MkdirAll(tempDir, 0700); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create keystore directory"})
        return
    }
    defer os.RemoveAll(tempDir) // Clean up after we're done

    // Create the keystore
    ks := keystore.NewKeyStore(tempDir, keystore.StandardScryptN, keystore.StandardScryptP)
    if ks == nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create keystore"})
        return
    }

    // Create a new account
    account, err := ks.NewAccount(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
        return
    }

    // Get the JSON representation
    jsonData, err := ks.Export(account, req.Password, req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export keystore JSON"})
        return
    }

    // Connect to Ethereum node to verify connectivity
    client, err := ethclient.Dial(h.Config.EthereumRPC)
    if err != nil {
        log.Printf("Warning: Failed to connect to Ethereum node: %v", err)
        // Continue with account creation even if node connection fails
    } else {
        defer client.Close()
        // Get balance just to verify account on blockchain
        balance, err := client.BalanceAt(context.Background(), account.Address, nil)
        if err == nil {
            log.Printf("Account created. Initial balance: %s wei", balance.String())
        }
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
        Email:        req.Email,
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
        Email:      req.Email,
        Address:      account.Address.Hex(),
        KeystoreJSON: jsonData,
        CreatedAt:    now,
    })
}

// RegisterUserHandler registers a new user without creating a wallet
func (h *Handler) RegisterUserHandler(c *gin.Context) {
    var req RegisterUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Please provide required fields."})
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

    // Check if email already exists
    var emailCount int64
    if result := h.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&emailCount); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email availability"})
        return
    }

    if emailCount > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
        return
    }

    // Password strength validation
    if len(req.Password) < 8 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters"})
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
        UUID:           userID,
        Username:       req.Username,
        Password:       string(hashedPassword),
        Email:          req.Email,
        FullName:       req.FullName,
        Phone:          req.Phone,
        CreatedAt:      now,
        UpdatedAt:      now,
        WalletAddress:  "",  // Empty wallet address - user can link wallet later
        TwoFactorEnabled: false,
    }

    // Save to database
    if result := h.DB.Create(&user); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user record"})
        return
    }

    // Return the user information without sensitive data
    c.JSON(http.StatusCreated, RegisterUserResponse{
        UUID:      userID.String(),
        Username:  req.Username,
        Email:     req.Email,
        FullName:  req.FullName,
        Phone:     req.Phone,
        CreatedAt: now,
    })
}

func randomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}