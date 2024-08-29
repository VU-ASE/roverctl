package enableservicespage

import (
	"fmt"
	"io"
	"strings"

	"github.com/VU-ASE/rover/src/configuration/asciitool"
	roverlock "github.com/VU-ASE/rover/src/lock"
	"github.com/VU-ASE/rover/src/roveryaml"
	"github.com/VU-ASE/rover/src/services"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	keys                     keyMap
	spinner                  spinner.Model
	list                     list.Model
	fetchServicesAction      tui.Action[[]services.FoundService]
	fetchConfigurationAction tui.Action[roveryaml.RoverConfig]
	error                    error // Can be shown to the user
}

// keyMap defines a set of keybindings. To work for help it must satisfy key.Map
type keyMap struct {
	Edit       key.Binding
	Delete     key.Binding
	MarkActive key.Binding
	New        key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.New, k.Edit, k.Delete, k.MarkActive}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var keys = keyMap{
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	MarkActive: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "set active"),
	),
	Delete: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "delete"),
	),
}

type item struct {
	service services.FoundService
	active  bool
}

func (i item) FilterValue() string { return i.service.Service.Name }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := i.service.Service.Name + " " + lipgloss.NewStyle().Foreground(style.AsePrimary).Render(i.service.Service.Version) + " " + lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("("+i.service.Path+")")
	if i.active {
		str += lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render(" (active)")
	}

	fn := lipgloss.NewStyle().Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().Bold(true).Render("> " + strings.Join(s, " "))
		}
	} else {
		str = "- " + str
	}

	fmt.Fprint(w, fn(str))
}

func servicesToListItem(services []services.FoundService, config *roveryaml.RoverConfig) []list.Item {
	items := make([]list.Item, 0)

	if services == nil {
		return items
	}

	for _, service := range services {
		items = append(items, item{
			service: service,
			active:  config.HasEnabled(service.Path),
		})
	}

	return items
}

