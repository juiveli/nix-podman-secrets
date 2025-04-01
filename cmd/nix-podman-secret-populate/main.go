package main

import (
	"flag"
	"os"
	"fmt"
	"github.com/juiveli/nix-podman-secrets/internal"
)

const (
	ENV_DEBUG = "NIX_PODMAN_SECRETS_DEBUG"
)

func main() {
	debug := false
	if os.Getenv(ENV_DEBUG) == "true" {
		debug = true
	}

	nonroot := flag.Bool("nonroot", false, "Force non-root usage even when run with root")
	flag.Parse()
	internal.WrapMain(func() {

		internal.InitEnvVars(*nonroot)

		secretDir := os.Getenv("NIX_SECRET_DIR")
		if secretDir == "" {
			panic(fmt.Errorf("NIX_SECRET_DIR is not set"))
		}
		
		mappingDir := os.Getenv("MAPPING_DIR")
		if mappingDir == "" {
			panic(fmt.Errorf("MAPPING_DIR is not set"))
		}

		internal.PopulatePodmanSecretsDB(
			secretDir,
			mappingDir,
			internal.DeletePodmanSecretImpl,
			internal.CreatePodmanSecretImpl,
			debug)
	})
}