package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// rpcURL is the Ethereum JSON-RPC endpoint URL
	rpcURL string
)

var rootCmd = &cobra.Command{
	Use:   "getho",
	Short: "A low-level Ethereum base-layer observability and debugging tool",
	Long: `getho is a low-level Ethereum base-layer observability and debugging tool
for inspecting execution-layer behavior.

It provides transaction inspection, gas analysis, execution tracing,
and calldata decoding capabilities.`,
	Version: "0.1.0",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rpcURL, "rpc", "", "Ethereum JSON-RPC endpoint URL (default: $GETHO_RPC_URL or http://localhost:8545)")
}

// GetRPCURL returns the configured RPC URL, falling back to environment variable or default.
func GetRPCURL() string {
	if rpcURL != "" {
		return rpcURL
	}
	return getRPCURLFromEnv()
}

func getRPCURLFromEnv() string {
	if url := os.Getenv("GETHO_RPC_URL"); url != "" {
		return url
	}
	return "http://localhost:8545"
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newTxCmd())
	rootCmd.AddCommand(newCalldataCmd())
	rootCmd.AddCommand(newGasCmd())
	rootCmd.AddCommand(newTraceCmd())
	rootCmd.AddCommand(newRLPCmd())
}

func er(msg interface{}) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", msg)
}
