package state

import (
	"github.com/VU-ASE/rover/src/configuration"
)

// This is the struct that implements the Bubbletea model interface. It contains all the app state.
// It is used as a singleton to keep track of the state of the app
type AppState struct {
	CurrentView      string                         // Keep track of the view that we are currently in
	RoverConnections configuration.RoverConnections // used to track state changes, if the connection state changes
}

var state *AppState = nil

// This initializes all actions and lists that can be used throughout the app. Both for connected and disconnected states
func initialize() *AppState {
	// Get the list of connections
	connections, err := configuration.ReadConnections()
	if err != nil {
		// todo: throw? Or show a warning?
	}

	return &AppState{
		CurrentView:      "",
		RoverConnections: connections,
	}
}

func Get() *AppState {
	if state == nil {
		state = initialize()
	}
	return state
}
