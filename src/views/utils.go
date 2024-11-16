package views

import (
	"github.com/VU-ASE/rover/src/components"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UtilitiesPage struct {
	// To select an action to perform with this utility
	actions list.Model // actions you can perform when connected to a Rover
	help    help.Model // to display a help footer
	spinner spinner.Model
}

func NewUtilitiesPage() UtilitiesPage {
	l := list.New([]list.Item{
		components.ActionItem{Name: "SSH", Desc: "Open an SSH terminal to your Rover"},
		components.ActionItem{Name: "Info", Desc: "Info about roverctl on your local system"},
	}, list.NewDefaultDelegate(), 0, 0)
	// If there are connections available, add the connected actions
	l.Title = lipgloss.NewStyle().Foreground(style.AsePrimary).Padding(0, 0).Render("Utilities")
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = style.TitleStyle
	l.Styles.PaginationStyle = style.PaginationStyle
	l.Styles.HelpStyle = style.HelpStyle

	sp := spinner.New()

	return UtilitiesPage{
		actions: l,
		help:    help.New(),
		spinner: sp,
	}
}

func (m UtilitiesPage) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m UtilitiesPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				case "Info":
					return RootScreen(state.Get()).SwitchScreen(NewInfoPage())
				}
				return m, tea.Quit
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	}

	var cmd tea.Cmd
	m.actions, cmd = m.actions.Update(msg)
	return m, cmd
}

func (m UtilitiesPage) View() string {
	return m.actions.View()
}
