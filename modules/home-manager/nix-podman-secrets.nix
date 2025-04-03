
{ lib, config, pkgs, options, ... }: 

let
    nix-podman-secrets = (pkgs.callPackage ../.. { }).nix-podman-secrets;

    serviceScript = toString (
        pkgs.writeShellScript "nix-podman-secret user"
        (
             ''
            echo "starting to generate secrets"
            ${
            nix-podman-secrets}/bin/nix-podman-secret-populate
            ''
        )
    );
in
{
    options.nix-podman-secrets = {
    podmanPackage = lib.mkOption {
        type = lib.types.package;
        default = pkgs.podman;
        description = "The podman package to use";
    };
    };

    config.systemd.user.services.nix-podman-secret = {
      Unit = {
        Description = "Populate podman secrets to user";
        After = [ "sops-nix.service" ];
        Requires = [ "sops-nix.service" ];
      };
      Service = {
        Type = "oneshot";
        Environment = "PATH=${config.nix-podman-secrets.podmanPackage}/bin:/run/wrappers/bin:${nix-podman-secrets}/bin:$PATH";
        ExecStart = serviceScript;
      };
      Install = {
        WantedBy = ["default.target"];
      };
    };

    config.home.activation.createPodmanSecretsDir = lib.hm.dag.entryAfter ["specialfs" "users" "groups" "setupSecrets"] 
        ''
            mkdir -p "$XDG_RUNTIME_DIR/containers/podman-secrets"
        '';
    


    config.home.packages = [ nix-podman-secrets ];
}

