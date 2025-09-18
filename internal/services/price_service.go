package services

import (
    "context"
    "errors"
    "fmt"
    "math/big"
    "time"

    "git.winteraccess.id/walanja/web3-tokensale-be/pkg/ethereum"
    "github.com/ethereum/go-ethereum/common"
)

type PriceService struct {
    client *ethereum.Client
    uniswap *ethereum.UniswapClient
}

const (
    ETH_ADDRESS = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" // WETH contract address
)

func NewPriceService(client *ethereum.Client, uniswap *ethereum.UniswapClient) *PriceService {
    return &PriceService{client: client, uniswap: uniswap}
}

func (ps *PriceService) GetTokenPrice(ctx context.Context, tokenAddress string) (*big.Float, error) {
    price, err := ps.uniswap.GetTokenPrice(ctx, tokenAddress)
    if err != nil {
        return nil, err
    }
    return price, nil
}

func (ps *PriceService) ConvertTokenPrice(ctx context.Context, tokenAddress string, amount *big.Float) (*big.Float, error) {
    price, err := ps.GetTokenPrice(ctx, tokenAddress)
    if err != nil {
        return nil, err
    }

    if price == nil {
        return nil, errors.New("price not found")
    }

    convertedAmount := new(big.Float).Mul(price, amount)
    return convertedAmount, nil
}

func (ps *PriceService) FetchPriceWithTimeout(ctx context.Context, tokenAddress string, timeout time.Duration) (*big.Float, error) {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    price, err := ps.GetTokenPrice(ctx, tokenAddress)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return nil, errors.New("request timed out")
        }
        return nil, err
    }

    return price, nil
}

func (ps *PriceService) ConvertTokenToEth(ctx context.Context, tokenAddress string, amount *big.Float) (*big.Float, error) {
    // First, get the token price in ETH
    tokenPrice, err := ps.GetTokenPriceInEth(ctx, tokenAddress)
    if err != nil {
        return nil, err
    }
    
    // Multiply amount by the token price to get the ETH equivalent
    ethAmount := new(big.Float).Mul(amount, tokenPrice)
    return ethAmount, nil
}

func (ps *PriceService) ConvertEthToToken(ctx context.Context, tokenAddress string, ethAmount *big.Float) (*big.Float, error) {
    // First, get the token price in ETH
    tokenPrice, err := ps.GetTokenPriceInEth(ctx, tokenAddress)
    if err != nil {
        return nil, err
    }
    
    // Ensure we don't divide by zero
    if tokenPrice.Cmp(new(big.Float).SetFloat64(0)) == 0 {
        return nil, errors.New("token price is zero")
    }
    
    // Divide ETH amount by the token price to get the token equivalent
    tokenAmount := new(big.Float).Quo(ethAmount, tokenPrice)
    return tokenAmount, nil
}

func (ps *PriceService) GetTokenPriceInEth(ctx context.Context, tokenAddress string) (*big.Float, error) {
    // Special case: if the token is ETH itself, return 1
    if tokenAddress == ps.uniswap.Config.WethAddress.Hex() {
        return big.NewFloat(1.0), nil
    }
    
    // Special case: if the token is our CIFO token
    if tokenAddress == ps.uniswap.Config.TokenAddress.Hex() {
        // Use the direct method you implemented in UniswapClient
        return ps.uniswap.GetTokenPrice(ctx, tokenAddress)
    }
    
    // Use Uniswap to get the price of the token in ETH
    tokenIn := common.HexToAddress(tokenAddress)
    tokenOut := ps.uniswap.Config.WethAddress
    amountIn := big.NewInt(1e18) // 1 token
    fee := uint32(3000) // 0.3%
    
    amountOut, err := ps.uniswap.QuoteExactInputSingle(ctx, tokenIn, tokenOut, fee, amountIn)
    if err != nil {
        return nil, err
    }
    
    // Convert to big.Float with proper decimal handling
    price := new(big.Float).SetInt(amountOut)
    divisor := new(big.Float).SetInt(big.NewInt(1e18))
    return new(big.Float).Quo(price, divisor), nil
}

func (ps *PriceService) GetWethAddress() common.Address {
    return ps.uniswap.Config.WethAddress
}

// CheckPoolLiquidity checks if a specific pool has liquidity
func (ps *PriceService) CheckPoolLiquidity(ctx context.Context, tokenInAddress, tokenOutAddress string, feeTier uint32) (string, error) {
    tokenIn := common.HexToAddress(tokenInAddress)
    tokenOut := common.HexToAddress(tokenOutAddress)
    amountIn := big.NewInt(1e18) // 1 token
    
    amount, err := ps.uniswap.QuoteExactInputSingleNoRetry(ctx, tokenIn, tokenOut, feeTier, amountIn)
    if err != nil {
        return "0", err
    }
    
    // Convert to readable format
    amountFloat := new(big.Float).SetInt(amount)
    divisor := new(big.Float).SetInt(big.NewInt(1e18))
    result := new(big.Float).Quo(amountFloat, divisor)
    
    // Convert to string with fixed precision
    resultStr := fmt.Sprintf("%.8f", result)
    
    return resultStr, nil
}

func (ps *PriceService) GetTokenAddress() common.Address {
    return ps.uniswap.Config.TokenAddress
}

func (ps *PriceService) GetUniswapClient() *ethereum.UniswapClient {
    return ps.uniswap
}

// SwapExactETHForTokens delegates to the UniswapClient
func (ps *PriceService) SwapExactETHForTokens(
    ctx context.Context,
    amountETH *big.Int,
    minAmountOut *big.Int,
    path []common.Address,
    to common.Address,
    deadline int64,
) (string, error) {
    // Delegate to the UniswapClient
    return ps.uniswap.SwapExactETHForTokens(
        ctx,
        amountETH,
        minAmountOut,
        path,
        to,
        deadline,
    )
}