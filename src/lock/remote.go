package roverlock

// Functions used to lock and unlock the Rover over a remote connection.
// This avoids multiple users from running locked actions at the same time.

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/VU-ASE/rover/src/configuration"
	"github.com/pkg/sftp"
)

// var lockfileName = configuration.RemoteConfigDir + "/.roverlock"
var lockfileName = ".roverlock"

// Functions to lock the Rover so that only one user can configure the Rover at a time

// Lock the Rover
func Lock(conn configuration.RoverConnection) error {
	// Create an ssh client
	client, err := conn.ToSSH()
	if err != nil {
		return err
	}
	defer client.Close()

	// Create an sftp client
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftp.Close()

	// Check if the lockfile exists and if it holds our hostname
	file, err := sftp.Open(lockfileName)
	if err == nil {
		defer file.Close()

		// See if the file holds our hostname
		b := make([]byte, 255)
		n, err := file.Read(b)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		// Check if the hostname matches
		hname, err := os.Hostname()
		if err != nil {
			return err
		}
		if string(b[:n]) != hname {
			return fmt.Errorf("The Rover is locked by '%s', so it cannot be locked by you.", strings.ReplaceAll((string(b[:n])), "\n", ""))
		}

		// The lockfile holds our hostname, so we locked already
		return nil
	}

	// Create the lockfile, and write the hostname to it
	f, err := sftp.Create(lockfileName)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the hostname to the lockfile
	hname, err := os.Hostname()
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(hname))
	return err
}

// Unlock the rover
func Unlock(conn configuration.RoverConnection) error {
	// Create an ssh client
	client, err := conn.ToSSH()
	if err != nil {
		return err
	}
	defer client.Close()

	// Create an sftp client
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftp.Close()

	// Check if the lockfile exists and if it holds our hostname
	file, err := sftp.Open(lockfileName)
	if err == nil {
		// See if the file holds our hostname
		b := make([]byte, 100)
		n, err := file.Read(b)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		// Check if the hostname matches
		hname, err := os.Hostname()
		if err != nil {
			return err
		}
		if string(b[:n]) != hname {
			return fmt.Errorf("The Rover is locked by '%s', so it cannot be unlocked by you.", strings.ReplaceAll((string(b[:n])), "\n", ""))
		}

		// Remove the lockfile
		err = sftp.Remove(lockfileName)
		return err
	}
	// No lockfile exists, so we are already unlocked
	return nil
}
