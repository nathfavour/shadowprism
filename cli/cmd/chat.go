package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nathfavour/shadowprism/cli/internal/agent"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Talk to ShadowPrism AI for privacy advice and system help",
	Run: func(cmd *cobra.Command, args []string) {
		pa := agent.NewPrismAgent()
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("🌐 ShadowPrism Conversational AI Online.")
		fmt.Println("Type 'exit' or 'quit' to leave the chat.")
		fmt.Println("")

		for {
			fmt.Print("👤 You: ")
			if !scanner.Scan() {
				break
			}
			input := scanner.Text()
			if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
				break
			}

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			resp, err := pa.Talk(ctx, input)
			cancel()

			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				continue
			}

			pa.DisplayResponse(resp)
			fmt.Println("")
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
