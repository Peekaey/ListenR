package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewBubbleteaModel() BubbleteaModel {
	ti := textinput.New()
	ti.Placeholder = "ws://host:port/hub"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60

	return BubbleteaModel{
		status:     "🔄  Connecting...",
		messages:   []string{},
		promptMode: true,
		ti:         ti,
	}
}

// Init runs when the program starts
// https://github.com/charmbracelet/bubbletea?tab=readme-ov-file#initialization
func (m BubbleteaModel) Init() tea.Cmd {

	if m.promptMode == true {
		return textinput.Blink
	}

	return func() tea.Msg {
		return connectToSignalR
	}
}

// Update handles messages and updates the BubbleteaModel
// https://github.com/charmbracelet/bubbletea?tab=readme-ov-file#the-update-method
func (m BubbleteaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.promptMode == true {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.url = m.ti.Value()
				m.promptMode = false
				m.status = "🔄 Connecting..."
				return m, func() tea.Msg {
					return connectToSignalR(m.url)
				}
			case "ctrl+c", "esc":
				return m, tea.Quit
			}
		}

		var cmd tea.Cmd
		m.ti, cmd = m.ti.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "r":
			return m, func() tea.Msg {
				returnResult := connectToSignalR(m.url)
				if returnResult.(connectionStatus).connected == true {
					m.messages = []string{}
				}
				return returnResult
			}
		case "i", "ctrl+i":
			m.promptMode = true

		}

	case connectionStatus:
		if msg.err != nil {
			m.status = "❌  Connection failed"
			m.err = msg.err
			m.connected = false
			return m, nil
		}

		m.connected = msg.connected
		m.status = "✅  Connected!"
		m.err = nil
		m.conn = nil

		// Start listening for messages
		return m, listenForMessages

	case SignalRResponse:
		m.messages = append(m.messages, msg.data)

		// Keep only last 10 messages
		if len(m.messages) > 10 {
			m.messages = m.messages[1:]
		}

		// Continue listening
		return m, listenForMessages

	case SignalRError:
		m.status = "⚠️  Error"
		m.err = msg.err
		m.connected = false
		return m, nil
	}

	return m, nil
}

// View renders the UI
// https://github.com/charmbracelet/bubbletea?tab=readme-ov-file#the-view-method
func (m BubbleteaModel) View() string {

	if m.quitting {
		return "👋  Adios!\n\n"
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1).
		Render("ListenR - Simple SignalR Listener")

	// For Prompt Screen
	if m.promptMode == true {
		msg := "Type the websocket URL of the SignalR hub and press Enter. Ctrl+C to quit."
		return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", titleStyle, msg, m.ti.View(), lipgloss.NewStyle().Faint(true).Render("Enter → connect"))
	}

	// For Normal Screen
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))

	if !m.connected {
		statusStyle = statusStyle.Foreground(lipgloss.Color("#FF6B6B"))
	}

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		MarginLeft(2)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B"))

	s := titleStyle
	s += "\n"
	s += statusStyle.Render(m.status)
	s += "\n\n"

	if m.err != nil {
		s += errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
		s += "\n\n"
	}

	if len(m.messages) > 0 {
		s += lipgloss.NewStyle().Bold(true).Render("Recent Updates:")
		s += "\n"
		for _, msg := range m.messages {
			s += messageStyle.Render("• " + msg)
			s += "\n"
		}
	} else if m.connected {
		s += lipgloss.NewStyle().Italic(true).Render("Waiting for presence updates...")
		s += "\n"
	}

	s += "\n"
	s += lipgloss.NewStyle().Faint(true).Render("Press to 'r' to reconnect, 'q' to quit, 'i' to change URL")

	return s
}
