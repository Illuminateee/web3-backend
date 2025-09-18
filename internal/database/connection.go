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

// Connect establishes a connection to the database
func Connect(cfg *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
    )

    // Configure logger based on environment
    logLevel := logger.Silent
    if cfg.Environment == "development" {
        logLevel = logger.Info
    }

    // Configure GORM
    gormConfig := &gorm.Config{
        Logger: logger.New(
            log.New(log.Writer(), "\r\n", log.LstdFlags),
            logger.Config{
                SlowThreshold:             time.Second,
                LogLevel:                  logLevel,
                IgnoreRecordNotFoundError: true,
                Colorful:                  true,
            },
        ),
    }

    // Connect to database
    db, err := gorm.Open(postgres.Open(dsn), gormConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Auto-migrate models
    err = autoMigrate(db)
    if err != nil {
        return nil, fmt.Errorf("failed to migrate database schema: %w", err)
    }

    log.Println("Database connected and migrated successfully")
    return db, nil
}

// autoMigrate automatically migrates the database schema
func autoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Recovery{},
        &models.Transaction{},
        &models.Wallet{},
        &models.ActivityLog{},
        // Add other models here as needed
    )
}