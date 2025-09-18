package blockchain

import (
    "math/big"
    "strings"

    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

// ERC20 represents an ERC20 token contract
type ERC20 struct {
    client   *ethclient.Client
    address  common.Address
    contract *bind.BoundContract
}

// NewERC20 creates a new ERC20 token contract instance
func NewERC20(address common.Address, client *ethclient.Client) (*ERC20, error) {
    // Get the ERC20 ABI
    parsedABI, err := abi.JSON(strings.NewReader(ERC20ABI))
    if err != nil {
        return nil, err
    }
    
    // Create a bound contract
    contract := bind.NewBoundContract(address, parsedABI, client, client, client)
    
    return &ERC20{
        client:   client,
        address:  address,
        contract: contract,
    }, nil
}

// BalanceOf gets the token balance of an account
func (e *ERC20) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
    var out []interface{}
    err := e.contract.Call(opts, &out, "balanceOf", account)
    if err != nil {
        return nil, err
    }
    
    if len(out) == 0 {
        return big.NewInt(0), nil
    }
    
    // Convert the first result to *big.Int
    return out[0].(*big.Int), nil
}

// Approve approves the spender to spend tokens
func (e *ERC20) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
    return e.contract.Transact(opts, "approve", spender, amount)
}

// Transfer transfers tokens to a recipient
func (e *ERC20) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
    return e.contract.Transact(opts, "transfer", recipient, amount)
}

// TransferFrom transfers tokens from one account to another
func (e *ERC20) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
    return e.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// Decimals gets the token decimals
func (e *ERC20) Decimals(opts *bind.CallOpts) (uint8, error) {
    var out []interface{}
    err := e.contract.Call(opts, &out, "decimals")
    if err != nil {
        return 0, err
    }
    
    if len(out) == 0 {
        return 0, nil
    }
    
    // Convert the first result to uint8
    return uint8(out[0].(uint8)), nil
}

// Symbol gets the token symbol
func (e *ERC20) Symbol(opts *bind.CallOpts) (string, error) {
    var out []interface{}
    err := e.contract.Call(opts, &out, "symbol")
    if err != nil {
        return "", err
    }
    
    if len(out) == 0 {
        return "", nil
    }
    
    // Convert the first result to string
    return out[0].(string), nil
}

// Name gets the token name
func (e *ERC20) Name(opts *bind.CallOpts) (string, error) {
    var out []interface{}
    err := e.contract.Call(opts, &out, "name")
    if err != nil {
        return "", err
    }
    
    if len(out) == 0 {
        return "", nil
    }
    
    // Convert the first result to string
    return out[0].(string), nil
}

// Allowance gets the remaining number of tokens that spender is allowed to spend
func (e *ERC20) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
    var out []interface{}
    err := e.contract.Call(opts, &out, "allowance", owner, spender)
    if err != nil {
        return nil, err
    }
    
    if len(out) == 0 {
        return big.NewInt(0), nil
    }
    
    // Convert the first result to *big.Int
    return out[0].(*big.Int), nil
}

// ERC20ABI is the standard ERC20 ABI
const ERC20ABI = `[
    {
        "constant": true,
        "inputs": [],
        "name": "name",
        "outputs": [{"name": "", "type": "string"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {"name": "_spender", "type": "address"},
            {"name": "_value", "type": "uint256"}
        ],
        "name": "approve",
        "outputs": [{"name": "", "type": "bool"}],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "totalSupply",
        "outputs": [{"name": "", "type": "uint256"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {"name": "_from", "type": "address"},
            {"name": "_to", "type": "address"},
            {"name": "_value", "type": "uint256"}
        ],
        "name": "transferFrom",
        "outputs": [{"name": "", "type": "bool"}],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "decimals",
        "outputs": [{"name": "", "type": "uint8"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [{"name": "_owner", "type": "address"}],
        "name": "balanceOf",
        "outputs": [{"name": "balance", "type": "uint256"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "symbol",
        "outputs": [{"name": "", "type": "string"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {"name": "_to", "type": "address"},
            {"name": "_value", "type": "uint256"}
        ],
        "name": "transfer",
        "outputs": [{"name": "", "type": "bool"}],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {"name": "_owner", "type": "address"},
            {"name": "_spender", "type": "address"}
        ],
        "name": "allowance",
        "outputs": [{"name": "", "type": "uint256"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "payable": true,
        "stateMutability": "payable",
        "type": "fallback"
    },
    {
        "anonymous": false,
        "inputs": [
            {"indexed": true, "name": "owner", "type": "address"},
            {"indexed": true, "name": "spender", "type": "address"},
            {"indexed": false, "name": "value", "type": "uint256"}
        ],
        "name": "Approval",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {"indexed": true, "name": "from", "type": "address"},
            {"indexed": true, "name": "to", "type": "address"},
            {"indexed": false, "name": "value", "type": "uint256"}
        ],
        "name": "Transfer",
        "type": "event"
    }
]`