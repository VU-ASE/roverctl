package initconnectionpage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/VU-ASE/rover/src/configuration"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	probing "github.com/prometheus-community/pro-bing"
)

type formValues struct {
	name     string
	index    string
	username string
	password string
}

// keyMap defines a set of keybindings. To work for help it must satisfy key.Map
type keyMap struct {
	Save  key.Binding
	New   key.Binding
	Back  key.Binding
	Retry key.Binding
	Quit  key.Binding
}

type model struct {
	form          *huh.Form
	help          help.Model
	spinner       spinner.Model
	routeExists   tui.Action[bool]
	authValid     tui.Action[bool]
	roverdVersion tui.Action[string]
	roverNumber   tui.Action[int32]
	isChecking    bool
	formValues    *formValues
	host          string // the ip or hostname of the rover to connect to
	error         error  // any errors that occurred
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Retry, k.Quit, k.Save, k.New}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var keys = keyMap{
	Back: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "back"),
	),
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

// Shown when the form is completed and the connection checks are successful
var successKeys = keyMap{
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "save and new"),
	),
	Back: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "back"),
	),
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

var failureKeys = keyMap{
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save anyway"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "save and new"),
	),
	Back: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "back"),
	),
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func InitialModel(val *formValues) model {
	s := spinner.New()
	s.Spinner = spinner.Line

	formValues := &formValues{
		name:     "",
		index:    "",
		username: "debix",
		password: "debix",
	}
	if val != nil {
		formValues = val
	}

	routeExistsAction := tui.NewAction[bool]("routeExists")
	authValidAction := tui.NewAction[bool]("authValid")
	roverdVersionAction := tui.NewAction[string]("roverdVersion")
	roverNumberAction := tui.NewAction[int32]("roverNumber")

	return model{
		spinner:       s,
		formValues:    formValues,
		host:          "",
		help:          help.New(),
		routeExists:   routeExistsAction,
		authValid:     authValidAction,
		roverdVersion: roverdVersionAction,
		roverNumber:   roverNumberAction,
		isChecking:    false,
		error:         nil,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Enter your Rover index (1-20, inclusive)").
					CharLimit(3).
					Prompt("> ").
					Value(&formValues.index).
					Validate(func(s string) error {
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
					Value((&formValues.username)),
				huh.NewInput().
					Title("Enter the authentication password").
					CharLimit(255).
					Prompt("> ").
					EchoMode(huh.EchoModePassword).
					Value((&formValues.password)),
				huh.NewInput().
					Title("Enter a name for this connection to find it back later").
					CharLimit(255).
					Prompt("> ").
					Value(&formValues.name).Validate(func(s string) error {
					if len(s) <= 0 {
						return fmt.Errorf("You cannot leave this field empty")
					}
					return nil
				}),
			),
		).WithTheme(style.FormTheme),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.form.State == huh.StateCompleted {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			// If we set a width on the help menu it can gracefully truncate
			// its view as needed.
			m.help.Width = msg.Width
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.Back):
				// Restore to the initial form, but recover the form values
				m = InitialModel(m.formValues)
				return m, tea.Batch(m.form.Init(), m.spinner.Tick)
			case key.Matches(msg, keys.Retry):
				// Retry the connection checks
				m.isChecking = true
				return m, tea.Batch(checkRoute(m), checkAuth(m), checkRoverdVersion(m), checkRoverNumber(m))
			}

			switch msg.String() {
			case "n", "enter":
				// Save the connection if all checks are successful
				if m.routeExists.IsDone() && m.authValid.IsDone() {
					// Save the connection
					state.Get().RoverConnections = state.Get().RoverConnections.Add(configuration.RoverConnection{
						Name:     m.formValues.name,
						Host:     m.host,
						Username: m.formValues.username,
						Password: m.formValues.password,
					})

					if msg.String() == "n" {
						m = InitialModel(nil)
						return m, tea.Batch(m.form.Init(), m.spinner.Tick)
					} else {
						state.Get().Route.Replace("connections")
						return m, tea.Quit
					}
				}
			}
		}
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tui.ActionInit[bool]:
		m.routeExists.ProcessInit(msg)
		m.authValid.ProcessInit(msg)
		return m, nil
	case tui.ActionInit[string]:
		m.roverdVersion.ProcessInit(msg)
		return m, nil
	case tui.ActionInit[int32]:
		m.roverNumber.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[bool]:
		m.authValid.ProcessResult(msg)
		m.routeExists.ProcessResult(msg)
		return m, nil
	case tui.ActionResult[string]:
		m.roverdVersion.ProcessResult(msg)
		return m, nil
	case tui.ActionResult[int32]:
		m.roverNumber.ProcessResult(msg)
		return m, nil
	default:
		cmds := []tea.Cmd{}
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			if f.State == huh.StateCompleted && !m.isChecking {
				m.isChecking = true

				index, err := strconv.Atoi(m.formValues.index)
				if err != nil || index < 1 || index > 20 {
					m.routeExists = tui.NewAction[bool]("routeExists")
					return m, cmd
				}
				// todo: change to 192.168.1 instead of 192.168.0
				// m.host = fmt.Sprintf("192.168.0.%d", index+100)
				m.host = "google.com"

				// We are optimistic, start all checks in parallel
				cmds = append(cmds, checkRoute(m), checkAuth(m), checkRoverdVersion(m), checkRoverNumber(m))
			}
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		} else {
			// Base command (put in this ugly nested else statement because we don't want to quit when a user is typing in a 'q')
			model, cmd := tui.Update(m, msg)
			if cmd != nil {
				return model, cmd
			}
		}
		return m, tea.Batch(cmds...)
	}
}

