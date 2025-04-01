package main

import (
	"flag"
	"fmt"
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
		deleteSecret(os.Getenv("MAPPING_DIR"), secretId)
	})
}

func deleteSecret(mappingDir, secretId string) {
	if err := os.Remove(filepath.Join(mappingDir, secretId)); err != nil {
		panic(fmt.Errorf("failed to remove mapping symlink for secret id %s: %w", secretId, err))
	}
}
