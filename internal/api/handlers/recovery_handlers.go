package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecoveryRequest represents the recovery request
type RecoveryRequest struct {
    Email        string `json:"email" binding:"required,email"`
    KeystoreJSON []byte `json:"keystore_json" binding:"required"`
}

// RecoveryVerifyRequest represents the verification request
type RecoveryVerifyRequest struct {
    Token string `json:"token" binding:"required"`
}

// RequestWalletRecoveryHandler handles wallet recovery requests
func (h *Handler) RequestWalletRecoveryHandler(c *gin.Context) {
    var req RecoveryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Email and keystore JSON required."})
        return
    }

    err := h.RecoveryService.RequestRecovery(req.Email, req.KeystoreJSON)
    log.Println("Recovery request processed for email:", req.Email)
    log.Printf("Error " , err)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process recovery request"})
        return
    }

    // Always return success to prevent email enumeration
    c.JSON(http.StatusOK, gin.H{
        "message": "If the email exists in our system, a recovery link has been sent.",
    })
}

// VerifyWalletRecoveryHandler verifies a recovery token
func (h *Handler) VerifyWalletRecoveryHandler(c *gin.Context) {
    var req RecoveryVerifyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Token required."})
        return
    }

    keystoreJSON, err := h.RecoveryService.VerifyRecoveryToken(req.Token)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Recovery successful",
        "keystore_json": keystoreJSON,
    })
}