package views

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/VU-ASE/rover/src/openapi"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ServicesUpdateKeyMap defines a set of keybindings. To work for help it must satisfy key.Map
type ServicesUpdateKeyMap struct {
	Retry   key.Binding
	Confirm key.Binding
	Quit    key.Binding
}

type ServicesUpdatePage struct {
	help           help.Model
	spinner        spinner.Model
	sourceList     tui.Action[[]openapi.SourcesGet200ResponseInner]
	serviceUpdates map[string]tui.Action[openapi.SourcesNamePost200Response]
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k ServicesUpdateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Retry, k.Confirm, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k ServicesUpdateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var errorFetchSourcesKeys = ServicesUpdateKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

var successFetchSourcesKeys = ServicesUpdateKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

// Shown when the services are being updated
var updateServicesKeys = ServicesUpdateKeyMap{
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func NewServicesUpdatePage() ServicesUpdatePage {
	s := spinner.New()
	s.Spinner = spinner.Line

	sourcesList := tui.NewAction[[]openapi.SourcesGet200ResponseInner]("sourcesList")
	servicesList := map[string]tui.Action[openapi.SourcesNamePost200Response]{}

	return ServicesUpdatePage{
		spinner:        s,
		help:           help.New(),
		sourceList:     sourcesList,
		serviceUpdates: servicesList,
	}
}

func (m ServicesUpdatePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, successFetchSourcesKeys.Confirm):
			if m.sourceList.IsSuccess() && len(m.serviceUpdates) <= 0 {
				cmds := []tea.Cmd{}
				for _, source := range *m.sourceList.Data {
					m.serviceUpdates[*source.Name] = tui.NewAction[openapi.SourcesNamePost200Response](*source.Name)
					cmds = append(cmds, updateService(m, *source.Name))
				}
				return m, tea.Batch(cmds...)
			}
			return m, nil
		case key.Matches(msg, errorFetchSourcesKeys.Retry):
			// Are there any sources being updated currently?
			updateOngoing := false
			for _, action := range m.serviceUpdates {
				if action.IsLoading() {
					updateOngoing = true
					break
				}
			}

			if m.sourceList.IsDone() && !updateOngoing {
				// Clear the service updates
				m.serviceUpdates = make(map[string]tui.Action[openapi.SourcesNamePost200Response])
				return m, fetchSources(m)
			}

			return m, nil
		// Restore to the initial form, but recover the form values
		// m = InitialModel(m.formValues)
		// return m, tea.Batch(m.form.Init(), m.spinner.Tick)
		// case key.Matches(msg, keys.Retry):
		// 	// Retry the connection checks
		// 	// m.isChecking = true
		// 	// return m, tea.Batch(checkRoute(m), checkAuth(m), checkRoverdVersion(m), checkRoverNumber(m))

		default:
			return m, nil
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tui.ActionInit[[]openapi.SourcesGet200ResponseInner]:
		m.sourceList.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[[]openapi.SourcesGet200ResponseInner]:
		m.sourceList.ProcessResult(msg)
		return m, nil
	case tui.ActionInit[openapi.SourcesNamePost200Response]:
		newServiceUpdates := make(map[string]tui.Action[openapi.SourcesNamePost200Response])
		for k, v := range m.serviceUpdates {
			newServiceUpdates[k] = v
		}

		if action, ok := newServiceUpdates[msg.Name]; ok {
			action.ProcessInit(msg)
			newServiceUpdates[msg.Name] = action
		}

		m.serviceUpdates = newServiceUpdates
		return m, nil
	case tui.ActionResult[openapi.SourcesNamePost200Response]:
		newServiceUpdates := make(map[string]tui.Action[openapi.SourcesNamePost200Response])
		for k, v := range m.serviceUpdates {
			newServiceUpdates[k] = v
		}

		if action, ok := newServiceUpdates[msg.Name]; ok {
			action.ProcessResult(msg)
			newServiceUpdates[msg.Name] = action // Reassign the updated action back
		}

		m.serviceUpdates = newServiceUpdates
		return m, nil
	default:
		return m, nil
	}
}

