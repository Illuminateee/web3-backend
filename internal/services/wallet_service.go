package services

import (
	"fmt"
	"math/big"
	"time"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/blockchain"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WalletService handles wallet-related operations
type WalletService struct {
    DB *gorm.DB
    StorageService *WalletStorageService
}

// NewWalletService creates a new wallet service
func NewWalletService(db *gorm.DB, storageService *WalletStorageService) *WalletService {
    return &WalletService{
        DB: db,
        StorageService: storageService,
    }
}

// CreateUserWallet creates a new wallet for a user
func (s *WalletService) CreateUserWallet(userID uuid.UUID, storeCredentials bool) (*models.Wallet, string, error) {
    // Check if user already has a wallet
    if s.DB == nil {
        return nil, "", fmt.Errorf("database connection not initialized")
    }
    
    var existingWallet models.Wallet
    if result := s.DB.Where("user_id = ?", userID).First(&existingWallet); result.Error == nil {
        return nil, "", fmt.Errorf("user already has a wallet")
    }

    // Create HD wallet
    hdWallet, err := blockchain.NewHDWallet()
    if err != nil {
        return nil, "", fmt.Errorf("failed to create wallet: %w", err)
    }

    // Derive first account
    account, err := hdWallet.DeriveAccount(0)
    if err != nil {
        return nil, "", fmt.Errorf("failed to derive account: %w", err)
    }

    // Create wallet record in main DB
    wallet := models.Wallet{
        UUID:          uuid.New(),
        UserID:        userID,
        WalletAddress: account.Address.Hex(),
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    if result := s.DB.Create(&wallet); result.Error != nil {
        return nil, "", fmt.Errorf("failed to save wallet: %w", result.Error)
    }

    // If storing credentials is requested
    if storeCredentials && s.StorageService != nil {
        // Store the mnemonic
        err = s.StorageService.StoreMnemonic(userID, account.Address.Hex(), hdWallet.Mnemonic, 0)
        if err != nil {
            // Log the error but don't fail the wallet creation
            fmt.Printf("Warning: Failed to store mnemonic: %v\n", err)
        }

        // Store private key and keystore
        privateKey := account.GetPrivateKeyHex()
        keystoreJSON, err := account.ExportKeystore("temporary_password") // Will be re-encrypted in the storage service
        if err == nil {
            err = s.StorageService.StorePrivateKey(userID, account.Address.Hex(), privateKey, keystoreJSON)
            if err != nil {
                fmt.Printf("Warning: Failed to store private key: %v\n", err)
            }
        }

        // Create a backup as well
        backupData, _ := account.ExportKeystore("backup_password")
        if backupData != nil {
            _ = s.StorageService.CreateBackup(userID, account.Address.Hex(), "keystore", backupData)
        }
    }

    // Return the wallet, mnemonic (for client-side backup)
    return &wallet, hdWallet.Mnemonic, nil
}

// ImportWalletFromMnemonic imports a wallet using a mnemonic
func (s *WalletService) ImportWalletFromMnemonic(userID uuid.UUID, mnemonic string, index uint32, storeCredentials bool) (*models.Wallet, error) {
    // Create HD wallet from mnemonic
    hdWallet, err := blockchain.NewHDWalletFromMnemonic(mnemonic)
    if err != nil {
        return nil, fmt.Errorf("invalid mnemonic: %w", err)
    }

    // Derive account at specified index
    account, err := hdWallet.DeriveAccount(index)
    if err != nil {
        return nil, fmt.Errorf("failed to derive account: %w", err)
    }

    // Create wallet record
    wallet := models.Wallet{
        UUID:          uuid.New(),
        UserID:        userID,
        WalletAddress: account.Address.Hex(),
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    if result := s.DB.Create(&wallet); result.Error != nil {
        return nil, fmt.Errorf("failed to save wallet: %w", result.Error)
    }

    // If storing credentials is requested
    if storeCredentials && s.StorageService != nil {
        // Store the mnemonic
        err = s.StorageService.StoreMnemonic(userID, account.Address.Hex(), mnemonic, index)
        if err != nil {
            fmt.Printf("Warning: Failed to store mnemonic: %v\n", err)
        }

        // Store private key
        privateKey := account.GetPrivateKeyHex()
        keystoreJSON, err := account.ExportKeystore("temporary_password")
        if err == nil {
            err = s.StorageService.StorePrivateKey(userID, account.Address.Hex(), privateKey, keystoreJSON)
            if err != nil {
                fmt.Printf("Warning: Failed to store private key: %v\n", err)
            }
        }
    }

    return &wallet, nil
}

func (s *WalletService) HasStoredCredentials(userID uuid.UUID, walletAddress string) bool {
    if s.StorageService == nil {
        return false
    }

    hasCredentials, _ := s.StorageService.HasStoredCredentials(userID, walletAddress)
    return hasCredentials
}

func (s *WalletService) RecoverWalletFromStorage(userID uuid.UUID, walletAddress string) (string, error) {
    if s.StorageService == nil {
        return "", fmt.Errorf("wallet storage service not initialized")
    }

    // Try to get mnemonic first (preferred)
    mnemonic, _, err := s.StorageService.GetMnemonic(userID, walletAddress)
    if err == nil {
        return mnemonic, nil
    }

    // If mnemonic not found, try private key
    privateKey, err := s.StorageService.GetPrivateKey(userID, walletAddress)
    if err == nil {
        return privateKey, nil
    }

    return "", fmt.Errorf("no wallet credentials found for this wallet")
}

// CreateSignedTransaction creates a signed transaction for token transfer
func (s *WalletService) CreateSignedTransaction(
    userID uuid.UUID, 
    mnemonic string, 
    toAddress common.Address, 
    amount *big.Int, 
    gasPrice *big.Int,
    nonce uint64,
) (*types.Transaction, error) {
    // Recreate the HD wallet
    hdWallet, err := blockchain.NewHDWalletFromMnemonic(mnemonic)
    if err != nil {
        return nil, fmt.Errorf("invalid mnemonic: %w", err)
    }

    // Derive account (typically the first one)
    account, err := hdWallet.DeriveAccount(0)
    if err != nil {
        return nil, fmt.Errorf("failed to derive account: %w", err)
    }

    // Validate that this is the user's wallet
    var wallet models.Wallet
    if result := s.DB.Where("user_id = ? AND wallet_address = ?", 
        userID, account.Address.Hex()).First(&wallet); result.Error != nil {
        return nil, fmt.Errorf("wallet not associated with user")
    }

    // Create transaction
    tx := types.NewTransaction(
        nonce,
        toAddress,
        amount,
        21000, // standard gas limit
        gasPrice,
        nil, // data
    )

    // Sign transaction
    chainID := big.NewInt(1) // mainnet
    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), account.PrivateKey)
    if err != nil {
        return nil, fmt.Errorf("failed to sign transaction: %w", err)
    }

    return signedTx, nil
}