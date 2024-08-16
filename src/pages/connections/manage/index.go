package manageconnectionspage

import (
	"fmt"
	"io"
	"strings"

	"github.com/VU-ASE/rover/src/configuration"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	keys keyMap
	list list.Model
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
	// Edit: key.NewBinding(
	// 	key.WithKeys("enter"),
	// 	key.WithHelp("enter", "edit"),
	// ),
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
	connection configuration.RoverConnection
	active     bool
}

func (i item) FilterValue() string { return i.connection.Name }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := i.connection.Name + lipgloss.NewStyle().Foreground(style.GrayPrimary).Render(fmt.Sprintf(" %s@%s", i.connection.Username, i.connection.Host))
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

func connectionsToListItems() []list.Item {
	// Get the connections from the state
	appState := state.Get()

	items := make([]list.Item, 0)

	for _, connection := range appState.RoverConnections.Available {
		items = append(items, item{
			connection: connection,
			active:     connection.Name == appState.RoverConnections.Active,
		})
	}

	return items
}

func InitialModel() model {
	l := list.New(connectionsToListItems(), itemDelegate{}, 0, 14)
	l.Title = "Manage Rover connections"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = style.TitleStyle
	l.Styles.PaginationStyle = style.PaginationStyle
	l.Styles.HelpStyle = style.HelpStyle
	l.AdditionalShortHelpKeys = keys.ShortHelp

	return model{
		list: l,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := style.Docstyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v) // leave some room for the header

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.MarkActive):
			if m.list.Index() >= 0 && m.list.Index() < len(m.list.Items()) {
				item := m.list.Items()[m.list.Index()].(item)
				state.Get().RoverConnections.Active = item.connection.Name
				m.list.SetItems(connectionsToListItems())
				return m, nil
			}
		case key.Matches(msg, keys.New):
			state.Get().CurrentView = "connect"
			return m, tea.Quit
		case key.Matches(msg, keys.Delete):
			if len(m.list.Items()) > 1 && m.list.Index() >= 0 && m.list.Index() < len(m.list.Items()) {
				item := m.list.Items()[m.list.Index()].(item)
				state.Get().RoverConnections = state.Get().RoverConnections.Remove(item.connection.Name)
				m.list.SetItems(connectionsToListItems())
				m.list.ResetSelected()
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
	return nil
}

func (m model) View() string {
	s := m.list.View()

	return style.Docstyle.Render(s)
}
