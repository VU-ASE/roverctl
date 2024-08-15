package view

import (
	"github.com/VU-ASE/rover/src/configuration"
	"github.com/VU-ASE/rover/src/style"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// This is the struct that implements the Bubbletea model interface. It contains all the app state.
type AppState struct {
	// To select an action to perform with this utility
	connectedActions    list.Model // actions you can perform when connected to a Rover
	disconnectedActions list.Model // actions you can perform when not connected to a Rover
	selectedAction      string
	roverConnections    configuration.RoverConnections // used to track state changes, if the connection state changes
	// To display a help footer
	help help.Model
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// A list of actions that can be performed with this utility
type actionItem struct {
	title, desc string
}

func (i actionItem) Title() string       { return i.title }
func (i actionItem) Description() string { return i.desc }
func (i actionItem) FilterValue() string { return i.title }

func InitialApp() AppState {
	a := initializeAppState()
	a.help = help.New()

	return a
}

// Implementation of the Bubbletea model interface
func (a AppState) Init() tea.Cmd {
	if a.selectedAction != "" {
		switch a.selectedAction {
		case "Connect":
			// Connect to a Rover
			return a.ConnectInit()
		}
	}

	return nil
}

// Define the root controls for the app (so that you can quit)
func (a AppState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c":
			return a, tea.Quit
		// If esc is pressed, go back to the home view
		case "esc", "q":
			if a.selectedAction == "" {
				return a, tea.Quit
			}

			a.selectedAction = ""
			return a, nil
		}
	}

	if a.selectedAction != "" {
		switch a.selectedAction {
		case "Configure":
			// Configure the pipeline
			return a.ConfigureUpdate(msg)
		case "Connect":
			// Connect to a Rover
			return a.ConnectUpdate(msg)
		}
	}

	// Use the home view, depending on if we are connected or not
	if len(a.roverConnections.Available) > 0 {
		return a.ConnectedHomeUpdate(msg)
	} else {
		return a.DisconnectedHomeUpdate(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	return a, nil
}

func (a AppState) View() string {
	if a.selectedAction != "" {
		switch a.selectedAction {
		case "Configure":
			// Configure the pipeline
			return a.ConfigureView()
		case "Connect":
			// Connect to a Rover
			return a.ConnectView()
		}
	}

	// Use just the home view for now
	s := a.DisconnectedHomeView()
	if len(a.roverConnections.Available) > 0 {
		s = a.ConnectedHomeView()
	}

	return docStyle.Render(s)
}

// Styling
var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(0)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(0).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

// A user can choose from these actions if they are connected to a Rover
var onConnectedActions = []list.Item{
	actionItem{title: "Start", desc: "Start your pipeline"}, // Should be "stop" when a pipeline is running
	actionItem{title: "Configure", desc: "Configure your pipeline"},
	actionItem{title: "Debug", desc: "Enable remote debugging for your pipeline"}, // Should not be available when no pipeline is running or disable when enabled
	actionItem{title: "Status", desc: "Watcdh module outputs and status logs"},    // Should not be available when no pipeline is running
	actionItem{title: "Update", desc: "Fetch the latest versions of all modules and install them"},
}

// A user can choose from these actions if they are not yet connected to a Rover
var onDisconnectedActions = []list.Item{
	actionItem{title: "Connect", desc: "Initialize a connection to a Rover"}, // Should be "stop" when a pipeline is running
}

// This initializes all actions and lists that can be used throughout the app. Both for connected and disconnected states
func initializeAppState() AppState {
	// Get the list of connections
	connections, err := configuration.ReadConnections()
	if err != nil {
		// todo: throw? Or show a warning?
	}

	// Initialize the list of actions when connected
	cl := list.New(onConnectedActions, list.NewDefaultDelegate(), 0, 0)
	// If there are connections available, add the connected actions
	cl.Title = lipgloss.NewStyle().Background(style.AsePrimary).Bold(true).Padding(0, 0).Render("VU ASE") + lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C")).Render(" - racing Rovers since 2024")
	cl.SetShowStatusBar(false)
	cl.SetFilteringEnabled(false)
	cl.Styles.Title = titleStyle
	cl.Styles.PaginationStyle = paginationStyle
	cl.Styles.HelpStyle = helpStyle

	// Initialize the list of actions when disconnected
	dl := list.New(onDisconnectedActions, list.NewDefaultDelegate(), 0, 0)
	// If there are connections available, add the connected actions
	dl.Title = lipgloss.NewStyle().Background(style.AsePrimary).Bold(true).Padding(0, 0).Render("VU ASE") + lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C")).Render(" - racing Rovers since 2024") + lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Render("\nYou did not connect to a Rover yet")
	dl.SetShowStatusBar(false)
	dl.SetFilteringEnabled(false)
	dl.Styles.Title = titleStyle
	dl.Styles.PaginationStyle = paginationStyle
	dl.Styles.HelpStyle = helpStyle

	return AppState{
		connectedActions:    cl,
		disconnectedActions: dl,
		selectedAction:      "",
		roverConnections:    connections,
	}
}
