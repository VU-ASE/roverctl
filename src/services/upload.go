package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/VU-ASE/rover/src/configuration"
	"github.com/VU-ASE/rover/src/serviceyaml"
	"github.com/VU-ASE/rover/src/state"
	"github.com/melbahja/goph"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Functions for uploading rover services to the Rover over a remote connection

// This will upload the current working directory to the Rover as a sservice
func UploadService() error {
	// Check if this service contains a service.yaml file
	_, err := os.Stat("./service.yaml")
	if err != nil {
		return fmt.Errorf("Cannot upload service: no service.yaml file found in the current directory. Initialize a service first.")
	}

	// Parse the service.yaml file
	yaml, err := serviceyaml.Parse("./service.yaml")
	if err != nil {
		return fmt.Errorf("Failed to parse service.yaml file: %v", err)
	}

	// Get the active connection
	activeConnection := state.Get().RoverConnections.GetActive()
	if activeConnection == nil {
		return fmt.Errorf("No active Rover connection found, do not know where to upload the service.")
	}

	// Create an SSH client
	sshconn, err := activeConnection.ToSSH()
	if err != nil {
		return err
	}
	defer sshconn.Close()

	// Connect over SSH
	auth := goph.Password(activeConnection.Password)
	sshclient, err := goph.NewConn(&goph.Config{
		User:     activeConnection.Username,
		Addr:     activeConnection.Host,
		Port:     22,
		Auth:     auth,
		Timeout:  goph.DefaultTimeout,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
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
