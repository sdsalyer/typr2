package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		log.Printf("logging to file: %s", f.Name())
		defer f.Close()
	}

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <keyboard-layout.json>")
	}

	p := tea.NewProgram(InitialModel(os.Args[1]), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
