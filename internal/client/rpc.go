package client

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// RPCClient is a JSON-RPC based implementation of the Client interface.
//
// It connects to an Ethereum node via standard JSON-RPC endpoints and
// provides execution-layer transaction inspection capabilities.
type RPCClient struct {
	client *ethclient.Client
	rpcURL string
}

// NewRPCClient creates a new RPC client connected to the specified RPC endpoint.
//
// The endpoint should be a full URL (e.g., "http://localhost:8545" or
// "https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY").
func NewRPCClient(ctx context.Context, rpcURL string) (*RPCClient, error) {
	if rpcURL == "" {
		return nil, errors.New("RPC URL cannot be empty")
	}

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, err
	}

	return &RPCClient{
		client: client,
		rpcURL: rpcURL,
	}, nil
}

// GetTransaction retrieves a transaction by hash.
func (c *RPCClient) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error) {
	tx, isPending, err := c.client.TransactionByHash(ctx, txHash)
	if err != nil {
		// Check if it's a "not found" error
		if errors.Is(err, ethereum.NotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return tx, isPending, nil
}

// GetTransactionReceipt retrieves a transaction receipt by hash.
func (c *RPCClient) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := c.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		// Check if it's a "not found" error
		if errors.Is(err, ethereum.NotFound) {
			return nil, nil
		}
		return nil, err
	}
	return receipt, nil
}

// GetBlockHeader retrieves a block header by number.
func (c *RPCClient) GetBlockHeader(ctx context.Context, blockNumber *big.Int) (*types.Header, error) {
	header, err := c.client.HeaderByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	return header, nil
}

// TraceTransaction generates an execution trace for a transaction.
//
// Note: This requires a debug-enabled node (e.g., Geth with --http.api eth,debug).
// The returned value is a raw trace result that should be parsed by the tracer package.
func (c *RPCClient) TraceTransaction(ctx context.Context, txHash common.Hash) (interface{}, error) {
	// For now, return an error indicating this needs to be implemented via
	// direct RPC calls to debug_traceTransaction.
	// We'll implement this properly when we add the tracer implementation.
	return nil, errors.New("trace transaction not yet implemented via RPC client")
}

// Close closes the client connection.
func (c *RPCClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

