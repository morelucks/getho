package decoder

// Decoder provides transaction and calldata decoding capabilities
type Decoder interface {
	// DecodeTransaction decodes RLP-encoded transaction data
	DecodeTransaction(data []byte) (*Transaction, error)

	// DecodeCalldata decodes function selector and arguments
	DecodeCalldata(calldata []byte) (*Calldata, error)
}

// Transaction represents a decoded Ethereum transaction
type Transaction struct {
}

// Calldata represents decoded calldata
type Calldata struct {
}

 
