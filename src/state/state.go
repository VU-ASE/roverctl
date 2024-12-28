package state

import (
	"math/rand"

	"github.com/VU-ASE/roverctl/src/configuration"
)

// This is the struct that implements the Bubbletea model interface. It contains all the app state.
// It is used as a singleton to keep track of the state of the app
type AppState struct {
	RoverConnections configuration.RoverConnections // used to track state changes, if the connection state changes
	// Window sizes
	WindowWidth  int
	WindowHeight int
	Quote        string
}

var state *AppState = nil

// This initializes all actions and lists that can be used throughout the app. Both for connected and disconnected states
func initialize() *AppState {
	// Get the list of connections
	connections, err := configuration.ReadConnections()
	if err != nil {
		// todo: throw? Or show a warning?
	}

	// Pick a random quote
	quote := quotes[rand.Intn(len(quotes))]

	return &AppState{
		RoverConnections: connections,
		Quote:            quote,
	}
}

func Get() *AppState {
	if state == nil {
		state = initialize()
	}
	return state
}

// List of quotes that are displayed on top
var quotes = []string{
	"racing Rovers since 2023",
	"configuring pipelines since 2023",
	"setting up services since 2023",
	"burning rubber since 2023",
	"exploring autonomously since 2023",
}
