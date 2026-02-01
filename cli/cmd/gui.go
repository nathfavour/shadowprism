package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nathfavour/shadowprism/cli/internal/sidecar"
	"github.com/nathfavour/shadowprism/cli/internal/ui"
	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch the ShadowPrism TUI",
	Run: func(cmd *cobra.Command, args []string) {
		manager := sidecar.NewManager(42069, "dev-token-123")
		
		fmt.Println("🚀 Starting ShadowPrism Core...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			fmt.Printf("❌ Failed to start core engine: %v\n", err)
			fmt.Println("Did you run 'cargo build' in the /core directory?")
			os.Exit(1)
		}
		defer manager.Stop()

		client := sidecar.NewCoreClient(cm.GetSocketPath(), authToken)
		p := bubbletea.NewProgram(ui.InitialModel(client), bubbletea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
}