package cli

import (
	"github.com/spf13/cobra"
)

func newTraceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trace [tx_hash]",
		Short: "Generate full execution trace",
		Long: `Generate an opcode-level execution trace for a transaction.
Tracks CALL, DELEGATECALL, STATICCALL, SSTORE, SLOAD operations
and identifies execution paths and state changes.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txHash := args[0]

			cmd.Printf("Tracing execution for transaction: %s\n", txHash)
			return nil
		},
	}

	return cmd
}
