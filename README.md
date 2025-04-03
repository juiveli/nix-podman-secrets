# nix-podman-secrets

This is a very simple program and flake to configure podman on a NixOS system
to use the secrets in `/run/podman-secrets` populated i.e. by [sops-nix](https://github.com/Mic92/sops-nix).

And for home-manager users `$XDG_RUNTIME_DIR/containers/podman-secrets`

## Adding to config

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
      ];

      # More specific sops installation please refer to their documentation.
      # This is here to show how to map path
      sops.secrets.yourkeyname = { 
        path = "/run/podman-secrets/exampleKeyRoot";
      };
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
      imports = [ nix-podman-secrets.homeManagerModules.nix-podman-secrets ];
            
      # More specific sops installation please refer to their documentation.
      # This is here to show how to map path
      sops.secrets.yourkeyname = { 
        %r will map to $XDG_RUNTIME_DIR
        path = "%r/containers/podman-secrets/exampleKeyUser";
      };
    }

    # The secrets are mapped to podman in a systemd user service called nix-podman-secret
    # so other services needing secrets must order after it:

    {
      systemd.user.services.mbsync.unitConfig.After = [ "nix-podman-secret.service" ];
      systemd.user.services.mbsync.unitConfig.Requires = [ "nix-podman-secret.service" ];
    }
  }
}
```

You can add it with both nixosModules for root, and home-manager for specific user, but then adding package should only be done once

For sops usage please refer to their documentation.
Sops could also be used with home-manager, but at the moment that is not supported, and above approach is needed

## Usage after nixos-rebuild switch
After configured you can use those secrets as an podman secret:

```
sudo podman secret ls

ID                         NAME        DRIVER      CREATED         UPDATED
c9jd6r0djdq934jfsdu5m02dl  exampleKeyRoot     shell       1 minutes ago  1 minutes ago

podman secret ls

ID                         NAME        DRIVER      CREATED         UPDATED
c760e707f228e6f8822fed6dc  exampleKeyUser     shell       1 minutes ago  1 minutes ago

sudo podman run --secret=exampleKeyRoot,type=env,target=MY_PASSWORD \
    registry.access.redhat.com/ubi9:latest \
    printenv MY_PASSWORD

ExampleRoot

podman run --secret=exampleKeyUser,type=env,target=MY_PASSWORD \
    registry.access.redhat.com/ubi9:latest \
    printenv MY_PASSWORD

ExampleUser

```

