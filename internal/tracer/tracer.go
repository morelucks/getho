package tracer

import "math/big"

// Tracer provides execution tracing capabilities.
//
// Implementations are expected to map client-specific trace formats into this
// normalized, execution-layer focused model.
type Tracer interface {
	// Trace generates an opcode-level execution trace for a transaction hash.
	Trace(txHash string) (*Trace, error)
}

// CallType represents the high-level kind of EVM call frame.
type CallType string

const (
	CallTypeCall         CallType = "CALL"
	CallTypeDelegateCall          = "DELEGATECALL"
	CallTypeStaticCall            = "STATICCALL"
	CallTypeCreate                = "CREATE"
	CallTypeCreate2               = "CREATE2"
)

// OpcodeStats aggregates counts of interesting opcodes within a frame.
type OpcodeStats struct {
	Total    uint64
	Calls    uint64 // CALL, DELEGATECALL, STATICCALL
	SLoads   uint64 // SLOAD
	SStores  uint64 // SSTORE
	Logs     uint64 // LOG0-LOG4
	Reverts  uint64 // REVERT
	Invalids uint64 // INVALID, BADJUMP, etc.
}

// CallFrame models a single EVM call frame within a transaction trace.
//
// This is a summary-level view – implementations may keep richer internal
// structures if needed, but getho surfaces this normalized representation.
type CallFrame struct {
	Type CallType

	From string // caller address (0x-prefixed)
	To   string // callee/created contract address (0x-prefixed)

	Value *big.Int

	// Depth of this frame in the call tree, where 0 is the root transaction.
	Depth int

	// Gas accounting.
	GasLimit uint64 // gas allocated to the frame
	GasUsed  uint64 // gas actually consumed

	// Summary of interesting opcodes executed in this frame.
	Opcodes OpcodeStats

	// Optional error string if this frame ended in REVERT or an error.
	Error string
}

// Trace represents an execution trace for a single transaction.
type Trace struct {
	TxHash string // 32-byte transaction hash (0x-prefixed)

	// Total gas used as reported by the execution client.
	TotalGasUsed uint64

	// Flattened list of all frames in depth-first order.
	Frames []CallFrame

	// Optional root-level error, if any.
	Error string
}
