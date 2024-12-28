package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/VU-ASE/roverctl/src/configuration"
	"github.com/VU-ASE/roverctl/src/state"
	"github.com/VU-ASE/roverctl/src/style"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConnectionsManagePage struct {
	list list.Model
}

// connectionsManageKeyMap defines a set of keybindings. To work for help it must satisfy key.Map
type connectionsManageKeyMap struct {
	Edit       key.Binding
	Delete     key.Binding
	MarkActive key.Binding
	New        key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k connectionsManageKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.New, k.Edit, k.Delete, k.MarkActive}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k connectionsManageKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var connectionsManageKeys = connectionsManageKeyMap{
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
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

func NewConnectionsManagePage() ConnectionsManagePage {
	l := list.New(connectionsToListItems(), itemDelegate{}, 0, 14)
	l.Title = "Manage Rover connections"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = style.TitleStyle
	l.Styles.PaginationStyle = style.PaginationStyle
	l.Styles.HelpStyle = style.HelpStyle
	l.AdditionalShortHelpKeys = connectionsManageKeys.ShortHelp

	return ConnectionsManagePage{
		list: l,
	}
}

func (m ConnectionsManagePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := style.Docstyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v) // leave some room for the header

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, connectionsManageKeys.MarkActive):
			if m.list.Index() >= 0 && m.list.Index() < len(m.list.Items()) {
				item := m.list.Items()[m.list.Index()].(item)
				state.Get().RoverConnections.Active = item.connection.Name
				m.list.SetItems(connectionsToListItems())
				return m, nil
			}
		case key.Matches(msg, connectionsManageKeys.New):
			return RootScreen(state.Get()).SwitchScreen(NewConnectionsInitPage(nil))
		case key.Matches(msg, connectionsManageKeys.Delete):
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
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ConnectionsManagePage) Init() tea.Cmd {
	return nil
}

func (m ConnectionsManagePage) View() string {
	s := m.list.View()

	return style.Docstyle.Render(s)
}
