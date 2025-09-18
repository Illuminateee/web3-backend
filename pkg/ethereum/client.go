package ethereum

import (
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
    RPCClient *rpc.Client     // Low-level RPC client
    Client    *ethclient.Client // High-level Ethereum client for contract interactions
}

func NewClient(url string) (*Client, error) {
    // Create the RPC client
    rpcClient, err := rpc.Dial(url)
    if err != nil {
        return nil, err
    }
    
    // Create the ethclient.Client from the RPC client
    ethClient := ethclient.NewClient(rpcClient)
    
    return &Client{
        RPCClient: rpcClient,
        Client:    ethClient,
    }, nil
}

func (c *Client) Close() {
    if c.RPCClient != nil {
        c.RPCClient.Close()
    }
    // The ethclient.Client doesn't need separate closing as it uses the underlying RPC connection
}