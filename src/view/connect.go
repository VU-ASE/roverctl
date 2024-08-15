package view

import (
	"fmt"
	"strconv"

	"github.com/VU-ASE/rover/src/components"
	"github.com/VU-ASE/rover/src/configuration"
	"github.com/VU-ASE/rover/src/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	wifiname "github.com/yelinaung/wifi-name"
	"gopkg.in/grignaak/tribool.v1"
)

var textfield = components.InitializeTextfield()
var spin = components.InitializeSpinner()

// Keep track of the current state of the connection
var newConnection = configuration.RoverConnection{
	Name:     "",
	Host:     "",
	Username: "debix",
	Password: "debix",
}

var roverIndex = ""

var form = huh.NewForm(
	huh.NewGroup(
		huh.NewInput().
			Title("Enter your Rover index (1-20, inclusive)").
			CharLimit(3).
			Prompt("> ").
			Value(&roverIndex).Validate(func(s string) error {
			index, err := strconv.Atoi(s)
			if err != nil || index < 1 || index > 20 {
				return fmt.Errorf("Please enter a valid Rover index between 1 and 20 (inclusive)")
			} else {
				newConnection.Host = fmt.Sprintf("192.168.1.%d", index+100)
				newConnection.Name = fmt.Sprintf("rover%d", index)
			}
			return nil
		}),
		huh.NewInput().
			Title("Enter the authentication username").
			CharLimit(255).
			Prompt("> ").
			Value(&newConnection.Username),
		huh.NewInput().
			Title("Enter the authentication password").
			CharLimit(255).
			Prompt("> ").
			Value(&newConnection.Password).
			EchoMode(huh.EchoModePassword),
		huh.NewInput().
			Title("Enter a name for this connection to find it back later").
			CharLimit(255).
			Prompt("> ").
			Value(&newConnection.Name),
	),
).WithTheme(style.FormTheme)

const defaultUsername = "debix"
const defaultPassword = "debix"

// To power the spinners that aare shown when testing the connection
var testingConnection = false
var routeExists = tribool.Maybe
var credentialsValid = tribool.Maybe
var connectionEstablished = tribool.Maybe

// The update view (where you can download the latest version of all modules) with the update method
func (a AppState) ConnectUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		case "enter":

		}
	}

	newSpinner, cmd := spin.Update(msg)
	spin = newSpinner
	// form, cmd := form.Update(msg)
	// if f, ok := form.(*huh.Form); ok {
	// 	form = f
	// }
	return a, cmd
}

// the update view with the view method
func (a AppState) connectEnterDetails() string {

	// Introduction
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Connect to a Rover") + "\n\n - Make sure you are connected to the ASElabs WiFi\n - Make sure that the Rover is powered on"

	// Get the current wifi name
	wifi := wifiname.WifiName()

	if wifi != "aselabs" {
		s += lipgloss.NewStyle().Bold(true).Foreground(style.WarningPrimary).Render("\n\nIt seems you are not connected to the ASElabs WiFi but to '" + wifi + "' instead. \nRead how to connect at: https://docs.ase.vu.nl/docs/tutorials/setting-up-your-workspace/accessing-the-network")
	}

	s += "\n\n" + form.View()

	return docStyle.Render(s)
}

func (a AppState) connectTestConnection() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Connecting to " + newConnection.Name)

	s += "\n\n" + spin.View()
	return docStyle.Render(s)
}

func (a AppState) ConnectInit() tea.Cmd {
	// return tea.Batch(form.Init(), spin.Tick)
	return spin.Tick
}

func (a AppState) ConnectView() string {
	return a.connectTestConnection()

	// if form.State == huh.StateCompleted {
	// 	// class := form.GetString("class")
	// 	// level := form.GetString("level")
	// 	return a.connectTestConnection()
	// } else {
	// 	return a.connectEnterDetails()
	// }
}
