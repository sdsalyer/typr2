package main

import (
	"github.com/charmbracelet/lipgloss"
)

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

// TODO: Add theme support
// type Theme struct {
//     Primary   lipgloss.Color
//     Secondary lipgloss.Color
//     Error     lipgloss.Color
//     // ... other colors
// }
//
// func ApplyTheme(theme Theme) {
//     // Update all styles with theme colors
// }
