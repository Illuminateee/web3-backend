package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/config"
	"git.winteraccess.id/walanja/web3-tokensale-be/pkg/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
    UNISWAP_V3_QUOTER_ADDRESS = "0x61fFE014bA17989E743c5F6cB21bF9697530B21e"
    UNISWAP_V2_ROUTER_ADDRESS = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
)

type UniswapClient struct {
    client    *rpc.Client
    Config    *config.Config  // Change to uppercase to make it exported
    ethClient *ethclient.Client
    router *contracts.UniswapV2Router02
    privateKey *ecdsa.PrivateKey
}

func NewUniswapClient(client *rpc.Client, cfg *config.Config) *UniswapClient {
    ethClient := ethclient.NewClient(client)

    // initialize the Uniswap router contract
    routerAddress := common.HexToAddress(UNISWAP_V2_ROUTER_ADDRESS)
    router, err := contracts.NewUniswapV2Router02(routerAddress, ethClient)
    if err != nil {
        log.Printf("Warning: Failed to initialize Uniswap router: %v", err)
    }
    // Parse private key from config
    var privateKey *ecdsa.PrivateKey
    if cfg.WalletPrivateKey != "" {
        privateKey, err = crypto.HexToECDSA(cfg.WalletPrivateKey)
        if err != nil {
            log.Printf("Warning: Failed to parse private key: %v", err)
        }
    } else {
        log.Printf("Warning: No wallet private key provided in config")
    }

    return &UniswapClient{
        client:     client,
        Config:     cfg,
        ethClient:  ethClient,
        router:     router,
        privateKey: privateKey,
    }
}

func (uc *UniswapClient) GetTokenPrice(ctx context.Context, tokenAddress string) (*big.Float, error) {
    // Use the fallback method directly which tries multiple paths
    return uc.GetTokenPriceWithFallback(ctx, tokenAddress)
}

// QuoteExactInputSingle queries the Uniswap V3 Quoter contract
func (uc *UniswapClient) QuoteExactInputSingle(
    ctx context.Context,
    tokenIn common.Address,
    tokenOut common.Address,
    fee uint32,
    amountIn *big.Int,
) (*big.Int, error) {
    // Function selector for quoteExactInputSingle
    functionSelector := "0xf7729d43" // keccak256("quoteExactInputSingle(address,address,uint24,uint256,uint160)")[0:4]

    // Encode parameters
    params := make([]byte, 0)

    // Encode tokenIn address (pad to 32 bytes)
    addressBytes := make([]byte, 32)
    copy(addressBytes[12:], tokenIn[:])
    params = append(params, addressBytes...)

    // Encode tokenOut address (pad to 32 bytes)
    addressBytes = make([]byte, 32)
    copy(addressBytes[12:], tokenOut[:])
    params = append(params, addressBytes...)

    // Encode fee (pad to 32 bytes)
    feeBytes := make([]byte, 32)
    feeBytes[31] = byte(fee)
    feeBytes[30] = byte(fee >> 8)
    feeBytes[29] = byte(fee >> 16)
    params = append(params, feeBytes...)

    // Encode amountIn (pad to 32 bytes)
    amountInBytes := make([]byte, 32)
    amountInBytes = amountIn.FillBytes(amountInBytes)
    params = append(params, amountInBytes...)

    // Encode sqrtPriceLimitX96 (0 for no limit, pad to 32 bytes)
    sqrtPriceLimit := make([]byte, 32)
    params = append(params, sqrtPriceLimit...)

    // Combine function selector and encoded params
    data := append([]byte(functionSelector), params...)

    // Prepare call params
    callArgs := map[string]interface{}{
        "to":   UNISWAP_V3_QUOTER_ADDRESS,
        "data": "0x" + hex.EncodeToString(data),
    }

    var result string
    err := uc.client.CallContext(ctx, &result, "eth_call", callArgs, "latest")
    if err != nil {
        // Try to provide a more helpful error message
        if strings.Contains(err.Error(), "execution reverted") {
            // Try with different fee tiers if the initial one fails
            if fee == 3000 { // If we tried 0.3% fee first
                log.Printf("Quoter call with fee tier 0.3%% reverted, trying with 1%% fee tier")
                return uc.QuoteExactInputSingle(ctx, tokenIn, tokenOut, 10000, amountIn) // Try with 1% fee
            } else if fee == 10000 { // If we tried 1% fee
                log.Printf("Quoter call with fee tier 1%% reverted, trying with 0.05%% fee tier")
                return uc.QuoteExactInputSingle(ctx, tokenIn, tokenOut, 500, amountIn) // Try with 0.05% fee
            }
            
            // If all fee tiers fail, return a more descriptive error
            return nil, fmt.Errorf("no liquidity found for %s/%s on any fee tier", 
                tokenIn.Hex(), tokenOut.Hex())
        }
        return nil, err
    }

    // Remove "0x" prefix if present
    result = strings.TrimPrefix(result, "0x")

    // Decode the result (32 bytes)
    resultBytes, err := hex.DecodeString(result)
    if err != nil {
        return nil, err
    }

    // Convert to big.Int
    amountOut := new(big.Int).SetBytes(resultBytes)
    return amountOut, nil
}

