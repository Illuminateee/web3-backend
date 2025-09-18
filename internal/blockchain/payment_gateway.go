package blockchain

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	// Import your generated bindings

	"git.winteraccess.id/walanja/web3-tokensale-be/internal/bindings/generated/fiattotokenpaymentgateway"
	"git.winteraccess.id/walanja/web3-tokensale-be/internal/bindings/generated/testtoken"
)

// PaymentStatus represents payment status in the contract
type PaymentStatus uint8

const (
    PaymentStatusPending   PaymentStatus = 0
    PaymentStatusCompleted PaymentStatus = 1
    PaymentStatusFailed    PaymentStatus = 2
    PaymentStatusRefunded  PaymentStatus = 3
)

// PaymentGatewayClient provides interaction with the payment gateway contract
type PaymentGatewayClient struct {
    client         *ethclient.Client
    contractAddr   common.Address
    contract       *fiattotokenpaymentgateway.FiatToTokenPaymentGateway
    privateKey     *ecdsa.PrivateKey
    tokenAddr      common.Address
    tokenContract  *testtoken.TestToken
}

type PaymentDetails struct {
    Buyer         common.Address
    TokenAmount   *big.Int
    FiatAmount    *big.Int
    Timestamp     *big.Int
    Gateway       string
    Status        PaymentStatus
    GasFundAmount *big.Int
    GasRefunded   bool
}

// NewPaymentGatewayClient creates a new client to interact with the payment gateway contract

func (c *PaymentGatewayClient) IsInitialized() bool {
    return c != nil && c.client != nil && c.privateKey != nil && c.contract != nil
}
func NewPaymentGatewayClient(rpcURL string, contractAddress string, privateKeyHex string) (*PaymentGatewayClient, error) {
    // Connect to Ethereum node
    client, err := ethclient.Dial(rpcURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
    }

    // Parse contract address
    contractAddr := common.HexToAddress(contractAddress)

    // Load the payment gateway contract
    contract, err := fiattotokenpaymentgateway.NewFiatToTokenPaymentGateway(contractAddr, client)
    if err != nil {
        return nil, fmt.Errorf("failed to instantiate payment gateway contract: %v", err)
    }

    // Parse private key
    privateKey, err := crypto.HexToECDSA(privateKeyHex)
    if err != nil {
        return nil, fmt.Errorf("failed to parse private key: %v", err)
    }

    // Get token address from contract
    tokenAddr, err := contract.Token(&bind.CallOpts{})
    if err != nil {
        return nil, fmt.Errorf("failed to get token address from contract: %v", err)
    }

    // Load the token contract
    tokenContract, err := testtoken.NewTestToken(tokenAddr, client)
    if err != nil {
        return nil, fmt.Errorf("failed to instantiate token contract: %v", err)
    }

    return &PaymentGatewayClient{
        client:        client,
        contractAddr:  contractAddr,
        contract:      contract,
        privateKey:    privateKey,
        tokenAddr:     tokenAddr,
        tokenContract: tokenContract,
    }, nil
}

// CreatePayment initializes a payment in the contract
func (c *PaymentGatewayClient) CreatePayment(ctx context.Context, paymentId string, tokenAmount *big.Int, fiatAmount *big.Int, gateway string, destinationWallet common.Address,gasDeposit *big.Int) (string, error) {
    // Create a keyed transactor
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return "",err
    }

    // Set value for gas deposit
    auth.Value = gasDeposit

    // Create payment
    tx, err := c.contract.CreatePayment(auth, paymentId, tokenAmount, fiatAmount, gateway,destinationWallet)
    if err != nil {
        return "",fmt.Errorf("failed to create payment: %v", err)
    }

    // Wait for the transaction to be mined
    receipt, err := bind.WaitMined(ctx, c.client, tx)
    if err != nil {
        return "",fmt.Errorf("transaction failed: %v", err)
    }

    if receipt.Status == 0 {
        return "", fmt.Errorf("transaction reverted")
    }

    txHash := receipt.TxHash.Hex()
    log.Printf("Payment created successfully: PaymentID=%s, TxHash=%s", paymentId, txHash)

    return txHash, nil
}

