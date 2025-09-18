package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gorm.io/gorm"
)

// RecoveryService handles wallet recovery functionality
type RecoveryService struct {
    DB              *gorm.DB
    SendGridAPIKey  string
    AppURL          string
    RecoveryExpiry  time.Duration
    FromEmail       string
    FromName        string
}
// NewRecoveryService creates a new recovery service
func NewRecoveryService(db *gorm.DB, sendgridAPIKey, appURL, fromEmail, fromName string) *RecoveryService {
    return &RecoveryService{
        DB:             db,
        SendGridAPIKey: sendgridAPIKey,
        AppURL:         appURL,
        RecoveryExpiry: 24 * time.Hour,
        FromEmail:      fromEmail,
        FromName:       fromName,
    }
}

// GenerateRecoveryToken creates a new recovery token
func (s *RecoveryService) GenerateRecoveryToken() (string, error) {
    tokenBytes := make([]byte, 32)
    _, err := rand.Read(tokenBytes)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(tokenBytes), nil
}

// RequestRecovery initiates a keystore recovery process
func (s *RecoveryService) RequestRecovery(email string, keystoreJSON []byte) error {
    var user models.User
    if result := s.DB.Where("email = ?", email).First(&user); result.Error != nil {
        // Don't reveal whether the email exists to prevent enumeration attacks
        return nil
    }

    // Generate a unique token
    token, err := s.GenerateRecoveryToken()
    if err != nil {
        return err
    }

    // Create recovery record
    recovery := models.Recovery{
        UUID:         uuid.New(),
        UserID:       user.UUID,
        Token:        token,
        KeystoreJSON: keystoreJSON,
        ExpiresAt:    time.Now().Add(s.RecoveryExpiry),
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    if result := s.DB.Create(&recovery); result.Error != nil {
        return result.Error
    }

    // Send recovery email
    return s.SendRecoveryEmail(user.Email, user.Username, token)
}

// SendRecoveryEmail sends a recovery link to the user's email

func (s *RecoveryService) SendRecoveryEmail(to, username, token string) error {
    from := mail.NewEmail(s.FromName, s.FromEmail)
    subject := "Your Wallet Recovery Link"
    toEmail := mail.NewEmail(username, to)
    
    recoveryLink := fmt.Sprintf("%s/recover?token=%s", s.AppURL, token)
    
    plainTextContent := fmt.Sprintf("Hello %s,\n\nYou requested to recover your wallet. "+
        "Click the link below to access your wallet:\n\n%s\n\n"+
        "This link will expire in 24 hours.\n\n"+
        "If you did not request this, please ignore this email.\n\n"+
        "Best regards,\nThe Web3 Tokensale Team", username, recoveryLink)
    
    htmlContent := fmt.Sprintf("<h2>Hello %s,</h2>"+
        "<p>You requested to recover your wallet. Click the button below to access your wallet:</p>"+
        "<p><a href=\"%s\" style=\"padding: 10px 20px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px;\">Recover My Wallet</a></p>"+
        "<p>Or copy this link: <a href=\"%s\">%s</a></p>"+
        "<p>This link will expire in 24 hours.</p>"+
        "<p>If you did not request this, please ignore this email.</p>"+
        "<p>Best regards,<br>The Web3 Tokensale Team</p>", username, recoveryLink, recoveryLink, recoveryLink)
    
    message := mail.NewSingleEmail(from, subject, toEmail, plainTextContent, htmlContent)
    client := sendgrid.NewSendClient(s.SendGridAPIKey)
    _, err := client.Send(message)
    return err
}

// VerifyRecoveryToken validates and processes a recovery token
func (s *RecoveryService) VerifyRecoveryToken(token string) ([]byte, error) {
    var recovery models.Recovery
    
    // Find the recovery record
    result := s.DB.Where("token = ? AND used = ? AND expires_at > ?", 
        token, false, time.Now()).First(&recovery)
    
    if result.Error != nil {
        return nil, fmt.Errorf("invalid or expired recovery token")
    }
    
    // Mark as used
    s.DB.Model(&recovery).Updates(map[string]interface{}{
        "used":       true,
        "updated_at": time.Now(),
    })
    
    return recovery.KeystoreJSON, nil
}


func (s *RecoveryService) SendRecoveryFAEmail(to, username, code string) error {
    from := mail.NewEmail(s.FromName, s.FromEmail)
    subject := "2FA Recovery Code"
    toEmail := mail.NewEmail(username, to)
    
    plainTextContent := fmt.Sprintf("Hello %s,\n\n"+
        "You requested to recover your Two-Factor Authentication (2FA).\n\n"+
        "Your recovery code is: %s\n\n"+
        "This code will expire in 30 minutes.\n\n"+
        "If you didn't request this recovery, please ignore this email or contact support immediately.\n\n"+
        "Best regards,\nThe Web3 Tokensale Team", username, code)
    
    htmlContent := fmt.Sprintf("<h2>Hello %s,</h2>"+
        "<p>We received a request to recover your Two-Factor Authentication (2FA).</p>"+
        "<p>Your recovery code is: <strong style=\"font-size: 18px; letter-spacing: 2px; background: #f8f9fa; padding: 5px 10px; border-radius: 3px;\">%s</strong></p>"+
        "<p style=\"color: #e74c3c; font-weight: bold;\">This code will expire in 30 minutes.</p>"+
        "<p>If you didn't request this recovery, please ignore this email or contact support immediately.</p>"+
        "<p style=\"margin-top: 20px; font-size: 12px; color: #7f8c8d;\">This is an automated message, please do not reply.</p>",
        username, code)
    
    message := mail.NewSingleEmail(from, subject, toEmail, plainTextContent, htmlContent)
    client := sendgrid.NewSendClient(s.SendGridAPIKey)
    response, err := client.Send(message)
    
    if err != nil {
        return fmt.Errorf("failed to send recovery email: %w", err)
    }
    
    if response.StatusCode >= 400 {
        return fmt.Errorf("failed to send recovery email: status code %d", response.StatusCode)
    }
    
    return nil
}