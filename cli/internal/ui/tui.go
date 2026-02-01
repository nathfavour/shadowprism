package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nathfavour/shadowprism/cli/internal/sidecar"
)

type sessionState int

const (
	stateDashboard sessionState = iota
	stateShield
	stateSettings
)

type statusMsg map[string]interface{}
type historyMsg []map[string]interface{}
type tickMsg time.Time

type model struct {
	state       sessionState
	cursor      int
	width       int
	height      int
	client      *sidecar.CoreClient
	lastStatus  map[string]interface{}
	lastHistory []map[string]interface{}
	err         error
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

	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	historyItemStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			MarginBottom(1)

	selectedItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)

	inactiveItemStyle = lipgloss.NewStyle().
			Padding(0, 1)
)

func InitialModel(client *sidecar.CoreClient) model {
	return model{
		state:  stateDashboard,
		cursor: 0,
		client: client,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchStatus(),
		m.fetchHistory(),
		m.tick(),
	)
}

func (m model) tick() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) fetchStatus() tea.Cmd {
	return func() tea.Msg {
		status, err := m.client.GetStatus()
		if err != nil {
			return err
		}
		return statusMsg(status)
	}
}

func (m model) fetchHistory() tea.Cmd {
	return func() tea.Msg {
		history, err := m.client.GetHistory()
		if err != nil {
			return err
		}
		return historyMsg(history)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.lastStatus = msg
	case historyMsg:
		m.lastHistory = msg
	case tickMsg:
		return m, tea.Batch(m.fetchStatus(), m.fetchHistory(), m.tick())
	case error:
		m.err = msg
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
	switch m.state {
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
	statusLine := "Core Engine: [OFFLINE]"
	if m.lastStatus != nil {
		statusLine = fmt.Sprintf("Core Engine: %s [ONLINE]", statusStyle.Render(fmt.Sprintf("%v", m.lastStatus["protocol"])))
	}

	var historyBuilder strings.Builder
	historyBuilder.WriteString(titleStyle.Render("Recent Transactions") + "\n")

	if len(m.lastHistory) == 0 {
		historyBuilder.WriteString("No transactions found.")
	} else {
		for _, tx := range m.lastHistory {
			item := fmt.Sprintf("ID: %s\nStatus: %s\nProvider: %s\nHash: %s",
				tx["id"], tx["status"], tx["provider"], tx["tx_hash"])
			historyBuilder.WriteString(historyItemStyle.Render(item) + "\n")
		}
	}

	return fmt.Sprintf(
		"%s\n%s\n\n%s",
		titleStyle.Render("System Dashboard"),
		statusLine,
		historyBuilder.String(),
	)
}

func (m model) renderShield() string {
	return titleStyle.Render("Shield Liquidity") + "\n\nPress [S] to initiate a privacy mix."
}

func (m model) renderSettings() string {
	return titleStyle.Render("Settings") + "\n\nEndpoint: " + m.client.Socket + "\nCompliance: Enabled"
}
