package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var reinstallCmd = &cobra.Command{
	Use:   "reinstall",
	Short: "Wipe all data and perform a fresh installation",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚠️  Preparing for Reinstallation...")
		
homeDir, _ := os.UserHomeDir()
		prismDir := fmt.Sprintf("%s/.shadowprism", homeDir)

		fmt.Printf("🗑️  Wiping application data at %s...\n", prismDir)
		if err := os.RemoveAll(prismDir); err != nil {
			fmt.Printf("❌ Failed to wipe data: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Data wiped. Starting fresh installation...")
		
		// Re-use update logic by calling the update command's Run function
		// or simply triggering the same shell logic
		updateCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(reinstallCmd)
}

