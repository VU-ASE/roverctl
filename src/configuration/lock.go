package configuration

import (
	"errors"
	"os"
	"strconv"
)

var lockfileName = RemoteConfigDir + "/.roverlock"

// Functions to lock the Rover so that only one user can configure the Rover at a time

// Lock the Rover
func Lock() error {
	// Does the lockfile already exist?
	if _, err := os.Stat(lockfileName); err == nil {
		return errors.New("The Rover utility is already in use. Close the other instance and try again.")
	}

	// Create the lockfile, and write the PID to it
	f, err := os.Create(lockfileName)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the PID to the lockfile
	_, err = f.WriteString(strconv.Itoa(os.Getpid()))
	if err != nil {
		return err
	}

	return nil
}

// Unlock the rover
func Unlock() error {
	// Check if the lockfile exists
	if _, err := os.Stat(lockfileName); err != nil {
		return errors.New("Tried to unlock Rover, but it was never locked.")
	}

	// Read the lockfile and check that the PID matches the current process
	f, err := os.Open(lockfileName)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read the PID from the lockfile
	b := make([]byte, 100)
	n, err := f.Read(b)
	if err != nil {
		return err
	}

	// Convert the PID to an integer
	pid, err := strconv.Atoi(string(b[:n]))
	if err != nil {
		return err
	}

	// Check that the PID matches the current process
	if pid != os.Getpid() {
		return errors.New("Cannot unlock Rover. It is locked by another process.")
	}

	// Remove the lockfile
	err = os.Remove(lockfileName)
	if err != nil {
		return err
	}

	return nil
}
