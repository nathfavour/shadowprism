package cmd

import (
	"fmt"
	"strconv"

	"github.com/nathfavour/shadowprism/cli/internal/sidecar"
	"github.com/spf13/cobra"
)

var payCmd = &cobra.Command{
	Use:   "pay [merchant_id] [amount]",
	Short: "Private payment to a Starpay merchant",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		merchant := args[0]
		amount, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			fmt.Println("❌ Invalid amount.")
			return
		}

		cm, _ := sidecar.NewConfigManager()
		socketPath := cm.GetSocketPath()
		client := sidecar.NewCoreClient(socketPath, "dev-token-123")

		fmt.Printf("💳 Sending private payment of %d lamports to %s...\n", amount, merchant)

		res, err := client.Pay(amount, merchant)
		if err != nil {
			fmt.Printf("❌ Payment failed: %v\n", err)
			return
		}

		fmt.Printf("✅ Payment Successful!\n")
		fmt.Printf("🔗 TX: %s\n", res["tx_hash"])
		fmt.Printf("🧾 Receipt ID: %s\n", res["receipt_id"])
	},
}

func init() {
	rootCmd.AddCommand(payCmd)
}
