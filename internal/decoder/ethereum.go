package decoder

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EthereumDecoder implements the Decoder interface for go-ethereum types.
type EthereumDecoder struct{}

// NewEthereumDecoder creates a new decoder for go-ethereum transaction types.
func NewEthereumDecoder() *EthereumDecoder {
	return &EthereumDecoder{}
}

// FromGoEthereumTransaction converts a go-ethereum types.Transaction into
// our internal Transaction model.
func (d *EthereumDecoder) FromGoEthereumTransaction(tx *types.Transaction, receipt *types.Receipt, from common.Address) (*Transaction, error) {
	if tx == nil {
		return nil, errors.New("transaction cannot be nil")
	}

	// Get transaction hash
	txHash := tx.Hash()

	// Determine transaction type
	txType := d.determineTransactionType(tx)

	// Get chain ID
	chainID := tx.ChainId()
	if chainID == nil {
		chainID = big.NewInt(0)
	}

	// Build access list
	accessList := d.buildAccessList(tx)

	// Get value
	value := tx.Value()
	if value == nil {
		value = big.NewInt(0)
	}

	// Build transaction
	result := &Transaction{
		Hash:                  txHash.Hex(),
		From:                  from.Hex(),
		To:                    d.getToAddress(tx),
		Nonce:                 tx.Nonce(),
		Value:                 value,
		GasLimit:              tx.Gas(),
		Type:                  txType,
		ChainID:               chainID,
		AccessList:            accessList,
		Input:                 tx.Data(),
		EstimatedIntrinsicGas: d.estimateIntrinsicGas(tx),
	}

	// Set gas price fields based on transaction type
	d.setGasPriceFields(result, tx, receipt)

	// Set blob gas fields for EIP-4844 transactions
	if tx.Type() == types.BlobTxType {
		d.setBlobGasFields(result, tx, receipt)
	}

	return result, nil
}

// determineTransactionType maps go-ethereum transaction types to our internal types.
func (d *EthereumDecoder) determineTransactionType(tx *types.Transaction) TransactionType {
	switch tx.Type() {
	case types.LegacyTxType:
		return TransactionTypeLegacy
	case types.AccessListTxType:
		return TransactionTypeAccessList
	case types.DynamicFeeTxType:
		return TransactionTypeDynamicFee
	case types.BlobTxType:
		return TransactionTypeBlob
	default:
		return TransactionTypeLegacy
	}
}

// buildAccessList converts go-ethereum access list to our format.
func (d *EthereumDecoder) buildAccessList(tx *types.Transaction) []AccessListEntry {
	accessList := tx.AccessList()
	if len(accessList) == 0 {
		return nil
	}

	result := make([]AccessListEntry, 0, len(accessList))
	for _, entry := range accessList {
		storageKeys := make([]string, 0, len(entry.StorageKeys))
		for _, key := range entry.StorageKeys {
			storageKeys = append(storageKeys, key.Hex())
		}
		result = append(result, AccessListEntry{
			Address:     entry.Address.Hex(),
			StorageKeys: storageKeys,
		})
	}
	return result
}

// getToAddress returns the recipient address or empty string for contract creation.
func (d *EthereumDecoder) getToAddress(tx *types.Transaction) string {
	to := tx.To()
	if to == nil {
		return "" // Contract creation
	}
	return to.Hex()
}

// setGasPriceFields sets gas price fields based on transaction type.
func (d *EthereumDecoder) setGasPriceFields(result *Transaction, tx *types.Transaction, receipt *types.Receipt) {
	switch tx.Type() {
	case types.LegacyTxType:
		result.GasPrice = tx.GasPrice()
		if receipt != nil {
			result.EffectiveGasPrice = receipt.EffectiveGasPrice
		} else {
			result.EffectiveGasPrice = tx.GasPrice()
		}
	case types.AccessListTxType, types.DynamicFeeTxType, types.BlobTxType:
		result.MaxFeePerGas = tx.GasFeeCap()
		result.MaxPriorityFeePerGas = tx.GasTipCap()
		if receipt != nil {
			result.EffectiveGasPrice = receipt.EffectiveGasPrice
		}
	}
}

// setBlobGasFields sets blob gas fields for EIP-4844 transactions.
func (d *EthereumDecoder) setBlobGasFields(result *Transaction, tx *types.Transaction, receipt *types.Receipt) {
	// For blob transactions, we can get blob fee cap from the transaction
	// Note: go-ethereum v1.14.0 may have different API for blob transactions
	// For now, we'll set these from the receipt if available
	if receipt != nil {
		result.BlobGasUsed = receipt.BlobGasUsed
		// Blob fee cap would need to be extracted from the transaction wrapper
		// This is a simplified version - full implementation would parse blob tx wrapper
		if receipt.BlobGasUsed > 0 {
			// Placeholder: blob fee cap would come from transaction data
			// For now, we leave it nil if not available
		}
	}
}

// estimateIntrinsicGas estimates the intrinsic gas cost of a transaction.
// This is a simplified version - full calculation would consider data, access list, etc.
func (d *EthereumDecoder) estimateIntrinsicGas(tx *types.Transaction) uint64 {
	// Base intrinsic gas
	gas := uint64(21000)

	// Add gas for data (4 gas per zero byte, 16 gas per non-zero byte)
	data := tx.Data()
	for _, b := range data {
		if b == 0 {
			gas += 4
		} else {
			gas += 16
		}
	}

	// Add gas for access list
	if tx.Type() == types.AccessListTxType || tx.Type() == types.DynamicFeeTxType {
		accessList := tx.AccessList()
		gas += uint64(len(accessList)) * 2400 // Per address
		for _, entry := range accessList {
			gas += uint64(len(entry.StorageKeys)) * 1900 // Per storage key
		}
	}

	// Contract creation costs extra
	if tx.To() == nil {
		gas += 32000
	}

	return gas
}

// GetSender extracts the sender address from a transaction.
// This uses EIP-155 replay protection to recover the sender.
func GetSender(tx *types.Transaction) (common.Address, error) {
	signer := types.LatestSignerForChainID(tx.ChainId())
	return types.Sender(signer, tx)
}
