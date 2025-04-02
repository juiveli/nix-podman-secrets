package internal

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	podmanBin = "podman"

	nixPodmanSecretsBin = "nix-podman-secret"
)

type DeletePodmanSecretFunc func(string) error
type CreatePodmanSecretFunc func(string) error

func listPodmanSecrets(mappingDirPath string) (secretNames []string, removedSecretIDs []string, err error) {

	files, err := os.ReadDir(mappingDirPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list entries in mapping dir: %w", err)
	}

	for _, secretFile := range files {
		secretPath := filepath.Join(mappingDirPath, secretFile.Name())
		actualSecretFile, err := evaluateNextSymlink(secretPath)
		if err != nil {
			removedSecretIDs = append(removedSecretIDs, secretFile.Name())
			continue
		}
		secretName := filepath.Base(actualSecretFile)

		secretNames = append(secretNames, strings.TrimSpace(secretName))
	}
	return
}

func evaluateNextSymlink(path string) (string, error) {
    // Get file or symlink information
    info, err := os.Lstat(path)
    if err != nil {
        return "", err
    }

    // Check if the path is a symlink
    if info.Mode()&os.ModeSymlink == 0 {
        return "", errors.New("provided path is not a symlink")
    }

    // Fetch the target of the symlink
    nextTarget, err := os.Readlink(path)
    if err != nil {
        return "", err
    }
    return nextTarget, nil
}


func DeletePodmanSecretImpl(secretName string) error {
	cmd := exec.Command(podmanBin, "secret", "rm", secretName)
	errBuf := &bytes.Buffer{}
	cmd.Stderr = errBuf
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to delete secret (%s): %w", errBuf.String(), err)
	}
	return nil
}

func CreatePodmanSecretImpl(secretName string) error {
    currentUser, err := user.Current()
    if err != nil {
        fmt.Printf("Error getting current user: %v\n", err)
        os.Exit(1) // Exit if user information cannot be retrieved
    }

	driverOptsDelete := "delete=" + nixPodmanSecretsBin + "-delete"
	driverOptsList := "list=" + nixPodmanSecretsBin + "-list"
	driverOptsLookup := "lookup=" + nixPodmanSecretsBin + "-lookup"
	driverOptsStore := "store=" + nixPodmanSecretsBin + "-store"

	var driverOptsFull string;
	if currentUser.Uid == "0" {
		driverOptsFull = driverOptsDelete + "," + driverOptsList + "," + driverOptsLookup + "," + driverOptsStore 
	} else {
		driverOptsFull = driverOptsDelete + " --nonroot," + driverOptsList + " --nonroot," + driverOptsLookup + " --nonroot," + driverOptsStore + " --nonroot"
	}


	cmd := exec.Command(podmanBin,
		"secret",
		"create",
		"--label", "source=nix",
		"--driver", "shell",
		"--driver-opts", driverOptsFull,
		secretName, "-")
	errBuff := &bytes.Buffer{}
	stdInBuff := bytes.NewBuffer([]byte(secretName))
	cmd.Stdin = stdInBuff
	cmd.Stderr = errBuff

	secretCreationError := cmd.Run()
	if secretCreationError != nil {
		return fmt.Errorf("failed to create secret (%s): %w", errBuff.String(), err)
	}
	return nil
}