// ProcessPaymentCallback processes a payment notification
func (c *PaymentGatewayClient) ProcessPaymentCallback(ctx context.Context, paymentId string, status uint8, signature []byte) (string, error) {
    
    if !c.IsInitialized() {
        return "", fmt.Errorf("payment gateway client not initialized")
    }
    // If no signature provided, create one
    if signature == nil {
        var err error
        signature, err = c.createPaymentSignature(paymentId, status)
        if err != nil {
            return "", fmt.Errorf("failed to create signature: %v", err)
        }
        log.Printf("Generated payment signature for transaction")
    }

    // Create a keyed transactor
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return "", fmt.Errorf("failed to create transactor: %v", err)
    }

    
    // Get payment details first to make sure it exists
    opts := &bind.CallOpts{Context: ctx}
    payment, err := c.contract.Payments(opts, paymentId)
    if err != nil {
        return "", fmt.Errorf("error getting payment details: %v", err)
    }

    log.Printf("Processing payment: ID=%s, buyer=%s, status=%d", 
    paymentId, payment.Buyer.Hex(), payment.Status)

    // If no signature is provided, create one
    if signature == nil {
        signature, err = c.createPaymentSignature(paymentId, status)
        if err != nil {
            return "", fmt.Errorf("failed to create signature: %v", err)
        }
    }

    // Process payment with signature
    tx, err := c.contract.ProcessPaymentCallback(auth, paymentId, status, signature)
    fmt.Printf("Processing payment callback: tx=%s", tx.Hash().Hex())
    
    if err != nil {
        return "", fmt.Errorf("failed to process payment callback: %v", err)
    }

    // Wait for the transaction to be mined
    receipt, err := bind.WaitMined(ctx, c.client, tx)
    if err != nil {
        return "", fmt.Errorf("transaction failed: %v", err)
    }

    if receipt.Status == 0 {
        return "", fmt.Errorf("transaction reverted")
    }

    return tx.Hash().Hex(),nil
}

