package analyzer

import "math/big"

// Analyzer provides gas and fee analysis capabilities.
//
// The Analyzer consumes decoded transaction and receipt data and produces a
// normalized GasAnalysis that focuses on execution-layer fee behavior.
type Analyzer interface {
	// AnalyzeGas analyzes gas usage and fees for a single transaction hash.
	AnalyzeGas(txHash string) (*GasAnalysis, error)
}

// GasComponent represents a labeled fee component (e.g. base, priority, blob).
type GasComponent struct {
	Label string   // human-readable label, e.g. "base", "priority", "blob"
	Value *big.Int // denominated in wei
}

// GasAnalysis represents a detailed gas and fee breakdown for a transaction.
//
// All monetary values are denominated in wei and represented as *big.Int to
// avoid precision loss. This model is intentionally execution-layer centric
// and does not attempt to interpret L2 semantics.
type GasAnalysis struct {
	TxHash      string // 32-byte transaction hash (0x-prefixed)
	BlockHash   string // containing block hash (0x-prefixed)
	BlockNumber uint64

	// Gas usage.
	GasUsed  uint64
	GasLimit uint64

	// Base fee context (per gas).
	BaseFeePerGas *big.Int

	// Transaction fee configuration (per gas).
	GasPrice             *big.Int // legacy gas price
	MaxFeePerGas         *big.Int // EIP-1559
	MaxPriorityFeePerGas *big.Int // EIP-1559

	// Effective fee actually paid (per gas).
	EffectiveGasPrice *big.Int

	// Aggregated fee amounts (in wei).
	TotalFeePaid *big.Int // EffectiveGasPrice * GasUsed
	BaseFeeBurnt *big.Int // portion burnt (base fee * GasUsed)
	PriorityFee  *big.Int // miner / proposer tip

	// Blob-related gas for EIP-4844 transactions.
	BlobGasUsed           uint64
	BlobGasPrice          *big.Int
	BlobGasFeeCap         *big.Int
	TotalBlobFeePaid      *big.Int
	TotalExecutionAndBlob *big.Int // TotalFeePaid + TotalBlobFeePaid

	// Convenience breakdown for presentation.
	Components []GasComponent

	// Any notes or warnings encountered during analysis
	// (e.g. missing base fee, pre-EIP-1559 block).
	Notes []string
}
