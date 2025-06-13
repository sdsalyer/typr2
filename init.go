package main

import (
	tea "github.com/charmbracelet/bubbletea"
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

// Messages
type ScreenChangeMsg struct {
	screen Screen
}

// Command mode state
type CommandMode int

const (
	NormalMode CommandMode = iota
	CommandModeActive
	SearchModeActive
)

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
