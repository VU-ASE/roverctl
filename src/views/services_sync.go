package views

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/VU-ASE/rover/src/openapi"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
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
	uploading     tui.Action[openapi.ServicesPost200Response]
	channel       chan fileChangeMsg
	watchDebounce time.Duration // debounce time for collecting changes
	spinner       spinner.Model
}

func NewServicesSyncPage() ServicesSyncPage {
	// Is there already a service.yaml file in the current directory?
	// _, err := os.Stat("./service.yaml")
	ch := make(chan fileChangeMsg)

	// Actions
	changes := tui.NewAction[[]fileChangeMsg]("fileChanges")
	uploading := tui.NewAction[openapi.ServicesPost200Response]("uploading")
	sp := spinner.New()

	model := ServicesSyncPage{
		isServiceDir:  true,
		channel:       ch,
		changes:       changes,
		watchDebounce: 500 * time.Millisecond,
		spinner:       sp,
		uploading:     uploading,
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
	switch msg := msg.(type) {
	case tui.ActionInit[openapi.ServicesPost200Response]:
		m.uploading.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[openapi.ServicesPost200Response]:
		m.uploading.ProcessResult(msg)
		if m.uploading.IsSuccess() {
			return m, m.collectChanges(nil)
		} else if m.uploading.IsError() {
			return m, m.collectChanges(m.changes.Data)
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
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// h, v := style.Docstyle.GetFrameSize()

	// Is it a key press?
	// case tea.KeyMsg:
	// 	// Cool, what was the actual key pressed?
	// 	switch msg.String() {
	// 	case "enter":
	// 		value := m.actions.SelectedItem().FilterValue()
	// 		if value != "" {
	// 			switch value {
	// 			case "Initialize":
	// 				return RootScreen(state.Get()).SwitchScreen(NewServiceInitPage())
	// 			case "Upload":
	// 				value = "service upload"
	// 			case "Update":
	// 				return RootScreen(state.Get()).SwitchScreen(NewServicesUpdatePage())
	// 			case "Download":
	// 				value = "service download"
	// 			}
	// 			// state.Get().Route.Push(value)
	// 			return m, tea.Quit
	// 		}
	// 	}
	// }

	// var cmd tea.Cmd
	// m.actions, cmd = m.actions.Update(msg)
	// return m, cmd

	return m, nil
}

func (m ServicesSyncPage) View() string {
	s := style.Title.Render("Synchronize service") + "\n\n"

	if m.uploading.IsLoading() {
		s += m.spinner.View() + " Syncing files..." + "\n\n"
	} else if m.uploading.IsError() {
		s += style.Error.Render("Failed to upload changes") + style.Gray.Render(" ("+m.uploading.Error.Error()+")") + "\n\n"
	}

	if m.changes.IsLoading() {
		s += m.spinner.View() + " Watching changes..." + "\n"
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

// Upload all collected changes to the Rover
func (m ServicesSyncPage) uploadChanges() tea.Cmd {
	return tui.PerformAction(&m.uploading, func() (*openapi.ServicesPost200Response, error) {
		// Create tmp zip file
		tmpZip := os.TempDir() + "/" + time.Now().Format("20060102150405") + ".zip"
		sourceDir := "."

		// Create the zip file
		zipFileWriter, err := os.Create(tmpZip)
		if err != nil {
			return nil, fmt.Errorf("failed to create zip file: %v", err)
		}
		defer zipFileWriter.Close()

		// Create a zip writer
		zipWriter := zip.NewWriter(zipFileWriter)
		defer zipWriter.Close()

		// Walk through the directory
		err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walking the path %q: %v", path, err)
			}

			// Get the relative path to store in the zip archive
			relPath, err := filepath.Rel(sourceDir, path)
			if err != nil {
				return err
			}

			// Skip the root directory itself
			if relPath == "." {
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
			return nil, err
		}

		// mock remove
		time.Sleep(2 * time.Second)

		if m.uploading.Attempt > 2 {
			return nil, fmt.Errorf("mocking an error here")
		}

		res := openapi.ServicesPost200Response{
			Name:    "service",
			Version: "1.0.0",
		}
		return &res, nil
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