{
  pkgs ? import <nixpkgs> { }
}:
let
  nix-podman-secrets = pkgs.callPackage ./pkgs/nix-podman-secrets/package.nix {};
in

{
  inherit nix-podman-secrets;
}