package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (a AppState) ConnectedHomeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		a.connectedActions.SetSize(msg.Width-h, msg.Height-v) // leave some room for the header

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		case "enter":
			value := a.connectedActions.SelectedItem().FilterValue()
			if value != "" {
				a.selectedAction = value
			}
			return a, nil
		}
	}

	var cmd tea.Cmd
	a.connectedActions, cmd = a.connectedActions.Update(msg)

	return a, cmd
}

func (a AppState) ConnectedHomeView() string {
	// Send the UI for rendering
	return a.connectedActions.View()
}
