package views

import (
	"fmt"
	"regexp"
	"time"

	"github.com/VU-ASE/rover/src/openapi"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lempiy/dgraph"
	"github.com/lempiy/dgraph/core"
)

//
// All keys
//

// Keys to navigate
type PipelineConfiguratorKeyMap struct {
	Retry   key.Binding
	Confirm key.Binding
	Switch  key.Binding // switch table focus
	Remove  key.Binding // remove service from pipeline
	Back    key.Binding // go back one level
	Save    key.Binding // save the pipeline
	Quit    key.Binding
	Refetch key.Binding
}

// Shown when the services are being updated
var pipelineConfiguratorKeysRegular = PipelineConfiguratorKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "refetch"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "go back one level"),
	),
	Switch: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch table"),
	),
	Remove: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "remove"),
	),
	Save: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "save"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

// When the "active" table is focussed
var pipelineConfiguratorKeysActiveTable = PipelineConfiguratorKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "refetch"),
	),
	Switch: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch table"),
	),
	Remove: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "remove from pipeline"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

// When the "remote" table is focussed on the "select an author" level
var pipelineConfiguratorKeysRemoteTableAuthor = PipelineConfiguratorKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "refetch"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view services"),
	),
	Switch: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch table"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

// When the "remote" table is focussed on the "select a service" level
var pipelineConfiguratorKeysRemoteTableService = PipelineConfiguratorKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "refetch"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view versions"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "view all authors"),
	),
	Switch: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch table"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

// When the "remote" table is focussed on the "select a version" level
var pipelineConfiguratorKeysRemoteTableVersion = PipelineConfiguratorKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "refetch"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "add to pipeline"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "view all services"),
	),
	Switch: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch table"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func (k PipelineConfiguratorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Retry, k.Confirm, k.Remove, k.Back, k.Switch, k.Quit}
}

