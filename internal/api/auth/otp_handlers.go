package auth

import (
    "fmt"
    "time"

    "github.com/pquerna/otp"
    "github.com/pquerna/otp/totp"
)

// TOTPService handles two-factor authentication
type TOTPService struct {
    Issuer     string
    SecretSize int
}

// NewTOTPService creates a new TOTP service
func NewTOTPService(issuer string) *TOTPService {
    return &TOTPService{
        Issuer:     issuer,
        SecretSize: 20,
    }
}

// GenerateSecret generates a new TOTP secret for a user
func (s *TOTPService) GenerateSecret(username string) (*otp.Key, error) {
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      s.Issuer,
        AccountName: username,
        SecretSize:  uint(s.SecretSize),
    })
    if err != nil {
        return nil, err
    }
    
    return key, nil
}

// ValidateCode validates a TOTP code
func (s *TOTPService) ValidateCode(secret string, code string) bool {
    return totp.Validate(code, secret)
}

// GenerateCode generates a TOTP code for testing
func (s *TOTPService) GenerateCode(secret string) (string, error) {
    return totp.GenerateCode(secret, time.Now())
}

// GetTOTPProvisioningURI returns the URI for QR code generation
func (s *TOTPService) GetTOTPProvisioningURI(username, secret string) string {
    return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
        s.Issuer, username, secret, s.Issuer)
}