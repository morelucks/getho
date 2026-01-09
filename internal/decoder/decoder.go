package decoder

import "math/big"

// Decoder provides transaction and calldata decoding capabilities.
//
// These models are intentionally execution-layer focused and mirror the
// semantics exposed by common Ethereum JSON-RPC and client internals.
type Decoder interface {
	// DecodeTransaction decodes RLP-encoded transaction data into a
	// normalized Transaction model.
	DecodeTransaction(data []byte) (*Transaction, error)

	// DecodeCalldata decodes function selector and arguments from raw
	// calldata bytes into a structured Calldata model.
	DecodeCalldata(calldata []byte) (*Calldata, error)
}

// TransactionType represents the high-level Ethereum transaction type.
//
// This is intentionally minimal – concrete mapping from raw tx data
// (legacy, access-list, EIP-1559, blob, etc.) happens inside the decoder.
type TransactionType uint8

const (
	TransactionTypeLegacy TransactionType = iota
	TransactionTypeAccessList
	TransactionTypeDynamicFee
	TransactionTypeBlob // EIP-4844
)

// AccessListEntry is a normalized representation of an access list item.
type AccessListEntry struct {
	Address     string   // 20-byte hex address (0x-prefixed)
	StorageKeys []string // 32-byte storage keys (0x-prefixed)
}

// Transaction is a decoded, execution-layer oriented Ethereum transaction.
//
// All monetary values are represented as *big.Int to avoid precision loss.
// Hex-encoded hashes/addresses are always 0x-prefixed.
type Transaction struct {
	Hash  string // 32-byte transaction hash (0x-prefixed)
	From  string // sender address (0x-prefixed)
	To    string // recipient address (0x-prefixed), empty for contract creation
	Nonce uint64

	// Value transferred with the transaction.
	Value *big.Int

	// Gas fields.
	GasLimit              uint64
	GasPrice              *big.Int // legacy gas price
	MaxFeePerGas          *big.Int // EIP-1559
	MaxPriorityFeePerGas  *big.Int // EIP-1559
	EffectiveGasPrice     *big.Int // derived from receipt when available
	Type                  TransactionType
	ChainID               *big.Int
	AccessList            []AccessListEntry
	Input                 []byte // raw calldata
	BlobGasUsed           uint64
	BlobGasFeeCap         *big.Int // EIP-4844 blob fee cap
	MaxFeePerBlobGas      *big.Int // EIP-4844 max fee per blob gas
	EstimatedIntrinsicGas uint64   // intrinsic gas cost (decoded, not executed)
}

// Argument represents a single decoded calldata argument.
type Argument struct {
	Name  string      // optional best-effort name (from ABI/metadata if available)
	Type  string      // canonical Solidity type, e.g. "uint256", "address[]"
	Value interface{} // decoded Go value (string, *big.Int, []byte, slices, etc.)
}

// Calldata represents decoded calldata, including selector and arguments.
type Calldata struct {
	// Raw calldata bytes as provided to the EVM.
	Raw []byte

	// 4-byte function selector as 0x-prefixed hex.
	Selector string

	// Optional resolved function name (best-effort from ABI or selector DB).
	FunctionName string

	// Decoded arguments. When ABI information is unavailable, this may be empty
	// and Unknown will be set to true.
	Arguments []Argument

	// Unknown indicates that we could not confidently decode this calldata
	// (e.g. missing ABI or malformed input).
	Unknown bool
}
