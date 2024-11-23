package views

import (
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
)

//
// All keys
//

// Keys to navigate
type PipelineConfiguratorKeyMap struct {
	Retry   key.Binding
	Confirm key.Binding
	Switch  key.Binding // switch table focus
	Quit    key.Binding
}

// Shown when the services are being updated
var pipelineConfiguratorKeysRegular = PipelineConfiguratorKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
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
	return []key.Binding{k.Retry, k.Confirm, k.Quit}
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
	pipeline tui.Action[PipelineOverviewSummary]
	authors  tui.Action[[]string]                                           // first part of FQN
	services tui.Action[[]string]                                           // second part of FQN
	versions tui.Action[openapi.ServicesAuthorServiceVersionGet200Response] // third part of FQN, with all information
	// Keep track of the focussed table (left/right)
	focussed int // 0 = active, 1 = remote
}

func NewPipelineConfiguratorPage() PipelineConfiguratorPage {
	// todo

	return PipelineConfiguratorPage{
		spinner:     spinner.New(),
		help:        help.New(),
		tableActive: table.New(),
		tableRemote: table.New(),
		pipeline:    tui.NewAction[PipelineOverviewSummary]("fetchActive"),
		authors:     tui.NewAction[[]string]("fetchAuthors"),
		services:    tui.NewAction[[]string]("fetchServices"),
		versions:    tui.NewAction[openapi.ServicesAuthorServiceVersionGet200Response]("fetchVersions"),
		focussed:    0,
	}
}

//
// Page model methods
//

func (m PipelineConfiguratorPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tui.ActionInit[PipelineOverviewSummary]:
		m.pipeline.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[PipelineOverviewSummary]:
		m.pipeline.ProcessResult(msg)
		if m.pipeline.IsSuccess() {
			m.tableActive = m.createActiveTable(*m.pipeline.Data)
		}
		return m, nil
	case tui.ActionInit[[]string]:
		m.authors.ProcessInit(msg)
		m.services.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[[]string]:
		m.authors.ProcessResult(msg)
		m.services.ProcessResult(msg)
		m.tableRemote = m.createServicesTable()
		return m, nil
	case tui.ActionInit[openapi.ServicesAuthorServiceVersionGet200Response]:
		m.versions.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[openapi.ServicesAuthorServiceVersionGet200Response]:
		m.versions.ProcessResult(msg)
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, pipelineConfiguratorKeysRegular.Quit):
			return m, tea.Quit
		case key.Matches(msg, pipelineConfiguratorKeysRegular.Retry):
			// todo:
			return m, nil
		case key.Matches(msg, pipelineConfiguratorKeysRegular.Confirm):
			// todo:
			return m, nil
		case key.Matches(msg, pipelineConfiguratorKeysRegular.Switch):
			m.focussed = (m.focussed + 1) % 2
			m.tableActive = m.createActiveTable(*m.pipeline.Data)
			m.tableRemote = m.createServicesTable()
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

func (m PipelineConfiguratorPage) View() string {
	s := style.Title.Render("Configure your pipeline") + "\n\n"

	// Calculate column width (subtract padding and borders)
	columnWidth := m.getColWidth()

	// Define styles for each column
	columnStyle := lipgloss.NewStyle().Width(columnWidth)

	graph := ""
	if m.pipeline.IsSuccess() {
		graph = "Pipeline loaded"
	} else if m.pipeline.IsError() {
		graph = style.Error.Render("Error loading pipeline") + style.Gray.Render(" "+m.pipeline.Error.Error())
	} else {
		graph = m.spinner.View() + " Loading pipeline"
	}
	graph += "\n\n"

	// Define the columns
	columnActive := columnStyle.Render("Active services") + "\n\n" + m.tableActive.View()
	columnRemote := columnStyle.Render("Remote services") + "\n\n" + m.tableRemote.View()

	row := lipgloss.JoinHorizontal(lipgloss.Top,
		columnStyle.Render(columnActive),
		" ",
		columnStyle.Render(columnRemote),
	)

	return s + graph + row
}

//
// Aux for table and col rendering
//

func (m PipelineConfiguratorPage) getColWidth() int {
	return (state.Get().WindowWidth - 4 - 6) / 2 // Adjust for padding and borders
}

func (m PipelineConfiguratorPage) colPct(pct int) int {
	total := m.getColWidth()
	return (total*pct)/100 - 1
}

//
// Actions
//

func (m PipelineConfiguratorPage) fetchPipeline() tea.Cmd {
	return tui.PerformAction(&m.pipeline, func() (*PipelineOverviewSummary, error) {
		// mock fetch
		// ! remove

		time.Sleep(1000 * time.Millisecond)
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

func (m PipelineConfiguratorPage) fetchAllAuthors() tea.Cmd {
	return tui.PerformAction(&m.authors, func() (*[]string, error) {
		// mock fetch
		// ! remove
		time.Sleep(5000 * time.Millisecond)

		return &[]string{
			"vu-ase",
			"ielaajez",
			"test",
		}, nil
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

func (m PipelineConfiguratorPage) createServicesTable() table.Model {
	// Retrieve the previously selected entry
	prev := m.tableActive.SelectedRow()

	columns := []table.Column{
		{Title: "Service", Width: m.colPct(30)},
		{Title: "Version", Width: m.colPct(30)},
		{Title: "Author", Width: m.colPct(39)},
	}

	rows := make([]table.Row, 0)

	if m.authors.HasData() {
		for _, author := range *m.authors.Data {
			rows = append(rows, table.Row{
				author,
				"n/a",
				"n/a",
			})
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
