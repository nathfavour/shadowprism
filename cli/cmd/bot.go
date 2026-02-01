package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nathfavour/shadowprism/cli/internal/sidecar"
	"github.com/spf13/cobra"
	tele "gopkg.in/telebot.v3"
)

var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Start the ShadowPrism Telegram Bot",
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("TELEGRAM_BOT_TOKEN")
		if token == "" {
			fmt.Println("❌ Error: TELEGRAM_BOT_TOKEN environment variable is not set.")
			os.Exit(1)
		}

		authToken := "dev-token-123"
		manager := sidecar.NewManager(42069, authToken)
		socketPath := "/tmp/shadowprism.sock"
		
		fmt.Println("🚀 Starting ShadowPrism Core for Bot Mode...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			fmt.Printf("❌ Failed to start core engine: %v\n", err)
			os.Exit(1)
		}
		defer manager.Stop()

		pref := tele.Settings{
			Token:  token,
			Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		}

		b, err := tele.NewBot(pref)
		if err != nil {
			fmt.Printf("❌ Failed to start Telegram Bot: %v\n", err)
			return
		}

		client := sidecar.NewCoreClient(socketPath, authToken)

		// Command Handlers
		b.Handle("/start", func(c tele.Context) error {
			return c.Send("🛡️ *Welcome to ShadowPrism*\nYour privacy-first Solana sidecar is active.\n\nCommands:\n/status - Check engine health\n/shield [amount] [address] - Anonymize SOL", tele.ModeMarkdown)
		})

		b.Handle("/status", func(c tele.Context) error {
			status, err := client.GetStatus()
			if err != nil {
				return c.Send("❌ Core Engine is unreachable.")
			}
			return c.Send(fmt.Sprintf("✅ *System Status*\nEngine: %v\nProtocol: %v\nStatus: %v", status["engine"], status["protocol"], status["status"])), tele.ModeMarkdown)
		})

		b.Handle("/shield", func(c tele.Context) error {
			// Basic parsing: /shield 1.0 ADDRESS
			args := c.Args()
			if len(args) < 2 {
				return c.Send("Usage: /shield [amount] [destination_address]")
			}

			c.Send("🕵️ *Initiating Privacy Shield...*\nRouting through Privacy Cash adapters.")

			payload := map[string]interface{}{
				"amount_lamports":  1000000000, // Hardcoded for demo
				"destination_addr": args[1],
				"strategy":         "mix_standard",
			}

			var result map[string]interface{}
			_, err := client.Http.R(). 
				SetBody(payload).
				SetResult(&result).
				Post("/v1/shield")

			if err != nil {
				return c.Send("❌ Shielding failed: Core communication error.")
			}

			return c.Send(fmt.Sprintf("✅ *Shield Success!*\n\n🔗 *TX:* `%v` \n🛡️ *Provider:* %v", result["tx_hash"], result["provider"])), tele.ModeMarkdown
		})

		fmt.Println("🤖 Telegram Bot is now online!")
		b.Start()
	},
}

func init() {
	rootCmd.AddCommand(botCmd)
}