func (k PipelineConfiguratorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

//
// The page model
//

type PipelineConfiguratorPage struct {
	help    help.Model
	spinner spinner.Model
	// tables that are shown next to each other
	tableActive table.Model // the pipeline as it is configured now, with the enabled services
	tableRemote table.Model // all remote services that still can be enabled
	// actions to fetch data to populate the tables
	pipeline         tui.Action[PipelineOverviewSummary]
	pipelineGraph    string               // preserved in the model to avoid re-rendering in .View()
	dependencyErrors []error              // errors in the pipeline configuration
	authors          tui.Action[[]string] // first part of FQN
	services         tui.Action[[]string] // second part of FQN
	versions         tui.Action[[]string] // third part of FQN
	savePipeline     tui.Action[bool]     // save the pipeline
	// Keep track of the focussed table (left/right)
	focussed int // 0 = active, 1 = remote
	// For querying remote services
	remoteAuthor  string
	remoteService string
}

func NewPipelineConfiguratorPage() PipelineConfiguratorPage {
	// todo

	return PipelineConfiguratorPage{
		spinner:       spinner.New(),
		help:          help.New(),
		tableActive:   table.New(),
		tableRemote:   table.New(),
		pipeline:      tui.NewAction[PipelineOverviewSummary]("fetchActive"),
		authors:       tui.NewAction[[]string]("fetchAuthors"),
		services:      tui.NewAction[[]string]("fetchServices"),
		versions:      tui.NewAction[[]string]("fetchVersions"),
		savePipeline:  tui.NewAction[bool]("savePipeline"),
		focussed:      0,
		remoteAuthor:  "",
		remoteService: "",
	}
}

//
// Page model methods
//

func (m PipelineConfiguratorPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.pipeline.HasData() {
			m.tableActive = m.createActiveTable(*m.pipeline.Data)
		} else {
			m.tableActive = m.createActiveTable(PipelineOverviewSummary{})
		}
		m.tableRemote = m.createRemoteTable()
		return m, nil
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tui.ActionInit[bool]:
		m.savePipeline.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[bool]:
		m.savePipeline.ProcessResult(msg)
		return m, nil
	case tui.ActionInit[PipelineOverviewSummary]:
		m.pipeline.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[PipelineOverviewSummary]:
		m.pipeline.ProcessResult(msg)
		if m.pipeline.IsSuccess() {
			if len(m.pipeline.Data.Pipeline.Enabled) <= 0 {
				m.focussed = 1 // nothing to navigate in the active table
				m.tableRemote = m.createRemoteTable()
			}

			// Create the pipeline graph based on enabled services
			nodes := make([]core.NodeInput, 0)
			for _, service := range m.pipeline.Data.Pipeline.Enabled {
				// Check if the service is selected, in this case unselect it
				if m.remoteService == service.Service.Name {
					m.remoteService = ""
				}

				nodes = append(nodes, core.NodeInput{
					Id: service.Service.Name,
					Next: func() []string {
						// Find services that depend on an output of this service
						found := make([]string, 0)
						for _, s := range m.pipeline.Data.Services {
							if s.Name != service.Service.Name {
								for _, input := range s.Configuration.Inputs {
									if input.Service == service.Service.Name {
										found = append(found, s.Name)
									}
								}
							}
						}

						return found
					}(),
				})
			}
			canvas, err := dgraph.DrawGraph(nodes)
			if len(nodes) <= 0 {
				m.pipelineGraph = style.Gray.Render("No pipeline to show! Go ahead and add some services.")
			} else if err != nil {
				m.pipelineGraph = "Failed to draw pipeline\n"
			} else {
				m.pipelineGraph = fmt.Sprintf("%s\n", canvas)
			}
			m.dependencyErrors = m.findDependencyErrors()
			m.tableActive = m.createActiveTable(*m.pipeline.Data)
			m.tableRemote = m.createRemoteTable()
		}
		m.savePipeline = tui.NewAction[bool]("savePipeline")
		return m, nil
	case tui.ActionInit[[]string]:
		m.authors.ProcessInit(msg)
		m.services.ProcessInit(msg)
		m.versions.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[[]string]:
		m.authors.ProcessResult(msg)
		m.services.ProcessResult(msg)
		m.versions.ProcessResult(msg)
		m.tableRemote = m.createRemoteTable()
		return m, nil
	case tea.KeyMsg:
		switch {
		// case key.Matches(msg, pipelineConfiguratorKeysRegular.Quit):
		// return m, tea.Quit
		case key.Matches(msg, pipelineConfiguratorKeysRegular.Remove):
			if m.focussed == 0 {
				return m.onActiveTableNavigation(msg)
			}

			// todo:
			return m, nil
		case key.Matches(msg, pipelineConfiguratorKeysRegular.Confirm):
			if m.focussed == 1 {
				return m.onRemoteTableNavigation(msg)
			}

			// todo:
			return m, nil
		case key.Matches(msg, pipelineConfiguratorKeysRegular.Back):
			if m.focussed == 1 {
				return m.onRemoteTableNavigation(msg)
			}

			// todo:
			return m, nil
		case key.Matches(msg, pipelineConfiguratorKeysRegular.Retry):
			// Redo all actions, reset to the initial state
			cmds := tea.Batch(
				m.fetchAllAuthors(),
				m.fetchPipeline(),
			)
			return m, cmds
		case key.Matches(msg, pipelineConfiguratorKeysRegular.Save):
			if len(m.dependencyErrors) <= 0 {
				return m, m.savePipelineRemote()
			}
		case key.Matches(msg, pipelineConfiguratorKeysRegular.Switch):
			m.focussed = (m.focussed + 1) % 2
			m.tableActive = m.createActiveTable(*m.pipeline.Data)
			m.tableRemote = m.createRemoteTable()
			return m, nil
		}
	}

	if m.focussed == 0 {
		m.tableActive, cmd = m.tableActive.Update(msg)
	} else if m.focussed == 1 {
		m.tableRemote, cmd = m.tableRemote.Update(msg)
	}

	return m, cmd
}

func (m PipelineConfiguratorPage) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchPipeline(), m.fetchAllAuthors())
}

