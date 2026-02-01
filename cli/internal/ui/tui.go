package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
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
type shieldResultMsg map[string]interface{}
type tickMsg time.Time

type model struct {
	state        sessionState
	cursor       int
	width        int
	height       int
	client       *sidecar.CoreClient
	lastStatus   map[string]interface{}
	lastHistory  []map[string]interface{}
	inputs       []textinput.Model
	focusedInput int
	isShielding  bool
	shieldResult string
	err          error
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

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func InitialModel(client *sidecar.CoreClient) model {
	m := model{
		state:  stateDashboard,
		cursor: 0,
		client: client,
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = focusedStyle
		t.CharLimit = 64

		switch i {
		case 0:
			t.Placeholder = "Amount (Lamports, e.g. 100000000)"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Destination Address (Base58)"
		}

		m.inputs[i] = t
	}

	return m
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

func (m model) runShield() tea.Cmd {
	return func() tea.Msg {
		amountStr := m.inputs[0].Value()
		dest := m.inputs[1].Value()

		// Simplified parsing: assume input is valid uint64 lamports
		var amount uint64
		fmt.Sscanf(amountStr, "%d", &amount)

		payload := map[string]interface{}{
			"amount_lamports":  amount,
			"destination_addr": dest,
			"strategy":         "mix_standard",
		}

		var result map[string]interface{}
		_, err := m.client.Http.R().
			SetBody(payload).
			SetResult(&result).
			Post("/v1/shield")

		if err != nil {
			return err
		}
		return shieldResultMsg(result)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.lastStatus = msg
	case historyMsg:
		m.lastHistory = msg
	case shieldResultMsg:
		m.isShielding = false
		m.shieldResult = fmt.Sprintf("Success! TX: %s", msg["tx_hash"])
		return m, m.fetchHistory()
	case tickMsg:
		return m, tea.Batch(m.fetchStatus(), m.fetchHistory(), m.tick())
	case error:
		m.err = msg
		m.isShielding = false
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			if m.state == stateDashboard {
				if msg.String() == "up" || msg.String() == "k" {
					if m.cursor > 0 {
						m.cursor--
					}
				} else {
					if m.cursor < 2 {
						m.cursor++
					}
				}
			} else if m.state == stateShield {
				s := msg.String()
				if s == "up" || s == "shift+tab" {
					m.focusedInput--
				} else {
					m.focusedInput++
				}

				if m.focusedInput < 0 {
					m.focusedInput = len(m.inputs) - 1
				} else if m.focusedInput >= len(m.inputs) {
					m.focusedInput = 0
				}

				cmds := make([]tea.Cmd, len(m.inputs))
				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.focusedInput {
						cmds[i] = m.inputs[i].Focus()
						m.inputs[i].PromptStyle = focusedStyle
						m.inputs[i].TextStyle = focusedStyle
						continue
					}
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = blurredStyle
					m.inputs[i].TextStyle = blurredStyle
				}
				return m, tea.Batch(cmds...)
			}

		case "enter":
			if m.state == stateDashboard {
				m.state = sessionState(m.cursor)
			} else if m.state == stateShield {
				m.isShielding = true
				m.shieldResult = ""
				return m, m.runShield()
			}

		case "esc":
			if m.state != stateDashboard {
				m.state = stateDashboard
				m.shieldResult = ""
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.state == stateShield {
		var cmd tea.Cmd
		m.inputs[m.focusedInput], cmd = m.inputs[m.focusedInput].Update(msg)
		return m, cmd
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
	var b strings.Builder

	b.WriteString(titleStyle.Render("Shield Liquidity") + "\n\n")
	b.WriteString("Enter the anonymization details below:\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View() + "\n")
	}

	b.WriteString("\n")
	if m.isShielding {
		b.WriteString(focusedStyle.Render("Processing privacy mix... Please wait."))
	} else if m.shieldResult != "" {
		b.WriteString(statusStyle.Render(m.shieldResult))
	} else {
		b.WriteString(blurredStyle.Render("Press [Tab] to switch fields • [Enter] to Mix • [Esc] for Dashboard"))
	}

	if m.err != nil {
		b.WriteString("\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return b.String()
}

func (m model) renderSettings() string {
	return titleStyle.Render("Settings") + "\n\nEndpoint: " + m.client.Socket + "\nCompliance: Enabled"
}