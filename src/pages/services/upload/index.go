package uploadservicepage

import (
	"github.com/VU-ASE/rover/src/services"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/grignaak/tribool.v1"
)

// Persistent global state (ugly, yes) to allow retrying of connection checks by discarding results with an attempt number lower than the current one
var attemptNumber = 1

// Action codes
const (
	transferAction = "transferfiles"
)

// Used to communicate the result of various tests
type resultMsg struct {
	action  string
	result  bool
	err     error
	attempt int
}

type model struct {
	spinner          spinner.Model
	filesTransferred tribool.Tribool
	error            error // any errors that occurred
}

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Line
	attemptNumber++

	return model{
		spinner:          s,
		filesTransferred: tribool.Maybe,
		error:            nil,
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
	case resultMsg:
		if msg.attempt < attemptNumber {
			return m, nil
		}
		switch msg.action {
		case transferAction:
			m.error = msg.err
			m.filesTransferred = tribool.FromBool(msg.result)
		}
		return m, nil
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

	if m.filesTransferred == tribool.True {
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
	return tea.Batch(m.spinner.Tick, uploadService(attemptNumber))
}

func (m model) View() string {
	if m.filesTransferred != tribool.Maybe {
		return style.Docstyle.Render(m.uploadResultsView())
	} else {
		return style.Docstyle.Render(m.uploadingView())
	}
}

func uploadService(a int) tea.Cmd {
	return func() tea.Msg {
		err := services.UploadService()
		return resultMsg{result: err == nil, err: err, action: transferAction, attempt: a}
	}
}
