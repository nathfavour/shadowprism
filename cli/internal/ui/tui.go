package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	stateDashboard sessionState = iota
	stateShield
	stateSettings
)

type model struct {
	state  sessionState
	cursor int
	width  int
	height int
}

var (
	sidebarStyle = lipgloss.NewStyle().
		Width(20).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("62"))

	mainStyle = lipgloss.NewStyle().
		Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	selectedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("62")).
		Padding(0, 1)

	inactiveItemStyle = lipgloss.NewStyle().
		Padding(0, 1)
)

func InitialModel() model {
	return model{
		state:  stateDashboard,
		cursor: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 2 {
				m.cursor++
			}
		case "enter":
			m.state = sessionState(m.cursor)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	sidebarItems := []string{"Dashboard", "Shield SOL", "Settings"}
	var sidebarBuilder strings.Builder

	sidebarBuilder.WriteString(titleStyle.Render("SHADOW PRISM"))
	sidebarBuilder.WriteString("\n\n")

	for i, item := range sidebarItems {
		if i == m.cursor {
			sidebarBuilder.WriteString(selectedItemStyle.Render(item) + "\n")
		} else {
			sidebarBuilder.WriteString(inactiveItemStyle.Render(item) + "\n")
		}
	}

	sidebar := sidebarStyle.Height(m.height - 4).Render(sidebarBuilder.String())

	var mainContent string
	sswitch m.state {
	case stateDashboard:
		mainContent = m.renderDashboard()
	case stateShield:
		mainContent = m.renderShield()
	case stateSettings:
		mainContent = m.renderSettings()
	}

	main := mainStyle.Width(m.width - 25).Render(mainContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
}

func (m model) renderDashboard() string {
	return fmt.Sprintf(
		"%s\n\n%s\n%s\n%s",
		titleStyle.Render("System Dashboard"),
		"Core Engine: [ONLINE]",
		"Privacy Level: [ULTRA]",
		"Recent Activity: No transactions yet.",
	)
}

func (m model) renderShield() string {
	return titleStyle.Render("Shield Liquidity") + "\n\nPress [S] to initiate a privacy mix."
}

func (m model) renderSettings() string {
	return titleStyle.Render("Settings") + "\n\nEndpoint: localhost:42069\nCompliance: Enabled"
}
