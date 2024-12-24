package views

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/VU-ASE/rover/src/openapi"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/VU-ASE/rover/src/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lempiy/dgraph"
	"github.com/lempiy/dgraph/core"
)

//
// Action responses
//

type PipelineOverviewServiceInfo struct {
	Name          string
	Version       string
	Author        string
	Configuration openapi.ServicesAuthorServiceVersionGet200Response
}

type PipelineOverviewSummary struct {
	// Basic pipeline GET request
	Pipeline openapi.PipelineGet200Response
	// Information about services specifically
	Services []PipelineOverviewServiceInfo
	// Status from roverd (for CPU and memory usage)
	Status openapi.StatusGet200Response
}

//
// All keys
//

// Keys to navigate
type PipelineOverviewKeyMap struct {
	Retry     key.Binding
	Toggle    key.Binding // start/stop pipeline
	Logs      key.Binding
	Details   key.Binding
	Configure key.Binding
	Quit      key.Binding
}

// Shown when the services are being updated
var pipelineOverviewKeysRegular = PipelineOverviewKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refetch"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

var pipelineOverviewKeysRunning = PipelineOverviewKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refetch"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "stop pipeline"),
	),
	Logs: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "logs"),
	),
	Details: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "details"),
	),
	Configure: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "configure"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

var pipelineOverviewKeysIdle = PipelineOverviewKeyMap{
	Retry: pipelineOverviewKeysRunning.Retry,
	Toggle: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start pipeline"),
	),
	Configure: pipelineOverviewKeysRunning.Configure,
	Quit:      pipelineOverviewKeysRunning.Quit,
}

func (k PipelineOverviewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Retry, k.Toggle, k.Logs, k.Details, k.Configure, k.Quit}
}

func (k PipelineOverviewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

//
// Possible tabs to select
//

type PipelineOverviewTab int

const (
	PipelineOverviewTabNone PipelineOverviewTab = iota
	PipelineOverviewTabServiceDetails
	PipelineOverviewTabLogs
)

//
// The page model
//

type PipelineOverviewPage struct {
	help          help.Model
	spinner       spinner.Model
	pipeline      tui.Action[PipelineOverviewSummary]
	pipelineGraph string // preserved in the model to avoid re-rendering in .View()
	progress      progress.Model
	table         table.Model
	openView      PipelineOverviewTab
	toggle        tui.Action[bool]
}

func NewPipelineOverviewPage() PipelineOverviewPage {
	// todo

	return PipelineOverviewPage{
		spinner:       spinner.New(),
		help:          help.New(),
		pipeline:      tui.NewAction[PipelineOverviewSummary]("pipelineFetch"),
		pipelineGraph: "",
		progress:      progress.New(progress.WithScaledGradient(string(style.AsePrimary), "#FFF")),
		table:         table.New(),
		openView:      PipelineOverviewTabNone,
	}
}

//
// Page model methods
//

func (m PipelineOverviewPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tui.ActionInit[bool]:
		m.toggle.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[bool]:
		m.toggle.ProcessResult(msg)
		if m.toggle.IsSuccess() {
			// we know data is stale now, need to refetch
			return m, m.fetchPipeline()
		}
	case tui.ActionInit[PipelineOverviewSummary]:
		m.pipeline.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[PipelineOverviewSummary]:
		m.pipeline.ProcessResult(msg)
		if m.pipeline.IsSuccess() {
			// Create the pipeline graph based on enabled services
			nodes := make([]core.NodeInput, 0)
			for _, service := range m.pipeline.Data.Pipeline.Enabled {
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
				m.pipelineGraph = style.Gray.Render("This pipeline is empty")
			} else if err != nil {
				m.pipelineGraph = "Failed to draw pipeline\n"
			} else {
				m.pipelineGraph = fmt.Sprintf("%s\n", canvas)
			}
			m.table = m.createServiceTable(*m.pipeline.Data)
		}

		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, pipelineOverviewKeysRegular.Quit):
			return m, tea.Quit
		case key.Matches(msg, pipelineOverviewKeysRegular.Retry):
			if !m.pipeline.IsLoading() {
				return m, m.fetchPipeline()
			}
			return m, nil
		case key.Matches(msg, pipelineOverviewKeysRunning.Logs):
			sel := m.table.SelectedRow()
			if sel != nil {
				service := sel[0]
				version := sel[1]
				author := sel[2]

				return RootScreen(state.Get()).SwitchScreen(NewPipelineLogsPage(
					service, author, version,
				))
			}
			return m, nil
		case key.Matches(msg, pipelineOverviewKeysRunning.Toggle):
			if !m.pipeline.IsLoading() && !m.toggle.IsLoading() {
				return m, m.toggleExecution()
			}
		case key.Matches(msg, pipelineOverviewKeysRunning.Details):
			if !m.pipeline.IsLoading() {
				var found *PipelineOverviewServiceInfo
				sel := m.table.SelectedRow()
				if sel != nil {
					for _, s := range m.pipeline.Data.Services {
						if s.Name == sel[0] {
							found = &s
							break
						}
					}

					if found != nil {
						return RootScreen(state.Get()).SwitchScreen(NewPipelineDetailsPage(*found))
					}
				}
			}
		case key.Matches(msg, pipelineOverviewKeysRunning.Configure):
			return RootScreen(state.Get()).SwitchScreen(NewPipelineConfiguratorPage())
		}
	case tea.WindowSizeMsg:
		m.progress.Width = (msg.Width - 4 - 6 - 6) / 3 // padding
		if m.pipeline.HasData() {
			m.table = m.createServiceTable(*m.pipeline.Data)
		} else {
			m.table = table.New()
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m PipelineOverviewPage) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchPipeline())
}