func (m PipelineConfiguratorPage) remoteTableView() string {
	note := ""
	if m.remoteAuthor == "" {
		if m.authors.IsSuccess() {
			if len(m.tableRemote.Rows()) <= 0 {
				note = style.Gray.Render(" No remote authors available. Upload a service first.")
			}

			return m.tableRemote.View() + note
		} else if m.authors.IsError() {
			return style.Error.Render("Error loading available authors") + style.Gray.Render(" "+m.authors.Error.Error())
		} else {
			return m.spinner.View() + style.Gray.Render(" Loading authors")
		}
	} else if m.remoteService == "" {
		if m.services.IsSuccess() {
			if len(m.tableRemote.Rows()) <= 0 {
				note = style.Gray.Render(" No unused services available for " + m.remoteAuthor)
			}

			return m.tableRemote.View() + note
		} else if m.services.IsError() {
			return style.Error.Render("Error loading services for "+m.remoteAuthor) + style.Gray.Render(" "+m.services.Error.Error())
		} else {
			return m.spinner.View() + style.Gray.Render(" Loading services for "+m.remoteAuthor)
		}
	} else if m.versions.IsSuccess() {
		if len(m.tableRemote.Rows()) <= 0 {
			note = style.Gray.Render(" No unused versions available for " + m.remoteAuthor + "/" + m.remoteService)
		}

		return m.tableRemote.View() + note
	} else if m.versions.IsError() {
		return style.Error.Render("Error loading versions for "+m.remoteAuthor+"/"+m.remoteService) + style.Gray.Render(" "+m.versions.Error.Error())
	} else {
		return m.spinner.View() + style.Gray.Render(" Loading versions for "+m.remoteAuthor+"/"+m.remoteService)
	}
}

func (m PipelineConfiguratorPage) View() string {
	s := style.Title.Render("Configure your pipeline") + "\n\n"

	// Calculate column width (subtract padding and borders)
	columnWidth := m.getColWidth()

	// Define styles for each column
	columnStyle := lipgloss.NewStyle().Width(columnWidth)

	graph := ""
	if m.pipeline.IsSuccess() {
		graph = m.postProcessGraph(m.pipelineGraph)

		if len(m.pipeline.Data.Services) > 0 && len(m.dependencyErrors) <= 0 {
			graph += "\n" + style.Success.Render("✓ This pipeline is valid") + " " + style.Gray.Render("- save it to your Rover with 's'")
		} else if len(m.pipeline.Data.Services) > 0 {
			for _, err := range m.dependencyErrors {
				graph += "\n" + style.Error.Render("✗ "+err.Error())
			}
		}

	} else if m.pipeline.IsError() {
		graph = style.Error.Render("Error loading pipeline") + style.Gray.Render(" "+m.pipeline.Error.Error())
	} else {
		graph = m.spinner.View() + " Loading pipeline"
	}

	if m.savePipeline.IsLoading() {
		graph += "\n" + m.spinner.View() + " Saving pipeline"
	} else if m.savePipeline.IsSuccess() {
		graph += "\n" + style.Success.Render("✓ Pipeline saved to Rover successfully")
	} else if m.savePipeline.IsError() {
		graph += "\n" + style.Error.Render("✗ Error saving pipeline") + style.Gray.Render(" "+m.savePipeline.Error.Error())
	}
	graph += "\n\n"

	// Define the columns
	columnActive := m.spinner.View() + style.Gray.Render(" Loading active services...")
	if m.pipeline.IsSuccess() {
		note := ""
		if len(m.tableActive.Rows()) <= 0 {
			note = style.Gray.Render(" There are no enabled services in this pipeline, yet.")
		}

		columnActive = m.tableActive.View() + note
	} else if m.pipeline.IsError() {
		columnActive = style.Error.Render("Error loading active services") + style.Gray.Render(" "+m.pipeline.Error.Error())
	}

	columnRemote := m.remoteTableView()

	row := lipgloss.JoinHorizontal(lipgloss.Top,
		columnStyle.Render(columnActive),
		" ",
		columnStyle.Render(columnRemote),
	)

	h := ""
	if m.focussed == 0 {
		h = m.tableActive.HelpView() + style.Gray.Render(" • ") + m.help.View(pipelineConfiguratorKeysActiveTable)
	} else {
		h = m.tableRemote.HelpView() + style.Gray.Render(" • ")
		if m.remoteService != "" {
			h += m.help.View(pipelineConfiguratorKeysRemoteTableVersion)
		} else if m.remoteAuthor != "" {
			h += m.help.View(pipelineConfiguratorKeysRemoteTableService)
		} else {
			h += m.help.View(pipelineConfiguratorKeysRemoteTableAuthor)
		}
	}

	return s + graph + row + "\n\n" + h
}

