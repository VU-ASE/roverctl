package view

import (
	"strings"

	"github.com/VU-ASE/rover/src/configuration/asciitool"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var index = 0

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

// The update view (where you can download the latest version of all modules) with the update method
func (a AppState) ConfigureUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {
		case key.Matches(msg, keys.Help):
			a.help.ShowAll = !a.help.ShowAll
		case key.Matches(msg, keys.Left):
			index--
		case key.Matches(msg, keys.Right):
			index++

		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return a, nil
}

// the update view with the view method
func (a AppState) ConfigureView() string {
	docStyle := lipgloss.NewStyle().Margin(1, 2)

	// The header
	s := "Pipeline schematic:\n\n"

	// Create a pipeline drawing with the mermaid-ascii tool
	pipeline := `graph LR
topcam --> controller
controller --> actuator`

	// Draw the pipeline
	pipelineDrawing, err := asciitool.Draw(pipeline)
	if err != nil {
		s = "Your pipeline could not be drawn. Please check the configuration."
	} else {
		s += pipelineDrawing
	}

	if index == 0 {
		s = strings.ReplaceAll(s, "topcam", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Render("topcam"))
	} else if index == 1 {
		s = strings.ReplaceAll(s, "controller", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Render("controller"))
	} else if index == 2 {
		s = strings.ReplaceAll(s, "actuator", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Render("actuator"))
	}

	s += "\nEnabled services:\n\n"

	helpView := a.help.View(keys)
	height := 20 - strings.Count(s, "\n") - strings.Count(helpView, "\n")

	s = "\n" + s + strings.Repeat("\n", height) + helpView

	return docStyle.Render(s)
}
