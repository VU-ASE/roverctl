package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

//
// All keys
//

// Keys to navigate
type PipelineLogsKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Retry   key.Binding
	Confirm key.Binding
	Quit    key.Binding
}

// Shown when the services are being updated
var PipelineLogsKeysRegular = PipelineLogsKeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "down"),
	),
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refetch"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func (k PipelineLogsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Retry, k.Confirm, k.Quit}
}

func (k PipelineLogsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

//
// The page model
//

type PipelineLogsPage struct {
	help     help.Model
	spinner  spinner.Model
	viewport viewport.Model
	logs     tui.Action[[]string]
	// Service FQN
	service string // name
	author  string
	version string
}

func NewPipelineLogsPage(service string, author string, version string) PipelineLogsPage {
	// todo

	return PipelineLogsPage{
		spinner: spinner.New(),
		help:    help.New(),
		logs:    tui.NewAction[[]string]("fetchLogs"),
		viewport: viewport.Model{
			Width:  state.Get().WindowWidth,
			Height: state.Get().WindowHeight - 4,
		},
		service: service,
		author:  author,
		version: version,
	}
}

//
// Page model methods
//

func (m PipelineLogsPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tui.ActionInit[[]string]:
		m.logs.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[[]string]:
		m.logs.ProcessResult(msg)
		if m.logs.IsSuccess() {
			m.viewport.YOffset = 100
			m.viewport.SetContent(strings.Join(*m.logs.Data, "\n"))
		}
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, PipelineLogsKeysRegular.Quit):
			return m, tea.Quit
		case key.Matches(msg, PipelineLogsKeysRegular.Retry):
			if !m.logs.IsLoading() {
				return m, m.fetchLogs()
			}
		case key.Matches(msg, PipelineLogsKeysRegular.Confirm):
			// todo:
			return m, nil
		}
	case tea.WindowSizeMsg:
		offset := m.viewport.YOffset
		m.viewport = viewport.New(msg.Width, msg.Height-2-2-2) // -2 for the header, -2 for the padding, -2 for the help
		if m.logs.IsSuccess() {
			m.viewport.YOffset = offset
			m.viewport.SetContent(strings.Join(*m.logs.Data, "\n"))
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m PipelineLogsPage) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchLogs())
}

func (m PipelineLogsPage) View() string {
	if m.logs.IsSuccess() {
		return style.Title.Render("Logs for "+m.author+"/"+m.service+":"+m.version) + style.Gray.Render(fmt.Sprintf(" %3.f%%", 100.0*m.viewport.ScrollPercent())) + "\n\n" + m.viewport.View() + "\n\n" + m.help.View(PipelineLogsKeysRegular)
	} else if m.logs.IsError() {
		return style.Error.Render("Failed to load logs for " + m.service + ".")
	} else {
		return m.spinner.View() + " Fetching logs..."
	}
}

