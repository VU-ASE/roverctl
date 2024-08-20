package tui

import (
	"github.com/VU-ASE/rover/src/state"
	tea "github.com/charmbracelet/bubbletea"
)

// Reusable tui components
func Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			// Go to the previous page
			state.Get().Route.Pop()
			return m, tea.Quit
		case "ctrl+c":
			// Exit the application immeediately
			state.Get().Route.Clear()
			return m, tea.Quit
		}
	}

	return m, nil
}
