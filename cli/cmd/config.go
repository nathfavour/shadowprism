package cmd

import (
	"fmt"
	"os"

	"github.com/nathfavour/shadowprism/cli/internal/sidecar"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure ShadowPrism settings",
}

var setBotTokenCmd = &cobra.Command{
	Use:   "set-bot-token [token]",
	Short: "Securely save the Telegram Bot token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cm, err := sidecar.NewConfigManager()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

		if err := cm.SaveSecret("tg_bot_token", args[0]); err != nil {
			fmt.Printf("❌ Failed to save token: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Telegram Bot token saved and encrypted in ~/.shadowprism")
	},
}

func init() {
	configCmd.AddCommand(setBotTokenCmd)
	rootCmd.AddCommand(configCmd)
}
