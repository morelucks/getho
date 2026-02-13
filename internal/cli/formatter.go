package cli

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/luckify/getho/internal/analyzer"
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
	b.WriteString("Value:       " + formatEth(tx.Value) + " ETH\n")
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
		b.WriteString("Gas Price:   " + formatGwei(tx.GasPrice) + " gwei\n")
	}
	if tx.MaxFeePerGas != nil {
		b.WriteString("Max Fee:     " + formatGwei(tx.MaxFeePerGas) + " gwei\n")
	}
	if tx.MaxPriorityFeePerGas != nil {
		b.WriteString("Max Priority: " + formatGwei(tx.MaxPriorityFeePerGas) + " gwei\n")
	}
	if tx.EffectiveGasPrice != nil {
		b.WriteString("Effective:   " + formatGwei(tx.EffectiveGasPrice) + " gwei\n")
		if receipt != nil {
			totalFee := new(big.Int).Mul(tx.EffectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
			b.WriteString("Total Fee:   " + formatEth(totalFee) + " ETH\n")
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
			b.WriteString("Max Fee/Blob:  " + formatGwei(tx.MaxFeePerBlobGas) + " gwei\n")
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

// FormatGasAnalysis displays a gas and fee breakdown in a human-readable format.
func FormatGasAnalysis(ga *analyzer.GasAnalysis) string {
	if ga == nil {
		return "no gas analysis available\n"
	}

	var b strings.Builder

	b.WriteString("Gas & Fee Analysis\n")
	b.WriteString(strings.Repeat("=", 80) + "\n\n")

	// Basic context
	b.WriteString("Tx Hash:      " + ga.TxHash + "\n")
	b.WriteString("Block:        " + formatUint64(ga.BlockNumber) + " (" + ga.BlockHash + ")\n")
	b.WriteString("Gas Used:     " + formatUint64(ga.GasUsed) + " / " + formatUint64(ga.GasLimit) + " (" + formatPercentage(ga.GasUsed, ga.GasLimit) + ")\n")
	b.WriteString("\n")

	// Per-gas pricing
	b.WriteString("Per-Gas Prices (gwei)\n")
	b.WriteString(strings.Repeat("-", 80) + "\n")
	if ga.BaseFeePerGas != nil {
		b.WriteString("Base fee:     " + formatGwei(ga.BaseFeePerGas) + "\n")
	}
	if ga.GasPrice != nil {
		b.WriteString("Gas price:    " + formatGwei(ga.GasPrice) + "\n")
	}
	if ga.MaxFeePerGas != nil {
		b.WriteString("Max fee:      " + formatGwei(ga.MaxFeePerGas) + "\n")
	}
	if ga.MaxPriorityFeePerGas != nil {
		b.WriteString("Max priority: " + formatGwei(ga.MaxPriorityFeePerGas) + "\n")
	}
	if ga.EffectiveGasPrice != nil {
		b.WriteString("Effective:    " + formatGwei(ga.EffectiveGasPrice) + "\n")
	}
	b.WriteString("\n")

	// Aggregated fees
	b.WriteString("Aggregated Fees (ETH)\n")
	b.WriteString(strings.Repeat("-", 80) + "\n")
	if ga.TotalFeePaid != nil {
		b.WriteString("Total execution fee: " + formatEth(ga.TotalFeePaid) + "\n")
	}
	if ga.BaseFeeBurnt != nil {
		b.WriteString("  Base fee burnt:    " + formatEth(ga.BaseFeeBurnt) + "\n")
	}
	if ga.PriorityFee != nil {
		b.WriteString("  Priority (tip):     " + formatEth(ga.PriorityFee) + "\n")
	}
	if ga.TotalBlobFeePaid != nil {
		b.WriteString("Blob fee:            " + formatEth(ga.TotalBlobFeePaid) + "\n")
	}
	if ga.TotalExecutionAndBlob != nil {
		b.WriteString("Total (exec+blob):   " + formatEth(ga.TotalExecutionAndBlob) + "\n")
	}
	b.WriteString("\n")

	// Components
	if len(ga.Components) > 0 {
		b.WriteString("Component Breakdown\n")
		b.WriteString(strings.Repeat("-", 80) + "\n")
		for _, c := range ga.Components {
			b.WriteString(fmt.Sprintf("  %-8s %s\n", c.Label+":", formatEth(c.Value)))
		}
		b.WriteString("\n")
	}

	// Blob gas
	if ga.BlobGasUsed > 0 {
		b.WriteString("Blob Gas (EIP-4844)\n")
		b.WriteString(strings.Repeat("-", 80) + "\n")
		b.WriteString("Blob Gas Used: " + formatUint64(ga.BlobGasUsed) + "\n")
		if ga.BlobGasPrice != nil {
			b.WriteString("Blob Gas Price: " + formatGwei(ga.BlobGasPrice) + " gwei\n")
		}
		if ga.BlobGasFeeCap != nil {
			b.WriteString("Blob Gas Cap:   " + formatGwei(ga.BlobGasFeeCap) + " gwei\n")
		}
		b.WriteString("\n")
	}

	// Notes
	if len(ga.Notes) > 0 {
		b.WriteString("Notes\n")
		b.WriteString(strings.Repeat("-", 80) + "\n")
		for _, n := range ga.Notes {
			b.WriteString(" - " + n + "\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}

// formatEth converts wei to a string representation in ETH (18 decimals).
func formatEth(wei *big.Int) string {
	if wei == nil {
		return "0"
	}
	// Convert to float using big.Float directly for better precision handling in display
	f := new(big.Float).SetInt(wei)
	f.Quo(f, big.NewFloat(1e18))
	// Format with up to 9 decimal places, trimming trailing zeros
	// This prevents 0.000000000000000000 but keeps precision for small amounts
	return strings.TrimRight(strings.TrimRight(f.Text('f', 9), "0"), ".")
}

// formatGwei converts wei to a string representation in gwei (9 decimals).
func formatGwei(wei *big.Int) string {
	if wei == nil {
		return "0"
	}
	f := new(big.Float).SetInt(wei)
	f.Quo(f, big.NewFloat(1e9))
	return strings.TrimRight(strings.TrimRight(f.Text('f', 9), "0"), ".")
}

// formatUint64 formats a uint64 with thousand separators.
func formatUint64(n uint64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) < 4 {
		return s
	}
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteRune(',')
		}
		b.WriteRune(r)
	}
	return b.String()
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
