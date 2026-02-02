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
	stateSwap
	statePay
	stateSettings
)

type statusMsg map[string]interface{}
type historyMsg []map[string]interface{}
type shieldResultMsg map[string]interface{}
type swapResultMsg map[string]interface{}
type payResultMsg map[string]interface{}
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
	isWorking    bool
	result       string
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
	}
	m.resetInputs()
	return m
}

func (m *model) resetInputs() {
	m.focusedInput = 0
	m.result = ""
	m.err = nil

	switch m.state {
	case stateShield:
		m.inputs = make([]textinput.Model, 2)
		for i := range m.inputs {
			t := textinput.New()
			t.Cursor.Style = focusedStyle
			if i == 0 {
				t.Placeholder = "Amount (Lamports)"
				t.Focus()
			} else {
				t.Placeholder = "Destination Address"
			}
			m.inputs[i] = t
		}
	case stateSwap:
		m.inputs = make([]textinput.Model, 3)
		for i := range m.inputs {
			t := textinput.New()
			t.Cursor.Style = focusedStyle
			switch i {
			case 0:
				t.Placeholder = "Amount (Lamports)"
				t.Focus()
			case 1:
				t.Placeholder = "From Token (e.g. SOL)"
			case 2:
				t.Placeholder = "To Token (e.g. USDC)"
			}
			m.inputs[i] = t
		}
	case statePay:
		m.inputs = make([]textinput.Model, 2)
		for i := range m.inputs {
			t := textinput.New()
			t.Cursor.Style = focusedStyle
			if i == 0 {
				t.Placeholder = "Amount (Lamports)"
				t.Focus()
			} else {
				t.Placeholder = "Merchant ID (Address)"
			}
			m.inputs[i] = t
		}
	default:
		m.inputs = nil
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

func (m model) runShield() tea.Cmd {
	return func() tea.Msg {
		var amount uint64
		fmt.Sscanf(m.inputs[0].Value(), "%d", &amount)
		dest := m.inputs[1].Value()

		res, err := m.client.Shield(amount, dest, "mix_standard", false)
		if err != nil {
			return err
		}
		return shieldResultMsg(res)
	}
}

func (m model) runSwap() tea.Cmd {
	return func() tea.Msg {
		var amount uint64
		fmt.Sscanf(m.inputs[0].Value(), "%d", &amount)
		from := m.inputs[1].Value()
		to := m.inputs[2].Value()

		res, err := m.client.Swap(amount, from, to)
		if err != nil {
			return err
		}
		return swapResultMsg(res)
	}
}

func (m model) runPay() tea.Cmd {
	return func() tea.Msg {
		var amount uint64
		fmt.Sscanf(m.inputs[0].Value(), "%d", &amount)
		merchant := m.inputs[1].Value()

		res, err := m.client.Pay(amount, merchant)
		if err != nil {
			return err
		}
		return payResultMsg(res)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.lastStatus = msg
	case historyMsg:
		m.lastHistory = msg
	case shieldResultMsg:
		m.isWorking = false
		note := ""
		if n, ok := msg["note"].(string); ok && n != "" {
			note = fmt.Sprintf("\n🔑 Note: %s", n)
		}
		m.result = fmt.Sprintf("Shield Success! TX: %s%s", msg["tx_hash"], note)
		return m, m.fetchHistory()

	case swapResultMsg:
		m.isWorking = false
		m.result = fmt.Sprintf("Swap Success! TX: %s\nRecieved: %v", msg["tx_hash"], msg["to_amount"])
		return m, m.fetchHistory()

	case payResultMsg:
		m.isWorking = false
		m.result = fmt.Sprintf("Payment Confirmed! TX: %s\nReceipt: %s", msg["tx_hash"], msg["receipt_id"])
		return m, m.fetchHistory()

	case tickMsg:
		return m, tea.Batch(m.fetchStatus(), m.fetchHistory(), m.tick())
	case error:
		m.err = msg
		m.isWorking = false
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
					if m.cursor < 4 { // dashboard, shield, swap, pay, settings
						m.cursor++
					}
				}
			} else if m.state == stateShield || m.state == stateSwap || m.state == statePay {
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
				for i := range m.inputs {
					if i == m.focusedInput {
						cmds[i] = m.inputs[i].Focus()
						continue
					}
					m.inputs[i].Blur()
				}
				return m, tea.Batch(cmds...)
			}

		case "enter":
			if m.state == stateDashboard {
				m.state = sessionState(m.cursor)
				m.resetInputs()
			} else if (m.state == stateShield || m.state == stateSwap || m.state == statePay) && !m.isWorking {
				m.isWorking = true
				m.result = ""
				m.err = nil
				switch m.state {
				case stateShield:
					return m, m.runShield()
				case stateSwap:
					return m, m.runSwap()
				case statePay:
					return m, m.runPay()
				}
			}

		case "esc":
			if m.state != stateDashboard {
				m.state = stateDashboard
				m.cursor = 0
				m.resetInputs()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.inputs != nil {
		var cmd tea.Cmd
		m.inputs[m.focusedInput], cmd = m.inputs[m.focusedInput].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	sidebarItems := []string{"Dashboard", "Shield SOL", "Private Swap", "Pay Merchant", "Settings"}
	var sidebarBuilder strings.Builder

	sidebarBuilder.WriteString(titleStyle.Render("SHADOW PRISM"))
	sidebarBuilder.WriteString("\n\n")

	for i, item := range sidebarItems {
		if i == m.cursor && m.state == stateDashboard {
			sidebarBuilder.WriteString(selectedItemStyle.Render(item) + "\n")
		} else if i == int(m.state) {
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
		mainContent = m.renderForm("Shield Liquidity", "Enter the anonymization details below:")
	case stateSwap:
		mainContent = m.renderForm("Private Swap", "Enter the swap details below:")
	case statePay:
		mainContent = m.renderForm("Pay Merchant", "Enter the payment details below:")
	case stateSettings:
		mainContent = m.renderSettings()
	}

	main := mainStyle.Width(m.width - 25).Render(mainContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
}

func (m model) renderForm(title, subtitle string) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(title) + "\n\n")
	b.WriteString(subtitle + "\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View() + "\n")
	}

	b.WriteString("\n")
	if m.isWorking {
		b.WriteString(focusedStyle.Render("Processing... Please wait."))
	} else if m.result != "" {
		b.WriteString(statusStyle.Render(m.result))
	} else {
		b.WriteString(blurredStyle.Render("Press [Tab] to switch fields • [Enter] to Execute • [Esc] for Dashboard"))
	}

	if m.err != nil {
		b.WriteString("\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return b.String()
}

func (m model) renderDashboard() string {
	statusLine := "Core Engine: [OFFLINE]"
	if m.lastStatus != nil {
		statusLine = fmt.Sprintf("Core Engine: %s [ONLINE]", statusStyle.Render(fmt.Sprintf("%v", m.lastStatus["engine"])))
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

func (m model) renderSettings() string {
	return titleStyle.Render("Settings") + "\n\nEndpoint: " + m.client.Socket + "\nCompliance: Enabled"
}