package client

import (
	"context"
	"os"
)

const (
	// DefaultRPCURL is the default Ethereum JSON-RPC endpoint.
	// This can be overridden via GETHO_RPC_URL environment variable or --rpc flag.
	DefaultRPCURL = "http://localhost:8545"
)

// GetRPCURL returns the RPC URL from environment variable or default.
func GetRPCURL() string {
	if url := os.Getenv("GETHO_RPC_URL"); url != "" {
		return url
	}
	return DefaultRPCURL
}

// NewClient creates a new client instance using the configured RPC URL.
func NewClient(ctx context.Context, rpcURL string) (Client, error) {
	if rpcURL == "" {
		rpcURL = GetRPCURL()
	}
	return NewRPCClient(ctx, rpcURL)
}

