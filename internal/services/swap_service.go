package services

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/blockchain"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SwapService handles token swapping functionality
type SwapService struct {
    DB              *gorm.DB
    EthClient       *ethclient.Client
    WalletService   *WalletService
    UniswapRouter   common.Address
    WrappedEthToken common.Address
    CifoToken       common.Address
}

// NewSwapService creates a new swap service
func NewSwapService(
    db *gorm.DB,
    ethClient *ethclient.Client,
    walletService *WalletService,
    uniswapRouter common.Address,
    wrappedEthToken common.Address,
    cifoToken common.Address,
) *SwapService {
    return &SwapService{
        DB:              db,
        EthClient:       ethClient,
        WalletService:   walletService,
        UniswapRouter:   uniswapRouter,
        WrappedEthToken: wrappedEthToken,
        CifoToken:       cifoToken,
    }
}

// SwapTokens performs a token swap
func (s *SwapService) SwapTokens(
    userID uuid.UUID,
    mnemonic string,
    fromToken string,
    toToken string,
    amountStr string,
    slippageTolerance float64,
) (string, error) {
    // Parse amount
    amount, ok := new(big.Int).SetString(amountStr, 10)
    if !ok {
        return "", fmt.Errorf("invalid amount")
    }

    // Validate user's wallet
    hdWallet, err := blockchain.NewHDWalletFromMnemonic(mnemonic)
    if err != nil {
        return "", fmt.Errorf("invalid mnemonic: %w", err)
    }

    // Derive account
    account, err := hdWallet.DeriveAccount(0)
    if err != nil {
        return "", fmt.Errorf("failed to derive account: %w", err)
    }

    // Verify user owns this wallet
    var wallet models.Wallet
    if result := s.DB.Where("user_id = ? AND wallet_address = ?", 
        userID, account.Address.Hex()).First(&wallet); result.Error != nil {
        return "", fmt.Errorf("wallet not associated with user")
    }

    // Get token addresses
    var fromTokenAddress, toTokenAddress common.Address
    
    // Map token symbols to addresses
    if fromToken == "ETH" {
        fromTokenAddress = s.WrappedEthToken // Use WETH for Uniswap
    } else if fromToken == "CIFO" {
        fromTokenAddress = s.CifoToken
    } else {
        return "", fmt.Errorf("unsupported from token: %s", fromToken)
    }
    
    if toToken == "ETH" {
        toTokenAddress = s.WrappedEthToken
    } else if toToken == "CIFO" {
        toTokenAddress = s.CifoToken
    } else {
        return "", fmt.Errorf("unsupported to token: %s", toToken)
    }

    // Get nonce for the user's address
    nonce, err := s.EthClient.PendingNonceAt(context.Background(), account.Address)
    if err != nil {
        return "", fmt.Errorf("failed to get nonce: %w", err)
    }

    // Get gas price
    gasPrice, err := s.EthClient.SuggestGasPrice(context.Background())
    if err != nil {
        return "", fmt.Errorf("failed to get gas price: %w", err)
    }

    // Prepare Uniswap swap parameters
    deadline := big.NewInt(time.Now().Add(15 * time.Minute).Unix())
    
    // Calculate minimum amount out based on slippage
    // This is a simplified example - in production you'd use the Uniswap SDK or API
    minAmountOut := new(big.Int).Div(
        new(big.Int).Mul(
            amount, 
            big.NewInt(int64((1-slippageTolerance/100)*1000)),
        ), 
        big.NewInt(1000),
    )

    // Create path array [fromToken, toToken]
    path := []common.Address{fromTokenAddress, toTokenAddress}

    // Create swap transaction
    // Note: This is a simplified example - you'd need the actual Uniswap ABI and interface
	uniswapContract, err := blockchain.NewUniswapRouter(s.UniswapRouter, s.EthClient)
    if err != nil {
        return "", fmt.Errorf("failed to create Uniswap contract: %w", err)
    }

    // Create transaction options
    auth := bind.NewKeyedTransactor(account.PrivateKey)
    auth.Nonce = big.NewInt(int64(nonce))
    auth.Value = big.NewInt(0)     // Set to amount if swapping ETH
    auth.GasLimit = uint64(300000) // Set appropriate gas limit
    auth.GasPrice = gasPrice

    // If from token is ETH, use swapExactETHForTokens
    var tx *types.Transaction
    if fromToken == "ETH" {
        auth.Value = amount // Send ETH with the transaction
        tx, err = uniswapContract.SwapExactETHForTokens(
            auth,
            minAmountOut,
            path,
            account.Address, // Recipient
            deadline,
        )
    } else {
        // For token to ETH or token to token, need to approve first
        // Approve the router to spend tokens
        tokenContract, err := blockchain.NewERC20(fromTokenAddress, s.EthClient)
        if err != nil {
            return "", fmt.Errorf("failed to create token contract: %w", err)
        }
        
        approveTx, err := tokenContract.Approve(auth, s.UniswapRouter, amount)
        if err != nil {
            return "", fmt.Errorf("failed to approve token: %w", err)
        }
        
        // Wait for approval to be mined
        _, err = bind.WaitMined(context.Background(), s.EthClient, approveTx)
        if err != nil {
            return "", fmt.Errorf("failed to wait for approval: %w", err)
        }
        
        // Increment nonce after approval
        auth.Nonce = big.NewInt(auth.Nonce.Int64() + 1)
        
        // If to token is ETH, use swapExactTokensForETH
        if toToken == "ETH" {
            tx, err = uniswapContract.SwapExactTokensForETH(
                auth,
                amount,
                minAmountOut,
                path,
                account.Address, // Recipient
                deadline,
            )
        } else {
            // Token to token swap
            tx, err = uniswapContract.SwapExactTokensForTokens(
                auth,
                amount,
                minAmountOut,
                path,
                account.Address, // Recipient
                deadline,
            )
        }
    }
    
    if err != nil {
        return "", fmt.Errorf("failed to create swap transaction: %w", err)
    }

    // Store transaction in database
    walletTx := models.WalletTransaction{
        UUID:          uuid.New(),
        UserID:        userID,
        WalletAddress: account.Address.Hex(),
        TxHash:        tx.Hash().Hex(),
        TxType:        "swap",
        Amount:        amount.String(),
        TokenSymbol:   fromToken,
        ToAddress:     s.UniswapRouter.Hex(),
        Status:        "pending",
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    if result := s.DB.Create(&walletTx); result.Error != nil {
        return "", fmt.Errorf("failed to save transaction: %w", result.Error)
    }

    return tx.Hash().Hex(), nil
}