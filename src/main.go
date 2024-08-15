package main

import (
	"os"

	"github.com/VU-ASE/rover/src/configuration"
	"github.com/VU-ASE/rover/src/view"
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
	defer configuration.Cleanup()

	// We start the app in a separate (full) screen
	app := view.InitialApp()
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
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
