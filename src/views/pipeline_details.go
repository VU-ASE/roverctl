package views

import (
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//
// All keys
//

// Keys to navigate
type PipelineDetailsKeyMap struct {
	Retry   key.Binding
	Confirm key.Binding
	Quit    key.Binding
}

// Shown when the services are being updated
var pipelineDetailsKeysRegular = PipelineDetailsKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func (k PipelineDetailsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Quit}
}

func (k PipelineDetailsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

//
// The page model
//

type PipelineDetailsPage struct {
	help    help.Model
	spinner spinner.Model
	service PipelineOverviewServiceInfo
}

func NewPipelineDetailsPage(s PipelineOverviewServiceInfo) PipelineDetailsPage {
	return PipelineDetailsPage{
		spinner: spinner.New(),
		help:    help.New(),
		service: s,
	}
}

//
// Page model methods
//

func (m PipelineDetailsPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, pipelineDetailsKeysRegular.Quit):
			return m, tea.Quit
		case key.Matches(msg, pipelineDetailsKeysRegular.Retry):
			// todo:
			return m, nil
		case key.Matches(msg, pipelineDetailsKeysRegular.Confirm):
			// todo:
			return m, nil
		}
	}

	return m, nil
}

func (m PipelineDetailsPage) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m PipelineDetailsPage) View() string {
	s := style.Title.Render("Details") + "\n\n"

	// Calculate column width (subtract padding and borders)
	columnWidth := (state.Get().WindowWidth - 4 - 6) / 3 // Adjust for padding and borders

	// Define styles for each column
	columnStyle := lipgloss.NewStyle().
		Width(columnWidth)

		// About
	about := style.Gray.Render("Service: ") + m.service.Name + "\n"
	about += style.Gray.Render("Author: ") + m.service.Author + "\n"
	about += style.Gray.Render("Version: ") + m.service.Version + "\n"

	inputs := ""
	if len(m.service.Configuration.Inputs) <= 0 {
		inputs += style.Gray.Render("This service does not depend on streams from other services")
	} else {
		inputs += style.Gray.Render("Inputs: ") + "\n"

		for _, i := range m.service.Configuration.Inputs {
			for _, s := range i.Streams {
				inputs += style.Gray.Render("- ") + s + style.Gray.Render(" from ") + i.Service + "\n"
			}
		}
	}

	outputs := ""
	if len(m.service.Configuration.Outputs) <= 0 {
		outputs += style.Gray.Render("This service does not provide streams to other services")
	} else {
		outputs += style.Gray.Render("Outputs: ") + "\n"

		for _, o := range m.service.Configuration.Outputs {
			outputs += style.Gray.Render("- ") + o + "\n"
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top,
		columnStyle.Render(about),
		columnStyle.Render(inputs),
		columnStyle.Render(outputs),
	)

	return s + row + "\n" + m.help.View(pipelineDetailsKeysRegular)

}
