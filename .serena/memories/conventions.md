# Conventions

## Go
- Exercise entry points: `func (e Exercise) One(instr string) (any, error)` and `Two` in `go/exercise.go`.
- Optional viz: `func (e Exercise) Vis(instr, outdir string) error` in `go/vis.go`.
- Match style of completed days in the same year.

## Rust
- Only `rs/src/solution.rs` and `Cargo.toml`/`Cargo.lock` are tracked.
- Entry: `pub fn part_one(input: &str) -> String` / `part_two`.

## Python
- `class Exercise(BaseExercise)` with `@staticmethod one(instr)/two(instr)` in `py/__init__.py`.

## General
- American English everywhere (color, gray, -ize).
- Never leave blank `info.json` test cases (`""` / `""`); fill from puzzle example or clear to `[]`.
- If example uses different parameters than real input, detect regime from input shape (not hardcoded).
- Visualizations: prefer Okabe-Ito palette; verify grayscale readability (`magick img -colorspace Gray out`).
- Prefer `fd`/`rg`/`sd`/`jq` over POSIX equivalents.
