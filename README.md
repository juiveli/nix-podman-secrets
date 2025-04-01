# nix-podman-secrets

This is a very simple program and flake to configure podman on a NixOS system
to use the secrets in `/run/secrets` populated i.e. by [sops-nix](https://github.com/Mic92/sops-nix).

And for home-manager users `$XDG_RUNTIME_DIR/containers/secrets`

To use this you can simply add this flake to your flake i.e.:

```
inputs = {
    nix-podman-secrets = {
      url = "github:juiveli/nix-podman-secrets?ref=latest";

      inputs.nixpkgs.follows = "nixpkgs";
    };
}
```

and add the module to you nixosSystem module list, i.e.

```
    nixosConfigurations = {
        podman-host = nixpkgs.lib.nixosSystem {
          system = "x86_64-linux";
          modules = [
            inputs.nix-podman-secrets.nixosModules.nix-podman-secrets
          ]
          environment.systemPackages = [nix-podman-secrets.packages.${pkgs.system}.nix-podman-secrets];
    }
  }
```

or with home-manager

```
    nixosConfigurations = {
        podman-host = nixpkgs.lib.nixosSystem {
          system = "x86_64-linux";

          home-manager.users.<user> =
          {

            uid = 1000; # Needed for sops part below
            imports = nix-podman-secrets.homeManagerModules.nix-podman-secrets
          }
          home.packages = [nix-podman-secrets.packages.${pkgs.system}.nix-podman-secrets];
    }

          # To map the key to the folder "$XDG_RUNTIME_DIR/containers/secrets"
            sops.secrets.exampleKey = {
              owner = config.users.users.<user>.name;
              path = "/run/user/${toString config.users.users.<user>.uid}/containers/secrets/something";
  }
}
```

You can add it with both nixosModules for root, and home-manager for specific user, but then adding package should only be done once

For sops usage please refer to their documentation.
Sops could also be used with home-manager, but at the moment that is not supported, and above approach is needed