func (uc *UniswapClient) GetTokenPriceWithFallback(ctx context.Context, tokenAddress string) (*big.Float, error) {
    // If the token address matches our token address, and we want price in WETH
    if common.HexToAddress(tokenAddress) == uc.Config.TokenAddress {
        // For CIFO, try all fee tiers and take the one with the best liquidity
        weth := uc.Config.WethAddress
        tokenIn := common.HexToAddress(tokenAddress)
        amountIn := big.NewInt(1e18) // 1 token
        
        // Try multiple fee tiers
        feeTiers := []uint32{3000, 10000, 500} // 0.3%, 1%, 0.05%
        
        var bestAmountOut *big.Int
        var bestFeeTier uint32
        
        for _, fee := range feeTiers {
            amountOut, err := uc.QuoteExactInputSingleNoRetry(ctx, tokenIn, weth, fee, amountIn)
            if err == nil {
                if bestAmountOut == nil || amountOut.Cmp(bestAmountOut) > 0 {
                    bestAmountOut = amountOut
                    bestFeeTier = fee
                }
            }
        }
        
        if bestAmountOut != nil {
            log.Printf("Found best liquidity for CIFO/WETH on fee tier %d", bestFeeTier)
            
            price := new(big.Float).SetInt(bestAmountOut)
            divisor := new(big.Float).SetInt(big.NewInt(1e18))
            return new(big.Float).Quo(price, divisor), nil
        }
        
        return nil, fmt.Errorf("no liquidity found for CIFO to WETH on any fee tier")
    }
    
    // For other tokens...
    // Try direct path first with different fee tiers
    feeTiers := []uint32{3000, 10000, 500} // 0.3%, 1%, 0.05%
    
    for _, fee := range feeTiers {
        tokenIn := common.HexToAddress(tokenAddress)
        tokenOut := uc.Config.StablecoinAddress
        amountIn := big.NewInt(1e18) // 1 token
        
        amountOut, err := uc.QuoteExactInputSingleNoRetry(ctx, tokenIn, tokenOut, fee, amountIn)
        if err == nil {
            price := new(big.Float).SetInt(amountOut)
            divisor := new(big.Float).SetInt(big.NewInt(1e18)) 
            return new(big.Float).Quo(price, divisor), nil
        }
    }
    
    // If direct path fails, try routing through WETH
    log.Printf("Direct price quote failed for %s, trying route through WETH", tokenAddress)
    
    // 1. Get price of token in WETH
    tokenIn := common.HexToAddress(tokenAddress)
    weth := uc.Config.WethAddress
    amountIn := big.NewInt(1e18) // 1 token
    
    var tokenWethAmount *big.Int
    var feeSuccess uint32
    
    for _, fee := range feeTiers {
        amt, err := uc.QuoteExactInputSingleNoRetry(ctx, tokenIn, weth, fee, amountIn)
        if err == nil {
            tokenWethAmount = amt
            feeSuccess = fee
            break
        }
    }
    
    if tokenWethAmount == nil {
        return nil, fmt.Errorf("cannot find liquidity for %s to WETH on any fee tier", tokenAddress)
    }
    
    log.Printf("Found liquidity for %s/WETH on fee tier %d", tokenAddress, feeSuccess)
    
    // 2. Get price of WETH in stablecoin
    wethStableAmount, err := uc.QuoteExactInputSingle(ctx, weth, uc.Config.StablecoinAddress, 3000, big.NewInt(1e18))
    if err != nil {
        return nil, fmt.Errorf("failed to get WETH/Stablecoin price: %w", err)
    }
    
    // 3. Calculate token price in stablecoin
    // token/stablecoin = (token/weth) * (weth/stablecoin)
    tokenWethPrice := new(big.Float).SetInt(tokenWethAmount)
    tokenWethPrice.Quo(tokenWethPrice, new(big.Float).SetInt(big.NewInt(1e18)))
    
    wethStablePrice := new(big.Float).SetInt(wethStableAmount)
    wethStablePrice.Quo(wethStablePrice, new(big.Float).SetInt(big.NewInt(1e18)))
    
    finalPrice := new(big.Float).Mul(tokenWethPrice, wethStablePrice)
    
    return finalPrice, nil
}

