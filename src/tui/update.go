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
		// These keys should exit the page.
		case "ctrl+c", "esc", "q":
			if state.Get().CurrentView == "home" {
				state.Get().CurrentView = ""
			} else if state.Get().CurrentView != "" {
				state.Get().CurrentView = "home"
			}
			return m, tea.Quit
		}
	}

	return m, nil
}
