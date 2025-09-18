package database

import (
    "fmt"
    "log"
    "time"

    "git.winteraccess.id/walanja/web3-tokensale-be/internal/config"
    "git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// ConnectWalletDB establishes a connection to the wallet database
func ConnectWalletDB(cfg *config.Config) (*gorm.DB, error) {
    dbLogger := logger.New(
        log.New(log.Writer(), "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             time.Second,
            LogLevel:                  logger.Info,
            IgnoreRecordNotFoundError: true,
            Colorful:                  true,
        },
    )

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
        cfg.WalletDB.Host, cfg.WalletDB.User, cfg.WalletDB.Password, 
        cfg.WalletDB.DBName, cfg.WalletDB.Port, 
        cfg.WalletDB.DBSSLMode) // Use the SSL mode from config

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: dbLogger,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to wallet database: %w", err)
    }

    // Auto-migrate wallet models only
    err = autoMigrateWalletDB(db)
    if err != nil {
        return nil, fmt.Errorf("failed to migrate wallet database schema: %w", err)
    }

    log.Println("Wallet database connected and migrated successfully")
    return db, nil
}

// autoMigrateWalletDB automatically migrates the wallet database schema
func autoMigrateWalletDB(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.EncryptedWallet{},
        &models.WalletMnemonic{},
        &models.WalletBackup{},
    )
}