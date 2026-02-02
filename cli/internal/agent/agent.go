package agent

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	agentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")),
			Bold(true)
	
hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")),
			Italic(true)
)

type PrismAgent struct {
	VibePath string
}

func NewPrismAgent() *PrismAgent {
	// Fallback to searching in common paths
	path := "/home/nathfavour/.local/bin/vibeaura"
	return &PrismAgent{VibePath: path}
}

func (a *PrismAgent) Talk(ctx context.Context, prompt string) (string, error) {
	// We want concise, conversational responses
	fullPrompt := "You are ShadowPrism AI, a privacy-first assistant for Solana. " +
		"Provide a very concise, professional, and slightly futuristic response (max 2 sentences). " +
		"Context: " + prompt

	args := []string{"direct", "--non-interactive"}
	cmd := exec.CommandContext(ctx, a.VibePath, args...)
	cmd.Stdin = strings.NewReader(fullPrompt)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("agent error: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func (a *PrismAgent) DisplayResponse(text string) {
	fmt.Printf("\n%s %s\n", agentStyle.Render("🤖 ShadowPrism:"), text)
}

func (a *PrismAgent) DisplayHint(hint string) {
	fmt.Printf("\n%s %s\n", hintStyle.Render("💡 Hint:"), hint)
}

func (a *PrismAgent) GetHint(ctx context.Context, action string) {
	go func() {
		// Non-blocking hint generation
		h, err := a.Talk(ctx, "Provide a quick privacy tip related to "+action)
		if err == nil && h != "" {
			a.DisplayHint(h)
		}
	}()
	// Small sleep to allow some output flow if desired, but usually we just let it go
}
