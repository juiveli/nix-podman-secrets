{ pkgs, lib }:

pkgs.buildGoModule {
  pname = "nix-podman-secrets";
  version = "0.3.0";

  src = ./.;
  doCheck = false;

  vendorHash = null;

  subPackages = [
    "cmd/nix-podman-secret-delete"
    "cmd/nix-podman-secret-list"
    "cmd/nix-podman-secret-lookup"
    "cmd/nix-podman-secret-populate"
    "cmd/nix-podman-secret-store"
  ];

  meta = {
    description = "Simple tool for podman secrets shell driver to access nix secrets";
    homepage = "https://github.com/juiveli/nix-podman-secrets";
  };
}
