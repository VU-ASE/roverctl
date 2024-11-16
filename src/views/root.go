package views

import (
	"github.com/VU-ASE/rover/src/components"
	startpageconnected "github.com/VU-ASE/rover/src/pages/start/connected"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// This is the main model that will be used to render all sub-models
type MainModel struct {
	current tea.Model
}

func RootScreen(s *state.AppState) MainModel {
	start := InitialModel()

	return MainModel{
		// current: startpagedisconnected.InitialModel(),
		current: &start, // needs to be a pointer so that the model state can be modified (see https://shi.foo/weblog/multi-view-interfaces-in-bubble-tea)
	}
}

func (m MainModel) Init() tea.Cmd {
	return m.current.Init()
}

func (m MainModel) View() string {
	return m.current.View()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global messages first
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Store the window dimensions
		state.Get().WindowWidth = msg.Width
		state.Get().WindowHeight = msg.Height

		// Forward the message to the current sub-model
		updatedModel, cmd := m.current.Update(msg)
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

type model struct {
	// To select an action to perform with this utility
	actions list.Model // actions you can perform when connected to a Rover
	help    help.Model // to display a help footer
}

func InitialModel() model {
	l := list.New([]list.Item{
		components.ActionItem{Name: "ServiceEEEE", Desc: "Create services"},
		components.ActionItem{Name: "Connect", Desc: "Initialize a connection to a Rover"}, // Should be "stop" when a pipeline is running
	}, list.NewDefaultDelegate(), 0, 0)
	// If there are connections available, add the connected actions
	l.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(style.AsePrimary).Bold(true).Padding(0, 0).Render("VU ASE") + lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C")).Render(" - racing Rovers since 2024")
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = style.TitleStyle
	l.Styles.PaginationStyle = style.PaginationStyle
	l.Styles.HelpStyle = style.HelpStyle

	return model{
		actions: l,
		help:    help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := style.Docstyle.GetFrameSize()
		m.actions.SetSize(msg.Width-h, msg.Height-v) // leave some room for the header

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		case "e":
			connected := startpageconnected.InitialModel()

			return RootScreen(state.Get()).SwitchScreen(&connected)
		case "enter":
			value := m.actions.SelectedItem().FilterValue()
			if value != "" {
				switch value {
				case "Connect":
					// state.Get().Route.Push("connection init")
				default:
					// state.Get().Route.Push(value)
				}
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.actions, cmd = m.actions.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return style.Docstyle.Render(m.actions.View())
}

func (m model) New() tea.Model {
	return InitialModel()
}
