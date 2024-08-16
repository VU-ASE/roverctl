package initservicepage

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	semver "github.com/Masterminds/semver/v3"
	"github.com/VU-ASE/rover/src/serviceyaml"
	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/grignaak/tribool.v1"
	"gopkg.in/src-d/go-git.v4"
)

// Persistent global state (ugly, yes) to allow retrying of connection checks by discarding results with an attempt number lower than the current one
var attemptNumber = 1

// Action codes
const (
	initializeService = "initializeService"
)

// Used to communicate the result of various tests
type resultMsg struct {
	action  string
	result  bool
	err     error
	attempt int
}

type model struct {
	form               *huh.Form
	spinner            spinner.Model
	serviceInitialized tribool.Tribool // download template and replace magic strings in it
	isInitializing     bool
	errors             []error // errors that occurred during the process
	service            *serviceyaml.RoverService
	selectedPreset     *string
}

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Line

	defaultAuthor, err := os.Hostname()
	if err != nil {
		defaultAuthor = ""
	}

	service := serviceyaml.RoverService{
		Name:          "",
		Author:        defaultAuthor,
		Source:        "github.com/username/repository",
		Version:       "0.0.1",
		Inputs:        []serviceyaml.RoverServiceInput{},
		Outputs:       []string{},
		Configuration: []serviceyaml.RoverServiceConfigurationOption{},
	}

	// We create some files based on the selected preset
	selectedPreset := "golang"

	return model{
		spinner:            s,
		service:            &service,
		selectedPreset:     &selectedPreset,
		serviceInitialized: tribool.Maybe,
		errors:             []error{},
		isInitializing:     false,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("What is the name of your service?").
					CharLimit(255).
					Prompt("> ").
					Value(&service.Name).
					Validate(func(s string) error {
						valid := regexp.MustCompile(`^[a-z0-9]*$`).MatchString(s)
						if !valid {
							return fmt.Errorf("Service names can only contain lowercase letters and numbers")
						}

						return nil
					}),
				huh.NewInput().
					Title("Who is the author of this service?").
					CharLimit(255).
					Prompt("> ").
					Value((&service.Author)).
					Validate(func(s string) error {
						valid := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(s)
						if !valid {
							return fmt.Errorf("Service names can only contain letters and numbers")
						}
						return nil
					}),
				huh.NewInput().
					Title("Where is this service published?").
					CharLimit(255).
					Prompt("> ").
					Value(&service.Source).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("Enter a valid source URL")
						}
						if strings.Contains(s, "username") || strings.Contains(s, "repository") {
							return fmt.Errorf("Please replace 'username' and 'repository' with your actual GitHub username/organization name and repository name")
						}
						if strings.Contains(s, "https://") || strings.Contains(s, "http://") || strings.Contains(s, "www.") {
							return fmt.Errorf("Do not include the protocol or 'www.' in the URL")
						}
						return nil
					}),
				huh.NewInput().
					Title("At what semantic version do you want to start?").
					CharLimit(255).
					Prompt("> ").
					Value(&service.Version).
					Validate(func(s string) error {
						// Try to parse the version
						_, err := semver.NewVersion(s)
						if err != nil {
							return fmt.Errorf("Please enter a valid semantic version (e.g. 0.0.1)")
						}
						return nil
					}),
				// Ask the user for a base burger and toppings.
				huh.NewSelect[string]().
					Title("Which programming language do you want to use?").
					Options(
						huh.NewOption("Go (Golang)", "golang"),
						huh.NewOption("Python", "python"),
						huh.NewOption("I will configure this myself", "none"),
					).
					Value(&selectedPreset),
			),
		).WithTheme(style.FormTheme),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.form.State == huh.StateCompleted {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			}
		}
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case resultMsg:
		if msg.attempt < attemptNumber {
			return m, nil
		}
		switch msg.action {
		case initializeService:
			m.serviceInitialized = tribool.FromBool(msg.result)
			if msg.err != nil {
				m.errors = append(m.errors, msg.err)
			}
		}
		return m, nil
	default:
		// Base command
		model, cmd := tui.Update(m, msg)
		if cmd != nil {
			return model, cmd
		}

		cmds := []tea.Cmd{}
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			if f.State == huh.StateCompleted && !m.isInitializing {
				attemptNumber++
				m.isInitializing = true
				cmds = append(cmds, initializeTemplate(m, attemptNumber))
			}
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}
}

