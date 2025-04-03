{ lib, config, pkgs, options, ... }: 

let
    nix-podman-secrets = (pkgs.callPackage ../.. { }).nix-podman-secrets;
in
{
    options.nix-podman-secrets = {
    podmanPackage = lib.mkOption {
        type = lib.types.package;
        default = pkgs.podman;
        description = "The podman package to use";
    };
    };

    config.home.activation.syncNixPodmanSecrets = lib.hm.dag.entryAfter ["specialfs" "users" "groups" "setupSecrets"]
    ''
    echo "Populating podman secrets from nix secrets..."
    # Optionally, check for something in your home config instead of /run/current-system
    [ -e "$XDG_RUNTIME_DIR/containers/secrets" ] || echo "secrets directory not found, continuing..."
    # Extend your PATH appropriately.
    PATH=$PATH:${lib.makeBinPath [
    config.nix-podman-secrets.podmanPackage
    nix-podman-secrets
    ]}
    ${
        nix-podman-secrets.outPath
    }/bin/nix-podman-secret-populate
    '';
}