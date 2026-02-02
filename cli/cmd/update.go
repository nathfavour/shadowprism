package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update ShadowPrism to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔄 Updating ShadowPrism...")

		// Construct the installation command
		// We use the local install.sh if we are in the repo, otherwise fetch from remote
		installCmd := "curl -fsSL https://raw.githubusercontent.com/nathfavour/shadowprism/main/install.sh | bash"
		
		if _, err := os.Stat("./install.sh"); err == nil {
			fmt.Println("📂 Local install.sh detected, using it for update.")
			installCmd = "./install.sh"
		}

		shellCmd := exec.Command("bash", "-c", installCmd)
		shellCmd.Stdout = os.Stdout
		shellCmd.Stderr = os.Stderr

		if err := shellCmd.Run(); err != nil {
			fmt.Printf("❌ Update failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✨ ShadowPrism updated successfully!")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