// QuoteExactInputSingleNoRetry - same as QuoteExactInputSingle but without retry logic
func (uc *UniswapClient) QuoteExactInputSingleNoRetry(
    ctx context.Context,
    tokenIn common.Address,
    tokenOut common.Address,
    fee uint32,
    amountIn *big.Int,
) (*big.Int, error) {
    // Function selector for quoteExactInputSingle
    functionSelector := "0xf7729d43"
    
    // Encode parameters
    params := make([]byte, 0)
    
    // Encode tokenIn address (pad to 32 bytes)
    addressBytes := make([]byte, 32)
    copy(addressBytes[12:], tokenIn[:])
    params = append(params, addressBytes...)
    
    // Encode tokenOut address (pad to 32 bytes)
    addressBytes = make([]byte, 32)
    copy(addressBytes[12:], tokenOut[:])
    params = append(params, addressBytes...)
    
    // Encode fee (pad to 32 bytes)
    feeBytes := make([]byte, 32)
    feeBytes[31] = byte(fee)
    feeBytes[30] = byte(fee >> 8)
    feeBytes[29] = byte(fee >> 16)
    params = append(params, feeBytes...)
    
    // Encode amountIn (pad to 32 bytes)
    amountInBytes := make([]byte, 32)
    amountInBytes = amountIn.FillBytes(amountInBytes)
    params = append(params, amountInBytes...)
    
    // Encode sqrtPriceLimitX96 (0 for no limit, pad to 32 bytes)
    sqrtPriceLimit := make([]byte, 32)
    params = append(params, sqrtPriceLimit...)
    
    // Combine function selector and encoded params
    data := append([]byte(functionSelector), params...)
    
    // Prepare call params
    callArgs := map[string]interface{}{
        "to":   UNISWAP_V3_QUOTER_ADDRESS,
        "data": "0x" + hex.EncodeToString(data),
    }
    
    var result string
    err := uc.client.CallContext(ctx, &result, "eth_call", callArgs, "latest")
    if err != nil {
        return nil, err
    }
    
    // Remove "0x" prefix if present
    result = strings.TrimPrefix(result, "0x")
    
    // Decode the result (32 bytes)
    resultBytes, err := hex.DecodeString(result)
    if err != nil {
        return nil, err
    }
    
    // Convert to big.Int
    amountOut := new(big.Int).SetBytes(resultBytes)
    return amountOut, nil
}

func (c *UniswapClient) createTransactor(ctx context.Context) (*bind.TransactOpts, error) {
    if c.privateKey == nil {
        return nil, fmt.Errorf("private key not initialized")
    }
    
    // Get the chainID
    chainID, err := c.ethClient.ChainID(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get chain ID: %v", err)
    }
    
    // Create the transaction signer
    auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
    if err != nil {
        return nil, fmt.Errorf("failed to create transactor: %v", err)
    }
    
    // Get the wallet address from private key
    publicKey := c.privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        return nil, fmt.Errorf("error casting public key to ECDSA")
    }
    
    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
    
    // Get the next nonce
    nonce, err := c.ethClient.PendingNonceAt(ctx, fromAddress)
    if err != nil {
        return nil, fmt.Errorf("failed to get nonce: %v", err)
    }
    
    auth.Nonce = big.NewInt(int64(nonce))
    auth.Value = big.NewInt(0)      // Default 0 ETH sent
    auth.GasLimit = uint64(3000000) // Default gas limit
    
    // Get gas price from the network
    gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get gas price: %v", err)
    }
    
    auth.GasPrice = gasPrice
    
    return auth, nil
}

func (c *UniswapClient) SwapExactETHForTokens(
    ctx context.Context,
    amountETH *big.Int,
    minAmountOut *big.Int,
    path []common.Address,
    to common.Address,
    deadline int64,
) (string, error) {
    // Check if router is initialized
    if c.router == nil {
        return "", fmt.Errorf("Uniswap router not initialized")
    }
    
    // Create a transactor with ETH value
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return "", err
    }
    
    // Set transaction value to ETH amount
    auth.Value = amountETH
    
    // Execute swap through router contract
    tx, err := c.router.SwapExactETHForTokens(
        auth,
        minAmountOut,
        path,
        to,
        big.NewInt(deadline),
    )
    if err != nil {
        return "", fmt.Errorf("swap failed: %v", err)
    }
    
    // Wait for transaction to be mined using ethClient instead of client
    receipt, err := waitForTransaction(ctx, c.ethClient, tx.Hash())
    if err != nil {
        return "", fmt.Errorf("waiting for swap transaction failed: %v", err)
    }
    
    if receipt.Status == 0 {
        return "", fmt.Errorf("swap transaction reverted")
    }
    
    return tx.Hash().Hex(), nil
}

// Helper function to wait for a transaction to be mined
func waitForTransaction(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
    for {
        receipt, err := client.TransactionReceipt(ctx, txHash)
        if err == nil {
            return receipt, nil
        }
        
        // Check if context is done
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
            // Continue polling
            time.Sleep(2 * time.Second)
        }
    }
}