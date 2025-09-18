package contracts

import (
    "math/big"
    "strings"

    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
)

// UniswapV2Router02 is a Go binding of the Uniswap V2 Router contract
type UniswapV2Router02 struct {
    UniswapV2Router02Caller     // Read-only binding to the contract
    UniswapV2Router02Transactor // Write-only binding to the contract
    UniswapV2Router02Filterer   // Log filterer for contract events
}

// UniswapV2Router02Caller is an auto generated read-only Go binding of the contract
type UniswapV2Router02Caller struct {
    contract *bind.BoundContract
}

// UniswapV2Router02Transactor is an auto generated write-only Go binding of the contract
type UniswapV2Router02Transactor struct {
    contract *bind.BoundContract
}

// UniswapV2Router02Filterer is an auto generated log filtering Go binding of the contract
type UniswapV2Router02Filterer struct {
    contract *bind.BoundContract
}

// NewUniswapV2Router02 creates a new instance of UniswapV2Router02, bound to a specific deployed contract
func NewUniswapV2Router02(address common.Address, backend bind.ContractBackend) (*UniswapV2Router02, error) {
    abi, err := abi.JSON(strings.NewReader(UniswapV2Router02ABI))
    if err != nil {
        return nil, err
    }
    contract := bind.NewBoundContract(address, abi, backend, backend, backend)
    return &UniswapV2Router02{
        UniswapV2Router02Caller:     UniswapV2Router02Caller{contract: contract},
        UniswapV2Router02Transactor: UniswapV2Router02Transactor{contract: contract},
        UniswapV2Router02Filterer:   UniswapV2Router02Filterer{contract: contract},
    }, nil
}

// SwapExactETHForTokens is a paid mutator transaction binding the contract method
func (t *UniswapV2Router02Transactor) SwapExactETHForTokens(opts *bind.TransactOpts, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
    return t.contract.Transact(opts, "swapExactETHForTokens", amountOutMin, path, to, deadline)
}

// SwapExactTokensForETH is a paid mutator transaction binding the contract method
func (t *UniswapV2Router02Transactor) SwapExactTokensForETH(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
    return t.contract.Transact(opts, "swapExactTokensForETH", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method
func (t *UniswapV2Router02Transactor) SwapExactTokensForTokens(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
    return t.contract.Transact(opts, "swapExactTokensForTokens", amountIn, amountOutMin, path, to, deadline)
}

// GetAmountsOut is a free caller binding the contract method
func (c *UniswapV2Router02Caller) GetAmountsOut(opts *bind.CallOpts, amountIn *big.Int, path []common.Address) ([]*big.Int, error) {
    var out []interface{}
    err := c.contract.Call(opts, &out, "getAmountsOut", amountIn, path)
    
    if err != nil {
        return nil, err
    }
    
    // Convert the result to []*big.Int
    amounts := make([]*big.Int, len(out))
    for i, val := range out {
        amounts[i] = val.(*big.Int)
    }
    
    return amounts, nil
}

// Then update your ABI string to include all methods
// UniswapV2Router02ABI is the input ABI used to generate the binding
const UniswapV2Router02ABI = `[
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "amountOutMin",
                "type": "uint256"
            },
            {
                "internalType": "address[]",
                "name": "path",
                "type": "address[]"
            },
            {
                "internalType": "address",
                "name": "to",
                "type": "address"
            },
            {
                "internalType": "uint256",
                "name": "deadline",
                "type": "uint256"
            }
        ],
        "name": "swapExactETHForTokens",
        "outputs": [
            {
                "internalType": "uint256[]",
                "name": "amounts",
                "type": "uint256[]"
            }
        ],
        "stateMutability": "payable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "amountIn",
                "type": "uint256"
            },
            {
                "internalType": "uint256",
                "name": "amountOutMin",
                "type": "uint256"
            },
            {
                "internalType": "address[]",
                "name": "path",
                "type": "address[]"
            },
            {
                "internalType": "address",
                "name": "to",
                "type": "address"
            },
            {
                "internalType": "uint256",
                "name": "deadline",
                "type": "uint256"
            }
        ],
        "name": "swapExactTokensForETH",
        "outputs": [
            {
                "internalType": "uint256[]",
                "name": "amounts",
                "type": "uint256[]"
            }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "amountIn",
                "type": "uint256"
            },
            {
                "internalType": "uint256",
                "name": "amountOutMin",
                "type": "uint256"
            },
            {
                "internalType": "address[]",
                "name": "path",
                "type": "address[]"
            },
            {
                "internalType": "address",
                "name": "to",
                "type": "address"
            },
            {
                "internalType": "uint256",
                "name": "deadline",
                "type": "uint256"
            }
        ],
        "name": "swapExactTokensForTokens",
        "outputs": [
            {
                "internalType": "uint256[]",
                "name": "amounts",
                "type": "uint256[]"
            }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "amountIn",
                "type": "uint256"
            },
            {
                "internalType": "address[]",
                "name": "path",
                "type": "address[]"
            }
        ],
        "name": "getAmountsOut",
        "outputs": [
            {
                "internalType": "uint256[]",
                "name": "amounts",
                "type": "uint256[]"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    }
]`