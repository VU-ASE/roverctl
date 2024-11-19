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
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
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
	Retry   key.Binding
	Confirm key.Binding
	Quit    key.Binding
}

// Shown when the services are being updated
var pipelineOverviewKeysRegular = PipelineOverviewKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func (k PipelineOverviewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Retry, k.Confirm, k.Quit}
}

func (k PipelineOverviewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

//
// The page model
//

type PipelineOverviewPage struct {
	help          help.Model
	spinner       spinner.Model
	pipeline      tui.Action[PipelineOverviewSummary]
	pipelineGraph string // preserved in the model to avoid re-rendering in .View()
	progress      progress.Model
}

func NewPipelineOverviewPage() PipelineOverviewPage {
	// todo

	return PipelineOverviewPage{
		spinner:       spinner.New(),
		help:          help.New(),
		pipeline:      tui.NewAction[PipelineOverviewSummary]("pipelineFetch"),
		pipelineGraph: "",
		progress:      progress.New(progress.WithScaledGradient(string(style.AsePrimary), "#FFF")),
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

				canvas, err := dgraph.DrawGraph(nodes)
				if err != nil {
					m.pipelineGraph = "Failed to draw pipeline\n"
				} else if len(nodes) > 0 {
					m.pipelineGraph = m.postProcessGraph(fmt.Sprintf("%s\n", canvas))
				}
			}
		}

		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, pipelineOverviewKeysRegular.Quit):
			return m, tea.Quit
		case key.Matches(msg, pipelineOverviewKeysRegular.Retry):
			// todo:
			return m, nil
		case key.Matches(msg, pipelineOverviewKeysRegular.Confirm):
			// todo:
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.progress.Width = (msg.Width - 4 - 6 - 6) / 3 // padding
	}

	return m, nil
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
	} else if m.pipeline.Data.Pipeline.Status == openapi.STARTED {
		status = style.Success.Bold(true).Render("Running")
	} else if m.pipeline.Data.Pipeline.Status == openapi.RESTARTING {
		status = style.Warning.Bold(true).Render("Restarting")
	}
	s := m.pipelineGraph
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

	cpu := ""
	if len(m.pipeline.Data.Status.Cpu) > 0 {
		cpu += style.Gray.Render("CPU") + "\n"
		for _, c := range m.pipeline.Data.Status.Cpu {
			cpu += m.progress.ViewAs(float64(c.Used)/float64(c.Total)) + "\n"
		}
	}
	mem := "\n" + style.Gray.Render("Memory") + "\n" + m.progress.ViewAs(float64(m.pipeline.Data.Status.Memory.Used)/float64(m.pipeline.Data.Status.Memory.Total)) + "\n"

	enabled := style.Gray.Render("Services") + "\n"
	for _, service := range m.pipeline.Data.Services {
		enabled += service.Author + "/" + style.Primary.Render(service.Name) + style.Gray.Render(""+service.Version) + "\n"
	}

	// Calculate column width (subtract padding and borders)
	columnWidth := (state.Get().WindowWidth - 4 - 6) / 3 // Adjust for padding and borders

	// Define styles for each column
	columnStyle := lipgloss.NewStyle().
		// BorderStyle(lipgloss.NormalBorder()).
		// Padding(1, 1).
		Width(columnWidth)

	// Join columns horizontally
	row := lipgloss.JoinHorizontal(lipgloss.Top,
		columnStyle.Render(status),
		columnStyle.Render(enabled),
		columnStyle.Render(cpu+mem),
	)

	return s + "\n\n" + row
}

func (m PipelineOverviewPage) View() string {
	s := style.Title.Render("Pipeline") + "\n\n"
	if m.pipeline.IsLoading() {
		s += m.spinner.View() + " Loading pipeline..."
	} else if m.pipeline.IsError() {
		s += style.Error.Render("Error loading pipeline") + " (" + m.pipeline.Error.Error() + ")"
	} else if m.pipeline.IsSuccess() {
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
		// mock fetch
		// ! remove

		time.Sleep(100 * time.Millisecond)
		// First roverd tells us what services are enabled, by reference (FQN)
		pipeline := openapi.PipelineGet200Response{
			Status:    openapi.STARTABLE,
			LastStart: openapi.PtrInt64(123456),
			LastStop:  openapi.PtrInt64(123456),
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
								Service: "controllertje",
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

// Clean up the graph to make it a bit more readable and compressed
func (m PipelineOverviewPage) postProcessGraph(s string) string {
	n := s

	// Remove empty lines
	n = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(n, "\n")

	return n
}
