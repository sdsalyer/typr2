package main

import (
	"fmt"
	"log"
	// "math"
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
	status := helpStyle.Render(fmt.Sprintf("Terminal: %dx%d", m.termWidth, m.termHeight))

	ui := lipgloss.JoinVertical(lipgloss.Left, title, menu, status)
	return m.centerContent(ui)
}

// Render main screen
func (m Model) renderMainScreen() string {
	log.Println("view.renderMainScreen()")
	// 	title := titleStyle.Render("📱 Main Application")
	//
	// 	content := contentStyle.Render(`This is the main application screen.
	//
	// Add your main application logic here.
	//
	// Features could include:
	// • Data display
	// • Interactive forms
	// • Real-time updates
	// • File operations`)
	//
	// 	help := helpStyle.Render("Press 'Esc' or 'b' to go back • 'q' or Ctrl+C to quit")
	//
	// 	ui := lipgloss.JoinVertical(lipgloss.Left, title, content, help)
	// 	return m.centerContent(ui)

	// title := titleStyle.Render("📱 Main Application")
	// Calculate dimensions
	// if h=24, then 8 rows for prompt and 16 for kb,
	// or approx 3 rows per "key" if 5 rows
	promptHeight := m.termHeight / 3
	keyboardHeight := m.termHeight - promptHeight

	// Build keyboard display
	promptSection := m.renderPrompt(promptHeight)
	keyboardSection := m.renderKeyboard(keyboardHeight)

	// return promptSection + "\n" + keyboardSection
	ui := lipgloss.JoinVertical(lipgloss.Left, promptSection, keyboardSection)
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
	availableHeight := m.termHeight - 1

	// Calculate content dimensions
	contentWidth := lipgloss.Width(content)
	contentHeight := lipgloss.Height(content)

	// Calculate padding needed for centering
	horizontalPadding := (m.termWidth - contentWidth) / 2
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
		Width(m.termWidth).
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
	rightInfo := fmt.Sprintf("[ %dx%d ]", m.termWidth, m.termHeight)

	// Calculate spacing
	totalUsed := len(leftInfo) + len(rightInfo) + 2
	spacing := max(m.termWidth-totalUsed, 0)

	statusContent := leftInfo + strings.Repeat(" ", spacing) + rightInfo

	return statusLineStyle.Width(m.termWidth).Render(statusContent)
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
		return commandErrorStyle.Width(m.termWidth).Render(fmt.Sprintf(" %s", m.commandError))
	}

	commandContent := fmt.Sprintf(" %s%s", prefix, m.commandInput)

	// Add cursor indicator
	if m.commandMode != NormalMode {
		commandContent += "█"
	}

	// Pad to full width
	if len(commandContent) < m.termWidth {
		commandContent += strings.Repeat(" ", m.termWidth-len(commandContent))
	}

	return commandLineStyle.Width(m.termWidth).Render(commandContent)
}

