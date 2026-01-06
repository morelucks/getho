package client

// Client provides an interface for interacting with Ethereum execution clients
type Client interface {
	// GetTransaction retrieves a transaction by hash
	GetTransaction(txHash string) ([]byte, error)

	// GetTransactionReceipt retrieves a transaction receipt
	GetTransactionReceipt(txHash string) ([]byte, error)

	// TraceTransaction generates an execution trace
	TraceTransaction(txHash string) ([]byte, error)
}
 
