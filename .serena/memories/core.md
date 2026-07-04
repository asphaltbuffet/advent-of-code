# Core

Advent of Code solutions repo: `github.com/asphaltbuffet/advent-of-code`.

## Layout

- `exercises/<year>/<day>-<slug>/` — one dir per puzzle; contains `info.json`, `input.txt`, `README.md`, `run-times.png`, `benchmark.json`, and per-language subdirs (`go/`, `rs/`, `py/`).
- `internal/` — shared Go internals.
- `pkg/` — exported Go packages.
- `lib/` — shared library code (e.g. `lib/bash/aoc.bash`).
- `docs/` — agent guidance (`docs/agents/`), ADRs.
- `CONTEXT.md` — domain doc (created lazily by `/grill-with-docs`).

## Key invariants

- Puzzle answers live only in `info.json` (`data.answers.a`/`.b`). Never in code comments.
- Exercise dir names must be valid Go import paths (no `?`/`!`/`'`).
- VCS is **jujutsu** (`jj`), not git. One commit per day solved.
- `elf` CLI (pinned in `flake.nix`) drives scaffold, solve, benchmark, analyze, visualize.

See `mem:tech_stack` for toolchain, `mem:conventions` for code style, `mem:task_completion` for done-criteria commands.
