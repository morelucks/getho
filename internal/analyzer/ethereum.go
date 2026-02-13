package analyzer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/luckify/getho/internal/client"
	"github.com/luckify/getho/internal/decoder"
)

// EthereumAnalyzer is a concrete Analyzer implementation backed by an
// execution client and the Ethereum decoder.
//
// It fetches the transaction, receipt, and block header and maps them into
// a normalized GasAnalysis focused on execution-layer behavior.
type EthereumAnalyzer struct {
	client client.Client
}

// NewEthereumAnalyzer creates a new Analyzer instance using the provided client.
func NewEthereumAnalyzer(c client.Client) *EthereumAnalyzer {
	return &EthereumAnalyzer{client: c}
}

// AnalyzeGas analyzes gas usage and fees for a single transaction hash.
func (a *EthereumAnalyzer) AnalyzeGas(txHash string) (*GasAnalysis, error) {
	if txHash == "" {
		return nil, fmt.Errorf("tx hash cannot be empty")
	}
	if len(txHash) != 66 || txHash[:2] != "0x" {
		return nil, fmt.Errorf("invalid transaction hash: %s", txHash)
	}

	ctx := context.Background()
	hash := common.HexToHash(txHash)

	// Fetch transaction
	tx, _, err := a.client.GetTransaction(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}
	if tx == nil {
		return nil, fmt.Errorf("transaction not found: %s", txHash)
	}

	// Fetch receipt
	receipt, err := a.client.GetTransactionReceipt(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch receipt: %w", err)
	}
	if receipt == nil {
		return nil, fmt.Errorf("receipt not found for transaction: %s", txHash)
	}

	// Fetch block header for base fee context
	header, err := a.client.GetBlockHeader(ctx, receipt.BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block header: %w", err)
	}

	// Decode transaction into normalized model so we can re-use fields like
	// EffectiveGasPrice, blob gas, etc.
	sender, err := decoder.GetSender(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract sender: %w", err)
	}

	dec := decoder.NewEthereumDecoder()
	decodedTx, err := dec.FromGoEthereumTransaction(tx, receipt, sender)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %w", err)
	}

	gasUsed := receipt.GasUsed
	gasLimit := decodedTx.GasLimit

	baseFeePerGas := header.BaseFee
	effectiveGasPrice := decodedTx.EffectiveGasPrice

	analysis := &GasAnalysis{
		TxHash:      decodedTx.Hash,
		BlockHash:   receipt.BlockHash.Hex(),
		BlockNumber: receipt.BlockNumber.Uint64(),

		GasUsed:  gasUsed,
		GasLimit: gasLimit,

		BaseFeePerGas:        baseFeePerGas,
		GasPrice:             decodedTx.GasPrice,
		MaxFeePerGas:         decodedTx.MaxFeePerGas,
		MaxPriorityFeePerGas: decodedTx.MaxPriorityFeePerGas,
		EffectiveGasPrice:    effectiveGasPrice,

		BlobGasUsed:   decodedTx.BlobGasUsed,
		BlobGasFeeCap: decodedTx.BlobGasFeeCap,

		Notes: make([]string, 0),
	}

	// Compute aggregated fee amounts.
	analysis.TotalFeePaid = mulUint64Big(gasUsed, effectiveGasPrice)

	if baseFeePerGas != nil && baseFeePerGas.Sign() > 0 {
		analysis.BaseFeeBurnt = mulUint64Big(gasUsed, baseFeePerGas)
	} else {
		analysis.Notes = append(analysis.Notes, "missing base fee (pre-EIP-1559 block or client did not return BaseFee)")
	}

	if analysis.TotalFeePaid != nil && analysis.BaseFeeBurnt != nil {
		priority := new(big.Int).Sub(analysis.TotalFeePaid, analysis.BaseFeeBurnt)
		if priority.Sign() < 0 {
			priority = big.NewInt(0)
		}
		analysis.PriorityFee = priority
	}

	// Blob-related fees (simplified – depends on client support).
	if decodedTx.BlobGasUsed > 0 && analysis.BlobGasPrice != nil {
		analysis.TotalBlobFeePaid = mulUint64Big(decodedTx.BlobGasUsed, analysis.BlobGasPrice)
	}
	if analysis.TotalFeePaid != nil && analysis.TotalBlobFeePaid != nil {
		analysis.TotalExecutionAndBlob = new(big.Int).Add(analysis.TotalFeePaid, analysis.TotalBlobFeePaid)
	}

	// Build convenience components for presentation.
	a.buildComponents(analysis)

	return analysis, nil
}

func (a *EthereumAnalyzer) buildComponents(ga *GasAnalysis) {
	components := make([]GasComponent, 0, 3)

	if ga.BaseFeeBurnt != nil {
		components = append(components, GasComponent{
			Label: "base",
			Value: ga.BaseFeeBurnt,
		})
	}
	if ga.PriorityFee != nil {
		components = append(components, GasComponent{
			Label: "priority",
			Value: ga.PriorityFee,
		})
	}
	if ga.TotalBlobFeePaid != nil && ga.TotalBlobFeePaid.Sign() > 0 {
		components = append(components, GasComponent{
			Label: "blob",
			Value: ga.TotalBlobFeePaid,
		})
	}

	ga.Components = components
}

func mulUint64Big(n uint64, x *big.Int) *big.Int {
	if x == nil {
		return nil
	}
	return new(big.Int).Mul(new(big.Int).SetUint64(n), x)
}


