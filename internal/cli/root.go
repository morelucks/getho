package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
