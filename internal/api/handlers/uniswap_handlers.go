package handlers

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
	"github.com/ethereum/go-ethereum/common"
)

func (h *Handler) processUniswapPurchase(transaction *models.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Minute)
	defer cancel()

	log.Println("Starting Uniswap purchase process for transaction: ", transaction.PaymentID)

	// Step 1: Convert ETH amount to wei (for sending to Uniswap)

	ethAmountWei := new(big.Float).Mul(
        big.NewFloat(transaction.EthAmount),
        new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
    )
    ethAmountInt, _ := ethAmountWei.Int(nil)

	minTokensWei := new(big.Float).Mul(
        big.NewFloat(transaction.MinTokenAmount),
        new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
    )
    minTokensInt, _ := minTokensWei.Int(nil)
    
    // Step 3: Get the destination wallet
    destWallet := common.HexToAddress(transaction.WalletAddress)
    
    // Step 4: Get token address - this is the CIFO token
    cifoTokenAddress := common.HexToAddress("0x7c74e84955891dfbdaf465be3d809f9605f93436")
    
    // Step 5: Get Uniswap router and WETH address
    uniswap := h.PriceService.GetUniswapClient()
    if uniswap == nil {
        return fmt.Errorf("Uniswap client not initialized")
    }

	    // Step 6: Prepare swap path (ETH -> CIFO)
    path := []common.Address{
        uniswap.Config.WethAddress, // WETH address
        cifoTokenAddress,         // CIFO token address
    }
    
    log.Printf("Executing Uniswap swap: %.8f ETH -> min %.8f CIFO tokens to %s",
        transaction.EthAmount, transaction.MinTokenAmount, destWallet.Hex())
    
    // Step 7: Execute the swap
    txHash, err := uniswap.SwapExactETHForTokens(
        ctx,
        ethAmountInt,
        minTokensInt,
        path,
        destWallet,
        time.Now().Add(20*time.Minute).Unix(), // 20 min deadline
    )

	if err != nil {
        transaction.Status = models.TransactionStatusFailed
        transaction.ErrorMessage = fmt.Sprintf("Uniswap swap failed: %v", err)
        h.DB.Save(&transaction)
        return fmt.Errorf("failed to execute Uniswap swap: %v", err)
    }
    
    // Step 8: Update transaction with swap details
    now := time.Now()
    transaction.Status = models.TransactionStatusCompleted
    transaction.BlockchainTxHash = txHash
    transaction.SwapTxHash = txHash
    transaction.BlockchainCompleted = true
    transaction.CompletedAt = &now
    transaction.UpdatedAt = now
    
    if err := h.DB.Save(&transaction).Error; err != nil {
        log.Printf("Warning: Failed to update transaction after successful swap: %v", err)
    }
    
    log.Printf("Successfully completed purchase of %.8f CIFO tokens via Uniswap, tx: %s",
        transaction.TokenAmount, txHash)
    
    return nil

}