// Rendered when the pipeline is successfully fetched
func (m PipelineOverviewPage) pipelineView() string {
	if len(m.pipeline.Data.Services) <= 0 {
		return style.Gray.Render("Your pipeline is empty. Start by adding services to it.")
	}

	status := style.Error.Bold(true).Render("Unknown")

	if m.pipeline.Data.Pipeline.Status == openapi.STARTABLE {
		status = style.Color(style.SuccessLight).Bold(true).Render("Startable")
		if m.toggle.IsLoading() {
			status = style.Warning.Bold(true).Render("Starting")
		}
	} else if m.pipeline.Data.Pipeline.Status == openapi.STARTED {
		status = style.Success.Bold(true).Render("Running")
		if m.toggle.IsLoading() {
			status = style.Warning.Bold(true).Render("Stopping")
		}
	} else if m.pipeline.Data.Pipeline.Status == openapi.RESTARTING {
		status = style.Warning.Bold(true).Render("Restarting")
		if m.toggle.IsLoading() {
			status = style.Warning.Bold(true).Render("Stopping")
		}
	}
	s := m.postProcessGraph(m.pipelineGraph)
	status = status + "\n"
	if m.pipeline.Data.Pipeline.LastStart != nil {
		status += style.Gray.Render("last started at: ") + time.Unix(*m.pipeline.Data.Pipeline.LastStart, 0).Format("2006-01-02 15:04:05") + "\n"
	}
	if m.pipeline.Data.Pipeline.LastStop != nil {
		status += style.Gray.Render("last stopped at: ") + time.Unix(*m.pipeline.Data.Pipeline.LastStop, 0).Format("2006-01-02 15:04:05") + "\n"
	}
	if m.pipeline.Data.Pipeline.LastRestart != nil {
		status += style.Gray.Render("last restarted at: ") + time.Unix(*m.pipeline.Data.Pipeline.LastRestart, 0).Format("2006-01-02 15:04:05") + "\n"
	}

	cpu := style.Gray.Render("Total CPU usage per core ") + "\n"
	if len(m.pipeline.Data.Status.Cpu) > 0 {
		for _, c := range m.pipeline.Data.Status.Cpu {
			cpu += m.progress.ViewAs(float64(c.Used)/float64(c.Total)) + "\n"
		}
	} else {
		cpu += style.Gray.Render("No CPU usage data available") + "\n"
	}
	mem := style.Gray.Render("Total memory usage") + "\n" + m.progress.ViewAs(float64(m.pipeline.Data.Status.Memory.Used)/float64(m.pipeline.Data.Status.Memory.Total)) + "\n"
	if m.pipeline.Data.Status.Memory.Total > 0 {
		mem += style.Gray.Render("The Rover uses ") + fmt.Sprintf("%d/%d MB", m.pipeline.Data.Status.Memory.Used, m.pipeline.Data.Status.Memory.Total) + "\n"
		mem += style.Gray.Render("Of which ") + fmt.Sprintf("%d MB", 10) + style.Gray.Render(" is used by this pipeline") + "\n"
	} else {
		mem += style.Gray.Render("No memory usage data available") + "\n"
	}

	// Calculate column width (subtract padding and borders)
	columnWidth := (state.Get().WindowWidth - 4 - 6) / 3 // Adjust for padding and borders

	// Define styles for each column
	columnStyle := lipgloss.NewStyle().
		Width(columnWidth)

	row := lipgloss.JoinHorizontal(lipgloss.Top,
		columnStyle.Render(status),
		columnStyle.Render(mem),
		columnStyle.Render(cpu),
	)

	proc_list := "\n" + m.table.View() + "\n\n"
	view := s + "\n\n" + row + proc_list

	if m.toggle.IsError() {
		view += style.Gray.Render("Could not toggle pipeline execution") + " (" + m.toggle.Error.Error() + ")" + "\n\n"
	}

	view += m.table.HelpView() + style.Gray.Render(" • ")
	if m.pipeline.Data.Pipeline.Status == openapi.STARTED {
		view += m.help.View(pipelineOverviewKeysRunning)
	} else {
		view += m.help.View(pipelineOverviewKeysIdle)
	}

	return view
}

