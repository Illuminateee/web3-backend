package services

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "io"
)

// EncryptionService handles encryption and decryption of sensitive data
type EncryptionService struct {
    Key []byte
}

// NewEncryptionService creates a new encryption service with the provided key
func NewEncryptionService(hexKey string) (*EncryptionService, error) {
    key, err := hex.DecodeString(hexKey)
    if err != nil {
        return nil, err
    }
    
    // Ensure key is 32 bytes (256 bits) for AES-256
    if len(key) != 32 {
        return nil, errors.New("encryption key must be 32 bytes (64 hex characters)")
    }
    
    return &EncryptionService{
        Key: key,
    }, nil
}

// Encrypt encrypts data using AES-GCM
func (s *EncryptionService) Encrypt(plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(s.Key)
    if err != nil {
        return nil, err
    }
    
    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    // Create a nonce
    nonce := make([]byte, aesGCM.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    // Encrypt and seal the data
    ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}

// Decrypt decrypts data using AES-GCM
func (s *EncryptionService) Decrypt(ciphertext []byte) ([]byte, error) {
    block, err := aes.NewCipher(s.Key)
    if err != nil {
        return nil, err
    }
    
    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    // Extract the nonce
    nonceSize := aesGCM.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }
    
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    
    // Decrypt the data
    plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }
    
    return plaintext, nil
}