package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "shadowprism",
	Short: "ShadowPrism is a privacy-first liquidity aggregator for Solana",
	Long:  `A secure sidecar that routes Solana transactions through privacy protocols.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Root flags can be defined here
}