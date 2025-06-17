package main

import (
	"log"

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
	termWidth     int
	termHeight    int
	ready         bool
	err           error
	menuSelection int // For navigating menu items
	commandMode   CommandMode
	commandInput  string
	commandError  string
	config        Config
	keyboard      Keyboard
	prompt        string
	userInput     string
	currentChar   int
	prompts       []string
	promptIndex   int
	pressedKeys   map[string]bool
}

// Initialize the application
func InitialModel(config string) Model {
	log.Println("init.InitialModel()")

	var prompts = []string{
		"Pack my box with five dozen liquor jugs.",
		"The quick brown fox jumps over the lazy dog.",
		"Waltz, bad nymph, for quick jigs vex.",
		"How vexingly quick daft zebras jump!",
		"Bright vixens jump; dozy fowl quack.",
	}

	kb, err := loadKeyboard(config)
	if err != nil {
		log.Fatalf("Failed to load keyboard: %v", err)
	}

	return Model{
		currentScreen: StartScreen,
		ready:         false,
		menuSelection: 0,
		commandMode:   NormalMode,
		commandInput:  "",
		commandError:  "",
		config:        DefaultConfig(),
		keyboard:      kb,
		prompts:       prompts,
		promptIndex:   0,
		prompt:        prompts[0],
		userInput:     "",
		currentChar:   0,
		pressedKeys:   make(map[string]bool),
	}
}

// Init method (required by tea.Model interface)
func (m Model) Init() tea.Cmd {
	log.Println("init.Init()")
	return nil
}
