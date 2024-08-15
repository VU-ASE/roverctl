package asciitool

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

var path = ""

// Functions used to draw ASCII graphs from mermaid content
func Init(folder string) error {
	// Check if the folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return err
	}

	// Check if the folder already contains the mermaid-ascii repository, if so set the path and build it
	if _, err := os.Stat(folder + "/mermaid-ascii"); err == nil {
		path = folder + "/mermaid-ascii"
	} else {
		// Clone the mermaid-ascii repository from github
		cmd := exec.Command("git", "clone", "https://github.com/AlexanderGrooff/mermaid-ascii.git")
		cmd.Dir = folder
		out, err := cmd.Output()
		log.Info().Msgf("Cloning ascii-mermaid: %s", out)
		if err != nil {
			return err
		}
	}

	// Build the mermaid-ascii repository
	return build(folder + "/mermaid-ascii")
}

func build(folder string) error {
	// Check if the path is set
	if folder == "" {
		return nil
	}

	// Check if the binary already exists
	if _, err := os.Stat(folder + "/mermaid-ascii"); err == nil {
		return nil
	}

	// Build the mermaid-ascii repository specified in the path
	cmd := exec.Command("go", "build")
	cmd.Dir = folder
	out, err := cmd.Output()
	log.Info().Msgf("Building ascii-mermaid: %s", out)
	return err
}

func Draw(mermaidContent string) (string, error) {
	// Check if the path is set
	if path == "" {
		return "", fmt.Errorf("The path to the mermaid-ascii repository is not set")
	}

	// Create a temporary file to store the mermaid content
	f, err := os.CreateTemp("", "mermaid-ascii-*.mmd")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())

	// Write the mermaid content to the temporary file
	if _, err := f.WriteString(mermaidContent); err != nil {
		return "", err
	}

	// Run the mermaid-ascii binary on the temporary file
	cmd := exec.Command(path+"/mermaid-ascii", "-f", f.Name())
	out, err := cmd.Output()
	return string(out), err
}