// the update view with the view method
func (m model) enterDetailsView() string {
	// Introduction
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Connect to a Rover")

	// // Get the current wifi name
	// wifi := wifiname.WifiName()

	// if wifi == "Could not get SSID" {
	// 	wifi = "unknown network"
	// }

	// if wifi != "aselabs" {
	// 	s += lipgloss.NewStyle().Foreground(style.WarningPrimary).Render("\n\nIt seems you are not connected to the ASElabs WiFi but to '" + wifi + "' instead. \nRead how to connect at: https://docs.ase.vu.nl/docs/tutorials/setting-up-your-workspace/accessing-the-network")
	// }

	s += "\n\n" + m.form.View()

	return style.Docstyle.Render(s)
}

func (m model) testConnectionView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Connecting to " + m.formValues.name)

	if m.routeExists.IsLoading() || m.authValid.IsLoading() || m.roverdVersion.IsLoading() || m.roverNumber.IsLoading() {
		s += "\n\n " + m.spinner.View() + " Performing connection checks..."

		s += "\n\n" + m.help.View(keys)
		return s
	}

	if !m.routeExists.IsSuccess() {
		s += "\n\n ✗ " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("No route could be established to the Rover. Are you sure it is powered on? (Tried "+m.host+")\n   Read how to connect at: https://docs.ase.vu.nl/docs/tutorials/setting-up-your-workspace/accessing-the-network")
		if m.routeExists.Error != nil {
			s += "\n   (" + m.routeExists.Error.Error() + ")"
		}
	} else {
		s += "\n\n ✓ " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Established route to Rover at "+m.host)
	}

	if m.routeExists.IsSuccess() {
		if !m.roverdVersion.IsSuccess() {
			s += "\n ✗ " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Could not determine roverd version")
			if m.roverdVersion.Error != nil {
				s += " (" + m.roverdVersion.Error.Error() + ")"
			}
		} else {
			s += "\n ✓ " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Found roverd version: "+*m.roverdVersion.Data)
		}

		index, _ := strconv.Atoi(m.formValues.index)
		if !m.roverNumber.IsSuccess() {
			s += "\n ✗ " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Could not determine rover number")
			if m.roverNumber.Error != nil {
				s += " (" + m.roverNumber.Error.Error() + ")"
			}
		} else if *m.roverNumber.Data != int32(index) {
			s += "\n ! " + lipgloss.NewStyle().Foreground(style.WarningPrimary).Render("This Rover presented itself as Rover "+strconv.Itoa(int(*m.roverNumber.Data))+" but you wanted to connect to Rover "+m.formValues.index)
		} else {
			s += "\n ✓ " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Rover number matches the index you entered ("+m.formValues.index+")")
		}

		if !m.authValid.IsSuccess() {
			s += "\n ✗ " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Authentication to the roverd endpoint failed. Please check your credentials")
			if m.authValid.Error != nil {
				s += " (" + m.authValid.Error.Error() + ")"
			}
		} else {
			s += "\n ✓ " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Authentication successful")
		}
	}

	if !m.routeExists.IsSuccess() || !m.authValid.IsSuccess() || !m.roverdVersion.IsSuccess() {
		s += "\n\n" + lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("This connection configuration is not valid and should not be saved.")
		s += "\n\n" + m.help.View(failureKeys)
	} else {
		s += "\n\n" + "You are all set!"
		s += "\n\n" + m.help.View(successKeys)
	}

	return s
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.form.Init(), m.spinner.Tick)
}

func (m model) View() string {
	if m.form.State == huh.StateCompleted {
		return style.Docstyle.Render(m.testConnectionView())
	} else {
		return m.enterDetailsView()
	}
}

func checkRoute(m model) tea.Cmd {
	return tui.PerformAction(&m.routeExists, func() (*bool, error) {
		ping, _ := probing.NewPinger(m.host)
		ping.Count = 3
		ping.Timeout = 10 * time.Second
		err := ping.Run()

		valid := ping.Statistics().PacketsRecv > 0
		if !valid {
			err = fmt.Errorf("No route to host")
		}
		return &valid, err
	})
}

func checkAuth(m model) tea.Cmd {
	return tui.PerformAction(&m.authValid, func() (*bool, error) {
		// Send a protected request to the roverd endpoint
		c := configuration.RoverConnection{
			Host:     m.host,
			Username: m.formValues.username,
			Password: m.formValues.password,
		}
		a := c.ToApiClient()

		_, _, err := a.PipelineAPI.PipelineGet(context.Background()).Execute()
		res := err == nil
		return &res, err
	})
}

func checkRoverdVersion(m model) tea.Cmd {
	return tui.PerformAction(&m.roverdVersion, func() (*string, error) {
		c := configuration.RoverConnection{
			Host:     m.host,
			Username: m.formValues.username,
			Password: m.formValues.password,
		}
		a := c.ToApiClient()

		res, _, err := a.HealthAPI.StatusGet(context.Background()).Execute()
		if err != nil {
			return nil, err
		}

		return res.Version, nil
	})
}

func checkRoverNumber(m model) tea.Cmd {
	return tui.PerformAction(&m.roverNumber, func() (*int32, error) {
		c := configuration.RoverConnection{
			Host:     m.host,
			Username: m.formValues.username,
			Password: m.formValues.password,
		}
		a := c.ToApiClient()

		res, _, err := a.HealthAPI.StatusGet(context.Background()).Execute()
		if err != nil {
			return nil, err
		}

		num := res.RoverId
		return num, err
	})
}
