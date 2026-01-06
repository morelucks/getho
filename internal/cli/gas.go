package cli

import (
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

			cmd.Printf("Analyzing gas for transaction: %s\n", txHash)
			return nil
		},
	}

	return cmd
}
