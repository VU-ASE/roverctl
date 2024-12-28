package views

import (
	"github.com/VU-ASE/roverctl/src/style"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

//
// All keys
//

// Keys to navigate
type TemplateKeyMap struct {
	Retry   key.Binding
	Confirm key.Binding
	Quit    key.Binding
}

// Shown when the services are being updated
var templateKeysRegular = TemplateKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func (k TemplateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Retry, k.Confirm, k.Quit}
}

func (k TemplateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

//
// The page model
//

type TemplatePage struct {
	help    help.Model
	spinner spinner.Model
}

func NewTemplatePage() TemplatePage {
	// todo

	return TemplatePage{
		spinner: spinner.New(),
		help:    help.New(),
	}
}

//
// Page model methods
//

func (m TemplatePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, templateKeysRegular.Quit):
			return m, tea.Quit
		case key.Matches(msg, templateKeysRegular.Retry):
			// todo:
			return m, nil
		case key.Matches(msg, templateKeysRegular.Confirm):
			// todo:
			return m, nil
		}
	}

	return m, nil
}

func (m TemplatePage) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m TemplatePage) View() string {
	return style.Title.Render("Template page")
}