func InitialModel() model {
	l := list.New([]list.Item{}, itemDelegate{}, 0, 14)
	l.Title = "Manage Rover services"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = style.TitleStyle
	l.Styles.PaginationStyle = style.PaginationStyle
	l.Styles.HelpStyle = style.HelpStyle
	l.AdditionalShortHelpKeys = keys.ShortHelp

	spin := spinner.New()
	spin.Spinner = spinner.Line

	fetchServicesAction := tui.NewAction[[]services.FoundService]("getServices")
	fetchServicesAction.Started = true

	fetchConfigAction := tui.NewAction[roveryaml.RoverConfig]("getConfig")
	fetchConfigAction.Started = true

	return model{
		keys:                     keys,
		list:                     l,
		spinner:                  spin,
		fetchServicesAction:      fetchServicesAction,
		fetchConfigurationAction: fetchConfigAction,
		error:                    nil,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.fetchServicesAction.IsSuccess() && m.fetchConfigurationAction.IsSuccess() {
		m.list.SetItems(servicesToListItem(*m.fetchServicesAction.Data, m.fetchConfigurationAction.Data))
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, _ := style.Docstyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height/3)
	case tui.ActionResult[[]services.FoundService]:
		if msg.IsFor(&m.fetchServicesAction) {
			m.fetchServicesAction.Result = msg.Result
			m.fetchServicesAction.Error = msg.Error
			m.fetchServicesAction.Finished = true
			m.fetchServicesAction.Data = msg.Data
			m.list.SetItems(servicesToListItem(*m.fetchServicesAction.Data, nil))
		}
	case tui.ActionResult[roveryaml.RoverConfig]:
		if msg.IsFor(&m.fetchConfigurationAction) {
			m.fetchConfigurationAction.Result = msg.Result
			m.fetchConfigurationAction.Error = msg.Error
			m.fetchConfigurationAction.Finished = true
			m.fetchConfigurationAction.Data = msg.Data
			m.list.SetItems(servicesToListItem(nil, m.fetchConfigurationAction.Data))
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.MarkActive):
			if m.list.Index() >= 0 && m.list.Index() < len(m.list.Items()) {
				item := m.list.Items()[m.list.Index()].(item)
				// We can only have one active service with this name, so if there is another service with the same name but a different path, show an error
				for _, other := range *m.fetchServicesAction.Data {
					if other.Service.Name == item.service.Service.Name && other.Path != item.service.Path && m.fetchConfigurationAction.Data.HasEnabled(other.Path) {
						m.error = fmt.Errorf("A service with the name '%s' is already active. Services must be unique.", item.service.Service.Name)
						return m, nil
					}
				}
				m.error = nil

				m.fetchConfigurationAction.Data.Toggle(item.service.Path)
				m.list.SetItems(servicesToListItem(*m.fetchServicesAction.Data, m.fetchConfigurationAction.Data))
				return m, nil
			}
		case key.Matches(msg, keys.New):
			state.Get().Route.Push("connection init")
			return m, tea.Quit
		case key.Matches(msg, keys.Delete):
			if len(m.list.Items()) > 1 && m.list.Index() >= 0 && m.list.Index() < len(m.list.Items()) {
				// todo:
				// item := m.list.Items()[m.list.Index()].(item)
				// state.Get().RoverConnections = state.Get().RoverConnections.Remove(item.connection.Name)
				// m.list.SetItems(connectionsToListItems())
				// m.list.ResetSelected()
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	// Base command
	model, cmd := tui.Update(m, msg)
	if cmd != nil {
		return model, cmd
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, getServices(m), getConfiguration(m))
}

func (m model) View() string {
	s := ""

	if m.fetchServicesAction.IsLoading() {
		s += m.spinner.View() + " Fetching services..."
	}

	if m.fetchConfigurationAction.IsLoading() {
		s += "\n" + m.spinner.View() + " Fetching configuration..."
	}

	if m.fetchServicesAction.IsSuccess() && m.fetchConfigurationAction.IsSuccess() {
		s += m.list.View()
	} else if m.fetchServicesAction.IsError() {
		s += "Failed to fetch services"
		if m.fetchServicesAction.Error != nil {
			s += "\n > " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render(m.fetchServicesAction.Error.Error())
		}
	} else if m.fetchConfigurationAction.IsError() {
		s += "Failed to fetch configuration"
		if m.fetchConfigurationAction.Error != nil {
			s += "\n > " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render(m.fetchConfigurationAction.Error.Error())
		}
	}

	if m.error != nil {
		s += "\n\n" + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render(m.error.Error())
	}

	if m.fetchServicesAction.IsSuccess() && m.fetchConfigurationAction.IsSuccess() {
		// From all the services, create a dot graph of the pipeline
		dotgraph := "graph LR"

		// For every service, add a connection if there is a service that depends on it
		for _, found := range *m.fetchServicesAction.Data {
			if !m.fetchConfigurationAction.Data.HasEnabled(found.Path) {
				continue
			}
			for _, outputStream := range found.Service.Outputs {
				// Go over all other services
				for _, other := range *m.fetchServicesAction.Data {
					if !m.fetchConfigurationAction.Data.HasEnabled(other.Path) || found.Path == other.Path {
						continue
					}

					// Does this service depend on the current service?
					for _, input := range other.Service.Inputs {
						for _, inputStream := range input.Streams {
							if inputStream == outputStream && input.Service == found.Service.Name {
								// Add a connection
								connection := found.Service.Name + " --> " + other.Service.Name
								if !strings.Contains(dotgraph, connection) {
									dotgraph += "\n" + connection
								}
							}
						}
					}
				}
			}
		}

		s += "\n\n"

		// Draw the pipeline
		pipelineDrawing, err := asciitool.Draw(dotgraph)
		if err != nil {
			s += "Your pipeline could not be drawn. Please check the configuration."
		} else {
			s += pipelineDrawing
		}
	}

	return style.Docstyle.Render(s)
}

func getServices(m model) tea.Cmd {
	return tui.PerformAction(&m.fetchServicesAction, func() (*[]services.FoundService, error) {
		conn := state.Get().RoverConnections.GetActive()
		if conn == nil {
			return nil, fmt.Errorf("Not connected to an active Rover")
		}

		found := []services.FoundService{}
		err := roverlock.WithLock(*conn, func() error {
			var err error
			// Get all the services
			found, err = services.Scan(*conn)
			return err
		})

		return &found, err
	})
}

func getConfiguration(m model) tea.Cmd {
	return tui.PerformAction(&m.fetchConfigurationAction, func() (*roveryaml.RoverConfig, error) {
		conn := state.Get().RoverConnections.GetActive()
		if conn == nil {
			return nil, fmt.Errorf("Not connected to an active Rover")
		}

		found := &roveryaml.RoverConfig{}
		err := roverlock.WithLock(*conn, func() error {
			var err error
			// Get the configuration
			found, err = roveryaml.Load(*conn)
			return err
		})

		return found, err
	})
}
