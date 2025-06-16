package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Update function
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

		// Check minimum dimensions
		if m.termWidth < MinWidth || m.termHeight < MinHeight {
			m.ready = false
			m.err = fmt.Errorf("terminal too small: need at least %dx%d, got %dx%d",
				MinWidth, MinHeight, m.termWidth, m.termHeight)
		} else {
			m.ready = true
			m.err = nil
		}

		return m, nil

	case tea.KeyMsg:
		if !m.ready {
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				return m, tea.Quit
			}
			return m, nil
		}

		// Handle command mode input
		if m.commandMode != NormalMode {
			return m.handleCommandMode(msg)
		}

		// Global navigation keys (only in normal mode)
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.currentScreen == StartScreen {
				return m, tea.Quit
			}
		case m.config.CommandKey:
			m.commandMode = CommandModeActive
			m.commandInput = ""
			m.commandError = ""
			return m, nil
		case m.config.SearchKey:
			m.commandMode = SearchModeActive
			m.commandInput = ""
			m.commandError = ""
			return m, nil
		}

		// Screen-specific navigation
		switch m.currentScreen {
		case StartScreen:
			return m.handleStartScreen(msg)
		case MainScreen:
			return m.handleMainScreen(msg)
		case ConfigScreen:
			return m.handleConfigScreen(msg)
		case ExtrasScreen:
			return m.handleExtrasScreen(msg)
		}

	case ScreenChangeMsg:
		m.currentScreen = msg.screen
		return m, nil
	}

	return m, nil
}

// Handle start screen input
func (m Model) handleStartScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuSelection > 0 {
			m.menuSelection--
		}
	case "down", "j":
		if m.menuSelection < 2 { // 3 menu items (0-2)
			m.menuSelection++
		}
	case "enter", " ":
		switch m.menuSelection {
		case 0:
			return m, func() tea.Msg { return ScreenChangeMsg{MainScreen} }
		case 1:
			return m, func() tea.Msg { return ScreenChangeMsg{ConfigScreen} }
		case 2:
			return m, func() tea.Msg { return ScreenChangeMsg{ExtrasScreen} }
		}
	case "1":
		m.menuSelection = 0
		return m, func() tea.Msg { return ScreenChangeMsg{MainScreen} }
	case "2":
		m.menuSelection = 1
		return m, func() tea.Msg { return ScreenChangeMsg{ConfigScreen} }
	case "3":
		m.menuSelection = 2
		return m, func() tea.Msg { return ScreenChangeMsg{ExtrasScreen} }
	}
	return m, nil
}

// Handle main screen input
func (m Model) handleMainScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// switch msg.String() {
	// case "esc", "b":
	// 	return m, func() tea.Msg { return ScreenChangeMsg{StartScreen} }
	// }
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "tab":
		// Next prompt
		m.promptIndex = (m.promptIndex + 1) % len(m.prompts)
		m.prompt = m.prompts[m.promptIndex]
		m.userInput = ""
		m.currentChar = 0
		m.pressedKeys = make(map[string]bool)

	case "backspace":
		if len(m.userInput) > 0 {
			m.userInput = m.userInput[:len(m.userInput)-1]
			if m.currentChar > 0 {
				m.currentChar--
			}
		}

	default:
		// Handle regular character input
		char := msg.String()
		if len(char) == 1 || char == "space" {
			if char == "space" {
				char = " "
			}

			// Simulate key press
			keyLabel := strings.ToUpper(char)
			if char == " " {
				keyLabel = "SPACE"
			}
			m.pressedKeys[keyLabel] = true

			// Clear pressed keys after a short delay
			go func() {
				time.Sleep(100 * time.Millisecond)
				delete(m.pressedKeys, keyLabel)
			}()

			m.userInput += char
			if m.currentChar < len(m.prompt) {
				m.currentChar++
			}

			// Check if prompt is completed
			if m.userInput == m.prompt {
				// Auto advance to next prompt after completion
				time.AfterFunc(1*time.Second, func() {
					m.promptIndex = (m.promptIndex + 1) % len(m.prompts)
					m.prompt = m.prompts[m.promptIndex]
					m.userInput = ""
					m.currentChar = 0
					m.pressedKeys = make(map[string]bool)
				})
			}
		}
	}
	return m, nil
}

// Handle config screen input
func (m Model) handleConfigScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		return m, func() tea.Msg { return ScreenChangeMsg{StartScreen} }
	}
	return m, nil
}

// Handle extras screen input
func (m Model) handleExtrasScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		return m, func() tea.Msg { return ScreenChangeMsg{StartScreen} }
	}
	return m, nil
}
