package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/juiveli/nix-podman-secrets/internal"
)

func main() {
	nonroot := flag.Bool("nonroot", false, "Force non-root usage even when run with root")
	flag.Parse()
	internal.WrapMain(func() {
		internal.InitEnvVars(*nonroot)
		secretId := os.Getenv("SECRET_ID")
		storeSecret(os.Stdin, secretId, os.Getenv("NIX_SECRET_DIR"), os.Getenv("MAPPING_DIR"))
	})
}

func storeSecret(in io.Reader, secretId, nixSecretDir, mappingDir string) {
	secretName, err := io.ReadAll(in) // Read nix secret name from stdin, because we give the name as secret content
	if err != nil {
		panic(fmt.Errorf("failed to read secret name data from stdin: %w", err))
	}
	if err := internal.EnsureMappingDirExists(mappingDir); err != nil {
		panic(fmt.Errorf("mapping dir does not exist: %w", err))
	}
	nixSecretPath := filepath.Join(nixSecretDir, string(secretName))
	targetPath := filepath.Join(mappingDir, secretId)
	if err := os.Symlink(nixSecretPath, targetPath); err != nil {
		panic(fmt.Errorf("failed to create symlink to nix secret: %w", err))
	}
}