// Render the onscreen keyboard
/*

TODO: This renders the keyboard in a vertical orientation and the keys appear to be too big.
It should fill the width of the terminal and be more compact.
*/
func (m Model) renderKeyboard(maxHeight int) string {
	if len(m.keyboard.Keys) == 0 {
		return "No keyboard layout loaded"
	}

	var info string = ""

	// Key styles
	// TODO: get colors from keys
	normalKeyStyle := lipgloss.NewStyle().
		// Border(lipgloss.RoundedBorder()).
		// BorderForeground(lipgloss.Color("240")).
		// Background(lipgloss.Color("236")).
		// Foreground(lipgloss.Color("255")).
		// Align(lipgloss.Center).
		// Padding(0, 1)
		Padding(0)

	pressedKeyStyle := lipgloss.NewStyle().
		// Border(lipgloss.RoundedBorder()).
		// BorderForeground(lipgloss.Color("46")).
		// Background(lipgloss.Color("46")).
		// Foreground(lipgloss.Color("0")).
		// Align(lipgloss.Center).
		// Padding(0, 1)
		Padding(0)

	specialKeyStyle := lipgloss.NewStyle().
		// Border(lipgloss.RoundedBorder()).
		// BorderForeground(lipgloss.Color("33")).
		// Background(lipgloss.Color("237")).
		// Foreground(lipgloss.Color("33")).
		// Align(lipgloss.Center).
		// Padding(0, 1)
		Padding(0)

	// Group keys by row (Y coordinate)
	rows := m.getKeyboardRows()

	/*
		DEBUG CODE - REMOVE
	*/
	// {
	// 	// info += fmt.Sprintf("# rows: %d\n", len(rows))
	// 	for r := range rows {
	// 		// for r := len(rows); r > 0; r-- {
	// 		// info += fmt.Sprintf(" > row[%d]:\n", r)
	// 		for range rows[r] {
	// 			info += fmt.Sprint("┌   ┐")
	// 		}
	// 		info += fmt.Sprint("\n")
	// 		for k := range rows[r] {
	// 			currentRow := rows[r][k]
	// 			label := currentRow.Labels[0]
	// 			if len(label) > 1 {
	// 				label = string(label[0])
	// 			} else if len(label) == 0 {
	// 				label = " "
	// 			}
	// 			//info += fmt.Sprintf("| %v |", label)
	// 			log.Printf("%v is [%.2fu x %.2fu] at (%d, %d)",
	// 				strings.Join(currentRow.Labels, ""),
	// 				currentRow.Width, currentRow.Height,
	// 				currentRow.X, currentRow.Y)
	// 		}
	// 		info += fmt.Sprint("\n")
	// 		for range rows[r] {
	// 			info += fmt.Sprint("└   ┘")
	// 		}
	// 		info += fmt.Sprint("\n")
	// 	}
	// }
	/*
		END DEBUG CODE
	*/

	var keyboardLines []string

	// Calculate available height for keyboard content
	// TODO: ??claude?? Account for border (2 lines), padding (2 lines), margin (2 lines), and info line (3 lines)
	//       maxHeight should already be 2/3 of term rows
	headerHeight := 0
	availableHeight := maxHeight - headerHeight
	// keyboardWidth := getKeyboardWidth(rows)
	// how many columns wide a 1u key can be
	// one unit should be 4 columns wide - this allows for dividing into fourths for certain key sizes

	keyUnitWidth := 4
	//maxUnitWidth := math.Round(float64(m.termWidth) / (keyboardWidth * float64(keyUnitWidth)))

	// Render each row, respecting height constraints
	renderedRows := 0
	// info += fmt.Sprintf("Keyboard layout (%d rows):\n", maxY+1)
	// info += fmt.Sprintf("maxHeight: %d\n", maxHeight)
	// info += fmt.Sprintf("headerHeight: %d\n", headerHeight)
	// info += fmt.Sprintf("availableHeight: %d\n", availableHeight)
	for y := 0; y <= len(rows) && renderedRows < availableHeight; y++ {
		// info += fmt.Sprintf("renderedRows: %d\n", renderedRows)
		if rowKeys, exists := rows[y]; exists {
			var row []string

			// info += fmt.Sprintf("\nRow %d: ", y)
			for _, key := range rowKeys {
				// // Determine key label
				// label := ""
				// if len(key.Labels) >= 5 {
				// 	label = key.Labels[4]
				// } else {
				// 	label = key.Labels[0]
				// }

				// Keys have up to 12 labels, in 3 columns and 3 rows, plus a "front face" row
				label1 := lipgloss.JoinHorizontal(lipgloss.Top, key.Labels[:3]...)
				label2 := lipgloss.JoinHorizontal(lipgloss.Top, key.Labels[3:6]...)
				label3 := lipgloss.JoinHorizontal(lipgloss.Top, key.Labels[6:9]...)
				// TODO: for now, ignoring front labels
				// label4 := lipgloss.JoinHorizontal(lipgloss.Top, key.Labels[9:]...)


				// Determine key width
				// width := max(3, int(key.Width)) //min(max(int(key.Width*4), 3), 12)
				// width := min(max(int(key.Width*4), 3), 12)

				// convert key unit width to actual terminal column width
				keyWidth := max(keyUnitWidth, int(key.Width*float64(keyUnitWidth)))
				log.Printf("%f key.Width -> %d term cols", key.Width, keyWidth)

				// Shorten labels that are too long
				if keyWidth < len(label1) {
					label1 = label1[:keyWidth]
				}
				if keyWidth < len(label2) {
					label2 = label2[:keyWidth]
				}
				if keyWidth < len(label1) {
					label3 = label3[:keyWidth]
				}

				// Join the label rows to form the keycap
				label := lipgloss.JoinVertical(lipgloss.Left, label1, label2, label3) //, label4)

				// Check if key is pressed
				// TODO: this now needs to check what the actual key is vs. what the label(s) might contain...
				keyPressed := m.pressedKeys[strings.ToUpper(label)]
				// || (label == "spaaaaaaaaaaaaaaaaaaaaaaaaaace" && m.pressedKeys["SPACE"])

				// Choose style
				var style lipgloss.Style
				if keyPressed {
					style = pressedKeyStyle
				} else if isSpecialKey(label) {
					style = specialKeyStyle
				} else {
					style = normalKeyStyle
				}

				// Apply custom colors if specified
				if key.Color != "" {
					style = style.Background(lipgloss.Color(key.Color))
				}
				if key.TextColor != "" {
					style = style.Foreground(lipgloss.Color(key.TextColor))
				}
/*
sizes...
... just for the keyboard

// Normal border
8*15=120 wide
5*5=25 tall
┌──────┐
│~     │
│      │
│`     │
└──────┘

// Rounded border
// same as normal border
╭──────╮
│~     │
│      │
│`     │
╰──────╯

// no border:
6x15=90 wide
3x5=15 tall
~     
      
`     


and then...

// 8 + 2 for the border = 10 cols for a 2u key
// but... it's 4 tall, should only be 3 tall
╭────────╮
│Backspac│
│e       │
│        │
│        │
╰────────╯

// 1u key is 4 wide + 2 border = 6 and 3+2=5 tall
// we could stand to pad it and make it 6 wide + 2 border = 8 columns but then it's 120 wide minimum
╭────╮
│~   │
│    │
│`   │
╰────╯
*/

				width := int(keyWidth)
				renderedKey := style.
					Border(lipgloss.RoundedBorder()).
					Width(width).
					Height(3).
					Render(label)
				// rowDisplay.WriteString(renderedKey)
				// rowDisplay.WriteString(" ")
				// row := lipgloss.JoinHorizontal(lipgloss.Top, renderedKey)
				log.Printf("\n%s\n", renderedKey)
				row = append(row, renderedKey)

			}

			rowDisplay := lipgloss.JoinHorizontal(lipgloss.Top, row...)

			// info += fmt.Sprintf("\nrowDisplay %v: ", rowDisplay.String())
			//keyboardLines = append(keyboardLines, rowDisplay.String())
			keyboardLines = append(keyboardLines, rowDisplay)
			renderedRows++
		}
		// TODO: Else for empty rows?
	}

	// Add truncation indicator if we couldn't fit all rows
	// if renderedRows < maxY+1 {
	// 	keyboardLines = append(keyboardLines, "... (keyboard truncated to fit)")
	// }

	// Join all rows
	// keyboard := strings.Join(keyboardLines, "\n")
	keyboard := lipgloss.JoinVertical(lipgloss.Center, keyboardLines...)

	// Add keyboard info
	// TODO: these fields could be empty
	// info += fmt.Sprintf("keyboard lines: %d\n", len(keyboardLines))
	// info += fmt.Sprintf("rows: %d\n", len(rows))
	// info += fmt.Sprintf("rendered rows: %d\n", renderedRows)
	if m.keyboard.Meta.Name != "" {
		info += fmt.Sprintf("Keyboard: %s", m.keyboard.Meta.Name)
	}
	if m.keyboard.Meta.Author != "" {
		info += fmt.Sprintf(" by %s", m.keyboard.Meta.Author)
	}

	// Apply height constraint to the final rendered output
	//finalContent := info + "\n\n" + keyboard
	styledKeyboard := lipgloss.NewStyle().
		// Border(lipgloss.RoundedBorder()).
		// BorderForeground(lipgloss.Color("62")).
		// Padding(1).
		// Margin(1).
		//Height(maxHeight).
		Render(keyboard)
		//Render(finalContent)

	/*
		DEBUG CODE - REMOVE
	*/
	// styledKeyboard = ""
	/*
		END DEBUG CODE
	*/

	// return styledKeyboard
	ui := lipgloss.JoinVertical(lipgloss.Left, info, styledKeyboard)
	return ui //m.centerContent(ui)
}

