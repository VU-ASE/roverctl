package startpageconnected

import (
	"github.com/VU-ASE/rover/src/components"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	// To select an action to perform with this utility
	actions list.Model // actions you can perform when connected to a Rover
	help    help.Model // to display a help footer
}

func InitialModel() model {
	l := list.New([]list.Item{
		components.ActionItem{Name: "Configure", Desc: "Configure your pipeline"},
		components.ActionItem{Name: "Services", Desc: "Create, upload, download and install services"},
		components.ActionItem{Name: "Connections", Desc: "Manage your Rover connections"},
		components.ActionItem{Name: "Utilities", Desc: "Various utilities to interact with your Rover"},
		// components.ActionItem{Name: "Debug", Desc: "Enable remote debugging for your pipeline"}, // Should not be available when no pipeline is running or disable when enabled
		// components.ActionItem{Name: "Status", Desc: "Watcdh module outputs and status logs"},    // Should not be available when no pipeline is running
		// components.ActionItem{Name: "Update", Desc: "Fetch the latest versions of all modules and install them"},
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

		case "enter":
			value := m.actions.SelectedItem().FilterValue()
			if value != "" {
				if value == "Configure" {
					state.Get().Route.Push("pipeline configure")
				} else {
					state.Get().Route.Push(value)
				}
				return m, tea.Quit
			}
		}
	}

	// Is there a main action to take?
	rootmodel, rootcmd := tui.Update(m, msg)
	if rootcmd != nil {
		return rootmodel, rootcmd
	}

	var cmd tea.Cmd
	m.actions, cmd = m.actions.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return style.Docstyle.Render(m.actions.View())
}