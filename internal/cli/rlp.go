package cli

import (
	"github.com/spf13/cobra"
)

func newRLPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rlp",
		Short: "RLP encoding/decoding utilities",
		Long:  "Decode raw RLP-encoded data",
	}

	decodeCmd := &cobra.Command{
		Use:   "decode [rlp_data]",
		Short: "Decode RLP-encoded data",
		Long:  "Decode raw RLP-encoded hexadecimal data",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rlpData := args[0]

			cmd.Printf("Decoding RLP data: %s\n", rlpData)
			return nil
		},
	}

	cmd.AddCommand(decodeCmd)
	return cmd
}
