package views

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/VU-ASE/rover/src/style"
	"github.com/VU-ASE/rover/src/tui"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	git "github.com/go-git/go-git/v5"
)

// Persistent global state (ugly, yes) to allow retrying of connection checks by discarding results with an attempt number lower than the current one
var attemptNumber = 1

type ServiceInitFormValues struct {
	Name    string
	Author  string
	Source  string
	Version string
}

type ServiceInitPage struct {
	serviceAlreadyExists bool
	form                 *huh.Form
	spinner              spinner.Model
	serviceInitialized   tui.Action[bool]
	isInitializing       bool
	errors               []error // errors that occurred during the process
	selectedPreset       *string
	service              ServiceInitFormValues
}

func NewServiceInitPage() ServiceInitPage {
	s := spinner.New()
	s.Spinner = spinner.Line

	defaultAuthor := ""
	userDir, err := os.UserHomeDir()
	if err == nil {
		// Get the last part of the user directory
		_, defaultAuthor = filepath.Split(userDir)
	}

	// Check if the service already exists, in which case we will not initialize it
	_, err = os.Stat("./service.yaml")
	serviceAlreadyExists := err == nil

	service := ServiceInitFormValues{
		Name:    "",
		Author:  defaultAuthor,
		Source:  "github.com/username/repository",
		Version: "0.0.1",
	}

	// We create some files based on the selected preset
	selectedPreset := "golang"

	return ServiceInitPage{
		spinner:              s,
		serviceAlreadyExists: serviceAlreadyExists,
		selectedPreset:       &selectedPreset,
		errors:               []error{},
		isInitializing:       false,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("What is the name of your service?").
					CharLimit(255).
					Prompt("> ").
					Value(&service.Name).
					Validate(func(s string) error {
						if len(s) < 3 {
							return fmt.Errorf("Service names must be at least 3 characters long")
						}

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
						if len(s) < 3 {
							return fmt.Errorf("Author names must be at least 3 characters long")
						}

						valid := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(s)
						if !valid {
							return fmt.Errorf("Author names can only contain letters and numbers")
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
						// _, err := semver.NewVersion(s)
						// if err != nil {
						// 	return fmt.Errorf("Please enter a valid semantic version (e.g. 0.0.1)")
						// }
						return nil
					}),
				// Ask the user for a base burger and toppings.
				huh.NewSelect[string]().
					Title("Which programming language do you want to use?").
					Options(
						huh.NewOption("Go", "golang"),
						huh.NewOption("Rust", "rust"),
						huh.NewOption("Python", "python"),
						huh.NewOption("C", "c"),
						huh.NewOption("I will configure this myself", "none"),
					).
					Value(&selectedPreset),
			),
		).WithTheme(style.FormTheme),
	}
}

func (m ServiceInitPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case tui.ActionInit[bool]:
		m.serviceInitialized.ProcessInit(msg)
		return m, nil
	case tui.ActionResult[bool]:
		m.serviceInitialized.ProcessResult(msg)
		return m, nil

	default:

		cmds := []tea.Cmd{}
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			if f.State == huh.StateCompleted && !m.isInitializing {
				attemptNumber++
				m.isInitializing = true
				cmds = append(cmds, m.initializeTemplate())
			}
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}
}

// the update view with the view method
func (m ServiceInitPage) enterDetailsView() string {
	// Introduction
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Create a new service")

	s += "\n\n" + m.form.View()

	return s
}

func (m ServiceInitPage) initializationView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Initializing your service")

	s += "\n\n" + m.spinner.View() + " Downloading template and setting up service"

	return s
}

func (m ServiceInitPage) initializedSuccessView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Service initialized!")

	s += "\n\n" + "Your service has been initialized successfully. Explore the files in your service folder to get started.\nBelow you can find a quick overview of the files that have been created for you."

	s += "\n\n" + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("service.yaml") + "\nThis is the most important file in your service. It helps the Rover understand what your service does and how to start it."
	s += "\n\n" + lipgloss.NewStyle().Foreground(style.SuccessPrimary).Render("Makefile") + "\nThis file contains commands to quickly build and run your service by using " + lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("make start") + " and " + lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("make build") + "."

	return s
}

func (m ServiceInitPage) initializedFailureView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Could not initialize service")

	s += "\n\nAn error occurred while initializing your service"
	if len(m.errors) > 0 {
		for _, err := range m.errors {
			s += lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("\n - " + err.Error())
		}
	}

	return s
}
func (m ServiceInitPage) serviceAlreadyExistsView() string {
	s := lipgloss.NewStyle().Foreground(style.AsePrimary).Render("Cannot initialize service")

	s += "\n\nYou already initialized a service in this folder. \nIf you want to initialize a new service, create a sibling folder and try again."

	return s
}

func (m ServiceInitPage) Init() tea.Cmd {
	return tea.Batch(m.form.Init(), m.spinner.Tick)
}

func (m ServiceInitPage) View() string {
	if m.serviceAlreadyExists {
		return m.serviceAlreadyExistsView()
	} else if m.form.State == huh.StateCompleted && !m.serviceInitialized.Started {
		return m.initializationView()
	} else if m.form.State == huh.StateCompleted && m.serviceInitialized.IsSuccess() {
		return m.initializedSuccessView()
	} else if m.form.State == huh.StateCompleted && m.serviceInitialized.IsError() {
		return m.initializedFailureView()
	} else {
		return m.enterDetailsView()
	}
}

func (m ServiceInitPage) initializeTemplate() tea.Cmd {
	return tui.PerformAction(&m.serviceInitialized, func() (*bool, error) {

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
			return nil, err
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

		return nil, nil
	})
}

// This function downloads a selected template from a repository and places it in the destination folder
func downloadTemplate(repository string, destination string) error {
	_, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL: repository,
	})
	if err != nil {
		return err
	}

	// Remove the .git folder from the template
	err = os.RemoveAll(filepath.Join(destination, ".git"))
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