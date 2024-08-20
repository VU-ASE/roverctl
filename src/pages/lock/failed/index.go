package lockfailedpage

import (
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	message string
}

func InitialModel(msg string) model {
	return model{
		message: msg,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			// Retry the lock
			return m, tea.Quit
		}
	}

	// Base command
	model, cmd := tui.Update(m, msg)
	return model, cmd
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	s := lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Failed to lock your Rover")

	s += "\n\nThe operation you were about to perform required exclusive access to your Rover\nbut this could not be obtained. The Rover reported the following error:"
	s += "\n > " + lipgloss.NewStyle().Foreground(style.WarningPrimary).Render(m.message)

	s += "\n\nIf no-one else is working on the Rover, rebooting the Rover by holding the power button for 10 seconds will release the lock. "

	s += lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("\n\nPress 'r' to retry, or 'q' to go back.")
	return style.Docstyle.Render(s)
}
