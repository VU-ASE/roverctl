package view

import (
	"fmt"
	"strconv"

	"github.com/VU-ASE/rover/src/components"
	"github.com/VU-ASE/rover/src/configuration"
	"github.com/VU-ASE/rover/src/style"
	tea "github.com/charmbracelet/bubbletea"
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
	Username: "",
	Password: "",
}

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
			txtfieldValue := textfield.Input.Value()
			textfield.Err = nil

			if newConnection.Host == "" {
				// Try to parse the input as a number between 1 and 65535
				// If it is not a number, set an error
				roverIndex, err := strconv.Atoi(txtfieldValue)
				if err != nil || roverIndex < 1 || roverIndex > 20 {
					textfield.Err = fmt.Errorf("Please enter a valid Rover index between 1 and 20 (inclusive)")
				} else {
					newConnection.Host = fmt.Sprintf("192.168.1.%d", roverIndex+100)
					textfield.Input.SetValue(defaultUsername)
					textfield.Input.SetCursor(len(defaultUsername))
				}
			} else if newConnection.Username == "" {
				newConnection.Username = txtfieldValue
				textfield.Input.SetValue(defaultPassword)
				textfield.Input.SetCursor(len(defaultPassword))
			} else if newConnection.Password == "" {
				newConnection.Password = txtfieldValue
				textfield.Input.SetValue("")
				textfield.Input.SetCursor(0)
			} else if newConnection.Name == "" {
				newConnection.Name = txtfieldValue
			} else {
				testingConnection = true
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	textfield.Input, cmd = textfield.Input.Update(msg)
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

	s += "\n\nEnter the number of the Rover you want to connect to (1-20):"
	if newConnection.Host == "" {
		s += "\n" + textfield.Input.View()
		if textfield.Err != nil {
			s += "\n" + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render(textfield.Err.Error())
		}

		return docStyle.Render(s)
	}
	s += "\n - " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Using IP "+newConnection.Host)

	s += "\n\nEnter the username used for authentication:"
	if newConnection.Username == "" {
		s += "\n" + textfield.Input.View()
		if textfield.Err != nil {
			s += "\n" + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render(textfield.Err.Error())
		}

		return docStyle.Render(s)
	}
	s += "\n - " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render(""+newConnection.Username)

	s += "\n\nEnter the password used for authentication:"
	if newConnection.Password == "" {
		s += "\n" + textfield.Input.View()
		if textfield.Err != nil {
			s += "\n" + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render(textfield.Err.Error())
		}

		return docStyle.Render(s)
	}
	s += "\n - " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Password set")

	s += "\n\nEnter a unique name for this connection to find it back later:"
	if newConnection.Name == "" {
		s += "\n" + textfield.Input.View()
		if textfield.Err != nil {
			s += "\n" + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render(textfield.Err.Error())
		}

		return docStyle.Render(s)
	}
	s += "\n - " + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render(newConnection.Name)

	s += "\n\n" + lipgloss.NewStyle().Bold(true).Foreground(style.SuccessPrimary).Render("You are all set! Press enter to save and connect.")

	return docStyle.Render(s)
}

func (a AppState) connectTestConnection() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Connecting to " + newConnection.Name)

	s += "\n\n" + spin.View()
	return docStyle.Render(s)
}

func (a AppState) ConnectInit() tea.Cmd {
	return spin.Tick
}

func (a AppState) ConnectView() string {
	if testingConnection {
		return a.connectTestConnection()
	} else {
		return a.connectEnterDetails()
	}
}
