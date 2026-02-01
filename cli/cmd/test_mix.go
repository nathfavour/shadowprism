package cmd

import (
	"fmt"
	"os"

	"github.com/nathfavour/shadowprism/cli/internal/sidecar"
	"github.com/spf13/cobra"
)

var testMixCmd = &cobra.Command{
	Use:   "test-mix",
	Short: "Send a test shielding request to the core engine",
	Run: func(cmd *cobra.Command, args []string) {
		token := "dev-token-123"
		client := sidecar.NewCoreClient(42069, token)

		fmt.Println("🧪 Sending test shielding request...")
		
		payload := map[string]interface{}{
			"amount_lamports":  1000000000,
			"destination_addr": "BuX...7z",
			"strategy":         "mix_standard",
		}

		var result map[string]interface{}
		resp, err := client.Http.R().
			SetBody(payload).
			SetResult(&result).
			Post("/v1/shield")

		if err != nil {
			fmt.Printf("❌ Request failed: %v\n", err)
			os.Exit(1)
		}

		if resp.IsError() {
			fmt.Printf("❌ Core returned error (%d): %s\n", resp.StatusCode(), resp.String())
			os.Exit(1)
		}

		fmt.Println("✅ Shielding Success!")
		fmt.Printf("🔗 Transaction Hash: %v\n", result["tx_hash"])
		fmt.Printf("🛡️ Provider Used: %v\n", result["provider"])
	},
}

func init() {
	rootCmd.AddCommand(testMixCmd)
}
