package views

import (
	"github.com/VU-ASE/rover/src/components"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type StartPage struct {
	// To select an action to perform with this utility
	actions list.Model // actions you can perform when connected to a Rover
	help    help.Model // to display a help footer
}

func NewStartPage() StartPage {
	d := style.DefaultListDelegate()
	l := list.New([]list.Item{
		components.ActionItem{Name: "Services", Desc: "Create services"},
		components.ActionItem{Name: "Connect", Desc: "Initialize a connection to a Rover"}, // Should be "stop" when a pipeline is running
	}, d, 0, 0)
	if len(state.Get().RoverConnections.Available) > 0 {
		l = list.New([]list.Item{
			components.ActionItem{Name: "Pipeline", Desc: "Manage your pipeline"},
			components.ActionItem{Name: "Services", Desc: "Create, upload, download and install services"},
			components.ActionItem{Name: "Connections", Desc: "Manage connections to roverd instances"},
			components.ActionItem{Name: "Utilities", Desc: "Various utilities to interact with your Rover"},
		}, d, 0, 0)
	}

	// If there are connections available, add the connected actions
	l.Title = "Roverctl"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = style.TitleStyle
	l.Styles.PaginationStyle = style.PaginationStyle
	l.Styles.HelpStyle = style.HelpStyle

	return StartPage{
		actions: l,
		help:    help.New(),
	}
}

func (m StartPage) Init() tea.Cmd {
	return nil
}

func (m StartPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := style.Docstyle.GetFrameSize()
		m.actions.SetSize(msg.Width-h, msg.Height-v) // leave some room for the header

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// case "e":
		// 	connected := startpageconnected.InitialStartPage()

		// 	return RootScreen(state.Get()).SwitchScreen(&connected)
		case "enter":
			value := m.actions.SelectedItem().FilterValue()
			if value != "" {
				switch value {
				case "Pipeline":
					return RootScreen(state.Get()).SwitchScreen(NewPipelineOverviewPage())
				case "Connect":
					return RootScreen(state.Get()).SwitchScreen(NewConnectionsInitPage(nil))
				case "Utilities":
					return RootScreen(state.Get()).SwitchScreen(NewUtilitiesPage())
				case "Connections":
					return RootScreen(state.Get()).SwitchScreen(NewConnectionsManagePage())
				case "Services":
					return RootScreen(state.Get()).SwitchScreen(NewServicesOverviewPage())
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

func (m StartPage) View() string {
	return m.actions.View()
}
