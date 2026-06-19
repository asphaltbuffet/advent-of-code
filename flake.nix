{
  description = "my coding practice solutions";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable-small";
    flake-utils.url = "github:numtide/flake-utils";
    nur.url = "github:nix-community/NUR";
    fenix = {
      url = "github:nix-community/fenix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    nixpkgs,
    flake-utils,
    nur,
    fenix,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [nur.overlays.default];
        };
        rustToolchain = fenix.packages.${system}.stable.withComponents [
          "cargo"
          "clippy"
          "rust-src"
          "rustc"
          "rustfmt"
        ];
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            jujutsu
            jjui
            ripgrep
            fd
            sd
            parallel
            gopls
            nixd
            uv
            nodejs
            gh

            go
            python3
            gfortran

            rustToolchain
            fenix.packages.${system}.rust-analyzer
          ];

          CGO_ENABLED = "0";
        };
      }
    );
}
