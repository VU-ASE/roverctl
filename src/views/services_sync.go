package views

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/VU-ASE/roverctl/src/openapi"
	"github.com/VU-ASE/roverctl/src/state"
	"github.com/VU-ASE/roverctl/src/style"
	"github.com/VU-ASE/roverctl/src/tui"
	"github.com/VU-ASE/roverctl/src/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/radovskyb/watcher"
)

type ServicesSyncPage struct {
	// To select an action to perform with this utility
	// actions list.Model // actions you can perform when connected to a Rover
	// help    help.Model // to display a help footer
	// Is the cwd a service directory?
	isServiceDir  bool
	changes       tui.Action[[]fileChangeMsg]
	uploading     tui.Action[openapi.FetchPost200Response]
	channel       chan fileChangeMsg
	watchDebounce time.Duration // debounce time for collecting changes
	spinner       spinner.Model
	help          help.Model
}

type ServicesSyncKeyMap struct {
	Retry key.Binding
	Quit  key.Binding
}

func (k ServicesSyncKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Retry, k.Quit}
}

func (k ServicesSyncKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var serviceSyncKeysFailure = ServicesSyncKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

var serviceSyncKeysRegular = ServicesSyncKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func NewServicesSyncPage() ServicesSyncPage {
	// Is there already a service.yaml file in the current directory?
	// _, err := os.Stat("./service.yaml")
	ch := make(chan fileChangeMsg)

	// Actions
	changes := tui.NewAction[[]fileChangeMsg]("fileChanges")
	uploading := tui.NewAction[openapi.FetchPost200Response]("uploading")
	sp := spinner.New()

	model := ServicesSyncPage{
		isServiceDir:  true,
		channel:       ch,
		changes:       changes,
		watchDebounce: 500 * time.Millisecond,
		spinner:       sp,
		uploading:     uploading,
		help:          help.New(),
	}
	go model.watch(".", ch)
	return model
}

func (m ServicesSyncPage) Init() tea.Cmd {
	return tea.Batch(
		m.uploadChanges(), // always upload the initial files first
		m.spinner.Tick,
	)
}

func (m ServicesSyncPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tui.ActionInit[openapi.FetchPost200Response]:
		m.uploading.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[openapi.FetchPost200Response]:
		m.uploading.ProcessResult(msg)
		if m.uploading.IsSuccess() {
			return m, m.collectChanges(nil)
		}
		return m, nil
	case tui.ActionInit[[]fileChangeMsg]:
		m.changes.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[[]fileChangeMsg]:
		m.changes.ProcessResult(msg)
		if m.changes.IsSuccess() {
			return m, m.uploadChanges()
		}
		return m, nil
	case key.Help:
		m.help, cmd = m.help.Update(msg)
		return m, cmd
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, serviceSyncKeysFailure.Retry):
			if m.uploading.IsError() {
				return m, m.uploadChanges()
			}
		}
	}

	return m, nil
}

func (m ServicesSyncPage) View() string {
	s := style.Title.Render("Synchronize service") + "\n\n"

	if m.uploading.IsLoading() {
		s += m.spinner.View() + " Syncing (uploading files)..." + "\n\n"
	} else if m.uploading.IsError() {
		s += style.Error.Render("✗ Could not sync files") + style.Gray.Render(" ("+m.uploading.Error.Error()+")") + "\n\n"
	} else if m.uploading.IsSuccess() {
		s += style.Success.Render("✓ Service is in sync")
	}

	if m.changes.IsLoading() && m.uploading.IsSuccess() {
		s += ", watching changes..." + "\n\n"
	} else if m.changes.IsSuccess() {
		for _, change := range *m.changes.Data {
			if change.action == created {
				s += style.Success.Render(" + " + change.filename)
			} else if change.action == modified {
				s += style.Warning.Render(" ~ " + change.filename)
			} else if change.action == deleted {
				s += style.Error.Render(" - " + change.filename)
			} else {
				s += style.Gray.Render(" ? " + change.filename)
			}
			s += "\n"
		}
		s += "\n"
	}

	if m.uploading.IsError() {
		s += m.help.View(serviceSyncKeysFailure)
	} else {
		s += m.help.View(serviceSyncKeysRegular)
	}

	return s
}

