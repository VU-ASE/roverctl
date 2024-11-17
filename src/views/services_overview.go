package views

import (
	"os"

	"github.com/VU-ASE/rover/src/components"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ServicesOverviewPage struct {
	// To select an action to perform with this utility
	actions list.Model // actions you can perform when connected to a Rover
	help    help.Model // to display a help footer
}

func NewServicesOverviewPage() ServicesOverviewPage {
	// Is there already a service.yaml file in the current directory?
	_, err := os.Stat("./service.yaml")

	listItems := []list.Item{}
	if err != nil {
		listItems = append(listItems, components.ActionItem{Name: "Initialize", Desc: "Initialize a new service in your current working directory"})
	} else {
		listItems = append(listItems, components.ActionItem{Name: "Sync", Desc: "Synchronize your local service with the Rover by watching for changes"})
	}
	listItems = append(listItems, []list.Item{
		components.ActionItem{Name: "Update", Desc: "Update official services from source onto your Rover"},
	}...)

	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	// If there are connections available, add the connected actions
	l.Title = lipgloss.NewStyle().Foreground(style.AsePrimary).Padding(0, 0).Render("Manage your services")
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = style.TitleStyle
	l.Styles.PaginationStyle = style.PaginationStyle
	l.Styles.HelpStyle = style.HelpStyle

	return ServicesOverviewPage{
		actions: l,
		help:    help.New(),
	}
}

func (m ServicesOverviewPage) Init() tea.Cmd {
	return nil
}

func (m ServicesOverviewPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := style.Docstyle.GetFrameSize()
		m.actions.SetSize(msg.Width-h, msg.Height-v) // leave some room for the header

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		case "enter":
			value := m.actions.SelectedItem().FilterValue()
			if value != "" {
				switch value {
				case "Initialize":
					return RootScreen(state.Get()).SwitchScreen(NewServiceInitPage())
				case "Sync":
					return RootScreen(state.Get()).SwitchScreen(NewServicesSyncPage())
				case "Update":
					return RootScreen(state.Get()).SwitchScreen(NewServicesUpdatePage())
				case "Download":
					value = "service download"
				}
				// state.Get().Route.Push(value)
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.actions, cmd = m.actions.Update(msg)
	return m, cmd
}

func (m ServicesOverviewPage) View() string {
	return m.actions.View()
}