//
// Aux for table and col rendering
//

func (m PipelineConfiguratorPage) getColWidth() int {
	return (state.Get().WindowWidth - 4 - 6) / 2 // Adjust for padding and borders
}

func (m PipelineConfiguratorPage) colPct(pct int) int {
	total := m.getColWidth() - 2
	return (total*pct)/100 - 1
}

//
// Actions
//

func (m PipelineConfiguratorPage) fetchPipeline() tea.Cmd {
	return tui.PerformAction(&m.pipeline, func() (*PipelineOverviewSummary, error) {
		// mock fetch
		// ! remove

		time.Sleep(200 * time.Millisecond)
		// First roverd tells us what services are enabled, by reference (FQN)
		pipeline := openapi.PipelineGet200Response{
			Status:    openapi.STARTED,
			LastStart: openapi.PtrInt64(123456),
			LastStop:  openapi.PtrInt64(123456),
			// LastRestart: openapi.PtrInt64(123456),
			Enabled: []openapi.PipelineGet200ResponseEnabledInner{
				{
					Service: openapi.PipelineGet200ResponseEnabledInnerService{
						Name:    "imaging",
						Version: "1.0.0",
						Author:  "vu-ase",
					},
				},
				{
					Service: openapi.PipelineGet200ResponseEnabledInnerService{
						Name:    "controller",
						Version: "1.0.0",
						Author:  "vu-ase",
					},
				},
				{
					Service: openapi.PipelineGet200ResponseEnabledInnerService{
						Name:    "transceiver",
						Version: "1.0.0",
						Author:  "vu-ase",
					},
				},
			},
		}

		// Then, for each service, we need to query the service for its actual configuration (inputs, outputs)
		services := make([]PipelineOverviewServiceInfo, 0)
		for _, enabled := range pipeline.Enabled {
			// mock fetch
			// ! remove

			if enabled.Service.Name == "imaging" {
				services = append(services, PipelineOverviewServiceInfo{
					Name:    enabled.Service.Name,
					Version: enabled.Service.Version,
					Author:  enabled.Service.Author,
					Configuration: openapi.ServicesAuthorServiceVersionGet200Response{
						Inputs: []openapi.ServicesAuthorServiceVersionGet200ResponseInputsInner{}, // no inputs
						Outputs: []string{
							"track",
						},
					},
				})
			} else if enabled.Service.Name == "controller" {
				services = append(services, PipelineOverviewServiceInfo{
					Name:    enabled.Service.Name,
					Version: enabled.Service.Version,
					Author:  enabled.Service.Author,
					Configuration: openapi.ServicesAuthorServiceVersionGet200Response{
						Inputs: []openapi.ServicesAuthorServiceVersionGet200ResponseInputsInner{
							{
								Service: "imaging",
								Streams: []string{
									"track",
								},
							},
						},
						Outputs: []string{}, // no outputs, last service
					},
				})
			} else if enabled.Service.Name == "transceiver" {
				services = append(services, PipelineOverviewServiceInfo{
					Name:    enabled.Service.Name,
					Version: enabled.Service.Version,
					Author:  enabled.Service.Author,
					Configuration: openapi.ServicesAuthorServiceVersionGet200Response{
						Inputs: []openapi.ServicesAuthorServiceVersionGet200ResponseInputsInner{
							{
								Service: "imaging",
								Streams: []string{
									"track",
								},
							},
							{
								Service: "controller yo",
								Streams: []string{
									"track",
								},
							},
						},
						Outputs: []string{}, // no outputs, last service
					},
				})
			}
		}

		// Then the status (mock data)
		status := openapi.StatusGet200Response{
			Cpu: []openapi.StatusGet200ResponseCpuInner{
				{
					Core:  0,
					Used:  5,
					Total: 10,
				},
				{
					Core:  1,
					Used:  2,
					Total: 10,
				},
			},
			Memory: openapi.StatusGet200ResponseMemory{
				Total: 100,
				Used:  50,
			},
		}

		// Combined response
		res := PipelineOverviewSummary{
			Pipeline: pipeline,
			Services: services,
			Status:   status,
		}

		return &res, nil
	})
}

