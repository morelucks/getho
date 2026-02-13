package cli

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/luckify/getho/internal/analyzer"
	"github.com/luckify/getho/internal/client"
	"github.com/spf13/cobra"
)

func newGasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gas [tx_hash]",
		Short: "Analyze gas and fee breakdown",
		Long: `Analyze gas usage and fee breakdown for a transaction.
Displays base fee, priority fee (tip), blob fee (EIP-4844),
and gas used vs gas limit.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txHash := args[0]

			// Validate and parse transaction hash (reuse logic style from tx.go).
			if len(txHash) < 2 || txHash[:2] != "0x" {
				return fmt.Errorf("invalid transaction hash: %s (must start with 0x)", txHash)
			}
			if len(txHash[2:]) != 64 {
				return fmt.Errorf("invalid transaction hash: %s (expected 64 hex characters after 0x, got %d)", txHash, len(txHash[2:]))
			}

			hash := common.HexToHash(txHash)
			if hash == (common.Hash{}) {
				return fmt.Errorf("invalid transaction hash: %s", txHash)
			}

			rpcURL := GetRPCURL()

			ctx := context.Background()
			c, err := client.NewClient(ctx, rpcURL)
			if err != nil {
				return fmt.Errorf("failed to connect to Ethereum node at %s: %w", rpcURL, err)
			}
			defer c.Close()

			a := analyzer.NewEthereumAnalyzer(c)

			ga, err := a.AnalyzeGas(txHash)
			if err != nil {
				return err
			}

			output := FormatGasAnalysis(ga)
			cmd.Print(output)

			return nil
		},
	}

	return cmd
}
