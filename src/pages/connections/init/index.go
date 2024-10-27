package initconnectionpage

import (
	"fmt"
	"strconv"
	"time"

	"github.com/VU-ASE/rover/src/configuration"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	probing "github.com/prometheus-community/pro-bing"
	wifiname "github.com/yelinaung/wifi-name"
)

type formValues struct {
	name     string
	index    string
	username string
	password string
}

type model struct {
	form          *huh.Form
	spinner       spinner.Model
	routeExists   tui.Action[bool]
	authValid     tui.Action[bool]
	roverdVersion tui.Action[string]
	roverNumber   tui.Action[int]
	isChecking    bool
	formValues    *formValues
	host          string // the ip or hostname of the rover to connect to
	error         error  // any errors that occurred
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
	roverNumberAction := tui.NewAction[int]("roverNumber")

	return model{
		spinner:       s,
		formValues:    formValues,
		host:          "",
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
		case tea.KeyMsg:
			switch msg.String() {
			case "n", "enter":
				// Save the connection if all checks are successful
				if m.routeExists.IsSuccess() && m.authValid.IsSuccess() {
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
			case "b":
				// Restore to the initial form, but recover the form values
				m = InitialModel(m.formValues)
				return m, tea.Batch(m.form.Init(), m.spinner.Tick)
			case "r":
				// Retry the connection checks
				m.isChecking = true
				return m, tea.Batch(checkRoute(m), checkAuth(m), checkRoverdVersion(m), checkRoverNumber(m))
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
	case tui.ActionInit[int]:
		m.roverNumber.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[bool]:
		m.authValid.ProcessResult(msg)
		m.routeExists.ProcessResult(msg)
		return m, nil
	case tui.ActionResult[string]:
		m.roverdVersion.ProcessResult(msg)
		return m, nil
	case tui.ActionResult[int]:
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

	// Get the current wifi name
	wifi := wifiname.WifiName()

	if wifi == "Could not get SSID" {
		wifi = "unknown network"
	}

	if wifi != "aselabs" {
		s += lipgloss.NewStyle().Foreground(style.WarningPrimary).Render("\n\nIt seems you are not connected to the ASElabs WiFi but to '" + wifi + "' instead. \nRead how to connect at: https://docs.ase.vu.nl/docs/tutorials/setting-up-your-workspace/accessing-the-network")
	}

	s += "\n\n" + m.form.View()

	return style.Docstyle.Render(s)
}

func (m model) testConnectionView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Connecting to " + m.formValues.name)

	s += "\n\n" + lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("Press 'b' to to back, 'r' to retry the connection checks, or 'q' to quit")

	if m.routeExists.IsLoading() || m.authValid.IsLoading() || m.roverdVersion.IsLoading() || m.roverNumber.IsLoading() {
		s += "\n\n " + m.spinner.View() + " Performing connection checks..."
		return s
	}

	if !m.routeExists.IsSuccess() {
		s += "\n\n ✗ " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("No route could be established to the Rover. Are you sure it is powered on? (Tried "+m.host+")")
	} else {
		s += "\n\n ✓ " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Established route to Rover at "+m.host)
	}

	if m.routeExists.IsSuccess() {
		if !m.roverdVersion.IsSuccess() {
			s += "\n ✗ " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Could not determine roverd version")
		} else {
			s += "\n ✓ " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Found roverd version: "+*m.roverdVersion.Data)
		}

		index, _ := strconv.Atoi(m.formValues.index)
		if !m.roverNumber.IsSuccess() {
			s += "\n ✗ " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Could not determine rover number")
		} else if *m.roverNumber.Data != index {
			s += "\n ! " + lipgloss.NewStyle().Foreground(style.WarningPrimary).Render("This Rover presented itself as Rover "+strconv.Itoa(*m.roverNumber.Data)+" but you wanted to connect to Rover "+m.formValues.index)
		} else {
			s += "\n ✓ " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Rover number matches the index you entered ("+m.formValues.index+")")
		}

		if !m.authValid.IsSuccess() {
			s += "\n ✗ " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Authentication to the roverd endpoint failed. Please check your credentials")
		} else {
			s += "\n ✓ " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Authentication successful")
		}
	}

	if !m.routeExists.IsSuccess() || !m.authValid.IsSuccess() || !m.roverdVersion.IsSuccess() {
		s += "\n\n" + lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("This connection configuration is not valid and cannot be saved.")
	} else {
		s += "\n\n" + "You are all set! Press enter to go start using your Rover, or press 'n' to add another connection."
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
	// todo: replace with actual authentication check
	return tui.PerformAction(&m.authValid, func() (*bool, error) {
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

func checkRoverdVersion(m model) tea.Cmd {
	return tui.PerformAction(&m.roverdVersion, func() (*string, error) {
		res := "linux 1234"
		return &res, nil
	})
}

func checkRoverNumber(m model) tea.Cmd {
	return tui.PerformAction(&m.roverNumber, func() (*int, error) {
		res := 123
		return &res, nil
	})
}
