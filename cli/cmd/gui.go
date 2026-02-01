package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/nathfavour/shadowprism/cli/internal/ui"
	"github.com/spf13/cobra"
)

var guiCmd = &cobra.RootCmd{
	Use:   "gui",
	Short: "Launch the ShadowPrism TUI",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
}
