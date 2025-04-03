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

    config.system.activationScripts.syncNixPodmanSecrets =
    (lib.stringAfter ([
        "specialfs"
        "users"
        "groups"
        "setupSecrets"
    ])) ''
        [ -e /run/current-system ] || echo "populating podman secrets from nix secrets"
        PATH=$PATH:${
        lib.makeBinPath [
            config.nix-podman-secrets.podmanPackage
            nix-podman-secrets
        ]
        } ${
        nix-podman-secrets.outPath
        }/bin/nix-podman-secret-populate
    '';

    config.environment.systemPackages = [ nix-podman-secrets ];
}