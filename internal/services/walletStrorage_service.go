package services

import (
    "fmt"
    "time"

    "git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// WalletStorageService handles secure storage of wallet credentials
type WalletStorageService struct {
    WalletDB        *gorm.DB
    MainDB          *gorm.DB
    EncryptionSvc   *EncryptionService
    ActivityLogger  *ActivityLoggerService
}

// NewWalletStorageService creates a new wallet storage service
func NewWalletStorageService(
    walletDB *gorm.DB,
    mainDB *gorm.DB,
    encryptionSvc *EncryptionService,
    activityLogger *ActivityLoggerService,
) *WalletStorageService {
    return &WalletStorageService{
        WalletDB:       walletDB,
        MainDB:         mainDB,
        EncryptionSvc:  encryptionSvc,
        ActivityLogger: activityLogger,
    }
}

// StoreMnemonic stores an encrypted mnemonic for a user's wallet
func (s *WalletStorageService) StoreMnemonic(userUUID uuid.UUID, walletAddress string, mnemonic string, pathIndex uint32) error {
    // Encrypt the mnemonic
    encMnemonic, err := s.EncryptionSvc.Encrypt([]byte(mnemonic))
    if err != nil {
        return fmt.Errorf("failed to encrypt mnemonic: %w", err)
    }

    // Create mnemonic record
    mnemonicRecord := models.WalletMnemonic{
        UUID:          uuid.New(),
        UserUUID:      userUUID,
        WalletAddress: walletAddress,
        EncMnemonic:   encMnemonic,
        PathIndex:     pathIndex,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    // Save to wallet database
    if result := s.WalletDB.Create(&mnemonicRecord); result.Error != nil {
        return fmt.Errorf("failed to store mnemonic: %w", result.Error)
    }

    return nil
}

// StorePrivateKey stores an encrypted private key for a user's wallet
func (s *WalletStorageService) StorePrivateKey(userUUID uuid.UUID, walletAddress string, privateKey string, keystoreJSON []byte) error {
    // Encrypt the private key and keystore JSON
    encPrivateKey, err := s.EncryptionSvc.Encrypt([]byte(privateKey))
    if err != nil {
        return fmt.Errorf("failed to encrypt private key: %w", err)
    }

    encKeystoreJSON, err := s.EncryptionSvc.Encrypt(keystoreJSON)
    if err != nil {
        return fmt.Errorf("failed to encrypt keystore JSON: %w", err)
    }

    // Create wallet record
    wallet := models.EncryptedWallet{
        UUID:           uuid.New(),
        UserUUID:       userUUID,
        WalletAddress:  walletAddress,
        EncPrivateKey:  encPrivateKey,
        EncKeystoreJSON: encKeystoreJSON,
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }

    // Save to wallet database
    if result := s.WalletDB.Create(&wallet); result.Error != nil {
        return fmt.Errorf("failed to store wallet credentials: %w", result.Error)
    }

    return nil
}

// GetMnemonic retrieves and decrypts a user's wallet mnemonic
func (s *WalletStorageService) GetMnemonic(userUUID uuid.UUID, walletAddress string) (string, uint32, error) {
    var mnemonicRecord models.WalletMnemonic

    // Find the mnemonic record
    if result := s.WalletDB.Where("user_uuid = ? AND wallet_address = ?", userUUID, walletAddress).
        First(&mnemonicRecord); result.Error != nil {
        return "", 0, fmt.Errorf("mnemonic not found: %w", result.Error)
    }

    // Decrypt the mnemonic
    mnemonicBytes, err := s.EncryptionSvc.Decrypt(mnemonicRecord.EncMnemonic)
    if err != nil {
        return "", 0, fmt.Errorf("failed to decrypt mnemonic: %w", err)
    }

    return string(mnemonicBytes), mnemonicRecord.PathIndex, nil
}

// GetPrivateKey retrieves and decrypts a user's private key
func (s *WalletStorageService) GetPrivateKey(userUUID uuid.UUID, walletAddress string) (string, error) {
    var wallet models.EncryptedWallet

    // Find the wallet record
    if result := s.WalletDB.Where("user_uuid = ? AND wallet_address = ?", userUUID, walletAddress).
        First(&wallet); result.Error != nil {
        return "", fmt.Errorf("wallet not found: %w", result.Error)
    }

    // Decrypt the private key
    privateKeyBytes, err := s.EncryptionSvc.Decrypt(wallet.EncPrivateKey)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt private key: %w", err)
    }

    return string(privateKeyBytes), nil
}

// GetKeystoreJSON retrieves and decrypts a user's keystore JSON
func (s *WalletStorageService) GetKeystoreJSON(userUUID uuid.UUID, walletAddress string) ([]byte, error) {
    var wallet models.EncryptedWallet

    // Find the wallet record
    if result := s.WalletDB.Where("user_uuid = ? AND wallet_address = ?", userUUID, walletAddress).
        First(&wallet); result.Error != nil {
        return nil, fmt.Errorf("wallet not found: %w", result.Error)
    }

    // Decrypt the keystore JSON
    keystoreBytes, err := s.EncryptionSvc.Decrypt(wallet.EncKeystoreJSON)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt keystore JSON: %w", err)
    }

    return keystoreBytes, nil
}

