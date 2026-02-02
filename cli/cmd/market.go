package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/nathfavour/shadowprism/cli/internal/agent"
	"github.com/nathfavour/shadowprism/cli/internal/sidecar"
	"github.com/spf13/cobra"
)

var marketCmd = &cobra.Command{
	Use:   "market",
	Short: "Get real-time market pricing from Encrypt.trade",
	Run: func(cmd *cobra.Command, args []string) {
		pa := agent.NewPrismAgent()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cm, _ := sidecar.NewConfigManager()
		socketPath := cm.GetSocketPath()
		client := sidecar.NewCoreClient(socketPath, "dev-token-123")

		res, err := client.GetMarket()
		if err != nil {
			fmt.Printf("❌ Failed to fetch market data: %v\n", err)
			return
		}

		fmt.Printf("📊 Market Data (via %s)\n", res["provider"])
		fmt.Printf("Asset: %s\n", res["asset"])
		fmt.Printf("Price: $%.2f USD\n", res["price_usd"])

		// Agent insight
		resp, _ := pa.Talk(ctx, fmt.Sprintf("The SOL price is $%.2f. Give a very brief market sentiment or tip.", res["price_usd"]))
		pa.DisplayResponse(resp)
	},
}

func init() {
	rootCmd.AddCommand(marketCmd)
}

