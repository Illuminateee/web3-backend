package blockchain

import (
    "context"
    "math/big"

    "git.winteraccess.id/walanja/web3-tokensale-be/pkg/contracts"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

// UniswapRouter represents a Uniswap V2 router contract interface
type UniswapRouter struct {
    client    *ethclient.Client
    address   common.Address
    contract  *contracts.UniswapV2Router02
}

// NewUniswapRouter creates a new instance of the Uniswap router
func NewUniswapRouter(address common.Address, client *ethclient.Client) (*UniswapRouter, error) {
    contract, err := contracts.NewUniswapV2Router02(address, client)
    if err != nil {
        return nil, err
    }
    
    return &UniswapRouter{
        client:   client,
        address:  address,
        contract: contract,
    }, nil
}

// SwapExactETHForTokens swaps an exact amount of ETH for tokens
func (u *UniswapRouter) SwapExactETHForTokens(
    opts *bind.TransactOpts,
    amountOutMin *big.Int,
    path []common.Address,
    to common.Address,
    deadline *big.Int,
) (*types.Transaction, error) {
    // The method is likely capitalized differently in the binding
    return u.contract.SwapExactETHForTokens(opts, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETH swaps an exact amount of tokens for ETH
func (u *UniswapRouter) SwapExactTokensForETH(
    opts *bind.TransactOpts,
    amountIn *big.Int,
    amountOutMin *big.Int,
    path []common.Address,
    to common.Address,
    deadline *big.Int,
) (*types.Transaction, error) {
    // Use the corrected method name - check the actual binding for exact casing
    // Try one of these:
    return u.contract.SwapExactTokensForETH(opts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokens swaps an exact amount of tokens for another token
func (u *UniswapRouter) SwapExactTokensForTokens(
    opts *bind.TransactOpts,
    amountIn *big.Int,
    amountOutMin *big.Int,
    path []common.Address,
    to common.Address,
    deadline *big.Int,
) (*types.Transaction, error) {
    return u.contract.SwapExactTokensForTokens(opts, amountIn, amountOutMin, path, to, deadline)
}

// GetAmountsOut calculates the expected output amounts for a swap
func (u *UniswapRouter) GetAmountsOut(
    amountIn *big.Int,
    path []common.Address,
) ([]*big.Int, error) {
    opts := &bind.CallOpts{Context: context.Background()}
    return u.contract.GetAmountsOut(opts, amountIn, path)
}