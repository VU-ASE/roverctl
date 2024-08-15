package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (a AppState) DisconnectedHomeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		a.disconnectedActions.SetSize(msg.Width-h, msg.Height-v) // leave some room for the header

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		case "enter":
			value := a.disconnectedActions.SelectedItem().FilterValue()
			if value != "" {
				a.selectedAction = value
			}
			return a, nil
		}
	}

	var cmd tea.Cmd
	a.disconnectedActions, cmd = a.disconnectedActions.Update(msg)

	return a, cmd
}

func (a AppState) DisconnectedHomeView() string {
	// Send the UI for rendering
	return a.disconnectedActions.View()
}
