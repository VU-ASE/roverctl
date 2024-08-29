package uploadservicepage

import (
	"fmt"

	roverlock "github.com/VU-ASE/rover/src/lock"
	"github.com/VU-ASE/rover/src/services"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	spinner        spinner.Model
	lockAction     tui.Action
	unlockAction   tui.Action
	transferAction tui.Action
	error          error // any errors that occurred
}

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Line

	return model{
		spinner:        s,
		transferAction: tui.NewAction("transfer"),
		error:          nil,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m = InitialModel()
			return m, m.Init()
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tui.ActionResult:
		actions := tui.Actions{&m.lockAction, &m.unlockAction, &m.transferAction}
		actions.ProcessResult(msg)
	default:
		// Base command
		model, cmd := tui.Update(m, msg)
		return model, cmd
	}

	return m, nil
}

// the update view with the view method
func (m model) uploadResultsView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Upload service")

	if m.transferAction.IsSuccess() {
		s += "\n\n" + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Files uploaded successfully")
	} else {
		s += "\n\n" + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Failed to upload files")
		if m.error != nil {
			s += "\n > " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render(m.error.Error())
		}
	}

	s += lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("\n\nPress 'r' to retry, or 'q' to quit")

	return s
}

func (m model) uploadingView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Upload service")

	s += "\n\n" + m.spinner.View() + " Uploading files..."

	return s
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, uploadService(m))
}

func (m model) View() string {
	if m.transferAction.IsLoading() {
		return style.Docstyle.Render(m.uploadResultsView())
	} else {
		return style.Docstyle.Render(m.uploadingView())
	}
}

func uploadService(m model) tea.Cmd {
	return tui.PerformAction(&m.transferAction, func() error {
		conn := state.Get().RoverConnections.GetActive()
		if conn == nil {
			return fmt.Errorf("Not connected to an active Rover")
		}

		// Lock the rover
		err := roverlock.Lock(*conn)
		if err != nil {
			return err
		}

		// Upload the service
		err = services.Upload(*conn)
		if err != nil {
			return err
		}

		// Unlock the rover
		err = roverlock.Unlock(*conn)
		return err
	})
}
