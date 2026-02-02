package cmd

import (
	"fmt"
	"strconv"

	"github.com/nathfavour/shadowprism/cli/internal/sidecar"
	"github.com/spf13/cobra"
)

var (
	shieldStrategy string
	shieldForce    bool
)

var shieldCmd = &cobra.Command{
	Use:   "shield [amount] [destination]",
	Short: "Anonymize SOL by routing through a privacy provider",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		amount, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			fmt.Println("❌ Invalid amount. Must be lamports (uint64)")
			return
		}
		dest := args[1]

		cm, _ := sidecar.NewConfigManager()
		socketPath := cm.GetSocketPath()
		client := sidecar.NewCoreClient(socketPath, "dev-token-123")

		fmt.Printf("🕵️  Initiating Privacy Shield for %d lamports...\n", amount)
		
		// Map strategy to correct provider name if needed
		// Here we just pass it through
		
		res, err := client.Shield(amount, dest, shieldStrategy, shieldForce)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}

		fmt.Printf("✅ Shield Success!\n")
		fmt.Printf("🔗 TX: %s\n", res["tx_hash"])
		fmt.Printf("🛡️  Provider: %v\n", res["provider"])
		if note, ok := res["note"].(string); ok && note != "" {
			fmt.Printf("🔑 Note: %s (Stored in local DB)\n", note)
		}
	},
}

func init() {
	shieldCmd.Flags().StringVarP(&shieldStrategy, "strategy", "s", "privacy_cash", "Privacy strategy to use (privacy_cash, radr_p2p)")
	shieldCmd.Flags().BoolVarP(&shieldForce, "force", "f", false, "Force transaction even if destination is high-risk")
	rootCmd.AddCommand(shieldCmd)
}
