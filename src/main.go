package main

import (
	"os"

	"github.com/VU-ASE/rover/src/configuration"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// func selectPage(s *state.AppState) tea.Model {
// 	switch strings.ToLower(s.Route.Peek()) {
// 	// SSH is different, it replaces the current process
// 	case "ssh":
// 		{
// 			// Get the active connection
// 			activeConnection := s.RoverConnections.GetActive()
// 			if activeConnection == nil {
// 				// This should never happen
// 				syscall.Exec("/bin/echo", []string{"error"}, os.Environ())
// 				return nil
// 			}

// 			ssh, lookErr := exec.LookPath("ssh")
// 			if lookErr != nil {
// 				panic(lookErr)
// 			}
// 			connectionString := fmt.Sprintf("%s@%s", activeConnection.Username, activeConnection.Host)
// 			syscall.Exec(ssh, []string{"ssh", connectionString, "-p", "22"}, os.Environ())
// 			return nil
// 		}
// 	case "connections":
// 		return manageconnectionspage.InitialModel()
// 	case "utilities":
// 		return utilitiespage.InitialModel()
// 	case "connection init":
// 		return initconnectionpage.InitialModel(nil)
// 	case "services":
// 		return servicespage.InitialModel()
// 	case "sources update":
// 		return updatesourcespage.InitialModel()
// 	// case "service init":
// 	// 	return initservicepage.InitialModel()
// 	// case "service upload":
// 	// 	return uploadservicepage.InitialModel()
// 	// case "pipeline configure":
// 	// 	return enableservicespage.InitialModel()
// 	default:
// 		{
// 			if len(s.RoverConnections.Available) > 0 {
// 				return
// 			} else {
// 				return startpagedisconnected.InitialModel()
// 			}
// 		}
// 	}

// }

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
	// for !appState.Route.IsEmpty() {
	p := tea.NewProgram(views.RootScreen(appState), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	// }

	// Save the connections to disk
	return state.Get().RoverConnections.Save()
}

func main() {
	// Clear the screen
	// fmt.Println("\033[2J")

	// Configure zerolog to output to stdout beautifully
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Run the program
	if err := run(); err != nil {
		log.Error().Err(err).Msg("An error occurred while running the program.")
		os.Exit(1)
	}
}