// ProcessRefund processes a refund for a payment
func (c *PaymentGatewayClient) ProcessRefund(ctx context.Context, paymentId string) error {
    // Create a keyed transactor
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return err
    }

    // Process refund
    tx, err := c.contract.ProcessRefund(auth, paymentId)
    if err != nil {
        return fmt.Errorf("failed to process refund: %v", err)
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

// GetPaymentStatus retrieves the status of a payment
func (c *PaymentGatewayClient) GetPaymentStatus(ctx context.Context, paymentId string) (PaymentStatus, error) {
    status, err := c.contract.GetPaymentStatus(&bind.CallOpts{Context: ctx}, paymentId)
    if err != nil {
        return 0, fmt.Errorf("failed to get payment status: %v", err)
    }
    return PaymentStatus(status), nil
}

// GetPaymentDetails retrieves the full details of a payment
func (c *PaymentGatewayClient) GetPaymentDetails(ctx context.Context, paymentId string) (*PaymentDetails, error) {
    payment, err := c.contract.Payments(&bind.CallOpts{Context: ctx}, paymentId)
    if err != nil {
        return nil, fmt.Errorf("failed to get payment details: %v", err)
    }
    
    // Convert the contract payment struct to our domain model
    return &PaymentDetails{
        Buyer:         payment.Buyer,
        TokenAmount:   payment.TokenAmount,
        FiatAmount:    payment.FiatAmount,
        Timestamp:     payment.Timestamp,
        Gateway:       payment.Gateway,
        Status:        PaymentStatus(payment.Status),
        GasFundAmount: payment.GasFundAmount,
        GasRefunded:   payment.GasRefunded,
    }, nil
}

// CalculateTokenAmount calculates token amount based on fiat amount
func (c *PaymentGatewayClient) CalculateTokenAmount(ctx context.Context, fiatAmount *big.Int) (*big.Int, error) {
    tokenAmount, err := c.contract.CalculateTokenAmount(&bind.CallOpts{Context: ctx}, fiatAmount)
    if err != nil {
        return nil, fmt.Errorf("failed to calculate token amount: %v", err)
    }
    return tokenAmount, nil
}

// WithdrawProcessingFees withdraws accumulated processing fees to owner
func (c *PaymentGatewayClient) WithdrawProcessingFees(ctx context.Context) error {
    // Create a keyed transactor
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return err
    }

    // Withdraw fees
    tx, err := c.contract.WithdrawProcessingFees(auth)
    if err != nil {
        return fmt.Errorf("failed to withdraw fees: %v", err)
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

// UpdateGasDepositRequirement updates the required gas deposit
func (c *PaymentGatewayClient) UpdateGasDepositRequirement(ctx context.Context, amount *big.Int) error {
    // Create a keyed transactor
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return err
    }

    // Update gas deposit requirement
    tx, err := c.contract.UpdateGasDepositRequirement(auth, amount)
    if err != nil {
        return fmt.Errorf("failed to update gas deposit requirement: %v", err)
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

// UpdateTokenPrice updates the token price
func (c *PaymentGatewayClient) UpdateTokenPrice(ctx context.Context, pricePerToken *big.Int) error {
    // Create a keyed transactor
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return err
    }

    // Update token price
    tx, err := c.contract.UpdateTokenPrice(auth, pricePerToken)
    if err != nil {
        return fmt.Errorf("failed to update token price: %v", err)
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

// UpdateGatewaySigner updates a gateway signer address
func (c *PaymentGatewayClient) UpdateGatewaySigner(ctx context.Context, gateway string, signer common.Address) error {
    // Create a keyed transactor
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return err
    }

    // Update gateway signer
    tx, err := c.contract.UpdateGatewaySigner(auth, gateway, signer)
    if err != nil {
        return fmt.Errorf("failed to update gateway signer: %v", err)
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

// MockPaymentCallback simulates a payment callback for testing
func (c *PaymentGatewayClient) MockPaymentCallback(ctx context.Context, paymentId string, status uint8) error {
    // Create a keyed transactor
    auth, err := c.createTransactor(ctx)
    if err != nil {
        return err
    }

    // Mock payment callback
    tx, err := c.contract.MockPaymentCallback(auth, paymentId, status)
    if err != nil {
        return fmt.Errorf("failed to mock payment callback: %v", err)
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

// Helper function to create a transactor
func (c *PaymentGatewayClient) createTransactor(ctx context.Context) (*bind.TransactOpts, error) {
    

    if c == nil {
        return nil, fmt.Errorf("payment gateway client is nil")
    }

    if c.client == nil {
        return nil, fmt.Errorf("ethereum client is nil, client not properly initialized")
    }

    if c.privateKey == nil {
        return nil, fmt.Errorf("private key is nil, check key initialization")
    }
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
    auth.GasLimit = uint64(3000000)
    auth.GasPrice = gasPrice
    auth.Context = ctx

    return auth, nil
}

// GetOwner returns the contract owner
func (c *PaymentGatewayClient) GetOwner(ctx context.Context) (common.Address, error) {
    return c.contract.Owner(&bind.CallOpts{Context: ctx})
}

// GetTokenAddress returns the address of the token contract
func (c *PaymentGatewayClient) GetTokenAddress(ctx context.Context) (common.Address, error) {
    return c.contract.Token(&bind.CallOpts{Context: ctx})
}

// GetRequiredGasDeposit returns the required gas deposit
func (c *PaymentGatewayClient) GetRequiredGasDeposit(ctx context.Context) (*big.Int, error) {
    return c.contract.RequiredGasDeposit(&bind.CallOpts{Context: ctx})
}

// GetPricePerToken returns the token price
func (c *PaymentGatewayClient) GetPricePerToken(ctx context.Context) (*big.Int, error) {
    return c.contract.PricePerToken(&bind.CallOpts{Context: ctx})
}

func (c *PaymentGatewayClient) CheckPaymentExists(ctx context.Context, paymentId string) (bool, error) {
    // Get payment details from contract
    payment, err := c.contract.Payments(&bind.CallOpts{Context: ctx}, paymentId)

    
    if err != nil {
        return false, fmt.Errorf("failed to check payment: %v", err)
    }
    
    // Log payment details for debugging
    log.Printf("CheckPaymentExists for ID=%s: tokenAmount=%s, status=%d", 
               paymentId, payment.TokenAmount.String(), payment.Status)
    
    // A payment exists if it has a non-zero token amount
    zero := big.NewInt(0)
    exists := payment.TokenAmount.Cmp(zero) > 0
    
    if exists {
        log.Printf("Payment %s found in contract", paymentId)
    } else {
        log.Printf("Payment %s NOT found in contract (token amount is zero)", paymentId)
    }
    
    return exists, nil
}

func (c *PaymentGatewayClient) createPaymentSignature(paymentId string, status uint8) ([]byte, error) {
    // Encode like abi.encodePacked(paymentId, status)
    var buf bytes.Buffer
    buf.WriteString(paymentId)
    buf.WriteByte(status)  // status is uint8

    // Get the hash
    packed := buf.Bytes()
    hash := crypto.Keccak256Hash(packed)

    // Apply Ethereum signed message prefix
    prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n32")
    prefixedHash := crypto.Keccak256Hash(append([]byte(prefix), hash.Bytes()...))

    // Sign
    signature, err := crypto.Sign(prefixedHash.Bytes(), c.privateKey)
    if err != nil {
        return nil, err
    }

    // Add 27 to the recovery ID if needed
    if signature[64] < 27 {
        signature[64] += 27
    }

    return signature, nil
}

func (c *PaymentGatewayClient) GetEthClient() *ethclient.Client {
    return c.client
}

func (c *PaymentGatewayClient) GetTokenBalance(ctx context.Context, address common.Address) (*big.Int, error) {
    // Use the already instantiated token contract instead of creating a new one
    callOpts := &bind.CallOpts{Context: ctx}
    
    // Call balanceOf function directly using the token contract instance
    balance, err := c.tokenContract.BalanceOf(callOpts, address)
    if err != nil {
        return nil, fmt.Errorf("failed to get token balance: %w", err)
    }
    
    return balance, nil
}