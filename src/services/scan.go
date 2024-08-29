package services

import (
	"io"
	"strings"

	"github.com/VU-ASE/rover/src/configuration"
	"github.com/VU-ASE/rover/src/serviceyaml"
	"github.com/pkg/sftp"
)

type FoundService struct {
	Path    string                   // the path (from root) to the service.yaml file found
	Service serviceyaml.RoverService // the actual service found
}

// Scan all downloaded services on the Rover, starting from the default services directory
func Scan(conn configuration.RoverConnection) ([]FoundService, error) {
	found := []FoundService{}

	// Set up SSH client
	sshconn, err := conn.ToSsh()
	if err != nil {
		return nil, err
	}

	// Set up SFTP client
	sshforsftp, err := conn.ToSshConnection()
	if err != nil {
		return nil, err
	}
	defer sshforsftp.Close()
	sftpclient, err := sftp.NewClient(sshforsftp)
	if err != nil {
		return nil, err
	}
	defer sftpclient.Close()

	// Find all paths to service.yaml (or service*.yaml paths) files
	res, err := sshconn.Run("find " + configuration.RemoteServiceDir + " -name 'service*.yaml'")
	if err != nil {
		return found, err
	}

	// Parse the paths, one per line
	paths := strings.Split(string(res), "\n")

	// Per path, read the service.yaml file and check if this is a valid service
	for _, path := range paths {
		if path == "" {
			continue
		}

		// Read the remote file
		file, err := sftpclient.Open(path)
		if err != nil {
			continue
		}
		contents, err := io.ReadAll(file)
		if err != nil {
			continue
		}

		// Parse the service.yaml file
		service, err := serviceyaml.Parse(contents)
		if err != nil {
			continue
		}

		// Add the service to the list
		found = append(found, FoundService{
			Path:    path,
			Service: *service,
		})
	}

	return found, nil
}