//
// Tables
//

func (m PipelineConfiguratorPage) createActiveTable(res PipelineOverviewSummary) table.Model {
	// Retrieve the previously selected entry
	prev := m.tableActive.SelectedRow()

	columns := []table.Column{
		{Title: "Service", Width: m.colPct(30)},
		{Title: "Version", Width: m.colPct(30)},
		{Title: "Author", Width: m.colPct(39)},
	}

	rows := make([]table.Row, 0)
	for _, service := range res.Services {
		rows = append(rows, table.Row{
			service.Name,
			service.Version,
			service.Author,
		})
	}

	// Find the previously selected row, if it exists
	cursor := 0
	if prev != nil {
		for i, row := range rows {
			if len(row) >= 3 && len(prev) >= 3 && row[0] == prev[0] && row[1] == prev[1] && row[2] == prev[2] {
				cursor = i
				break
			}
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)
	if len(rows) < 7 {
		t.SetHeight(len(rows) + 1)
	}
	t.SetCursor(cursor)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	if m.focussed == 0 {
		t.Focus()
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("FFF")).
			Background(style.AsePrimary).
			Bold(false)
	} else {
		t.Blur()
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("FFF")).
			Background(style.GrayPrimary).
			Bold(false)
	}

	t.SetStyles(s)

	return t
}

func (m PipelineConfiguratorPage) createRemoteTable() table.Model {
	// Retrieve the previously selected entry
	prev := m.tableRemote.SelectedRow()

	columns := []table.Column{
		{Title: "Author", Width: m.colPct(100)},
	}
	rows := make([]table.Row, 0)
	// Go from most fine-grained to least fine-grained
	if m.remoteService != "" {
		columns = []table.Column{
			{Title: fmt.Sprintf("Versions (%s/%s)", m.remoteAuthor, m.remoteService), Width: m.colPct(100)},
		}
		if m.versions.HasData() {
			for _, version := range *m.versions.Data {
				// Does this version already exist in the pipeline?
				exists := false
				if m.pipeline.HasData() {
					for _, enabled := range m.pipeline.Data.Pipeline.Enabled {
						if enabled.Service.Name == m.remoteService && enabled.Service.Version == version {
							exists = true
							break
						}
					}
				}

				if !exists {
					rows = append(rows, table.Row{
						version,
					})
				}
			}
		}
	} else if m.remoteAuthor != "" {
		columns = []table.Column{
			{Title: fmt.Sprintf("Services (%s)", m.remoteAuthor), Width: m.colPct(100)},
		}
		if m.services.HasData() {
			for _, service := range *m.services.Data {
				// Does this service already exist in the pipeline?
				exists := false
				if m.pipeline.HasData() {
					for _, enabled := range m.pipeline.Data.Pipeline.Enabled {
						if enabled.Service.Name == service {
							exists = true
							break
						}
					}
				}

				if !exists {
					rows = append(rows, table.Row{
						service,
					})
				}
			}
		}

	} else {
		if m.authors.HasData() {
			for _, author := range *m.authors.Data {
				rows = append(rows, table.Row{
					author,
				})
			}
		}
	}

	// Find the previously selected row, if it exists
	cursor := 0
	if prev != nil {
		for i, row := range rows {
			if len(row) >= 1 && len(prev) >= 1 && row[0] == prev[0] {
				cursor = i
				break
			}
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)
	if len(rows) < 7 {
		t.SetHeight(len(rows) + 1)
	}
	t.SetCursor(cursor)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	if m.focussed == 1 {
		t.Focus()
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("FFF")).
			Background(style.AsePrimary).
			Bold(false)
	} else {
		t.Blur()
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("FFF")).
			Background(style.GrayPrimary).
			Bold(false)
	}

	t.SetStyles(s)

	return t
}

