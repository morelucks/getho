package cli

import (
	"github.com/spf13/cobra"
)

func newCalldataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calldata [tx_hash]",
		Short: "Decode calldata from a transaction",
		Long: `Decode function selectors and parse arguments from transaction calldata.
Highlights unknown or malformed calldata.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txHash := args[0]

			cmd.Printf("Decoding calldata for transaction: %s\n", txHash)
			return nil
		},
	}

	return cmd
}
