package cli

import (
	"github.com/spf13/cobra"
)

func newTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx [tx_hash]",
		Short: "Inspect a transaction",
		Long: `Inspect a transaction by its hash. Displays sender, recipient,
value, nonce, type, and decodes RLP-encoded transaction data.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txHash := args[0]

			cmd.Printf("Inspecting transaction: %s\n", txHash)
			return nil
		},
	}

	return cmd
}
