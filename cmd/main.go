package main

import (
	"os"

	"github.com/luckify/getho/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
