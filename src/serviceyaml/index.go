package serviceyaml

import (
	"os"

	"gopkg.in/yaml.v3"
)

// This should all become included in the roverlib package, that handles parsing and writing the service.yaml files

type RoverService struct {
	Name          string                            `yaml:"name"`
	Author        string                            `yaml:"author"`
	Source        string                            `yaml:"source"`
	Version       string                            `yaml:"version"`
	Inputs        []RoverServiceInput               `yaml:"inputs"`
	Outputs       []string                          `yaml:"outputs"`
	Configuration []RoverServiceConfigurationOption `yaml:"configuration"`
}

type RoverServiceInput struct {
	Service string   `yaml:"service"`
	Streams []string `yaml:"streams"`
}

type RoverServiceConfigurationOption struct {
	Name    string `yaml:"name"`
	Value   string `yaml:"value"`
	Type    string `yaml:"type"` // string, int, float or undefined for autoparsing
	Tunable bool   `yaml:"tunable"`
}

// Parse the YAML content
// NB: having a custom Parse function allows us to set default values for the struct
func Parse(content []byte) (*RoverService, error) {
	service := &RoverService{}
	err := yaml.Unmarshal(content, service)

	return service, err
}

// Parse the YAML file at the given path
func ParseFrom(path string) (*RoverService, error) {
	// Read the file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Parse(content)
}