// Collect all changes so far reported on the channel
func (m ServicesSyncPage) collectChanges(initial *[]fileChangeMsg) tea.Cmd {
	start := time.Now()
	return tui.PerformAction(&m.changes, func() (*[]fileChangeMsg, error) {
		all := make([]fileChangeMsg, 0)
		changed := false // We only want to return if there are new changes, not if the initial list is the same
		if initial != nil {
			all = append(all, *initial...)
		}

		for {
			select {
			case value, ok := <-m.channel:
				if !ok {
					// Channel is closed
					return &all, nil
				}

				// Is this change already in the list?
				var found *fileChangeMsg
				for _, v := range all {
					if v.filename == value.filename {
						found = &v
						break
					}
				}

				// If it is, update the action
				if found != nil && value.action != unknown {
					found.action = value.action
				} else if found == nil {
					all = append(all, value)
				}
				start = time.Now()
				changed = true

			default:
				// No more values to read (debounced)
				if len(all) > 0 && changed && time.Since(start) > m.watchDebounce {
					return &all, nil
				}
			}
		}
	})
}

func createZipFromDirectory(zipPath, sourceDir string) error {
	// Create the zip file
	tmpZip, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create temp zip file: %v", err)
	}
	defer tmpZip.Close()

	// Create a zip writer
	zipWriter := zip.NewWriter(tmpZip)
	defer func() {
		if closeErr := zipWriter.Close(); closeErr != nil {
			fmt.Printf("Error closing zip writer: %v\n", closeErr)
		}
	}()

	// Walk through the source directory
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking the path %q: %v", path, err)
		}

		// Get the relative path to store in the zip archive
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %v", err)
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Skip the "roverctl" binary
		if relPath == "roverctl" {
			return nil
		}

		if info.IsDir() {
			// Add a directory entry to the zip file
			_, err := zipWriter.Create(relPath + "/")
			if err != nil {
				return fmt.Errorf("failed to create directory entry: %v", err)
			}
			return nil
		}

		// Add a file entry to the zip file
		fileWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return fmt.Errorf("failed to create file entry: %v", err)
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()

		// Copy the file content to the zip file
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			return fmt.Errorf("failed to write file content: %v", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Upload all collected changes to the Rover
func (m ServicesSyncPage) uploadChanges() tea.Cmd {
	return tui.PerformAction(&m.uploading, func() (*openapi.FetchPost200Response, error) {
		remote := state.Get().RoverConnections.GetActive()
		if remote == nil {
			return nil, fmt.Errorf("No active rover connection")
		}

		// Create a temp zip file
		zipPath := os.TempDir() + "/rover-sync-" + time.Now().Format("20060102150405") + ".zip"
		sourceDir := "."

		err := createZipFromDirectory(zipPath, sourceDir)
		if err != nil {
			return nil, fmt.Errorf("Error creating zip: %v\n", err)
		}

		// Open the zip file
		zipFile, err := os.Open(zipPath)
		if err != nil {
			return nil, fmt.Errorf("Failed to open temp zip file: %v", err)
		}
		defer zipFile.Close()

		// Upload the zip file
		api := remote.ToApiClient()
		req := api.ServicesAPI.UploadPost(
			context.Background(),
		)
		req = req.Content(zipFile)

		// req.Content(zipFile)
		res, htt, err := req.Execute()

		if err != nil && htt != nil {
			return nil, utils.ParseHTTPError(err, htt)
		}

		return res, err
	})
}

// Enum for possible file change actions, using iota
type fileChangeAction int

const (
	created fileChangeAction = iota
	modified
	deleted
	unknown
)

type fileChangeMsg struct {
	filename string
	action   fileChangeAction
}

// Start the watcher and let it report changes on the specified channel (this is a long running function)
func (m ServicesSyncPage) watch(path string, c chan fileChangeMsg) {
	w := watcher.New()

	// Only notify rename and move events
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Create, watcher.Remove, watcher.Write)

	// Ignore .git directory
	err := w.Ignore("./.git")
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			select {
			case event := <-w.Event:
				if event.Op == watcher.Write {
					c <- fileChangeMsg{filename: event.Path, action: modified}
				} else if event.Op == watcher.Create {
					c <- fileChangeMsg{filename: event.Path, action: created}
				} else if event.Op == watcher.Remove {
					c <- fileChangeMsg{filename: event.Path, action: deleted}
				} else {
					c <- fileChangeMsg{filename: event.Path, action: unknown}
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch recursively for changes
	if err := w.AddRecursive(path); err != nil {
		log.Fatalln(err)
	}

	// Start the watching process - it'll check for changes every 100ms
	if err := w.Start(time.Millisecond * 500); err != nil {
		log.Fatalln(err)
	}
}
