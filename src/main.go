package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/VU-ASE/rover/src/configuration"
	initconnectionpage "github.com/VU-ASE/rover/src/pages/connections/init"
	startpageconnected "github.com/VU-ASE/rover/src/pages/start/connected"
	startpagedisconnected "github.com/VU-ASE/rover/src/pages/start/disconnected"
	"github.com/VU-ASE/rover/src/state"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func selectPage(s *state.AppState) tea.Model {
	switch s.CurrentView {
	// SSH is different, it replaces the current process
	case "SSH":
		{
			// Get the active connection
			activeConnection := s.RoverConnections.GetActive()
			if activeConnection == nil {
				// This should never happen
				syscall.Exec("/bin/echo", []string{"error"}, os.Environ())
				return nil
			}

			ssh, lookErr := exec.LookPath("ssh")
			if lookErr != nil {
				panic(lookErr)
			}
			connectionString := fmt.Sprintf("%s@%s", activeConnection.Username, activeConnection.Host)
			syscall.Exec(ssh, []string{"ssh", connectionString, "-p", "22"}, os.Environ())
			return nil
		}
	case "Connect":
		return initconnectionpage.InitialModel(nil)
	default:
		{
			if len(s.RoverConnections.Available) > 0 {
				return startpageconnected.InitialModel()
			} else {
				return startpagedisconnected.InitialModel()
			}
		}
	}

}

func run() error {
	// Initialize the app
	err := configuration.Initialize()
	if err != nil {
		return err
	}
	defer configuration.Cleanup()

	// Create the app state
	appState := state.Get()
	defer appState.RoverConnections.Save()

	// We start the app in a separate (full) screen
	firsttime := true
	for firsttime || appState.CurrentView != "" {
		firsttime = false
		page := selectPage(appState)
		p := tea.NewProgram(page, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// Configure zerolog to output to stdout beautifully
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Run the program
	if err := run(); err != nil {
		log.Error().Err(err).Msg("An error occurred while running the program.")
		os.Exit(1)
	}
}
