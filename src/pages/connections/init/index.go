package initconnectionpage

import (
	"fmt"
	"strconv"

	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	wifiname "github.com/yelinaung/wifi-name"
)

type model struct {
	form    *huh.Form
	spinner spinner.Model
}

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Line

	return model{
		spinner: s,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Enter your Rover index (1-20, inclusive)").
					CharLimit(3).
					Prompt("> ").
					Key("roverIndex").Validate(func(s string) error {
					index, err := strconv.Atoi(s)
					if err != nil || index < 1 || index > 20 {
						return fmt.Errorf("Please enter a valid Rover index between 1 and 20 (inclusive)")
					}
					return nil
				}),
				huh.NewInput().
					Title("Enter the authentication username").
					CharLimit(255).
					Prompt("> ").
					Key("username"),
				huh.NewInput().
					Title("Enter the authentication password").
					CharLimit(255).
					Prompt("> ").
					Key("password").
					EchoMode(huh.EchoModePassword),
				huh.NewInput().
					Title("Enter a name for this connection to find it back later").
					CharLimit(255).
					Prompt("> ").
					Key("name"),
			),
		).WithTheme(style.FormTheme),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Base command
	model, command := tui.Update(m, msg)
	if command != nil {
		return model, command
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	if cmd != nil {
		return m, cmd
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// the update view with the view method
func (m model) enterDetailsView() string {

	// Introduction
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Connect to a Rover") + "\n\n - Make sure you are connected to the ASElabs WiFi\n - Make sure that the Rover is powered on"

	// Get the current wifi name
	wifi := wifiname.WifiName()

	if wifi != "aselabs" {
		s += lipgloss.NewStyle().Bold(true).Foreground(style.WarningPrimary).Render("\n\nIt seems you are not connected to the ASElabs WiFi but to '" + wifi + "' instead. \nRead how to connect at: https://docs.ase.vu.nl/docs/tutorials/setting-up-your-workspace/accessing-the-network")
	}

	s += "\n\n" + m.form.View()

	return style.Docstyle.Render(s)
}

func (m model) testConnectionView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Connecting to " + m.form.GetString("name"))

	s += "\n\n" + m.spinner.View() + " checking if a route to Rover exists"
	s += "\n" + m.spinner.View() + " checking if a route to Rover exists"

	return style.Docstyle.Render(s)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.form.Init(), m.spinner.Tick)
}

func (m model) View() string {
	if m.form.State == huh.StateCompleted {
		return m.testConnectionView()
	} else {
		return m.enterDetailsView()
	}
}
