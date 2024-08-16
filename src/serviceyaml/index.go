package serviceyaml

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
