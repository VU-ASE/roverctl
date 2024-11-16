package views

import (
	"runtime"

	"github.com/VU-ASE/rover/src/configuration"
	"github.com/VU-ASE/rover/src/style"
	tea "github.com/charmbracelet/bubbletea"
)

var version = "UNSET"

type InfoPage struct {
	// To select an action to perform with this utility
}

func NewInfoPage() InfoPage {
	return InfoPage{}
}

func (m InfoPage) Init() tea.Cmd {
	return nil
}

func (m InfoPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m InfoPage) View() string {
	s := `
 _    ____  __       ___   _____ ______
| |  / / / / /      /   | / ___// ____/
| | / / / / /      / /| | \__ \/ __/   
| |/ / /_/ /      / ___ |___/ / /___   
|___/\____/      /_/  |_/____/_____/        
                      `

	s += "\n\n" + style.Gray.Render("Roverctl version: ") + version + "\n"
	s += style.Gray.Render("Configuration location: ") + configuration.LocalConfigDir() + "\n"
	s += style.Gray.Render("Architecture: ") + runtime.GOOS + "/" + runtime.GOARCH + "\n"

	return s
}