func (m PipelineLogsPage) fetchLogs() tea.Cmd {
	return tui.PerformAction(&m.logs, func() (*[]string, error) {
		// mock log fetching
		// !remove

		time.Sleep(2 * time.Second)

		logs := []string{
			"[2024-11-22 12:00:01] [INFO] Starting the program.",
			"[2024-11-22 12:00:02] [INFO] Loading configuration files.",
			"[2024-11-22 12:00:03] [WARN] Configuration file 'config.yaml' is missing optional fields.",
			"[2024-11-22 12:00:04] [INFO] Connecting to the database.",
			"[2024-11-22 12:00:05] [ERROR] Failed to connect to the database: timeout.",
			"[2024-11-22 12:00:06] [INFO] Retrying database connection.",
			"[2024-11-22 12:00:07] [INFO] Database connection established successfully.",
			"[2024-11-22 12:00:08] [INFO] Starting the main service loop.",
			"[2024-11-22 12:00:09] [DEBUG] Checking health of service 'serviceA'.",
			"[2024-11-22 12:00:10] [INFO] Service 'serviceA' is running normally.",
			"[2024-11-22 12:00:11] [WARN] High memory usage detected: 85%.",
			"[2024-11-22 12:00:12] [INFO] Performing garbage collection.",
			"[2024-11-22 12:00:13] [DEBUG] Garbage collection completed in 50ms.",
			"[2024-11-22 12:00:14] [INFO] Received request from client 192.168.1.10.",
			"[2024-11-22 12:00:15] [ERROR] Failed to process client request: invalid input data.",
			"[2024-11-22 12:00:16] [INFO] Sending error response to client.",
			"[2024-11-22 12:00:17] [DEBUG] Cache hit for key 'user123'.",
			"[2024-11-22 12:00:18] [INFO] Successfully processed request for user 'user123'.",
			"[2024-11-22 12:00:19] [INFO] Performing scheduled task 'dailyCleanup'.",
			"[2024-11-22 12:00:20] [DEBUG] Cleaning temporary files.",
			"[2024-11-22 12:00:21] [INFO] Task 'dailyCleanup' completed successfully.",
			"[2024-11-22 12:00:22] [WARN] Disk space low: 5% remaining.",
			"[2024-11-22 12:00:23] [INFO] Archiving old logs to free up disk space.",
			"[2024-11-22 12:00:24] [DEBUG] Archived 1000 log files.",
			"[2024-11-22 12:00:25] [INFO] Disk space now at 20%.",
			"[2024-11-22 12:00:26] [DEBUG] Checking connection pool status.",
			"[2024-11-22 12:00:27] [INFO] Connection pool healthy: 10 connections active.",
			"[2024-11-22 12:00:28] [INFO] Service 'serviceB' started successfully.",
			"[2024-11-22 12:00:29] [WARN] Service 'serviceB' responding slowly: 2s latency.",
			"[2024-11-22 12:00:30] [INFO] Adjusting thread pool size for 'serviceB'.",
			"[2024-11-22 12:00:31] [DEBUG] New thread pool size: 20 threads.",
			"[2024-11-22 12:00:32] [INFO] Sending heartbeat to monitoring server.",
			"[2024-11-22 12:00:33] [INFO] Heartbeat acknowledged by monitoring server.",
			"[2024-11-22 12:00:34] [DEBUG] Preparing to send batch job 'job123'.",
			"[2024-11-22 12:00:35] [INFO] Batch job 'job123' started successfully.",
			"[2024-11-22 12:00:36] [WARN] Batch job 'job123' running longer than expected.",
			"[2024-11-22 12:00:37] [INFO] Completed batch job 'job123'.",
			"[2024-11-22 12:00:38] [ERROR] Error in worker thread: null pointer exception.",
			"[2024-11-22 12:00:39] [DEBUG] Restarting failed worker thread.",
			"[2024-11-22 12:00:40] [INFO] Worker thread restarted successfully.",
			"[2024-11-22 12:00:41] [INFO] Server received shutdown signal.",
			"[2024-11-22 12:00:42] [INFO] Gracefully shutting down services.",
			"[2024-11-22 12:00:43] [INFO] Service 'serviceA' stopped successfully.",
			"[2024-11-22 12:00:44] [INFO] Service 'serviceB' stopped successfully.",
			"[2024-11-22 12:00:45] [INFO] Closing database connections.",
			"[2024-11-22 12:00:46] [DEBUG] All database connections closed.",
			"[2024-11-22 12:00:47] [INFO] Program shutdown completed.",
			"[2024-11-22 12:01:00] [INFO] System maintenance starting.",
			"[2024-11-22 12:01:05] [DEBUG] Checking disk integrity.",
			"[2024-11-22 12:01:10] [INFO] Disk integrity check passed.",
			"[2024-11-22 12:01:15] [WARN] High CPU usage detected: 95%.",
			"[2024-11-22 12:01:20] [INFO] Scaling down background processes.",
			"[2024-11-22 12:01:25] [DEBUG] CPU usage now at 50%.",
			"[2024-11-22 12:01:30] [INFO] Service 'serviceC' started successfully.",
			"[2024-11-22 12:01:35] [INFO] Running scheduled report generation.",
			"[2024-11-22 12:01:40] [INFO] Report generation completed successfully.",
			"[2024-11-22 12:01:45] [INFO] Sending email notifications.",
			"[2024-11-22 12:01:50] [DEBUG] Email sent to 'user@example.com'.",
			"[2024-11-22 12:01:55] [INFO] All email notifications sent successfully.",
			"[2024-11-22 12:02:00] [INFO] Running nightly backup.",
			"[2024-11-22 12:02:05] [DEBUG] Backup process started.",
			"[2024-11-22 12:02:10] [INFO] Backup completed: 1GB data saved.",
			"[2024-11-22 12:02:15] [INFO] System check complete, no issues detected.",
			"[2024-11-22 12:02:20] [INFO] Starting program cleanup.",
			"[2024-11-22 12:02:25] [DEBUG] Removing temporary files.",
			"[2024-11-22 12:02:30] [INFO] Cleanup completed.",
			"[2024-11-22 12:02:35] [INFO] Running pre-shutdown checks.",
			"[2024-11-22 12:02:40] [DEBUG] Verifying service dependencies.",
			"[2024-11-22 12:02:45] [INFO] All dependencies verified.",
			"[2024-11-22 12:02:50] [INFO] Program ready for shutdown.",
			"[2024-11-22 12:02:55] [INFO] Initiating shutdown sequence.",
			"[2024-11-22 12:03:00] [INFO] Shutdown sequence completed.",
			"[2024-11-22 12:03:05] [INFO] Program terminated.",
		}

		return &logs, nil
	})

}
