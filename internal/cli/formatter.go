package cli

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/luckify/getho/internal/decoder"
)

// FormatTransaction displays a decoded transaction in a human-readable format.
func FormatTransaction(tx *decoder.Transaction, receipt *types.Receipt, isPending bool) string {
	var b strings.Builder

	// Header
	b.WriteString("Transaction Details\n")
	b.WriteString(strings.Repeat("=", 80) + "\n\n")

	// Basic Information
	b.WriteString("Hash:        " + tx.Hash + "\n")
	b.WriteString("Status:      ")
	if isPending {
		b.WriteString("PENDING\n")
	} else if receipt != nil {
		if receipt.Status == types.ReceiptStatusSuccessful {
			b.WriteString("SUCCESS\n")
		} else {
			b.WriteString("FAILED (reverted)\n")
		}
	} else {
		b.WriteString("UNKNOWN\n")
	}
	b.WriteString("\n")

	// From/To
	b.WriteString("From:        " + tx.From + "\n")
	if tx.To != "" {
		b.WriteString("To:          " + tx.To + "\n")
	} else {
		b.WriteString("To:          [Contract Creation]\n")
	}
	b.WriteString("\n")

	// Value
	b.WriteString("Value:       " + formatWei(tx.Value) + " ETH\n")
	b.WriteString("Nonce:       " + fmt.Sprintf("%d", tx.Nonce) + "\n")
	b.WriteString("\n")

	// Transaction Type
	b.WriteString("Type:        " + formatTransactionType(tx.Type) + "\n")
	if tx.ChainID != nil && tx.ChainID.Sign() > 0 {
		b.WriteString("Chain ID:    " + tx.ChainID.String() + "\n")
	}
	b.WriteString("\n")

	// Gas Information
	b.WriteString("Gas Information\n")
	b.WriteString(strings.Repeat("-", 80) + "\n")
	b.WriteString("Gas Limit:   " + formatUint64(tx.GasLimit) + "\n")
	if receipt != nil {
		b.WriteString("Gas Used:    " + formatUint64(receipt.GasUsed) + " (" + formatPercentage(receipt.GasUsed, tx.GasLimit) + ")\n")
	}
	b.WriteString("Intrinsic:   " + formatUint64(tx.EstimatedIntrinsicGas) + "\n")
	b.WriteString("\n")

	// Fee Information
	b.WriteString("Fee Information\n")
	b.WriteString(strings.Repeat("-", 80) + "\n")
	if tx.GasPrice != nil {
		b.WriteString("Gas Price:   " + formatWei(tx.GasPrice) + " gwei\n")
	}
	if tx.MaxFeePerGas != nil {
		b.WriteString("Max Fee:     " + formatWei(tx.MaxFeePerGas) + " gwei\n")
	}
	if tx.MaxPriorityFeePerGas != nil {
		b.WriteString("Max Priority: " + formatWei(tx.MaxPriorityFeePerGas) + " gwei\n")
	}
	if tx.EffectiveGasPrice != nil {
		b.WriteString("Effective:   " + formatWei(tx.EffectiveGasPrice) + " gwei\n")
		if receipt != nil {
			totalFee := new(big.Int).Mul(tx.EffectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
			b.WriteString("Total Fee:   " + formatWei(totalFee) + " ETH\n")
		}
	}
	b.WriteString("\n")

	// Access List (if present)
	if len(tx.AccessList) > 0 {
		b.WriteString("Access List\n")
		b.WriteString(strings.Repeat("-", 80) + "\n")
		for i, entry := range tx.AccessList {
			b.WriteString(fmt.Sprintf("  [%d] Address: %s\n", i+1, entry.Address))
			if len(entry.StorageKeys) > 0 {
				for j, key := range entry.StorageKeys {
					if j < 3 {
						b.WriteString(fmt.Sprintf("         Key[%d]: %s\n", j, key))
					} else if j == 3 {
						b.WriteString(fmt.Sprintf("         ... and %d more keys\n", len(entry.StorageKeys)-3))
					}
				}
			}
		}
		b.WriteString("\n")
	}

	// Blob Gas (EIP-4844)
	if tx.Type == decoder.TransactionTypeBlob {
		b.WriteString("Blob Gas (EIP-4844)\n")
		b.WriteString(strings.Repeat("-", 80) + "\n")
		if tx.BlobGasUsed > 0 {
			b.WriteString("Blob Gas Used: " + formatUint64(tx.BlobGasUsed) + "\n")
		}
		if tx.MaxFeePerBlobGas != nil {
			b.WriteString("Max Fee/Blob:  " + formatWei(tx.MaxFeePerBlobGas) + " gwei\n")
		}
		b.WriteString("\n")
	}

	// Input Data
	b.WriteString("Input Data\n")
	b.WriteString(strings.Repeat("-", 80) + "\n")
	if len(tx.Input) == 0 {
		b.WriteString("(no input data)\n")
	} else {
		b.WriteString(fmt.Sprintf("Length: %d bytes\n", len(tx.Input)))
		if len(tx.Input) <= 10 {
			b.WriteString("Data:   0x" + fmt.Sprintf("%x", tx.Input) + "\n")
		} else {
			b.WriteString("Data:   0x" + fmt.Sprintf("%x", tx.Input[:10]) + "...\n")
			b.WriteString(fmt.Sprintf("         (truncated, full: %d bytes)\n", len(tx.Input)))
		}
	}
	b.WriteString("\n")

	// Block Information (if receipt available)
	if receipt != nil {
		b.WriteString("Block Information\n")
		b.WriteString(strings.Repeat("-", 80) + "\n")
		b.WriteString("Block Number: " + formatUint64(receipt.BlockNumber.Uint64()) + "\n")
		b.WriteString("Block Hash:   " + receipt.BlockHash.Hex() + "\n")
		b.WriteString("Tx Index:     " + formatUint64(uint64(receipt.TransactionIndex)) + "\n")
		if receipt.ContractAddress != (common.Address{}) {
			b.WriteString("Contract:     " + receipt.ContractAddress.Hex() + "\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}

// formatWei converts wei to a human-readable string (ETH or gwei).
func formatWei(wei *big.Int) string {
	if wei == nil {
		return "0"
	}

	// For gwei display, divide by 1e9
	gwei := new(big.Int).Div(wei, big.NewInt(1e9))
	remainder := new(big.Int).Mod(wei, big.NewInt(1e9))

	if remainder.Sign() == 0 {
		return gwei.String()
	}

	// For ETH display, divide by 1e18
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
	return eth.Text('f', 18)
}

// formatUint64 formats a uint64 with thousand separators.
func formatUint64(n uint64) string {
	return fmt.Sprintf("%d", n)
}

// formatPercentage calculates and formats a percentage.
func formatPercentage(part, total uint64) string {
	if total == 0 {
		return "0%"
	}
	percent := float64(part) / float64(total) * 100
	return fmt.Sprintf("%.2f%%", percent)
}

// formatTransactionType returns a human-readable transaction type string.
func formatTransactionType(t decoder.TransactionType) string {
	switch t {
	case decoder.TransactionTypeLegacy:
		return "Legacy (0x0)"
	case decoder.TransactionTypeAccessList:
		return "Access List (0x1)"
	case decoder.TransactionTypeDynamicFee:
		return "EIP-1559 (0x2)"
	case decoder.TransactionTypeBlob:
		return "Blob (EIP-4844) (0x3)"
	default:
		return fmt.Sprintf("Unknown (%d)", t)
	}
}
