package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Update function
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Check minimum dimensions
		if m.width < MinWidth || m.height < MinHeight {
			m.ready = false
			m.err = fmt.Errorf("terminal too small: need at least %dx%d, got %dx%d",
				MinWidth, MinHeight, m.width, m.height)
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
	switch msg.String() {
	case "esc", "b":
		return m, func() tea.Msg { return ScreenChangeMsg{StartScreen} }
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
