package cmd

import (
	"fmt"
	"strconv"

	"github.com/nathfavour/shadowprism/cli/internal/sidecar"
	"github.com/spf13/cobra"
)

var (
	fromToken string
	toToken   string
)

var swapCmd = &cobra.Command{
	Use:   "swap [amount]",
	Short: "Private token exchange via SilentSwap",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		amount, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			fmt.Println("❌ Invalid amount. Must be lamports (uint64)")
			return
		}

		cm, _ := sidecar.NewConfigManager()
		socketPath := cm.GetSocketPath()
		authToken := "dev-token-123" // In production, this is managed by the sidecar manager

		client := sidecar.NewCoreClient(socketPath, authToken)

		fmt.Printf("🔄 Initiating Private Swap: %d %s -> %s...\n", amount, fromToken, toToken)
		
		res, err := client.Swap(amount, fromToken, toToken)
		if err != nil {
			fmt.Printf("❌ Swap failed: %v\n", err)
			return
		}

		fmt.Printf("✅ Swap Confirmed!\n")
		fmt.Printf("🔗 TX: %s\n", res["tx_hash"])
		fmt.Printf("💰 Received: %v %s\n", res["to_amount"], toToken)
	},
}

func init() {
	swapCmd.Flags().StringVar(&fromToken, "from", "SOL", "Token to swap from")
	swapCmd.Flags().StringVar(&toToken, "to", "USDC", "Token to swap to")
	rootCmd.AddCommand(swapCmd)
}
