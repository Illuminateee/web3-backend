package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/bindings/generated/testtoken"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TokenClient provides interaction with the TestToken contract
type TokenClient struct {
    client       *ethclient.Client
    contractAddr common.Address
    contract     *testtoken.TestToken
    privateKey   *ecdsa.PrivateKey
}

// NewTokenClient creates a new client to interact with the TestToken contract
func NewTokenClient(rpcURL string, contractAddress string, privateKeyHex string) (*TokenClient, error) {
    // Connect to Ethereum node
    client, err := ethclient.Dial(rpcURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
    }

    // Parse contract address
    contractAddr := common.HexToAddress(contractAddress)

    // Load the token contract
    contract, err := testtoken.NewTestToken(contractAddr, client)
    if err != nil {
        return nil, fmt.Errorf("failed to instantiate token contract: %v", err)
    }

    // Parse private key if provided
    var privateKey *ecdsa.PrivateKey
    if privateKeyHex != "" {
        privateKey, err = crypto.HexToECDSA(privateKeyHex)
        if err != nil {
            return nil, fmt.Errorf("failed to parse private key: %v", err)
        }
    }

    return &TokenClient{
        client:       client,
        contractAddr: contractAddr,
        contract:     contract,
        privateKey:   privateKey,
    }, nil
}

// GetName returns the name of the token
func (c *TokenClient) GetName(ctx context.Context) (string, error) {
    return c.contract.Name(&bind.CallOpts{Context: ctx})
}

// GetSymbol returns the symbol of the token
func (c *TokenClient) GetSymbol(ctx context.Context) (string, error) {
    return c.contract.Symbol(&bind.CallOpts{Context: ctx})
}

// GetDecimals returns the number of decimals for the token
func (c *TokenClient) GetDecimals(ctx context.Context) (uint8, error) {
    return c.contract.Decimals(&bind.CallOpts{Context: ctx})
}

// GetTotalSupply returns the total supply of the token
func (c *TokenClient) GetTotalSupply(ctx context.Context) (*big.Int, error) {
    return c.contract.TotalSupply(&bind.CallOpts{Context: ctx})
}

// GetBalance returns the balance of the given address
func (c *TokenClient) GetBalance(ctx context.Context, address common.Address) (*big.Int, error) {
    return c.contract.BalanceOf(&bind.CallOpts{Context: ctx}, address)
}

// GetAllowance returns the amount of tokens that the spender is allowed to spend on behalf of the owner
func (c *TokenClient) GetAllowance(ctx context.Context, owner, spender common.Address) (*big.Int, error) {
    return c.contract.Allowance(&bind.CallOpts{Context: ctx}, owner, spender)
}

// Transfer transfers tokens to the given address
func (c *TokenClient) Transfer(ctx context.Context, to common.Address, amount *big.Int) error {
    if c.privateKey == nil {
        return fmt.Errorf("private key not provided")
    }

    // Create a keyed transactor
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return err
    }

    // Transfer tokens
    tx, err := c.contract.Transfer(auth, to, amount)
    if err != nil {
        return fmt.Errorf("failed to transfer tokens: %v", err)
    }

    // Wait for the transaction to be mined
    receipt, err := bind.WaitMined(ctx, c.client, tx)
    if err != nil {
        return fmt.Errorf("transaction failed: %v", err)
    }

    if receipt.Status == 0 {
        return fmt.Errorf("transaction reverted")
    }

    return nil
}

// Approve approves the spender to spend the given amount of tokens
func (c *TokenClient) Approve(ctx context.Context, spender common.Address, amount *big.Int) error {
    if c.privateKey == nil {
        return fmt.Errorf("private key not provided")
    }

    // Create a keyed transactor
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return err
    }

    // Approve tokens
    tx, err := c.contract.Approve(auth, spender, amount)
    if err != nil {
        return fmt.Errorf("failed to approve tokens: %v", err)
    }

    // Wait for the transaction to be mined
    receipt, err := bind.WaitMined(ctx, c.client, tx)
    if err != nil {
        return fmt.Errorf("transaction failed: %v", err)
    }

    if receipt.Status == 0 {
        return fmt.Errorf("transaction reverted")
    }

    return nil
}

// TransferFrom transfers tokens from one address to another
func (c *TokenClient) TransferFrom(ctx context.Context, from, to common.Address, amount *big.Int) error {
    if c.privateKey == nil {
        return fmt.Errorf("private key not provided")
    }

    // Create a keyed transactor
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return err
    }

    // Transfer tokens
    tx, err := c.contract.TransferFrom(auth, from, to, amount)
    if err != nil {
        return fmt.Errorf("failed to transfer tokens: %v", err)
    }

    // Wait for the transaction to be mined
    receipt, err := bind.WaitMined(ctx, c.client, tx)
    if err != nil {
        return fmt.Errorf("transaction failed: %v", err)
    }

    if receipt.Status == 0 {
        return fmt.Errorf("transaction reverted")
    }

    return nil
}

// Mint mints new tokens and assigns them to the given address
// Only the contract owner can call this function
func (c *TokenClient) Mint(ctx context.Context, to common.Address, amount *big.Int) error {
    if c.privateKey == nil {
        return fmt.Errorf("private key not provided")
    }

    // Create a keyed transactor
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return err
    }

    // Mint tokens
    tx, err := c.contract.Mint(auth, to, amount)
    if err != nil {
        return fmt.Errorf("failed to mint tokens: %v", err)
    }

    // Wait for the transaction to be mined
    receipt, err := bind.WaitMined(ctx, c.client, tx)
    if err != nil {
        return fmt.Errorf("transaction failed: %v", err)
    }

    if receipt.Status == 0 {
        return fmt.Errorf("transaction reverted")
    }

    return nil
}

// GetOwner returns the contract owner
func (c *TokenClient) GetOwner(ctx context.Context) (common.Address, error) {
    return c.contract.Owner(&bind.CallOpts{Context: ctx})
}

// Helper function to create a transactor
func (c *TokenClient) createTransactor(ctx context.Context) (*bind.TransactOpts, error) {
    // Get the current nonce
    nonce, err := c.client.PendingNonceAt(ctx, crypto.PubkeyToAddress(c.privateKey.PublicKey))
    if err != nil {
        return nil, fmt.Errorf("failed to get nonce: %v", err)
    }

    // Get gas price
    gasPrice, err := c.client.SuggestGasPrice(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to suggest gas price: %v", err)
    }

    // Get chain ID
    chainID, err := c.client.NetworkID(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get chain ID: %v", err)
    }

    // Create auth
    auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
    if err != nil {
        return nil, fmt.Errorf("failed to create transactor: %v", err)
    }

    auth.Nonce = big.NewInt(int64(nonce))
    auth.Value = big.NewInt(0)
    auth.GasLimit = uint64(300000)
    auth.GasPrice = gasPrice
    auth.Context = ctx

    return auth, nil
}