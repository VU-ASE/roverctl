package configuration

import (
	"os"

	"github.com/VU-ASE/rover/src/configuration/asciitool"
	// "github.com/rs/zerolog/log"
)

// This is where the Rover configuration files are saved both remotely
// and on the local machine of the user.
// Directories will be created if they do not exist

const LocalConfigDir = "/etc/rover"
const RemoteConfigDir = "/etc/rover"

// Initialize the configuration directory
func Initialize() error {
	// Create the configuration directory if it does not exist
	if err := os.MkdirAll(LocalConfigDir, 0755); err != nil {
		return err
	}

	// // Lock the Rover
	// if err := Lock(); err != nil {
	// 	return err
	// }

	// Initialize the ascii tool that we need
	return asciitool.Init(LocalConfigDir)
}

// Cleanup the configuration directory
func Cleanup() {
	// Unlock the Rover
	// if err := Unlock(); err != nil {
	// 	log.Error().Err(err).Msg("An error occurred while unlocking the Rover.")
	// }
}
