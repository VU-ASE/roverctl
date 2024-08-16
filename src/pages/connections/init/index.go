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
	"github.com/melbahja/goph"
	probing "github.com/prometheus-community/pro-bing"
	wifiname "github.com/yelinaung/wifi-name"
	"gopkg.in/grignaak/tribool.v1"
)

// Persistent global state (ugly, yes) to allow retrying of connection checks by discarding results with an attempt number lower than the current one
var attemptNumber = 1

type formValues struct {
	name     string
	index    string
	username string
	password string
}

// Action codes
const (
	routePingAction  = "routePing"
	authCheckAction  = "authCheck"
	debixCheckAction = "debixCheck"
)

// Used to communicate the result of various tests
type resultMsg struct {
	action  string
	result  bool
	err     error
	attempt int
}

type model struct {
	form        *huh.Form
	spinner     spinner.Model
	routeExists tribool.Tribool
	authValid   tribool.Tribool
	debixValid  tribool.Tribool
	isChecking  bool
	formValues  *formValues
	host        string // the IP of the rover to use
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

	return model{
		spinner:     s,
		formValues:  formValues,
		host:        "",
		routeExists: tribool.Maybe,
		authValid:   tribool.Maybe,
		debixValid:  tribool.Maybe,
		isChecking:  false,
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
					Value(&formValues.name),
			),
		).WithTheme(style.FormTheme),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.form.State == huh.StateCompleted {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				// Save the connection if all checks are successful
				if m.routeExists == tribool.True && m.authValid == tribool.True && m.debixValid == tribool.True {
					// Save the connection
					state.Get().RoverConnections = state.Get().RoverConnections.Add(configuration.RoverConnection{
						Name:     m.formValues.name,
						Host:     m.host,
						Username: m.formValues.username,
						Password: m.formValues.password,
					})
					state.Get().CurrentView = "home"
					return m, tea.Quit
				}
			case "n":
				if m.routeExists == tribool.True && m.authValid == tribool.True && m.debixValid == tribool.True {
					// Save the connection
					state.Get().RoverConnections = state.Get().RoverConnections.Add(configuration.RoverConnection{
						Name:     m.formValues.name,
						Host:     m.host,
						Username: m.formValues.username,
						Password: m.formValues.password,
					})
				}
			case "b":
				// Restore to the initial form, but recover the form values
				m = InitialModel(m.formValues)
				attemptNumber++
				return m, tea.Batch(m.form.Init(), m.spinner.Tick)
			case "r":
				// Retry the connection checks
				attemptNumber++
				m.isChecking = true
				m.authValid = tribool.Maybe
				m.routeExists = tribool.Maybe
				m.debixValid = tribool.Maybe
				return m, tea.Batch(checkRoute(m, attemptNumber), checkAuth(m, attemptNumber), checkDebix(m, attemptNumber))
			}
		}
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case resultMsg:
		if msg.attempt < attemptNumber {
			return m, nil
		}
		switch msg.action {
		case routePingAction:
			m.routeExists = tribool.FromBool(msg.result)
		case authCheckAction:
			m.authValid = tribool.FromBool(msg.result)
		case debixCheckAction:
			m.debixValid = tribool.FromBool(msg.result)
		}
		return m, nil
	default:
		// Base command
		model, cmd := tui.Update(m, msg)
		if cmd != nil {
			return model, cmd
		}

		cmds := []tea.Cmd{}
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			if f.State == huh.StateCompleted && !m.isChecking {
				m.isChecking = true

				index, err := strconv.Atoi(m.formValues.index)
				if err != nil || index < 1 || index > 20 {
					m.routeExists = tribool.False
					return m, cmd
				}
				m.host = fmt.Sprintf("192.168.1.%d", index+100)
				m.host = "www.google.com"
				m.host = "localhost"

				cmds = append(cmds, checkRoute(m, attemptNumber), checkAuth(m, attemptNumber), checkDebix(m, attemptNumber))
			}
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
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

	if m.routeExists == tribool.Maybe {
		s += "\n\n " + m.spinner.View() + " checking if a route to Rover exists"
		return s
	} else if m.routeExists == tribool.False {
		s += "\n\n - " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("No route could be established to the Rover at "+m.host+" Are you sure it is powered on? ")
	} else {
		s += "\n\n - " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Established route to Rover at "+m.host)
	}

	if m.authValid == tribool.Maybe {
		s += "\n " + m.spinner.View() + " checking if authentication is valid"
		return s
	} else if m.authValid == tribool.False {
		s += "\n - " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Authentication failed. Please check your credentials")
	} else {
		s += "\n - " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Authentication successful")
	}

	if m.debixValid == tribool.Maybe {
		s += "\n " + m.spinner.View() + " checking Debix state"
		return s
	} else if m.debixValid == tribool.False {
		s += "\n - " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Debix is not properly configured")
	} else {
		s += "\n - " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Debix is properly configured")
	}

	if m.routeExists == tribool.False || m.authValid == tribool.False || m.debixValid == tribool.False {
		s += "\n\n" + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("This Rover configuration is not valid and cannot be saved.")
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

func checkRoute(m model, a int) tea.Cmd {
	return func() tea.Msg {
		ping, _ := probing.NewPinger(m.host)
		ping.Count = 3
		ping.Timeout = 10 * time.Second
		err := ping.Run()
		if ping.Statistics().PacketsRecv > 0 {
			return resultMsg{result: err == nil, err: nil, action: routePingAction, attempt: a}
		} else {
			return resultMsg{result: false, err: err, action: routePingAction, attempt: a}
		}
	}
}

func checkAuth(m model, a int) tea.Cmd {
	return func() tea.Msg {
		// Start new ssh connection with private key.
		auth := goph.Password(m.formValues.password)
		client, err := goph.New(m.formValues.username, m.host, auth)
		if err != nil {
			return resultMsg{result: false, err: err, action: authCheckAction, attempt: a}
		}
		defer client.Close()

		// Check if the connection is working by running a simple command
		_, err = client.Run("ls /tmp/")
		return resultMsg{result: err == nil, err: nil, action: authCheckAction, attempt: a}
	}
}

func checkDebix(m model, a int) tea.Cmd {
	return func() tea.Msg {
		// Start new ssh connection with private key.
		auth := goph.Password(m.formValues.password)
		client, err := goph.New(m.formValues.username, m.host, auth)
		if err != nil {
			return resultMsg{result: false, err: err, action: debixCheckAction, attempt: a}
		}
		defer client.Close()

		// Check if the connection is working by running a simple command
		_, err = client.Run("ls ./debix") // todo: change to /home/debix
		return resultMsg{result: err == nil, err: nil, action: debixCheckAction, attempt: a}
	}
}
