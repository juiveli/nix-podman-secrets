{
  description = "Use nix secrets in podman";

  inputs = { nixpkgs = { url = "github:nixos/nixpkgs/nixos-24.11"; }; };

  outputs = inputs@{ nixpkgs, self, ... }:
    let
      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
      ];

      # Helper to provide system-specific attributes
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });

    in {

        packages = forAllSystems ({ pkgs }: {
          nix-podman-secrets = pkgs.callPackage ./pkgs/nix-podman-secrets/package.nix {};

        });

    overlays.default = 
      final: prev: 
      {
      nix-podman-secrets = final.callPackage ./pkgs/nix-podman-secrets/package.nix {};
      };

      nixosModules = {
        nix-podman-secrets = ./modules/nix-podman-secrets;
        default = self.nix-podman-secrets;
      };

      homeManagerModules = {
        nix-podman-secrets = ./modules/home-manager/nix-podman-secrets.nix;
        default = self.nix-podman-secrets;
      };
    };
}
