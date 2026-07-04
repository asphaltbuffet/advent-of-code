# Tech Stack

- **Go 1.25.4** — primary solution language; module `github.com/asphaltbuffet/advent-of-code`.
- **Rust** (stable via fenix flake) — secondary solution language.
- **Python** — tertiary solution language (numpy available only via `nix develop --command elf ...`).
- **elf v0.4.4** — pinned in `flake.nix` (use flake input, not overlay). Drives scaffold/solve/benchmark/analyze/visualize.
- **Nix flakes** — dev environment; `nix develop` gives Rust toolchain + Python with extras.
- Notable Go deps: `testify`, `dominikbraun/graph`, `kettek/apng`, `golang.org/x/image`, `caarlos0/log`, `fatih/color`.
