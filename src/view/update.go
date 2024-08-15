package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

// The update view (where you can download the latest version of all modules) with the update method
func (a AppState) UpdateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return a, tea.Quit
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return a, nil
}

// the update view with the view method
func (a AppState) UpdateView() string {
	// The header
	s := "What should we buy at the market? This is the update view\n\n"

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
