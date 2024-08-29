package roveryaml

import (
	"io"
	"os"
	"slices"

	"github.com/VU-ASE/rover/src/configuration"
	"github.com/pkg/sftp"
	"gopkg.in/yaml.v3"
)

//
// Everything for working with rover.yaml files and their Golang representation
//

type RoverConfig struct {
	Downloaded []ServiceDownload `yaml:"downloaded"`
	// Every string is a path to a service on the Rover. The path should lead to a service.yaml file
	// or to a directory containing a service.yaml file.
	Enabled []string `yaml:"enabled"`
}

// Describes a downloaded service from a given source
type ServiceDownload struct {
	Name    string `yaml:"name"`
	Source  string `yaml:"source"`
	Version string `yaml:"version"`
}

func (c *RoverConfig) Enable(path string) {
	if c == nil {
		return
	}

	c.Enabled = append(c.Enabled, path)
}

func (c *RoverConfig) Disable(path string) {
	if c == nil {
		return
	}

	c.Enabled = slices.DeleteFunc(
		c.Enabled,
		func(p string) bool {
			return p == path
		},
	)
}

func (c *RoverConfig) Toggle(path string) {
	if c.HasEnabled(path) {
		c.Disable(path)
	} else {
		c.Enable(path)
	}
}

func (c *RoverConfig) HasEnabled(path string) bool {
	if c == nil {
		return false
	}

	for _, p := range c.Enabled {
		if p == path {
			return true
		}
	}

	return false
}

var defaultPath = configuration.RemoteConfigDir + "/rover.yaml"

// Parse the YAML file at the given path
func Parse(path string) (*RoverConfig, error) {
	// Read the file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse YAML
	config := &RoverConfig{}
	err = yaml.Unmarshal(content, config)

	return config, err
}

// Load the RoverConfig from the default save location on a Rover
func Load(conn configuration.RoverConnection) (*RoverConfig, error) {
	res := &RoverConfig{
		Downloaded: []ServiceDownload{},
		Enabled:    []string{},
	}

	// Establish an SSH connection
	sshconn, err := conn.ToSshConnection()
	if err != nil {
		return res, err
	}
	defer sshconn.Close()

	// Create an SFTP client
	client, err := sftp.NewClient(sshconn)
	if err != nil {
		return res, err
	}
	defer client.Close()

	// Check if the file exists
	file, err := client.Open(defaultPath)
	if err != nil {
		return res, nil
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		return res, err
	}

	// Parse the YAML content
	err = yaml.Unmarshal(contents, res)
	return res, err
}

// Save the RoverConfig to the default save location on a Rover
func (c *RoverConfig) Save(conn configuration.RoverConnection) error {
	// Establish an SSH connection
	sshconn, err := conn.ToSshConnection()
	if err != nil {
		return err
	}
	defer sshconn.Close()

	// Create an SFTP client
	client, err := sftp.NewClient(sshconn)
	if err != nil {
		return err
	}
	defer client.Close()

	// Marshal the config to YAML
	content, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// Write the file, overwriting the existing one
	file, err := client.Create(defaultPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	return err
}
