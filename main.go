package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Minimum terminal dimensions
const (
	MinWidth  = 80
	MinHeight = 24
)

// Screen types
type Screen int

const (
	StartScreen Screen = iota
	MainScreen
	ConfigScreen
	ExtrasScreen
)

// Command mode state
type CommandMode int

const (
	NormalMode CommandMode = iota
	CommandModeActive
	SearchModeActive
)

// Application configuration
type Config struct {
	CommandKey string // Key to enter command mode (default ":")
	SearchKey  string // Key to enter search mode (default "/")
}

// Default configuration
func DefaultConfig() Config {
	return Config{
		CommandKey: ":",
		SearchKey:  "/",
	}
}

// Main application model
type Model struct {
	currentScreen Screen
	width         int
	height        int
	ready         bool
	err           error
	menuSelection int // For navigating menu items
	commandMode   CommandMode
	commandInput  string
	commandError  string
	config        Config
}

// Initialize the application
func InitialModel() Model {
	return Model{
		currentScreen: StartScreen,
		ready:         false,
		menuSelection: 0,
		commandMode:   NormalMode,
		commandInput:  "",
		commandError:  "",
		config:        DefaultConfig(),
	}
}

// Init method (required by tea.Model interface)
func (m Model) Init() tea.Cmd {
	return nil
}

// Messages
type ScreenChangeMsg struct {
	screen Screen
}

type SizeCheckMsg struct {
	width  int
	height int
}

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

// Styles
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1)
		
	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196")).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("196")).
		Padding(1, 2).
		Margin(1, 0)
		
	menuStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))
		
	menuItemStyle = lipgloss.NewStyle().
		Padding(0, 2)
		
	selectedMenuItemStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("63")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 2).
		Bold(true)
		
	contentStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Margin(1, 0)
		
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Margin(1, 0)
		
	statusLineStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("255")).
		Padding(0, 1).
		Bold(true)
		
	commandLineStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("234")).
		Foreground(lipgloss.Color("255")).
		Padding(0, 1)
		
	commandErrorStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("196")).
		Foreground(lipgloss.Color("255")).
		Padding(0, 1).
		Bold(true)
)

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

func main() {
	p := tea.NewProgram(InitialModel(), tea.WithAltScreen())
	
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