// the update view with the view method
func (m model) enterDetailsView() string {
	// Introduction
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Create a new service")

	s += "\n\n" + m.form.View()

	return style.Docstyle.Render(s)
}

func (m model) initializationView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Initializing your service")

	s += "\n\n" + m.spinner.View() + " Downloading template and setting up service"

	return s
}

func (m model) initializedSuccessView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Service initialized!")

	s += "\n\n" + "Your service has been initialized successfully. Explore the files in your service folder to get started.\nBelow you can find a quick overview of the files that have been created for you."

	s += "\n\n" + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("service.yaml") + "\nThis is the most important file in your service. It helps the Rover understand what your service does and how to start it."
	s += "\n\n" + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Makefile") + "\nThis file contains commands to quickly build and run your service by using " + lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("make start") + " and " + lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("make build") + "."

	return s
}

func (m model) initializedFailureView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Could not initialize service")

	s += lipgloss.NewStyle().Foreground(style.WarningPrimary).Render("\n\nAn error occurred while initializing your service")
	if len(m.errors) > 0 {
		for _, err := range m.errors {
			s += lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("\n - " + err.Error())
		}
	}

	return s
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.form.Init(), m.spinner.Tick)
}

func (m model) View() string {
	if m.form.State == huh.StateCompleted && m.serviceInitialized == tribool.Maybe {
		return style.Docstyle.Render(m.initializationView())
	} else if m.form.State == huh.StateCompleted && m.serviceInitialized == tribool.True {
		return style.Docstyle.Render(m.initializedSuccessView())
	} else if m.form.State == huh.StateCompleted && m.serviceInitialized == tribool.False {
		return style.Docstyle.Render(m.initializedFailureView())
	} else {
		return m.enterDetailsView()
	}
}

func initializeTemplate(m model, a int) tea.Cmd {
	return func() tea.Msg {
		// Based on the programming language chosen, download a specific template and replace the magic strings in it
		templateRepo := "unsupported"
		switch *m.selectedPreset {
		case "golang":
			templateRepo = "https://github.com/VU-ASE/service-template-go"
		case "python":
			templateRepo = "unsupported"
		}

		err := downloadTemplate(templateRepo, ".")
		if err != nil {
			return resultMsg{result: false, err: err, action: initializeService, attempt: a}
		}

		// Strings to be replaced
		toreplace := map[string]string{
			"$SERVICE_NAME":    m.service.Name,
			"$SERVICE_AUTHOR":  m.service.Author,
			"$SERVICE_VERSION": m.service.Version,
			"$SERVICE_SOURCE":  m.service.Source,
		}

		// Replace the magic strings in the template
		_ = replaceMagicStrings("service.yaml", toreplace)
		_ = replaceMagicStrings("Makefile", toreplace)
		_ = replaceMagicStrings("go.mod", toreplace)

		return resultMsg{result: true, err: nil, action: initializeService, attempt: a}
	}
}

// This function downloads a selected template from a repository and places it in the destination folder
func downloadTemplate(repository string, destination string) error {
	_, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL: repository,
	})
	return err
}

func replaceMagicStrings(filepath string, replacements map[string]string) error {
	// Read the file
	content, err := os.ReadFile(filepath)
	// If the file does not exist, we don't return an error, we just skip
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	// Replace the magic strings
	for key, value := range replacements {
		content = []byte(strings.ReplaceAll(string(content), key, value))
	}

	// Write the file back
	err = os.WriteFile(filepath, content, 0644)
	return err
}