// CreateBackup creates a backup of wallet data
func (s *WalletStorageService) CreateBackup(userUUID uuid.UUID, walletAddress string, backupType string, data []byte) error {
    // Encrypt the backup data
    encData, err := s.EncryptionSvc.Encrypt(data)
    if err != nil {
        return fmt.Errorf("failed to encrypt backup data: %w", err)
    }

    // Create backup record
    backup := models.WalletBackup{
        UUID:         uuid.New(),
        UserUUID:     userUUID,
        WalletAddress: walletAddress,
        BackupType:   backupType,
        BackupData:   encData,
        CreatedAt:    time.Now(),
    }

    // Save to wallet database
    if result := s.WalletDB.Create(&backup); result.Error != nil {
        return fmt.Errorf("failed to create backup: %w", result.Error)
    }

    return nil
}

// HasStoredCredentials checks if a user has stored wallet credentials
func (s *WalletStorageService) HasStoredCredentials(userUUID uuid.UUID, walletAddress string) (bool, error) {
    var count int64

    // Check for mnemonic
    if err := s.WalletDB.Model(&models.WalletMnemonic{}).
        Where("user_uuid = ? AND wallet_address = ?", userUUID, walletAddress).
        Count(&count).Error; err != nil {
        return false, err
    }

    if count > 0 {
        return true, nil
    }

    // Check for private key
    if err := s.WalletDB.Model(&models.EncryptedWallet{}).
        Where("user_uuid = ? AND wallet_address = ?", userUUID, walletAddress).
        Count(&count).Error; err != nil {
        return false, err
    }

    return count > 0, nil
}

// DeleteWalletCredentials deletes stored wallet credentials
func (s *WalletStorageService) DeleteWalletCredentials(userUUID uuid.UUID, walletAddress string) error {
    // Delete mnemonic records
    if err := s.WalletDB.Where("user_uuid = ? AND wallet_address = ?", userUUID, walletAddress).
        Delete(&models.WalletMnemonic{}).Error; err != nil {
        return fmt.Errorf("failed to delete mnemonic: %w", err)
    }

    // Delete wallet records
    if err := s.WalletDB.Where("user_uuid = ? AND wallet_address = ?", userUUID, walletAddress).
        Delete(&models.EncryptedWallet{}).Error; err != nil {
        return fmt.Errorf("failed to delete wallet: %w", err)
    }

    // Delete backup records
    if err := s.WalletDB.Where("user_uuid = ? AND wallet_address = ?", userUUID, walletAddress).
        Delete(&models.WalletBackup{}).Error; err != nil {
        return fmt.Errorf("failed to delete backups: %w", err)
    }

    return nil
}

// ValidateUser ensures the user exists in the main database
func (s *WalletStorageService) ValidateUser(userUUID uuid.UUID) (bool, error) {
    var user models.User
    result := s.MainDB.Where("uuid = ?", userUUID).First(&user)
    return result.Error == nil, result.Error
}