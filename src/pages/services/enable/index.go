package enableservicespage

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/VU-ASE/rover/src/configuration"
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
	"github.com/lempiy/dgraph"
	"github.com/lempiy/dgraph/core"
)

type model struct {
	keys                     keyMap
	spinner                  spinner.Model
	list                     list.Model
	fetchServicesAction      tui.Action[[]services.FoundService]
	fetchConfigurationAction tui.Action[roveryaml.RoverConfig]
	saveConfigurationAction  tui.Action[any]
	error                    error // Can be shown to the user
}

// keyMap defines a set of keybindings. To work for help it must satisfy key.Map
type keyMap struct {
	MarkActive key.Binding
	Save       key.Binding
	Reload     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.MarkActive, k.Save, k.Reload}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var keys = keyMap{
	MarkActive: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle service"),
	),
	Save: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "save to Rover"),
	),
	Reload: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reload"),
	),
}

// List item to render
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

	shortPath := strings.Replace(i.service.Path, configuration.RemoteServiceDir+"/", "", 1)

	str := i.service.Service.Name + " " + lipgloss.NewStyle().Foreground(style.AsePrimary).Render(i.service.Service.Version) + " " + lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("("+shortPath+")")

	fn := lipgloss.NewStyle().Render
	if index == m.Index() {
		fn = func(s ...string) string {
			if i.active {
				return lipgloss.NewStyle().Bold(true).Foreground(style.SuccessPrimary).Render("> " + strings.Join(s, " "))
			} else {
				return lipgloss.NewStyle().Bold(true).Render("> " + strings.Join(s, " "))
			}
		}
	} else if i.active {
		str = style.RenderColor("✓ ", style.SuccessPrimary) + str
	} else {
		str = "- " + str
	}

	fmt.Fprint(w, fn(str))
}

type serviceDependency struct {
	service string
	stream  string
}

func getUnmetDependencies(service services.FoundService, enabled []services.FoundService) []serviceDependency {
	dependencies := make([]serviceDependency, 0)
	for _, dependency := range service.Service.Inputs {
		for _, stream := range dependency.Streams {
			dependencies = append(dependencies, serviceDependency{
				service: dependency.Service,
				stream:  stream,
			})
		}
	}

	for _, dependency := range service.Service.Inputs {
		// Go over all other service
		for _, other := range enabled {
			// Is this the service we are looking for?
			if dependency.Service == other.Service.Name {
				// Are all the streams available?
				for _, stream := range dependency.Streams {
					if slices.Contains(other.Service.Outputs, stream) {
						// Remove the dependency
						dependencies = slices.DeleteFunc(dependencies, func(d serviceDependency) bool {
							return d.service == dependency.Service && d.stream == stream
						})
					}
				}
			}
		}
	}

	return dependencies
}

