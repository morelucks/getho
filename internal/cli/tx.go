package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/luckify/getho/internal/client"
	"github.com/luckify/getho/internal/decoder"
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
			txHashStr := args[0]

			// Validate and parse transaction hash
			if !common.IsHexAddress(txHashStr) && len(txHashStr) != 66 {
				return fmt.Errorf("invalid transaction hash: %s (expected 0x-prefixed 64-character hex string)", txHashStr)
			}

			txHash := common.HexToHash(txHashStr)
			if txHash == (common.Hash{}) {
				return fmt.Errorf("invalid transaction hash: %s", txHashStr)
			}

			// Get RPC URL
			rpcURL := GetRPCURL()

			// Create client
			ctx := context.Background()
			ethClient, err := client.NewClient(ctx, rpcURL)
			if err != nil {
				return fmt.Errorf("failed to connect to Ethereum node at %s: %w", rpcURL, err)
			}
			defer ethClient.Close()

			// Fetch transaction
			tx, isPending, err := ethClient.GetTransaction(ctx, txHash)
			if err != nil {
				return fmt.Errorf("failed to fetch transaction: %w", err)
			}

			if tx == nil {
				return fmt.Errorf("transaction not found: %s", txHashStr)
			}

			// Get sender address
			sender, err := decoder.GetSender(tx)
			if err != nil {
				return fmt.Errorf("failed to extract sender address: %w", err)
			}

			// Fetch receipt if transaction is not pending
			var receipt *types.Receipt
			if !isPending {
				receipt, err = ethClient.GetTransactionReceipt(ctx, txHash)
				if err != nil {
					// Receipt might not be available yet, continue without it
					fmt.Fprintf(os.Stderr, "Warning: could not fetch receipt: %v\n", err)
				}
			}

			// Decode transaction
			dec := decoder.NewEthereumDecoder()
			decodedTx, err := dec.FromGoEthereumTransaction(tx, receipt, sender)
			if err != nil {
				return fmt.Errorf("failed to decode transaction: %w", err)
			}

			// Store decoded transaction for output formatting
			// We'll display it in the next commit
			_ = decodedTx

			if isPending {
				cmd.Printf("Transaction is pending\n")
			}

			return nil
		},
	}

	return cmd
}
