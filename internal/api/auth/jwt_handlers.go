package auth

import (
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

// JWT configuration settings
type JWTConfig struct {
    SecretKey     string
    TokenDuration time.Duration
}

// Claims struct for JWT token
type Claims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    Address  string `json:"address"`
    Email string `json:"email"`
    WalletAddress string `json:"wallet_address"`
    TokenType     string `json:"token_type,omitempty"` // "access" or "refresh"
    jwt.RegisteredClaims
}

// TokenService handles JWT token generation and validation
type TokenService struct {
    config JWTConfig
    blacklist *TokenBlacklist
}

// NewTokenService creates a new token service
func NewTokenService(config JWTConfig) *TokenService {
    return &TokenService{
        config:    config,
        blacklist: NewTokenBlacklist(),
    }
}

// Add new method to get config
func (s *TokenService) GetConfig() JWTConfig {
    return s.config
}

// GenerateToken creates a new JWT token for a user
func (s *TokenService) GenerateToken(userID uuid.UUID, username, walletAddress string) (string, error) {
    // Set expiration time
    expirationTime := time.Now().Add(s.config.TokenDuration)

    // Create claims with user data
    claims := &Claims{
        UserID:   userID.String(),
        Username: username,
        Address:  walletAddress,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "web3-tokensale",
            Subject:   username,
            ID:        uuid.New().String(),
        },
    }

    // Create token with claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Generate signed token
    tokenString, err := token.SignedString([]byte(s.config.SecretKey))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// ValidateToken validates and parses a JWT token
func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
    // Check if token is blacklisted
    if s.blacklist.IsBlacklisted(tokenString) {
        return nil, fmt.Errorf("token is blacklisted")
    }
    
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Validate signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(s.config.SecretKey), nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}

func (s *TokenService) GenerateTempToken(userID uuid.UUID, username string) (string, error) {
    // Create token with claims
    claims := jwt.MapClaims{
        "user_id":  userID.String(),
        "username": username,
        "temp":     true,                          // Mark this as a temporary token
        "exp":      time.Now().Add(5 * time.Minute).Unix(), // Short expiration (5 minutes)
        "iat":      time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Generate signed token
    tokenString, err := token.SignedString([]byte(s.config.SecretKey))
    if err != nil {
        return "", err
    }
    
    return tokenString, nil
}

// Add new function for refresh tokens
func (s *TokenService) GenerateRefreshToken(userID uuid.UUID, username, walletAddress string) (string, error) {
    // Set expiration time (longer than access token)
    expirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 days
    
    claims := &Claims{
        UserID:        userID.String(),
        Username:      username,
        WalletAddress: walletAddress,
        TokenType:     "refresh",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "web3-tokensale",
            Subject:   username,
            ID:        uuid.New().String(),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(s.config.SecretKey))
    if err != nil {
        return "", err
    }
    
    return tokenString, nil
}

func (s *TokenService) BlacklistToken(tokenString string) error {
    // Parse without validating to get expiry
    claims := &Claims{}
    _, _, err := jwt.NewParser().ParseUnverified(tokenString, claims)
    if err != nil {
        return err
    }
    
    // Get expiry time
    expiry := claims.ExpiresAt.Time
    
    // Add to blacklist
    s.blacklist.Add(tokenString, expiry)
    return nil
}