package views

import (
	"fmt"

	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// This is the main model that will be used to render all sub-models
type MainModel struct {
	current tea.Model
}

func RootScreen(s *state.AppState) MainModel {
	start := NewStartPage()

	return MainModel{
		current: &start, // needs to be a pointer so that the model state can be modified (see https://shi.foo/weblog/multi-view-interfaces-in-bubble-tea)
	}
}

func (m MainModel) Init() tea.Cmd {
	return m.current.Init()
}

func (m MainModel) View() string {
	// Define the header style
	headerStyle := lipgloss.NewStyle().
		Width(state.Get().WindowWidth). // Set the width of the header to the window width
		Align(lipgloss.Center).         // Center-align the text
		Background(style.AsePrimary)    // Set the background color

	// Define the URL and the text
	url := "https://ase.vu.nl"
	text := "read the docs"

	// Hyperlink escape sequence
	link := fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, text)

	header := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Padding(0, 0).Render("VU ASE") + lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(style.AsePrimary).Bold(false).Padding(0, 0).Render(", "+state.Get().Quote+" | "+link)
	fullScreen := lipgloss.NewStyle().Padding(1, 2).Width(state.Get().WindowWidth).Height(state.Get().WindowHeight - 1) // leave room for the header

	return fullScreen.Render(m.current.View()) + "\n" + headerStyle.Render(header)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global messages first
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Store the window dimensions
		state.Get().WindowWidth = msg.Width
		state.Get().WindowHeight = msg.Height

		passedMsg := tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - 2, // leave room for the header
		}

		// Forward the message to the current sub-model
		updatedModel, cmd := m.current.Update(passedMsg)
		m.current = updatedModel
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":

			// Return based on the current route
			// todo ...

			return m, tea.Quit

		case "ctrl+c":
			return m, tea.Quit
		}
	}

	// Delegate other messages to the current sub-model
	updatedModel, cmd := m.current.Update(msg)
	m.current = updatedModel
	return m, cmd
}

// This function is used to switch between screens, the caller should supply the "route" taken so far to get to this screen, so that users can return to the previous screen
func (m MainModel) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	m.current = model

	// Notify the new model of the current window size
	windowSizeMsg := tea.WindowSizeMsg{
		Width:  state.Get().WindowWidth,
		Height: state.Get().WindowHeight,
	}

	// Initialize the new model and send the size message
	initCmd := m.current.Init()
	sizeCmd := func() tea.Cmd {
		return func() tea.Msg {
			return windowSizeMsg
		}
	}

	return m.current, tea.Batch(initCmd, sizeCmd())
}