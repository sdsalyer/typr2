package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	prompt        string
	userInput     string
	onScreenKeys  []string
	highlightedKey string
	spinner       spinner.Model
	err           error
	score         int
	started       bool
	finished      bool
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		onScreenKeys: []string{"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P", "A", "S", "D", "F", "G", "H", "J", "K", "L", "Z", "X", "C", "V", "B", "N", "M"},
		spinner:      s,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c": //, "q":
			return m, tea.Quit
		case "enter":
			if !m.started {
				m.started = true
				m.prompt = m.generatePrompt()
			} else if !m.finished {
				if m.userInput == m.prompt {
					m.score++
					m.err = nil
					m.prompt = m.generatePrompt()
					m.userInput = ""
				} else {
					m.err = fmt.Errorf("Incorrect. Try again.")
					m.userInput = ""
				}
			} else {
				return m, tea.Quit
			}
		default:
			if len(m.userInput) < len(m.prompt) {
				m.userInput += msg.String()
				m.highlightedKey = msg.String()
			}
		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("\n   Typr2 - Let's Go\n\n")

	if !m.started {
		b.WriteString("   Press Enter to start.\n\n")
	} else {
		b.WriteString(fmt.Sprintf("   Prompt: %s\n\n   ", m.prompt))

		for i, key := range m.onScreenKeys {
			if key == m.highlightedKey {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(key))
			} else {
				b.WriteString(key)
			}
			if (i+1)%10 == 0 {
				b.WriteString("\n   ")
			} else {
				b.WriteString(" ")
			}
		}

		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("   Your input: %s\n\n", m.userInput))

		if m.err != nil {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(m.err.Error()))
			b.WriteString("\n\n")
		}

		b.WriteString(fmt.Sprintf("   Score: %d\n\n", m.score))
	}

	b.WriteString(m.spinner.View())
	b.WriteString("\n\n")

	return b.String()
}

func (m *model) generatePrompt() string {
	// Generate a random prompt for the user to type
	rand.Seed(time.Now().UnixNano())
	words := []string{
		"The five boxing wizards jump quickly.",
		"Pack my box with five dozen liquor jugs.",
		"When zombies arrive, quickly fax Judge Pat.",
		"Amazingly few discotheques provide jukeboxes.",
		"The quick onyx goblin jumps over the lazy dwarf.",
	}
	return words[rand.Intn(len(words))]
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

