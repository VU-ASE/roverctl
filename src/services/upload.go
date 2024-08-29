package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/VU-ASE/rover/src/configuration"
	"github.com/VU-ASE/rover/src/serviceyaml"
	"github.com/pkg/sftp"
)

// Functions for uploading rover services to the Rover over a remote connection

// This will upload the current working directory as a service to the Rover
func Upload(conn configuration.RoverConnection) error {
	// Check if this service contains a service.yaml file
	_, err := os.Stat("./service.yaml")
	if err != nil {
		return fmt.Errorf("Cannot upload service: no service.yaml file found in the current directory. Initialize a service first.")
	}

	// Parse the service.yaml file
	yaml, err := serviceyaml.ParseFrom("./service.yaml")
	if err != nil {
		return fmt.Errorf("Failed to parse service.yaml file: %v", err)
	}

	// Create an SSH client
	sshconn, err := conn.ToSshConnection()
	if err != nil {
		return err
	}
	defer sshconn.Close()

	// Connect over SSH
	sshclient, err := conn.ToSsh()
	if err != nil {
		return err
	}
	defer sshclient.Close()

	// Create an SFTP client
	client, err := sftp.NewClient(sshconn)
	if err != nil {
		return err
	}
	defer client.Close()

	// Get the folder that we are in
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// The path to upload to
	remotePath := configuration.RemoteServiceDir + "/" + yaml.Author + "/" + yaml.Name + "/" + yaml.Version
	localPath := wd

	// Delete the remote directory if it exists (we will overwrite it)
	// we will ignore the error if the directory does not exist
	// (we can't use the sftp.RemoveDirectory function, which gives SSH_FX_ERROR)
	_, err = sshclient.Run("rm -rf " + remotePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove remote directory %s: %v", remotePath, err)
	}

	// Create the service directory if it does not exist yet
	_, err = client.Stat(remotePath)
	if err != nil {
		err = client.MkdirAll(remotePath)
		if err != nil {
			return err
		}
	}

	// Copy the local working directory to the remote directory
	err = filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create the corresponding remote path
		relPath := path[len(localPath):]
		remotePath := filepath.Join(remotePath, relPath)

		if info.IsDir() {
			// Create the remote directory
			err = client.MkdirAll(remotePath)
			if err != nil {
				return fmt.Errorf("failed to create remote directory %s: %v", remotePath, err)
			}
		} else if path == "rover" {
			// Excluded
			return nil
		} else {
			// Upload the file
			srcFile, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open local file %s: %v", path, err)
			}
			defer srcFile.Close()

			dstFile, err := client.Create(remotePath)
			if err != nil {
				return fmt.Errorf("failed to create remote file %s: %v", remotePath, err)
			}
			defer dstFile.Close()

			_, err = srcFile.Seek(0, 0) // rewind to start
			if err != nil {
				return fmt.Errorf("failed to seek local file %s: %v", path, err)
			}

			_, err = dstFile.ReadFrom(srcFile)
			if err != nil {
				return fmt.Errorf("failed to copy file %s to %s: %v", path, remotePath, err)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walk local directory: %v", err)
	}

	return nil
}
