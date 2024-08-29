package configuration

import (
	"fmt"
	"os"

	"github.com/VU-ASE/rover/src/configuration/asciitool"
)

// This is where the Rover configuration files are saved both remotely
// and on the local machine of the user.
// Directories will be created if they do not exist

func LocalConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/home/.rover"
	}

	return home + "/.rover"
}

const RemoteConfigDir = "/home/debix/.rover/config"
const RemoteServiceDir = "/home/debix/.rover/services" // this directory holds all service folders. Each subfolder represents a service

// Initialize the configuration directory
func Initialize() error {
	// Check if we are root
	uid := os.Geteuid()
	if uid == 0 {
		return fmt.Errorf("You should not run this utility as root.")
	}

	// Create the configuration directory if it does not exist
	if err := os.MkdirAll(LocalConfigDir(), 0755); err != nil {
		return err
	}

	// // Lock the Rover
	// if err := Lock(); err != nil {
	// 	return err
	// }

	// Initialize the ascii tool that we need
	return asciitool.Init(LocalConfigDir())
}

// Cleanup the configuration directory
func Cleanup() {
	// Unlock the Rover
	// if err := Unlock(); err != nil {
	// 	log.Error().Err(err).Msg("An error occurred while unlocking the Rover.")
	// }

}
