package client

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Client provides an interface for interacting with Ethereum execution clients.
//
// Implementations should map client-specific formats (JSON-RPC, direct Geth
// internals, etc.) into these normalized, execution-layer focused types.
type Client interface {
	// GetTransaction retrieves a transaction by hash.
	// Returns nil, nil if the transaction is not found.
	GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)

	// GetTransactionReceipt retrieves a transaction receipt by hash.
	// Returns nil, nil if the receipt is not found.
	GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)

	// GetBlockHeader retrieves a block header by number.
	// This is needed for base fee and other block-level context.
	GetBlockHeader(ctx context.Context, blockNumber *big.Int) (*types.Header, error)

	// TraceTransaction generates an execution trace for a transaction.
	// The format of the returned data is implementation-specific (e.g., Geth's
	// trace format), but should be parseable into the tracer.Trace model.
	TraceTransaction(ctx context.Context, txHash common.Hash) (interface{}, error)

	// Close closes the client connection and releases resources.
	Close()
}

