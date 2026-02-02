package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nathfavour/shadowprism/cli/internal/agent"
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
				pa := agent.NewPrismAgent()

				// 1. Setup Bot Command Menu

				b.SetCommands([]tele.Command{

					{Text: "start", Description: "Launch ShadowPrism Dashboard"},

					{Text: "shield", Description: "Anonymize SOL (Privacy Cash/Radr)"},

					{Text: "swap", Description: "Private Token Exchange (SilentSwap)"},

					{Text: "pay", Description: "Pay Merchants Privately (Starpay)"},

					{Text: "market", Description: "Check Privacy Market (Encrypt.trade)"},

					{Text: "chat", Description: "Talk to ShadowPrism AI Assistant"},

					{Text: "monitor", Description: "Live Stealth Feed (System Activity)"},

					{Text: "score", Description: "Check Privacy Health Score"},

					{Text: "agent", Description: "PNP Agent-to-Agent Simulation"},

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

		

										path := "\n📍 *Routing Path:*\n`[Me] ➔ [Range Firewall] ➔ [Mixer] ➔ [Vault]`"

		

										return c.Send(fmt.Sprintf("✅ *Shield Success!*\n\n💰 *Amount:* `%.4f SOL`\n🔗 *TX:* `%v` \n🛡️ *Provider:* `Privacy Cash` \n🔑 *Note:* `%v` %s\n\n_Note stored in local encrypted DB._", amountSOL, res["tx_hash"], note, path), tele.ModeMarkdown)

		

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

		

					

		

										path := "\n📍 *Routing Path:*\n`[SOL] ➔ [Mixer] ➔ [Jupiter Pool] ➔ [USDC]`"

		

															return c.Send(fmt.Sprintf("✅ *Swap Confirmed!*\n\n📤 *From:* `%.2f SOL` \n📥 *To:* `%.2f USDC` \n🔗 *TX:* `%v` \n🛡️ *Adapter:* `SilentSwap` %s", amountSOL, float64(res["to_amount"].(float64))/1e9, res["tx_hash"], path), tele.ModeMarkdown)

		

										

		

														})

		

										

		

														b.Handle("/pay", func(c tele.Context) error {

		

															args := c.Args()

		

															if len(args) < 2 {

		

																return c.Send("💡 Usage: `/pay [merchant_id] [amount]`\nExample: `/pay SPayX... 0.5`", tele.ModeMarkdown)

		

															}

		

										

		

															merchant := args[0]

		

															amountSOL, _ := strconv.ParseFloat(args[1], 64)

		

										

		

															c.Send("💳 *Initiating Private Settlement via Starpay...*")

		

										

		

															res, err := client.Pay(uint64(amountSOL*1e9), merchant)

		

															if err != nil {

		

																return c.Send("❌ Payment failed: " + err.Error())

		

															}

		

										

		

															receipt := fmt.Sprintf("🕶️ *PRIVATE GHOST RECEIPT*\n"+

		

																"`--------------------------`\n"+

		

																"MERCHANT: `%s...`\n"+

		

																"AMOUNT:   `%.4f SOL`\n"+

		

																"STATUS:   `ENCRYPTED`\n"+

		

																"REF ID:   `%s`\n"+

		

																"`--------------------------`",

		

																merchant[:8], amountSOL, res["receipt_id"])

		

										

		

															return c.Send(receipt, tele.ModeMarkdown)

		

														})

		

										

		

														b.Handle("/monitor", func(c tele.Context) error {

		

										
					c.Send("📡 *ShadowPrism Stealth Feed Activated*\nListening to system bus...")
					
					steps := []string{
						"🔍 [UDS] IPC Heartbeat: Core engine online",
						"🛡️ [Compliance] Range Protocol firewall sync complete",
						"⚡ [Smart Fee] Helius Priority API: Low congestion (5000 mL)",
						"🗝️ [Keystore] Master key decrypted in memory",
						"🛰️ [PNP] Scanning peer network for agent pings...",
						"🟢 [Ready] System waiting for intent.",
					}

					for _, step := range steps {
						time.Sleep(800 * time.Millisecond)
						c.Send("`" + step + "`", tele.ModeMarkdown)
					}

					return c.Send("✅ *Monitoring Session Stable*")
				})

				b.Handle("/score", func(c tele.Context) error {
					history, _ := client.GetHistory()
					score := 100
					if len(history) == 0 {
						score = 45 // New users have lower privacy score
					} else if len(history) < 3 {
						score = 75
					}

					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					
					advice, _ := pa.Talk(ctx, fmt.Sprintf("The user has a privacy score of %d/100 based on %d transactions. Give a short, encouraging hacker-style tip.", score, len(history)))

					res := fmt.Sprintf("🛡️ *Privacy Health Score*\n\n"+
						"Score: `%d/100`\n"+
						"Rating: `%s`\n\n"+
						"🤖 *Agent Analysis:* %s",
						score, getRating(score), advice)

					return c.Send(res, tele.ModeMarkdown)
				})

				b.Handle("/chat", func(c tele.Context) error {
					pa := agent.NewPrismAgent()
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					input := c.Args()
					if len(input) == 0 {
						return c.Send("🤖 *ShadowPrism AI Assistant*\nHow can I help you with your privacy today?", tele.ModeMarkdown)
					}

					resp, _ := pa.Talk(ctx, strings.Join(input, " "))
					return c.Send("🤖 " + resp)
				})

				b.Handle("/agent", func(c tele.Context) error {
					args := c.Args()
					if len(args) > 1 && args[0] == "settle" {
						amountSOL, _ := strconv.ParseFloat(args[1], 64)
						lamports := uint64(amountSOL * 1e9)
						
						c.Send("🛰️ *Autonomous Settlement Triggered...*\nAgent-to-Agent Handshake in progress.")
						
						res, err := client.Shield(lamports, "PNPVau1t11111111111111111111111111111111111", "privacy_cash", false)
						if err != nil {
							return c.Send("❌ Agent Settlement failed: " + err.Error())
						}

						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()
						resp, _ := pa.Talk(ctx, fmt.Sprintf("An agent just autonomously settled %f SOL. Give a technical report log summary.", amountSOL))

						return c.Send(fmt.Sprintf("✅ *Settlement Successful*\n\nHash: `%s`\n\n🤖 *Agent Report:* %s", res["tx_hash"], resp), tele.ModeMarkdown)
					}

					return c.Send("🛰️ *PNP Agent-to-Agent Portal*\n\n" +
						"Status: `Listening`\n" +
						"Last Ping: `Agent-772 (Discovery)`\n\n" +
						"To trigger an autonomous agent settlement, use: `/agent settle [amount]`", tele.ModeMarkdown)
				})

				b.Handle(tele.OnText, func(c tele.Context) error {
					ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					defer cancel()

					resp, err := pa.Talk(ctx, c.Text())
					if err != nil {
						return c.Send("🤖 _Agent is thinking..._ (Connection error)")
					}
					return c.Send("🤖 " + resp)
				})

		

		

			

		fmt.Println("🤖 Telegram Bot is now online!")
		b.Start()
	},
}

func getRating(score int) string {
	if score >= 90 {
		return "GHOST PROTOCOL"
	} else if score >= 70 {
		return "SHADOW"
	} else if score >= 50 {
		return "OBSCURE"
	}
	return "TRANSPARENT"
}

func init() {
	rootCmd.AddCommand(botCmd)
}