func (m PipelineOverviewPage) View() string {
	s := style.Title.Render("Pipeline")
	// We're doing optimistic updates, so we want to show an indicator without disrupting the view
	if (m.pipeline.IsLoading() && m.pipeline.HasData()) || m.toggle.IsLoading() {
		s += " " + m.spinner.View()
	}
	s += "\n\n"

	if m.pipeline.IsLoading() && !m.pipeline.HasData() {
		s += m.spinner.View() + " Loading pipeline..."
	} else if m.pipeline.IsError() {
		s += style.Error.Render("Error loading pipeline") + " (" + m.pipeline.Error.Error() + ")"
	} else if m.pipeline.HasData() {
		s += m.pipelineView()
	}
	s += "\n"

	return s
}

//
// Actions
//

func (m PipelineOverviewPage) fetchPipeline() tea.Cmd {
	return tui.PerformAction(&m.pipeline, func() (*PipelineOverviewSummary, error) {
		remote := state.Get().RoverConnections.GetActive()
		if remote == nil {
			return nil, fmt.Errorf("No active rover connection")
		}

		api := remote.ToApiClient()

		// First, fetch all services and the status of the current pipeline
		pipeline, htt, err := api.PipelineAPI.PipelineGet(
			context.Background(),
		).Execute()

		if err != nil && htt != nil {
			return nil, utils.ParseHTTPError(err, htt)
		}

		// Then, for each service, we need to query the service for its actual configuration (inputs, outputs)
		services := make([]PipelineOverviewServiceInfo, 0)
		for _, enabled := range pipeline.Enabled {
			configuration, htt, err := api.ServicesAPI.ServicesAuthorServiceVersionGet(
				context.Background(),
				enabled.Service.Author,
				enabled.Service.Name,
				enabled.Service.Version,
			).Execute()

			if err != nil && htt != nil {
				return nil, utils.ParseHTTPError(err, htt)
			}

			services = append(services, PipelineOverviewServiceInfo{
				Name:          enabled.Service.Name,
				Version:       enabled.Service.Version,
				Author:        enabled.Service.Author,
				Configuration: *configuration,
			})
		}

		// Then the Rover status
		status, htt, err := api.HealthAPI.StatusGet(
			context.Background(),
		).Execute()

		if err != nil && htt != nil {
			return nil, utils.ParseHTTPError(err, htt)
		}

		// Combined response
		res := PipelineOverviewSummary{
			Pipeline: *pipeline,
			Services: services,
			Status:   *status,
		}

		return &res, nil
	})
}

// Clean up the graph to make it a bit more readable and compressed
func (m PipelineOverviewPage) postProcessGraph(s string) string {
	n := s

	// Remove empty lines
	n = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(n, "\n")

	// Highlight the currently selected service
	sel := m.table.SelectedRow()
	if sel != nil {
		// The first item is always the service name
		name := sel[0]

		// Find the service in the graph
		n = regexp.MustCompile(fmt.Sprintf(`\b%s\b`, name)).ReplaceAllString(n, style.Primary.Bold(true).Render(name))
	}

	return n
}

// Converts a percentage to a table column width in characters
func pct(pct int) int {
	total := state.Get().WindowWidth - 4 - 6 - 6 // padding
	return int(float64(total)*float64(pct)/100.0) - 1
}