//
// Actions
//

func (m PipelineConfiguratorPage) fetchAllAuthors() tea.Cmd {
	return tui.PerformAction(&m.authors, func() (*[]string, error) {
		// mock fetch
		// ! remove
		time.Sleep(200 * time.Millisecond)

		return &[]string{
			"vu-ase",
			"ielaajezdev",
			"maxgallup",
		}, nil
	})
}

func (m PipelineConfiguratorPage) fetchServicesForAuthor(author string) tea.Cmd {
	return tui.PerformAction(&m.services, func() (*[]string, error) {
		// mock fetch
		// ! remove
		time.Sleep(200 * time.Millisecond)

		res := []string{
			"lux",
			"controller",
			"actuator",
		}

		return &res, nil
	})
}

func (m PipelineConfiguratorPage) fetchVersionsForService(author string, service string) tea.Cmd {
	return tui.PerformAction(&m.versions, func() (*[]string, error) {
		// mock fetch
		// ! remove
		time.Sleep(200 * time.Millisecond)

		res := []string{
			"1.0.0",
			"1.0.1",
			"1.0.2",
		}

		return &res, nil
	})
}

func (m PipelineConfiguratorPage) onActiveTableNavigation(pressedKey tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.focussed != 0 {
		return m, nil
	}

	if key.Matches(pressedKey, pipelineConfiguratorKeysRegular.Remove) {
		// If there is no value, nothing we can do
		sel := m.tableActive.SelectedRow()
		if len(sel) <= 0 {
			return m, nil
		}

		// Remove service from pipeline
		return m, m.removeServiceFromPipeline(sel[2], sel[0], sel[1])
	}

	return m, nil
}

func (m PipelineConfiguratorPage) onRemoteTableNavigation(pressedKey tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.focussed != 1 {
		return m, nil
	}

	// Hitting "enter" makes the lookup go one level deeper, unless we are at the deepest level (a specific version), which will then insert the service into the pipeline
	if key.Matches(pressedKey, pipelineConfiguratorKeysRegular.Confirm) {
		// If there is no value, nothing we can do
		sel := m.tableRemote.SelectedRow()
		if len(sel) <= 0 {
			return m, nil
		}

		if m.remoteAuthor == "" {
			m.remoteAuthor = sel[0]
			return m, m.fetchServicesForAuthor(m.remoteAuthor)
		} else if m.remoteService == "" {
			m.remoteService = sel[0]
			return m, m.fetchVersionsForService(m.remoteAuthor, m.remoteService)
		} else {
			// Insert service into pipeline
			return m, m.addServiceToPipeline(m.remoteAuthor, m.remoteService, sel[0])
		}
	}

	// Hitting "backspace" goes one level up
	if key.Matches(pressedKey, pipelineConfiguratorKeysRegular.Back) {
		if m.remoteService != "" {
			m.remoteService = ""
		} else if m.remoteAuthor != "" {
			m.remoteAuthor = ""
		}
	}

	m.tableRemote = m.createRemoteTable()
	return m, nil
}

