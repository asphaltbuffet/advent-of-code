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
    elf = {
      url = "github:asphaltbuffet/elf/v0.6.3";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    nixpkgs,
    flake-utils,
    nur,
    fenix,
    elf,
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
            elf.packages.${system}.default

            go
            golangci-lint
            (python3.withPackages (ps: [ps.numpy]))
            gfortran
            (lua5_2.withPackages (ps: [ps.dkjson]))

            dotnet-sdk_8

            rustToolchain
            fenix.packages.${system}.rust-analyzer
          ];

          CGO_ENABLED = "0";
          DOTNET_CLI_TELEMETRY_OPTOUT = "1";
          DOTNET_NOLOGO = "1";
        };
      }
    );
}
