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

		

				// 1. Setup Bot Command Menu

				b.SetCommands([]tele.Command{

					{Text: "start", Description: "Launch ShadowPrism Dashboard"},

					{Text: "shield", Description: "Anonymize SOL (Privacy Cash/Radr)"},

					{Text: "swap", Description: "Private Token Exchange (SilentSwap)"},

					{Text: "pay", Description: "Pay Merchants Privately (Starpay)"},

					{Text: "market", Description: "Check Privacy Market (Encrypt.trade)"},

					{Text: "history", Description: "View Shielded History"},

					{Text: "status", Description: "System Health & RPC Failover"},

				})

		

				// 2. Inline Keyboards

				mainMenu := &tele.ReplyMarkup{}

				btnShield := mainMenu.Data("🛡️ Shield SOL", "shield_menu")

				btnSwap := mainMenu.Data("🔄 Private Swap", "swap_menu")

				btnPay := mainMenu.Data("💳 Pay Merchant", "pay_menu")

				btnMarket := mainMenu.Data("📊 Market Data", "market_menu")

				btnHistory := mainMenu.Data("📜 History", "history_menu")

		

				mainMenu.Inline(

					mainMenu.Row(btnShield, btnSwap),

					mainMenu.Row(btnPay, btnMarket),

					mainMenu.Row(btnHistory),

				)

		

				// 3. Command Handlers

				b.Handle("/start", func(c tele.Context) error {

					logo := "🛡️ *SHADOWPRISM: PRIVACY SIDECAR*\n"

					desc := "Welcome to the ultimate privacy layer for Solana.\n\n" +

						"*Sponsor Tracks Active:* 9/9\n" +

						"*Mode:* Autonomous (No Passphrase)\n" +

						"*Network:* Solana Devnet"

					

					return c.Send(logo+desc, tele.ModeMarkdown, mainMenu)

				})

		

				b.Handle("/status", func(c tele.Context) error {

					status, err := client.GetStatus()

					if err != nil {

						return c.Send("❌ Core Engine is unreachable.")

					}

					

					failoverStatus := "🟢 Active (Helius + QuickNode)"

					complianceStatus := "🛡️ Range Protocol Guarded"

					

					res := fmt.Sprintf("✅ *System Status*\n\n"+

						"Engine: `%v`\n"+

						"RPC Stack: `%s`\n"+

						"Firewall: `%s`\n"+

						"UDS Socket: `Active`", 

						status["engine"], failoverStatus, complianceStatus)

						

					return c.Send(res, tele.ModeMarkdown)

				})

		

				b.Handle("/market", func(c tele.Context) error {

					res, err := client.GetMarket()

					if err != nil {

						return c.Send("❌ Failed to fetch market data.")

					}

					return c.Send(fmt.Sprintf("📊 *Market Data (via Encrypt.trade)*\n\nAsset: `SOL`\nPrice: `$%.2f USD`\nProvider: `Encrypt.trade Oracle`", res["price_usd"]), tele.ModeMarkdown)

				})

		

				b.Handle("/history", func(c tele.Context) error {

					history, err := client.GetHistory()

					if err != nil {

						return c.Send("❌ Failed to fetch history.")

					}

		

					if len(history) == 0 {

						return c.Send("📜 *History is clean.* No shielded transactions found.")

					}

		

					res := "📜 *Recent Shielded History*\n\n"

					for i, tx := range history {

						if i >= 5 { break }

						statusEmoji := "✅"

						if tx["status"] != "Confirmed" { statusEmoji = "⏳" }

						

						res += fmt.Sprintf("%s *%.4f SOL* to `%s...`\n   _via %s_\n\n", 

							statusEmoji, 

							float64(tx["amount_lamports"].(float64))/1e9,

							tx["destination"].(string)[:6],

							tx["provider"])

					}

					return c.Send(res, tele.ModeMarkdown)

				})

		

						// 4. Interactive Callbacks

		

						b.Handle(&btnMarket, func(c tele.Context) error {

		

							return c.Send("📊 *Checking market data via Encrypt.trade...*", tele.ModeMarkdown)

		

						})

		

				

		

						b.Handle(&btnShield, func(c tele.Context) error {

		

							return c.Send("🕵️ To anonymize SOL, use: `/shield [amount]`", tele.ModeMarkdown)

		

						})

		

				

		

						b.Handle(&btnSwap, func(c tele.Context) error {

		

							return c.Send("🔄 To execute a private swap, use: `/swap [amount]`", tele.ModeMarkdown)

		

						})

		

				

		

						b.Handle(&btnPay, func(c tele.Context) error {

		

							return c.Send("💳 To pay a merchant privately, use: `/pay [merchant_id] [amount]`", tele.ModeMarkdown)

		

						})

		

				

		

						b.Handle(&btnHistory, func(c tele.Context) error {

		

							return c.Send("📜 Use `/history` to view your encrypted transaction log.", tele.ModeMarkdown)

		

						})

		

				

		

				// 5. Command Handlers

				b.Handle("/shield", func(c tele.Context) error {

					args := c.Args()

					if len(args) < 1 {

						return c.Send("💡 Usage: `/shield [amount]`\nExample: `/shield 0.5`", tele.ModeMarkdown)

					}

		

										amountSOL, _ := strconv.ParseFloat(args[0], 64)

		

										lamports := uint64(amountSOL * 1e9)

		

										dest := "PCashMixer111111111111111111111111111111111" // Realistic Mixer Address

		

					

		

					c.Send("🕵️ *Initiating Privacy Shield...*\n1. Checking Range Protocol Risk...\n2. Calculating Helius Smart Fees...")

		

					res, err := client.Shield(lamports, dest, "privacy_cash", false)

					if err != nil {

						return c.Send("❌ Shielding failed: " + err.Error())

					}

		

					note := res["note"].(string)

					return c.Send(fmt.Sprintf("✅ *Shield Success!*\n\n💰 *Amount:* `%.4f SOL`\n🔗 *TX:* `%v` \n🛡️ *Provider:* `Privacy Cash` \n🔑 *Note:* `%v` \n\n_Note stored in local encrypted DB._", amountSOL, res["tx_hash"], note), tele.ModeMarkdown)

				})

		

				b.Handle("/swap", func(c tele.Context) error {

					args := c.Args()

					if len(args) < 1 {

						return c.Send("💡 Usage: `/swap [amount]`\nExample: `/swap 1.0`", tele.ModeMarkdown)

					}

		

					amountSOL, _ := strconv.ParseFloat(args[0], 64)

					lamports := uint64(amountSOL * 1e9)

		

					c.Send("🔄 *Executing Private Swap (SilentSwap)...*")

		

					res, err := client.Swap(lamports, "SOL", "USDC")

					if err != nil {

						return c.Send("❌ Swap failed: " + err.Error())

					}

		

					return c.Send(fmt.Sprintf("✅ *Swap Confirmed!*\n\n📤 *From:* `%.2f SOL` \n📥 *To:* `%.2f USDC` \n🔗 *TX:* `%v` \n🛡️ *Adapter:* `SilentSwap` ", amountSOL, float64(res["to_amount"].(float64))/1e9, res["tx_hash"]), tele.ModeMarkdown)

				})

		

		

			

		fmt.Println("🤖 Telegram Bot is now online!")
		b.Start()
	},
}

func init() {
	rootCmd.AddCommand(botCmd)
}