// This adds a service to a pipeline *locally*. It will only be checked by the server when the pipeline is saved.
func (m PipelineConfiguratorPage) addServiceToPipeline(author string, service string, version string) tea.Cmd {
	return tui.PerformAction(&m.pipeline, func() (*PipelineOverviewSummary, error) {
		// There should already be a pipeline in the model
		if !m.pipeline.IsSuccess() {
			return nil, fmt.Errorf("Cannot add a service to a non-fetched pipeline")
		}

		// Fetch the specific service data
		// mock fetch
		// ! remove

		res := openapi.ServicesAuthorServiceVersionGet200Response{
			BuiltAt: openapi.PtrInt64(123456),
			Inputs: []openapi.ServicesAuthorServiceVersionGet200ResponseInputsInner{
				{
					Service: "imaging",
					Streams: []string{
						"track",
						"nonex",
					},
				},
			},
			Outputs: []string{
				"lux",
			},
		}

		// Add this service to the pipeline
		pipeline := *m.pipeline.Data

		// Check if the service is already in the pipeline
		for _, enabled := range pipeline.Pipeline.Enabled {
			if enabled.Service.Name == service {
				return &pipeline, nil
			}
		}

		pipeline.Services = append(pipeline.Services, PipelineOverviewServiceInfo{
			Name:          service,
			Author:        author,
			Version:       version,
			Configuration: res,
		})
		pipeline.Pipeline.Enabled = append(pipeline.Pipeline.Enabled, openapi.PipelineGet200ResponseEnabledInner{
			Service: openapi.PipelineGet200ResponseEnabledInnerService{
				Name:    service,
				Version: version,
				Author:  author,
			},
		})

		return &pipeline, nil
	})
}

// This removes a service from a pipeline *locally*. It will only be checked by the server when the pipeline is saved.
func (m PipelineConfiguratorPage) removeServiceFromPipeline(author string, service string, version string) tea.Cmd {
	return tui.PerformAction(&m.pipeline, func() (*PipelineOverviewSummary, error) {
		// There should already be a pipeline in the model
		if !m.pipeline.IsSuccess() {
			return nil, fmt.Errorf("Cannot remove a service from a non-fetched pipeline")
		}

		// Remove this service from the pipeline
		pipeline := *m.pipeline.Data
		newServices := make([]PipelineOverviewServiceInfo, 0)
		for _, s := range pipeline.Services {
			if s.Name != service {
				newServices = append(newServices, s)
			}
		}
		pipeline.Services = newServices
		newEnabled := make([]openapi.PipelineGet200ResponseEnabledInner, 0)
		for _, enabled := range pipeline.Pipeline.Enabled {
			if enabled.Service.Name != service {
				newEnabled = append(newEnabled, enabled)
			}
		}
		pipeline.Pipeline.Enabled = newEnabled
		return &pipeline, nil
	})
}

func (m PipelineConfiguratorPage) findDependencyErrors() []error {
	errors := make([]error, 0)
	if !m.pipeline.IsSuccess() {
		return errors
	}

	// For each service, check if it has unmet dependencies with other services
	for _, service := range m.pipeline.Data.Services {
		for _, input := range service.Configuration.Inputs {
			for _, stream := range input.Streams {
				found := false
				for _, other := range m.pipeline.Data.Services {
					if other.Name == input.Service {
						for _, output := range other.Configuration.Outputs {
							if output == stream {
								found = true
								break
							}
						}
					}
				}

				if !found {
					errors = append(errors, fmt.Errorf("Service '%s' depends on unresolved stream '%s' from service '%s'", service.Name, stream, input.Service))
				}

			}
		}
	}

	return errors
}

func (m PipelineConfiguratorPage) savePipelineRemote() tea.Cmd {
	return tui.PerformAction(&m.savePipeline, func() (*bool, error) {

		// mock save
		// ! remove

		time.Sleep(200 * time.Millisecond)

		return openapi.PtrBool(true), nil
	})
}

//
// Aux methods for views
//

// Clean up the graph to make it a bit more readable and compressed
func (m PipelineConfiguratorPage) postProcessGraph(s string) string {
	n := s

	// Remove empty lines
	n = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(n, "\n")

	// Highlight the currently selected service
	sel := m.tableActive.SelectedRow()
	if sel != nil {
		// The first item is always the service name
		name := sel[0]

		// Find the service in the graph
		if m.focussed == 0 {
			n = regexp.MustCompile(fmt.Sprintf(`\b%s\b`, name)).ReplaceAllString(n, style.Primary.Bold(true).Render(name))
		} else {
			n = regexp.MustCompile(fmt.Sprintf(`\b%s\b`, name)).ReplaceAllString(n, lipgloss.NewStyle().Bold(true).Render(name))
		}
	}

	return n
}