// Returns the errors in the configuration, if none are found, the configuration is valid
func configurationValid(config *roveryaml.RoverConfig, allservices []services.FoundService) []error {
	enabledServices := make([]services.FoundService, 0)
	for _, service := range allservices {
		if config.HasEnabled(service.Path) {
			enabledServices = append(enabledServices, service)
		}
	}

	errors := make([]error, 0)
	for _, service := range enabledServices {
		unmet := getUnmetDependencies(service, enabledServices)
		for _, dep := range unmet {
			errors = append(errors, fmt.Errorf("Service '%s' depends on service '%s' for stream '%s' which is not enabled", service.Service.Name, dep.service, dep.stream))
		}
	}

	return errors
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
	l.Title = "Configure Rover pipeline"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = style.TitleStyle
	l.Styles.PaginationStyle = style.PaginationStyle
	l.Styles.HelpStyle = style.HelpStyle
	l.AdditionalShortHelpKeys = keys.ShortHelp

	spin := spinner.New()
	spin.Spinner = spinner.Line

	fetchServicesAction := tui.NewAction[[]services.FoundService]("getServices")
	fetchServicesAction.Start()

	fetchConfigAction := tui.NewAction[roveryaml.RoverConfig]("getConfig")
	fetchConfigAction.Start()

	return model{
		keys:                     keys,
		list:                     l,
		spinner:                  spin,
		fetchServicesAction:      fetchServicesAction,
		fetchConfigurationAction: fetchConfigAction,
		saveConfigurationAction:  tui.NewAction[any]("saveConfig"),
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
			allservices := msg.Data
			if allservices != nil {
				m.fetchServicesAction.Data = msg.Data
				m.list.SetItems(servicesToListItem(*allservices, m.fetchConfigurationAction.Data))
			}
		}
	case tui.ActionResult[roveryaml.RoverConfig]:
		if msg.IsFor(&m.fetchConfigurationAction) {
			m.fetchConfigurationAction.Result = msg.Result
			m.fetchConfigurationAction.Error = msg.Error
			m.fetchConfigurationAction.Finished = true

			config := msg.Data
			if config != nil {
				m.fetchConfigurationAction.Data = config
				m.list.SetItems(servicesToListItem(nil, config))
			}
		}
	case tui.ActionResult[any]:
		if msg.IsFor(&m.saveConfigurationAction) {
			m.saveConfigurationAction.Result = msg.Result
			m.saveConfigurationAction.Error = msg.Error
			m.saveConfigurationAction.Finished = true
			m.error = nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Reload):
			if !m.fetchServicesAction.IsLoading() && !m.fetchConfigurationAction.IsLoading() && !m.saveConfigurationAction.IsLoading() {
				m = InitialModel()
				return m, m.Init()
			}
		case key.Matches(msg, keys.MarkActive):
			if m.list.Index() >= 0 && m.list.Index() < len(m.list.Items()) {
				m.saveConfigurationAction.Reset()
				m.error = nil

				item := m.list.Items()[m.list.Index()].(item)
				// We can only have one active service with this name, so if there is another service with the same name but a different path, show an error
				for _, other := range *m.fetchServicesAction.Data {
					if other.Service.Name == item.service.Service.Name && other.Path != item.service.Path && m.fetchConfigurationAction.Data.HasEnabled(other.Path) {
						m.error = fmt.Errorf("A service with the name '%s' is already active. Services must be unique.", item.service.Service.Name)
						return m, nil
					}
				}

				m.fetchConfigurationAction.Data.Toggle(item.service.Path)
				m.list.SetItems(servicesToListItem(*m.fetchServicesAction.Data, m.fetchConfigurationAction.Data))
				return m, nil
			}
		case key.Matches(msg, keys.Save):
			if m.saveConfigurationAction.IsLoading() {
				return m, nil
			}

			var configErrors []error
			if m.fetchServicesAction.IsSuccess() && m.fetchConfigurationAction.IsSuccess() {
				configErrors = configurationValid(m.fetchConfigurationAction.Data, *m.fetchServicesAction.Data)
			} else {
				m.error = fmt.Errorf("Cannot save configuration, services or configuration not loaded")
				return m, nil
			}

			// Remove enabled services that are not in the list of available services
			for _, enabled := range m.fetchConfigurationAction.Data.Enabled {
				if !slices.ContainsFunc(*m.fetchServicesAction.Data, func(f services.FoundService) bool {
					return f.Path == enabled
				}) {
					m.fetchConfigurationAction.Data.Disable(enabled)
				}
			}

			if len(configErrors) > 0 {
				errorString := "Cannot save configuration because it is invalid:"
				maxHeight := 3
				for _, err := range configErrors[:min(len(configErrors), maxHeight)] {
					errorString += "\n > " + err.Error()
				}
				if len(configErrors) > maxHeight {
					errorString += "\n > ..."
				}

				m.error = fmt.Errorf(errorString)
				return m, nil
			}
			m.error = nil
			m.saveConfigurationAction.Start()
			return m, saveConfiguration(m)
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

func configureView(m model) string {
	s := m.list.View()
	if m.error != nil {
		s += "\n" + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render(m.error.Error()) + "\n"
	}

	if m.saveConfigurationAction.IsSuccess() {
		s += "\n" + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Configuration saved successfully to Rover") + "\n"
	}

	// From all the services, create a dot graph of the pipeline
	nodes := make([]core.NodeInput, 0)

	// Shorthands
	config := m.fetchConfigurationAction.Data
	allservices := *m.fetchServicesAction.Data
	enabledservices := make([]services.FoundService, 0)
	for _, service := range allservices {
		if config.HasEnabled(service.Path) {
			enabledservices = append(enabledservices, service)
		}
	}

	// For every service, add a connection if there is a service that depends on it
	for _, found := range enabledservices {
		newNode := core.NodeInput{
			Id:   found.Service.Name,
			Next: make([]string, 0),
		}

		for _, outputStream := range found.Service.Outputs {
			// Go over all other services
			for _, other := range enabledservices {
				if found.Path == other.Path {
					continue
				}

				// Does this service depend on the current service?
				for _, input := range other.Service.Inputs {
					for _, inputStream := range input.Streams {
						if inputStream == outputStream && input.Service == found.Service.Name {
							// Add a connection
							newNode.Next = append(newNode.Next, other.Service.Name)
						}
					}
				}
			}
		}

		nodes = append(nodes, newNode)
	}

	// Draw the pipeline
	canvas, err := dgraph.DrawGraph(nodes)
	canvasView := ""
	if err != nil {
		canvasView += "Failed to draw pipeline"
	} else if len(nodes) > 0 {
		canvasView += fmt.Sprintf("\n%s\n", canvas)
	}

	// Add arrows
	canvasView = strings.ReplaceAll(canvasView, "─┤", ">┤")

	// Show unmet dependencies in the pipeline with a red > symbol
	for _, service := range enabledservices {
		unmet := getUnmetDependencies(service, enabledservices)
		if len(unmet) > 0 {
			canvasView = strings.Replace(canvasView, ">┤ "+service.Service.Name, style.RenderColor(">┤ ", style.ErrorPrimary)+service.Service.Name, -1)
			canvasView = strings.Replace(canvasView, "│ "+service.Service.Name, style.RenderColor("> ", style.ErrorPrimary)+service.Service.Name, -1)
		}
	}

	// If the current list item is active, highlight it
	currItem := m.list.Items()[m.list.Index()].(item)
	if currItem.active {
		unmet := getUnmetDependencies(currItem.service, enabledservices)
		color := style.SuccessPrimary
		if len(unmet) > 0 {
			color = style.ErrorPrimary
		}
		canvasView = strings.Replace(canvasView, currItem.service.Service.Name, lipgloss.NewStyle().Foreground(color).Bold(true).Render(currItem.service.Service.Name), -1)
	}

	s += lipgloss.NewStyle().Margin(0, 1).Render(canvasView)
	return s
}

func loadErrorView(m model) string {
	s := ""

	if m.fetchServicesAction.IsError() {
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

	return s
}

func loadindView(m model) string {
	s := ""

	if m.fetchServicesAction.IsLoading() {
		s += m.spinner.View() + " Fetching services...\n"
	}

	if m.fetchConfigurationAction.IsLoading() {
		s += m.spinner.View() + " Fetching configuration..."
	}

	return s
}

func saveErrorView(m model) string {
	s := ""

	if m.saveConfigurationAction.IsError() {
		s += "Failed to save configuration"
		if m.saveConfigurationAction.Error != nil {
			s += "\n > " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render(m.saveConfigurationAction.Error.Error())
		}
	}

	s += "\n\n" + style.RenderColor("Press 'r' to reload the configuration or press 'q' to quit", style.GrayPrimary)

	return s
}

func savingView(m model) string {
	return m.spinner.View() + " Saving configuration..."
}

func (m model) View() string {
	s := ""

	if m.fetchServicesAction.IsLoading() || m.fetchConfigurationAction.IsLoading() {
		s = loadindView(m)
	} else if m.fetchServicesAction.IsError() || m.fetchConfigurationAction.IsError() {
		s = loadErrorView(m)
	} else if m.saveConfigurationAction.IsLoading() {
		s = savingView(m)
	} else if m.saveConfigurationAction.IsError() {
		s = saveErrorView(m)
	} else {
		s = configureView(m)
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

func saveConfiguration(m model) tea.Cmd {
	return tui.PerformAction(&m.saveConfigurationAction, func() (*any, error) {
		// Shorthand
		config := m.fetchConfigurationAction.Data
		if config == nil {
			return nil, fmt.Errorf("No configuration loaded")
		}

		conn := state.Get().RoverConnections.GetActive()
		if conn == nil {
			return nil, fmt.Errorf("Not connected to an active Rover")
		}

		err := roverlock.WithLock(*conn, func() error {
			// Check if all enabled services still exist
			current, err := services.Scan(*conn)
			if err != nil {
				return err
			}
			for _, enabled := range config.Enabled {
				if !slices.ContainsFunc(current, func(f services.FoundService) bool {
					return f.Path == enabled
				}) {
					return fmt.Errorf("Service '%s' does not exist anymore", enabled)
				}
			}

			// Try to save the configuration
			return config.Save(*conn)
		})

		return nil, err
	})
}
