package cmd

import (
	"fmt"
	"time"

	"github.com/nathfavour/shadowprism/cli/internal/sidecar"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent-listen",
	Short: "Start the Autonomous PNP Payment Agent",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🤖 ShadowPrism PNP Agent starting...")
		fmt.Println("🛰️ Listening for autonomous payment requests via PNP Protocol...")

		cm, _ := sidecar.NewConfigManager()
		socketPath := cm.GetSocketPath()
		client := sidecar.NewCoreClient(socketPath, "dev-token-123")

		// Simulation loop for the hackathon demo
		ticker := time.NewTicker(8 * time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Println("⏳ [Agent] Heartbeat: PNP network connected. Standing by...")
				
				// Simulate a triggered event
				if time.Now().Unix() % 3 == 0 {
					fmt.Println("🔔 [Agent] ALERT: Incoming payment request from AI-Agent-7")
					fmt.Println("📜 [Agent] Details: 50,000,000 Lamports for 'Private Inference Fee'")
					
					fmt.Println("🛡️ [Agent] Executing Auto-Shielded Transfer...")
					res, err := client.Shield(50000000, "PNP-Vault-11111111111111111111111111111111", "pnp_autonomous", false)
					if err != nil {
						fmt.Printf("❌ [Agent] Failed to fulfill request: %v\n", err)
					} else {
						fmt.Printf("✅ [Agent] Request Fulfilled! TX: %s\n", res["tx_hash"])
						fmt.Printf("🔑 [Agent] Privacy Note stored for Agent audit.\n")
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
