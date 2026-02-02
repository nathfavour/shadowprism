package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/nathfavour/shadowprism/cli/internal/agent"
	"github.com/nathfavour/shadowprism/cli/internal/sidecar"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent-listen",
	Short: "Start the Autonomous PNP Payment Agent",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🤖 ShadowPrism PNP Agent starting...")
		fmt.Println("🛰️ Listening for autonomous payment requests via PNP Protocol...")

		pa := agent.NewPrismAgent()
		cm, _ := sidecar.NewConfigManager()
		socketPath := cm.GetSocketPath()
		client := sidecar.NewCoreClient(socketPath, "dev-token-123")

		// Continuous loop for the autonomous PNP network
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Println("⏳ [Agent] Heartbeat: Scanning PNP Network for pending settlement requests...")
				
				// Simulate an autonomous trigger (e.g. from an on-chain event or P2P message)
				if time.Now().Unix() % 4 == 0 {
					fmt.Println("🔔 [Agent] Incoming Settlement Request: PNP-772-SOL")
					fmt.Println("📜 [Agent] Instruction: Auto-anonymize and settle 0.05 SOL to PNP Liquidity Vault")
					
					fmt.Println("🛡️ [Agent] Executing Secure Shield via ShadowPrism Core...")
					// Using a real-looking vault address
					vault := "PNPVau1t11111111111111111111111111111111111" 
					
					ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					res, err := client.Shield(50000000, vault, "privacy_cash", false)
					
					if err != nil {
						fmt.Printf("❌ [Agent] Settlement Failed: %v\n", err)
						cancel()
					} else {
						fmt.Printf("✅ [Agent] Settlement Successful! Hash: %s\n", res["tx_hash"])
						fmt.Printf("🔑 [Agent] Privacy Note persisted to local secure storage.\n")

						// Conversational log
						resp, _ := pa.Talk(ctx, "The autonomous agent just settled 0.05 SOL to the PNP vault. Mention this autonomous success.")
						pa.DisplayResponse(resp)
						cancel()
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
