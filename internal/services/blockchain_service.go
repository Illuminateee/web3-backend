package services

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/blockchain"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/config"
	"github.com/ethereum/go-ethereum/common"
)

// BlockchainService handles interactions with blockchain
type BlockchainService struct {
    config         *config.Config
    PaymentGateway *blockchain.PaymentGatewayClient
}

// NewBlockchainService creates a new blockchain service
func NewBlockchainService(cfg *config.Config, paymentGateway *blockchain.PaymentGatewayClient) (*BlockchainService, error) {
    // Check if payment gateway is initialized if provided
    if paymentGateway != nil && !paymentGateway.IsInitialized() {
        return nil, fmt.Errorf("payment gateway client was not properly initialized")
    }

    return &BlockchainService{
        config:         cfg,
        PaymentGateway: paymentGateway,
    }, nil
}
func (s *BlockchainService) ProcessPayment(ctx context.Context, paymentID string, isSuccess bool, gateway string, transactionData map[string]interface{}) error {
    log.Printf("Processing blockchain payment %s, success: %v, gateway: %s", paymentID, isSuccess, gateway)

    // Check if payment gateway is initialized
    if s.PaymentGateway == nil {
        return fmt.Errorf("payment gateway client is not initialized")
    }

    exists, err := s.PaymentGateway.CheckPaymentExists(ctx, paymentID)
    log.Printf("Payment ID %s exists: %v", paymentID, exists)
    if err != nil {
        log.Printf("Warning: Error checking payment existence: %v", err)
    }

    if !exists {
        log.Printf("Payment ID %s does not exist in contract, creating it now", paymentID)
        
        // Get required gas deposit
        gasDeposit, err := s.PaymentGateway.GetRequiredGasDeposit(ctx)
        if err != nil {
            return fmt.Errorf("failed to get required gas deposit: %v", err)
        }

        // Get token and fiat amounts from transaction data
        var tokenAmount *big.Int
        var fiatAmount *big.Int
        var destinationWallet common.Address

        // Extract transaction data if available
        if transactionData != nil {
            // Get token amount
            if tokenAmountValue, ok := transactionData["token_amount"]; ok {
                if tokenAmountBig, ok := tokenAmountValue.(*big.Int); ok {
                    tokenAmount = tokenAmountBig
                }
            }

            // Get fiat amount
            if fiatAmountValue, ok := transactionData["fiat_amount"]; ok {
                if fiatAmountBig, ok := fiatAmountValue.(*big.Int); ok {
                    fiatAmount = fiatAmountBig
                }
            }

            // Get destination wallet
            if walletValue, ok := transactionData["destination_wallet"]; ok {
                if walletStr, ok := walletValue.(string); ok && common.IsHexAddress(walletStr) {
                    destinationWallet = common.HexToAddress(walletStr)
                }
            }
        }

        // Use default values if not provided
        if tokenAmount == nil {
            tokenAmount = big.NewInt(1999815000000) // Example token amount
        }
        
        if fiatAmount == nil {
            fiatAmount = big.NewInt(26569891000000) // Example fiat amount
        }
        
        // If no destination wallet, use the default account (for testing only)
        if destinationWallet == (common.Address{}) {
            log.Printf("Warning: No destination wallet provided, using default")
            // In production, you should handle this error properly
            destinationWallet = common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
        }

        log.Printf("Creating payment with destination wallet: %s", destinationWallet.Hex())

        // Create the payment with destination wallet
        _ , err = s.PaymentGateway.CreatePayment(
            ctx, 
            paymentID, 
            tokenAmount, 
            fiatAmount, 
            gateway,
            destinationWallet, // Pass the destination wallet
            gasDeposit,
        )
        
        if err != nil {
            return fmt.Errorf("failed to create payment: %v", err)
        }
        
        log.Printf("Successfully created payment %s in contract for wallet %s", 
            paymentID, destinationWallet.Hex())
    }

    // Convert boolean success to uint8 status code
    var statusCode uint8
    if isSuccess {
        statusCode = 1 // Success status code
    } else {
        statusCode = 2 // Failure status code
    }

    // Call the payment gateway contract to process the payment
    err = s.PaymentGateway.MockPaymentCallback(ctx, paymentID, statusCode)
    if err != nil {
        return fmt.Errorf("failed to process payment callback: %v", err)
    }

    return nil
}

// GetPaymentStatus retrieves the status of a payment
func (s *BlockchainService) GetPaymentStatus(ctx context.Context, paymentID string) (blockchain.PaymentStatus, error) {
    return s.PaymentGateway.GetPaymentStatus(ctx, paymentID)
}

// GetPaymentDetails retrieves full details of a payment
func (s *BlockchainService) GetPaymentDetails(ctx context.Context, paymentID string) (*blockchain.PaymentDetails, error) {
    return s.PaymentGateway.GetPaymentDetails(ctx, paymentID)
}