func (m ServicesUpdatePage) fetchSourcesView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Update sources")

	if m.sourceList.IsLoading() {
		s += "\n\n " + m.spinner.View() + " Fetching sources..."
		return s
	}

	if m.sourceList.IsSuccess() {
		if len(*m.sourceList.Data) == 0 {
			s += "\n\n This rover has no enabled sources, nothing to update"
		} else {
			s += "\n\n The following " + strconv.Itoa(len(*m.sourceList.Data)) + " source" + func() string {
				if len(*m.sourceList.Data) > 1 {
					return "s"
				}
				return ""
			}() + " will be updated:\n"

			for _, source := range *m.sourceList.Data {
				s += "\n - " + lipgloss.NewStyle().Bold(true).Render(*source.Name) + " " + lipgloss.NewStyle().Foreground(style.AsePrimary).Render(*source.Url) + " " + lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("(now at v"+*source.Version+")")
			}

			s += "\n\n" + m.help.View(successFetchSourcesKeys)
		}
	} else {
		s += "\n\n " + lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Failed to fetch sources")
		s += "\n\n" + m.help.View(errorFetchSourcesKeys)
	}

	return s
}

func (m ServicesUpdatePage) testConnectionView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Updating sources") + "\n"

	errorUpdates := 0
	successUpdates := 0

	// Convert the map to a slice, alphabetically sorted by key name
	keys := make([]string, 0, len(m.serviceUpdates))
	for k := range m.serviceUpdates {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		s += "\n - " + lipgloss.NewStyle().Bold(true).Render(name) + " \n  "
		action := m.serviceUpdates[name]

		if action.IsLoading() {
			s += " " + m.spinner.View() + lipgloss.NewStyle().Foreground(style.GrayPrimary).Render(" Updating...")
		} else if action.IsSuccess() {
			successUpdates++

			// Find the old version
			oldVersion := "unknown"
			if source := m.sourceList.Data; source != nil {
				for _, s := range *source {
					if *s.Name == name {
						oldVersion = *s.Version
						break
					}
				}
			}

			if oldVersion != *action.Data.Version {
				s += lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render(" ✓ Updated v" + oldVersion + " -> v" + *action.Data.Version)
			} else {
				s += lipgloss.NewStyle().Foreground(style.AsePrimary).Render(" ✓ Unchanged at v" + *action.Data.Version)
			}

		} else if action.IsError() {
			errorUpdates++
			s += lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render(" ✗ Failed to update")
			if action.Error != nil {
				s += style.Gray.Render(" (" + action.Error.Error() + ")")
			}
		}
	}

	if errorUpdates+successUpdates == len(m.serviceUpdates) {
		s += "\n\n" + m.help.View(updateServicesKeys)
	}

	return s
}

func (m ServicesUpdatePage) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchSources(m))
}

func (m ServicesUpdatePage) View() string {
	if len(m.serviceUpdates) > 0 {
		return m.testConnectionView()
	}

	return m.fetchSourcesView()
}

func fetchSources(m ServicesUpdatePage) tea.Cmd {
	return tui.PerformAction(&m.sourceList, func() (*[]openapi.SourcesGet200ResponseInner, error) {
		// Wait 10 seconds
		time.Sleep(1 * time.Second)
		// return nil, errors.New("Failed to connect")

		// Mock sources
		//! remove

		sources := []openapi.SourcesGet200ResponseInner{
			{
				Name:    openapi.PtrString("source1"),
				Version: openapi.PtrString("1.0.0"),
				Url:     openapi.PtrString("https://example.com/source1"),
				Sha:     openapi.PtrString("1234567890"),
			},
			{
				Name:    openapi.PtrString("source2a"),
				Version: openapi.PtrString("1.0.0"),
				Url:     openapi.PtrString("https://example.com/source2"),
				Sha:     openapi.PtrString("1234567890"),
			},
			{
				Name:    openapi.PtrString("source3"),
				Version: openapi.PtrString("1.0.1"),
				Url:     openapi.PtrString("https://example.com/source1"),
				Sha:     openapi.PtrString("1234567890"),
			},
		}

		return &sources, nil
	})
}

func updateService(m ServicesUpdatePage, name string) tea.Cmd {
	update := m.serviceUpdates[name]
	return tui.PerformAction(&update, func() (*openapi.SourcesNamePost200Response, error) {
		// Wait for a random duration between 1 and 10 seconds
		time.Sleep(time.Duration(1+rand.Intn(10)) * time.Second)

		// Mock sources
		//! remove
		if len(name)%2 == 0 {
			res := openapi.SourcesNamePost200Response{
				Version: openapi.PtrString("1.0.1"),
				New:     openapi.PtrBool(true),
			}
			return &res, fmt.Errorf("Failed to update %s", name)
		}

		res := openapi.SourcesNamePost200Response{
			Version: openapi.PtrString("1.0.1"),
			New:     openapi.PtrBool(true),
		}
		return &res, nil
	})
}
