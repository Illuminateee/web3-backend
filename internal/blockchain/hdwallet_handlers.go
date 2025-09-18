package blockchain

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

// HDWallet represents a hierarchical deterministic wallet
type HDWallet struct {
    Mnemonic  string
    Seed      []byte
    MasterKey *hdkeychain.ExtendedKey
}

// Account represents an Ethereum account
type Account struct {
    Address    common.Address
    PrivateKey *ecdsa.PrivateKey
}

// NewHDWallet creates a new wallet with a random mnemonic
func NewHDWallet() (*HDWallet, error) {
    // Generate entropy for mnemonic
    entropy, err := bip39.NewEntropy(128) // 128 bits = 12 words
    if err != nil {
        return nil, fmt.Errorf("failed to generate entropy: %v", err)
    }

    // Generate mnemonic
    mnemonic, err := bip39.NewMnemonic(entropy)
    if err != nil {
        return nil, fmt.Errorf("failed to generate mnemonic: %v", err)
    }

    return NewHDWalletFromMnemonic(mnemonic)
}

// NewHDWalletFromMnemonic creates a wallet from an existing mnemonic
func NewHDWalletFromMnemonic(mnemonic string) (*HDWallet, error) {
    // Validate mnemonic
    if !bip39.IsMnemonicValid(mnemonic) {
        return nil, fmt.Errorf("invalid mnemonic")
    }

    // Generate seed from mnemonic
    seed := bip39.NewSeed(mnemonic, "")

    // Generate master key
    masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
    if err != nil {
        return nil, fmt.Errorf("failed to generate master key: %v", err)
    }

    return &HDWallet{
        Mnemonic:  mnemonic,
        Seed:      seed,
        MasterKey: masterKey,
    }, nil
}

// DeriveAccount derives a new account from the wallet using the standard Ethereum derivation path
func (w *HDWallet) DeriveAccount(index uint32) (*Account, error) {
    path := fmt.Sprintf("m/44'/60'/0'/0/%d", index)
    return w.DeriveAccountFromPath(path)
}

// DeriveAccountFromPath derives a new account using a custom path
func (w *HDWallet) DeriveAccountFromPath(path string) (*Account, error) {
    // Parse derivation path
    parsedPath, err := accounts.ParseDerivationPath(path)
    if err != nil {
        return nil, fmt.Errorf("invalid derivation path: %v", err)
    }

    var key = w.MasterKey
    for _, n := range parsedPath {
        key, err = key.Derive(n)
        if err != nil {
            return nil, fmt.Errorf("failed to derive key: %v", err)
        }
    }

    // Get private key
    privateKey, err := key.ECPrivKey()
    if err != nil {
        return nil, fmt.Errorf("failed to get private key: %v", err)
    }

    // Convert to Ethereum private key
    privateKeyECDSA := privateKey.ToECDSA()
    
    // Derive Ethereum address
    publicKey := privateKeyECDSA.PublicKey
    address := crypto.PubkeyToAddress(publicKey)

    return &Account{
        Address:    address,
        PrivateKey: privateKeyECDSA,
    }, nil
}

// GetPrivateKeyHex returns the hex representation of the private key
func (a *Account) GetPrivateKeyHex() string {
    privateKeyBytes := crypto.FromECDSA(a.PrivateKey)
    return hexutil.Encode(privateKeyBytes)[2:] // Remove 0x prefix
}

// ExportKeystore creates a keystore JSON for the account
func (a *Account) ExportKeystore(password string) ([]byte, error) {
    key := &keystore.Key{
        Address:    a.Address,
        PrivateKey: a.PrivateKey,
    }
    return keystore.EncryptKey(key, password, keystore.StandardScryptN, keystore.StandardScryptP)
}