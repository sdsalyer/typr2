package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View function
func (m Model) View() string {
	if !m.ready {
		return m.renderWithStatusLine(m.renderErrorScreen())
	}

	var content string
	switch m.currentScreen {
	case StartScreen:
		content = m.renderStartScreen()
	case MainScreen:
		content = m.renderMainScreen()
	case ConfigScreen:
		content = m.renderConfigScreen()
	case ExtrasScreen:
		content = m.renderExtrasScreen()
	default:
		content = "Unknown screen"
	}

	return m.renderWithStatusLine(content)
}

// Render error screen for insufficient terminal size
func (m Model) renderErrorScreen() string {
	var content string

	if m.err != nil {
		content = errorStyle.Render(fmt.Sprintf("❌ %s", m.err.Error()))
	} else {
		content = errorStyle.Render("❌ Terminal size check failed")
	}

	help := helpStyle.Render("Please resize your terminal and try again. Press 'q' or Ctrl+C to quit.")

	ui := lipgloss.JoinVertical(lipgloss.Left, content, help)
	return m.centerContent(ui)
}

// Render start screen with menu
func (m Model) renderStartScreen() string {
	title := titleStyle.Render("🚀 My TUI Application")

	// Menu items with selection highlighting
	menuItems := []string{
		"Main Application",
		"Configuration",
		"Extras",
	}

	var renderedItems []string
	for i, item := range menuItems {
		if i == m.menuSelection {
			renderedItems = append(renderedItems, selectedMenuItemStyle.Render(fmt.Sprintf("▶ %s", item)))
		} else {
			renderedItems = append(renderedItems, menuItemStyle.Render(fmt.Sprintf("  %s", item)))
		}
	}

	menuContent := lipgloss.JoinVertical(lipgloss.Left,
		"Welcome! Choose an option:",
		"",
		renderedItems[0],
		renderedItems[1],
		renderedItems[2],
		"",
		"Navigation:",
		"• Use ↑/↓ or j/k to navigate",
		"• Press Enter or Space to select",
		"• Use number keys for direct access",
		"• Press 'q' or Ctrl+C to quit",
	)

	menu := menuStyle.Render(menuContent)
	status := helpStyle.Render(fmt.Sprintf("Terminal: %dx%d", m.width, m.height))

	ui := lipgloss.JoinVertical(lipgloss.Left, title, menu, status)
	return m.centerContent(ui)
}

// Render main screen
func (m Model) renderMainScreen() string {
	title := titleStyle.Render("📱 Main Application")

	content := contentStyle.Render(`This is the main application screen.

Add your main application logic here.

Features could include:
• Data display
• Interactive forms  
• Real-time updates
• File operations`)

	help := helpStyle.Render("Press 'Esc' or 'b' to go back • 'q' or Ctrl+C to quit")

	ui := lipgloss.JoinVertical(lipgloss.Left, title, content, help)
	return m.centerContent(ui)
}

// Render configuration screen
func (m Model) renderConfigScreen() string {
	title := titleStyle.Render("⚙️  Configuration")

	msg := fmt.Sprintf(`	Configuration settings:

• Theme settings
• User preferences  
• Data sources
• Export options
• Keyboard shortcuts

Current Configuration:
• Command key: '%[1]s' (use '%[1]shelp' for commands)
• Search key: '%[1]s' (reserved for future search)

Try these commands:
• %[1]sset commandkey ; (change to semicolon)
• %[1]sset commandkey : (change to colon - default)
• %[1]sset searchkey ? (change search key)`, m.config.CommandKey)
	content := contentStyle.Render(msg)

	help := helpStyle.Render("Press 'Esc' or 'b' to go back • 'q' or Ctrl+C to quit")

	ui := lipgloss.JoinVertical(lipgloss.Left, title, content, help)
	return m.centerContent(ui)
}

// Render extras screen
func (m Model) renderExtrasScreen() string {
	title := titleStyle.Render("✨ Extras")

	content := contentStyle.Render(`Additional features:

• Help documentation
• About information
• Debug tools
• Export data
• Import settings

[Extra features would go here]`)

	help := helpStyle.Render("Press 'Esc' or 'b' to go back • 'q' or Ctrl+C to quit")

	ui := lipgloss.JoinVertical(lipgloss.Left, title, content, help)
	return m.centerContent(ui)
}

// Center content both horizontally and vertically
func (m Model) centerContent(content string) string {
	// Reserve space for status line (subtract 1 from height)
	availableHeight := m.height - 1

	// Calculate content dimensions
	contentWidth := lipgloss.Width(content)
	contentHeight := lipgloss.Height(content)

	// Calculate padding needed for centering
	horizontalPadding := (m.width - contentWidth) / 2
	verticalPadding := (availableHeight - contentHeight) / 2

	// Ensure padding is not negative
	if horizontalPadding < 0 {
		horizontalPadding = 0
	}
	if verticalPadding < 0 {
		verticalPadding = 0
	}

	// Apply centering (using available height minus status line)
	centeredStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(availableHeight).
		Align(lipgloss.Center, lipgloss.Center)

	return centeredStyle.Render(content)
}

// Render content with status line
func (m Model) renderWithStatusLine(content string) string {
	statusLine := m.renderStatusLine()

	// If we're in command mode, show command line instead of just status
	if m.commandMode != NormalMode {
		commandLine := m.renderCommandLine()
		return lipgloss.JoinVertical(lipgloss.Left, content, commandLine)
	}

	return lipgloss.JoinVertical(lipgloss.Left, content, statusLine)
}

// Render vim-style status line
func (m Model) renderStatusLine() string {
	// Get screen name
	var screenName string
	switch m.currentScreen {
	case StartScreen:
		screenName = "START"
	case MainScreen:
		screenName = "MAIN"
	case ConfigScreen:
		screenName = "CONFIG"
	case ExtrasScreen:
		screenName = "EXTRAS"
	}

	// Left side: screen info
	leftInfo := fmt.Sprintf("| %s |", screenName)

	// Right side: terminal dimensions and help
	rightInfo := fmt.Sprintf("[ %dx%d ]", m.width, m.height)

	// Calculate spacing
	totalUsed := len(leftInfo) + len(rightInfo) + 2
	spacing := m.width - totalUsed
	if spacing < 0 {
		spacing = 0
	}

	statusContent := leftInfo + strings.Repeat(" ", spacing) + rightInfo

	return statusLineStyle.Width(m.width).Render(statusContent)
}

// Render command line (when in command mode)
func (m Model) renderCommandLine() string {
	var prefix string
	switch m.commandMode {
	case CommandModeActive:
		prefix = m.config.CommandKey
	case SearchModeActive:
		prefix = m.config.SearchKey
	}

	// Show error if there is one, otherwise show command input
	if m.commandError != "" {
		return commandErrorStyle.Width(m.width).Render(fmt.Sprintf(" %s", m.commandError))
	}

	commandContent := fmt.Sprintf(" %s%s", prefix, m.commandInput)

	// Add cursor indicator
	if m.commandMode != NormalMode {
		commandContent += "█"
	}

	// Pad to full width
	if len(commandContent) < m.width {
		commandContent += strings.Repeat(" ", m.width-len(commandContent))
	}

	return commandLineStyle.Width(m.width).Render(commandContent)
}
