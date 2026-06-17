{
  description = "my coding practice solutions";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable-small";
    flake-utils.url = "github:numtide/flake-utils";
    nur.url = "github:nix-community/NUR";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    nur,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [nur.overlays.default];
        };
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            jujutsu
            jjui
            ripgrep
            fd
            sd
            gopls
            nixd
            uv
            nodejs
            gh
          ];

          CGO_ENABLED = "0";
        };
      }
    );
}
