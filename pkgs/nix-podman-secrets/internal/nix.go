package internal

import (
    "fmt"
    "os"
    "os/user"
    "path/filepath"
)

const (
    NIX_SECRET_DIR_ROOT  = "/run/podman-secrets"
    MAPPING_DIR_ROOT     = "/var/lib/containers/storage/secrets/nix-mapping"
    NIX_SECRET_DIR_NONROOT = "$XDG_RUNTIME_DIR/containers/podman-secrets"
    MAPPING_DIR_NONROOT  = "$HOME/.local/share/containers/storage/secrets/nix-mappings"
)

func InitEnvVars(forceNonroot bool) {
    // Get the current user
    currentUser, err := user.Current()
    if err != nil {
        fmt.Printf("Error getting current user: %v\n", err)
        os.Exit(1) // Exit if user information cannot be retrieved
    }

    var nixSecretDir, mappingDir string

	// forceNonroot is a workaround, as podman secret creation pass Uid of 0 too, even when running as regular user
    if currentUser.Uid == "0" && !forceNonroot { // Root user
        nixSecretDir = NIX_SECRET_DIR_ROOT
        mappingDir = MAPPING_DIR_ROOT
    } else { // Non-root user, expand paths
        nixSecretDir = os.ExpandEnv(NIX_SECRET_DIR_NONROOT)
        mappingDir = os.ExpandEnv(MAPPING_DIR_NONROOT)

        // Check if environment variables are set
        if nixSecretDir == "$XDG_RUNTIME_DIR/containers/podman-secrets" || mappingDir == "$HOME/.local/share/containers/storage/secrets/nix-mappings" {
            fmt.Println("Error: Environment variables XDG_RUNTIME_DIR or HOME are not set.")
            os.Exit(1) // Exit if variables are not properly expanded
        }
    }

    // Set environment variables
    err = os.Setenv("NIX_SECRET_DIR", nixSecretDir)
    if err != nil {
        fmt.Printf("Error setting NIX_SECRET_DIR environment variable: %v\n", err)
        os.Exit(1)
    }
    err = os.Setenv("MAPPING_DIR", mappingDir)
    if err != nil {
        fmt.Printf("Error setting MAPPING_DIR environment variable: %v\n", err)
        os.Exit(1)
    }
}


func EnsureMappingDirExists(mappingDirPath string) error {
	if stat, err := os.Stat(mappingDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(mappingDirPath, 0700); err != nil {
			return fmt.Errorf("failed to create mapping dir: %w", err)
		}
	} else if !stat.IsDir() {
		return fmt.Errorf("mapping dir path %s exists, but is not a directory", mappingDirPath)
	}
	return nil
}

func ListNixSecrets(secretsDir string) (secretNames []string, err error) {
	secretsDir, err = filepath.EvalSymlinks(secretsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve secrets dir: %w", err)
	}

	secretFiles, err := os.ReadDir(secretsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets dir: %w", err)
	}


	for _, secretFile := range secretFiles {

		if !secretFile.IsDir() {
			secretNames = append(secretNames, secretFile.Name())
		}
	}
 
	return
}