func (m Model) renderPrompt(maxHeight int) string {
	// Style definitions
	promptStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Margin(1)

	correctStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))                                    // Green
	incorrectStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))                                 // Red
	currentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Background(lipgloss.Color("240")) // Yellow bg
	futureStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))                                    // Gray

	// Build prompt display
	var promptDisplay strings.Builder
	for i, char := range m.prompt {
		if i < len(m.userInput) {
			// Character has been typed
			if i < len(m.userInput) && rune(m.userInput[i]) == char {
				promptDisplay.WriteString(correctStyle.Render(string(char)))
			} else {
				promptDisplay.WriteString(incorrectStyle.Render(string(char)))
			}
		} else if i == m.currentChar {
			// Current character to type
			promptDisplay.WriteString(currentStyle.Render(string(char)))
		} else {
			// Future characters
			promptDisplay.WriteString(futureStyle.Render(string(char)))
		}
	}

	// Progress info
	progress := fmt.Sprintf("Progress: %d/%d characters | Prompt %d/%d",
		len(m.userInput), len(m.prompt), m.promptIndex+1, len(m.prompts))

	instructions := "Tab: Next prompt | Ctrl+C/Q: Quit"

	return promptStyle.Render(
		fmt.Sprintf("Type: %s\n\n%s\n\n%s\n%s",
			m.prompt,
			promptDisplay.String(),
			progress,
			instructions))
		//Height(maxHeight).

}

func (m Model) getKeyboardRows() map[int][]Key {
	rows := make(map[int][]Key)
	maxY := 0
	for _, key := range m.keyboard.Keys {
		y := key.Y
		rows[y] = append(rows[y], key)
		if y > maxY {
			maxY = y
		}
	}
	return rows
}

// Find the total u-width of the longest row
func getKeyboardWidth(kbRows map[int][]Key) float64 {
	maxWidth := 0.0
	for row := range(kbRows) {
		rowWidth := 0.0
		for key := range(row) {
			rowWidth += kbRows[row][key].Width
		}
		if rowWidth > maxWidth {
			maxWidth = rowWidth
		}
	}
	return maxWidth
}

func isSpecialKey(label string) bool {
	// TODO: why is this checking specific keys loaded from JSON which likely won't exist?
	specialKeys := []string{
		"Tab", "Caps Lock", "Shift", "Enter", "Backspace", "Space",
		"CMD", "Alt", "FN", "win", "menü", "spaaaaaaaaaaaaaaaaaaaaaaaaaace",
	}

	for _, special := range specialKeys {
		if strings.EqualFold(label, special) || strings.Contains(strings.ToLower(label), strings.ToLower(special)) {
			return true
		}
	}
	return false
}
