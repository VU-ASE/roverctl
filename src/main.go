package main

import (
	"os"

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
	if s.CurrentView == "Connect" {
		return initconnectionpage.InitialModel()
	}

	switch s.CurrentView {
	case "Connect":
		return initconnectionpage.InitialModel()
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