// Create a nicely formatted table based on input data
func (m PipelineOverviewPage) createServiceTable(res PipelineOverviewSummary) table.Model {
	// Depending on the state of the pipeline, we want to show different columns
	var columns []table.Column
	var rows []table.Row

	// Pipeline is currently running
	if res.Pipeline.Status == openapi.STARTED {
		columns = []table.Column{
			{Title: "Service", Width: pct(10)},
			{Title: "Version", Width: pct(10)},
			{Title: "Author", Width: pct(10)},
			{Title: "Faults", Width: pct(10)},
			{Title: "Uptime", Width: pct(10)},
			{Title: "PID", Width: pct(10)},
			{Title: "CPU", Width: pct(10)},
			{Title: "Memory", Width: pct(30)},
		}

		rows = []table.Row{
			// {"imaging", "1.0.1", "vu-ase", "0", "1h 23m 45s", "1234", "5%", "50MB"},
			// {"controller", "1.1.1", "vu-ase", "0", "1h 23m 45s", "1234", "10%", "150MB"},
			// {"transceiver", "1.2.2", "vu-ase", "0", "1h 23m 45s", "1234", "15%", "250MB"},
		}

		for _, e := range res.Pipeline.Enabled {
			row := []string{
				e.Service.Name,
				e.Service.Version,
				e.Service.Author,
				fmt.Sprintf("%d", *e.Service.Faults),
			}

			if e.Process != nil {
				row = append(row,
					utils.FormatDuration(e.Process.Uptime),
				)
				row = append(row, fmt.Sprintf("%d", e.Process.Pid))
				row = append(row, fmt.Sprintf("%d%%", e.Process.Cpu))
				row = append(row, fmt.Sprintf("%dMB", e.Process.Memory))
			}

			rows = append(rows, row)
		}
	} else if res.Pipeline.Status != openapi.STARTED && res.Pipeline.LastStart == nil {
		// This pipeline is not running, and has not been started before
		columns = []table.Column{
			{Title: "Service", Width: pct(10)},
			{Title: "Version", Width: pct(5)},
			{Title: "Author", Width: pct(85)},
		}

		rows = []table.Row{
			// {"imaging", "1.0.1", "vu-ase"},
			// {"controller", "1.1.1", "vu-ase"},
			// {"transceiver", "1.2.2", "vu-ase"},
		}

		for _, e := range res.Pipeline.Enabled {
			row := []string{
				e.Service.Name,
				e.Service.Version,
				e.Service.Author,
			}

			rows = append(rows, row)
		}
	} else {
		// This pipeline is not running, but has been started before
		columns = []table.Column{
			{Title: "Service", Width: pct(10)},
			{Title: "Version", Width: pct(5)},
			{Title: "Author", Width: pct(10)},
			{Title: "Faults", Width: pct(5)},
			{Title: "Uptime", Width: pct(70)},
		}

		rows = []table.Row{
			// {"imaging", "1.0.1", "vu-ase", "0", "1h 23m 45s"},
			// {"controller", "1.1.1", "vu-ase", "0", "1h 23m 45s"},
		}

		for _, e := range res.Pipeline.Enabled {
			row := []string{
				e.Service.Name,
				e.Service.Version,
				e.Service.Author,
				fmt.Sprintf("%d", *e.Service.Faults),
			}

			if e.Process != nil {
				row = append(row,
					utils.FormatDuration(e.Process.Uptime),
				)
			}

			rows = append(rows, row)
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

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("FFF")).
		Background(style.AsePrimary).
		Bold(false)
	t.SetStyles(s)

	return t
}

func (m PipelineOverviewPage) toggleExecution() tea.Cmd {
	return tui.PerformAction(&m.toggle, func() (*bool, error) {
		remote := state.Get().RoverConnections.GetActive()
		if remote == nil {
			return nil, fmt.Errorf("No active rover connection")
		}

		api := remote.ToApiClient()
		var err error
		if m.pipeline.Data.Pipeline.Status == openapi.STARTED {
			htt, err := api.PipelineAPI.PipelineStopPost(
				context.Background(),
			).Execute()

			if err != nil && htt != nil {
				return nil, utils.ParseHTTPError(err, htt)
			}
		} else {
			htt, err := api.PipelineAPI.PipelineStartPost(
				context.Background(),
			).Execute()

			if err != nil && htt != nil {
				return nil, utils.ParseHTTPError(err, htt)
			}
		}

		// todo: also need to call API build endpoint
		return openapi.PtrBool(true), err
	})
}
