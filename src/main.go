package main

import (
	"os"

	"github.com/VU-ASE/roverctl/src/configuration"
	"github.com/VU-ASE/roverctl/src/state"
	"github.com/VU-ASE/roverctl/src/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func run() error {
	// Initialize the app
	err := configuration.Initialize()
	if err != nil {
		return err
	}

	// Create the app state
	appState := state.Get()

	// We start the app in a separate (full) screen
	p := tea.NewProgram(views.RootScreen(appState), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	// Save the connections to disk
	return state.Get().RoverConnections.Save()
}

func main() {
	// Configure zerolog to output to stdout beautifully
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) //

	// Run the program
	if err := run(); err != nil {
		log.Error().Err(err).Msg("An error occurred while running the program.")
		os.Exit(1)
	}
}
