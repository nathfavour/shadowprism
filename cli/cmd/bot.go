package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nathfavour/shadowprism/cli/internal/sidecar"
	"github.com/spf13/cobra"
	tele "gopkg.in/telebot.v3"
)

var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Start the ShadowPrism Telegram Bot",
	Run: func(cmd *cobra.Command, args []string) {
		cm, err := sidecar.NewConfigManager()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

		token, err := cm.LoadSecret("tg_bot_token")
		if err != nil {
			token = os.Getenv("TELEGRAM_BOT_TOKEN")
		}

		if token == "" {
			fmt.Println("❌ Error: Telegram Bot token not found in config or environment.")
			fmt.Println("Run: shadowprism config set-bot-token <your_token>")
			os.Exit(1)
		}

		authToken := "dev-token-123"

		manager := sidecar.NewManager(42069, authToken)

		
		socketPath := cm.GetSocketPath()
		
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

					return c.Send(fmt.Sprintf("✅ *System Status*\nEngine: %v\nProtocol: %v\nStatus: %v", status["engine"], status["protocol"], status["status"]), tele.ModeMarkdown)

				})

		

				b.Handle("/market", func(c tele.Context) error {

					res, err := client.GetMarket()

					if err != nil {

						return c.Send("❌ Failed to fetch market data.")

					}

					return c.Send(fmt.Sprintf("📊 *Market Data (via Encrypt.trade)*\nAsset: %v\nPrice: $%v USD", res["asset"], res["price_usd"]), tele.ModeMarkdown)

				})

		

				b.Handle("/shield", func(c tele.Context) error {

					args := c.Args()

					if len(args) < 2 {

						return c.Send("Usage: /shield [amount] [destination_address]")

					}

		

					c.Send("🕵️ *Initiating Privacy Shield...*\nRouting through Privacy Cash adapters.")

		

					amount, _ := strconv.ParseUint(args[0], 10, 64)

					res, err := client.Shield(amount, args[1], "privacy_cash", false)

		

					if err != nil {

						return c.Send("❌ Shielding failed: " + err.Error())

					}

		

					note := "N/A"

					if n, ok := res["note"].(string); ok {

						note = n

					}

		

					return c.Send(fmt.Sprintf("✅ *Shield Success!*\n\n🔗 *TX:* `%v` \n🛡️ *Provider:* %v\n🔑 *Note:* `%v`", res["tx_hash"], res["provider"], note), tele.ModeMarkdown)

				})

		

				b.Handle("/swap", func(c tele.Context) error {

					args := c.Args()

					if len(args) < 3 {

						return c.Send("Usage: /swap [amount] [from] [to]")

					}

		

					amount, _ := strconv.ParseUint(args[0], 10, 64)

					res, err := client.Swap(amount, args[1], args[2])

					if err != nil {

						return c.Send("❌ Swap failed: " + err.Error())

					}

		

					return c.Send(fmt.Sprintf("🔄 *Private Swap Confirmed!*\n\n🔗 *TX:* `%v` \n💰 *Received:* %v %v", res["tx_hash"], res["to_amount"], args[2]), tele.ModeMarkdown)

				})

		

				b.Handle("/pay", func(c tele.Context) error {

					args := c.Args()

					if len(args) < 2 {

						return c.Send("Usage: /pay [merchant_id] [amount]")

					}

		

					amount, _ := strconv.ParseUint(args[1], 10, 64)

					res, err := client.Pay(amount, args[0])

					if err != nil {

						return c.Send("❌ Payment failed: " + err.Error())

					}

		

					return c.Send(fmt.Sprintf("💳 *Private Payment Sent!*\n\n🔗 *TX:* `%v` \n🧾 *Receipt:* `%v`", res["tx_hash"], res["receipt_id"]), tele.ModeMarkdown)

				})

		

			

		fmt.Println("🤖 Telegram Bot is now online!")
		b.Start()
	},
}

func init() {
	rootCmd.AddCommand(botCmd)
}
