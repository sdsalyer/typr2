package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Handle command mode input
func (m Model) handleCommandMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Execute command
		cmd := m.executeCommand(m.commandInput)
		m.commandMode = NormalMode
		m.commandInput = ""
		return m, cmd
	case "esc":
		// Cancel command mode
		m.commandMode = NormalMode
		m.commandInput = ""
		m.commandError = ""
		return m, nil
	case "backspace":
		if len(m.commandInput) > 0 {
			m.commandInput = m.commandInput[:len(m.commandInput)-1]
		}
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	default:
		// Add character to command input
		if len(msg.String()) == 1 {
			m.commandInput += msg.String()
		}
		return m, nil
	}
}

// Execute command and return appropriate tea.Cmd
func (m Model) executeCommand(command string) tea.Cmd {
	m.commandError = "" // Clear previous errors

	switch command {
	case "q", "quit":
		return tea.Quit
	case "start", "home":
		return func() tea.Msg { return ScreenChangeMsg{StartScreen} }
	case "main":
		return func() tea.Msg { return ScreenChangeMsg{MainScreen} }
	case "config", "settings":
		return func() tea.Msg { return ScreenChangeMsg{ConfigScreen} }
	case "extras":
		return func() tea.Msg { return ScreenChangeMsg{ExtrasScreen} }
	case "resize":
		// Force a window size check (useful for debugging)
		return func() tea.Msg {
			return tea.WindowSizeMsg{Width: m.width, Height: m.height}
		}
	case "help":
		m.commandError = "Commands: q|quit, start|home, main, config|settings, extras, resize, set"
		return nil
	default:
		// Handle 'set' commands for configuration
		if len(command) > 4 && command[:4] == "set " {
			return m.handleSetCommand(command[4:])
		}

		if command == "" {
			return nil
		}
		m.commandError = fmt.Sprintf("Unknown command: %s (try 'help')", command)
		return nil
	}
}

// Handle 'set' commands for configuration
func (m Model) handleSetCommand(args string) tea.Cmd {
	parts := strings.Fields(args)
	if len(parts) < 2 {
		m.commandError = "Usage: set <option> <value> (try: set commandkey :)"
		return nil
	}

	option := parts[0]
	value := parts[1]

	switch option {
	case "commandkey":
		if len(value) == 1 {
			m.config.CommandKey = value
			m.commandError = fmt.Sprintf("Command key set to '%s'", value)
		} else {
			m.commandError = "Command key must be a single character"
		}
	case "searchkey":
		if len(value) == 1 {
			m.config.SearchKey = value
			m.commandError = fmt.Sprintf("Search key set to '%s'", value)
		} else {
			m.commandError = "Search key must be a single character"
		}
	default:
		m.commandError = fmt.Sprintf("Unknown option: %s (try: commandkey, searchkey)", option)
	}

	return nil
}
