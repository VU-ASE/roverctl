package services

import (
	"github.com/VU-ASE/rover/src/configuration"
)

// Clear all services from the Rover
func ClearAll(conn configuration.RoverConnection) error {
	// Create an SSH connection
	sshconn, err := conn.ToSsh()
	if err != nil {
		return err
	}
	defer sshconn.Close()

	// Remove the directory recursively
	_, err = sshconn.Run("rm -rf " + configuration.RemoteServiceDir)
	return err
}
