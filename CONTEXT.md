# Context

A personal archive of solutions to programming-practice exercises ("ProgPrac"),
primarily Advent of Code, plus Project Euler (TBD) and Exercism tracks. AoC
exercises are downloaded, run, tested, and benchmarked with **elf** — an external
CLI (https://github.com/asphaltbuffet/elf), not part of this repo. This repo holds
the exercise data, the polyglot solutions, and shared Go utilities (`pkg/`).

This file is the project glossary. When naming a domain concept in an issue,
refactor, test, or hypothesis, use the term as defined here. Extend it via
`/grill-with-docs` when a new term is resolved — don't drift to synonyms.

## Glossary

- **elf** — the external CLI used to manage AoC exercises
  (https://github.com/asphaltbuffet/elf), supporting AoC only. Key commands:
  `elf download` (fetch puzzle + input), `elf solve <year>/<day>` (run a
  solution), `elf test` (run test cases), `elf visualize` (write a Vis output to
  disk), `elf benchmark` (time implementations), `elf analyze` (graph benchmark
  results), and `elf config init|check|update-token`. Configured via `elf.toml`
  (in `$XDG_CONFIG_HOME/elf/` or the cwd) and env vars `ELF_ADVENT_TOKEN`,
  `ELF_LANGUAGE`. Relevant settings: `settings.language` (default `go`),
  `settings.input-file` (default `input.txt`), `settings.advent.dir` (default
  `exercises`), `settings.task-timeout` (default 2 min). Language toolchains must
  be installed separately.

- **Exercise** — one puzzle (e.g. AoC 2023 day 1, "Trebuchet?!"). Identified by
  an `id` like `2023-01`. Each exercise has a directory under
  `exercises/<year>/<NN-slug>/` containing `info.json`, `README.md`, optional
  `benchmark.json`, and one subdirectory per language implementation.

- **Part** — an exercise has two parts. In Go code and the exercise interface
  they are **One** and **Two** (methods `One(instr string)` / `Two(instr string)`).
  In `info.json` the *test cases* are keyed `one`/`two` but the *final answers*
  are keyed `a`/`b`. Treat "Part One" / "Part Two" as the canonical names; `a`/`b`
  is a storage detail of the answers block only.

- **Exercise interface** — each Go solution is an `Exercise` struct embedding
  `common.BaseExercise` and implementing `One`, `Two`, and optionally `Vis`
  (visualization). `instr` is the raw puzzle input string. See
  `internal/common/adventofcode.go`.

- **Implementation** — a single-language solution of an exercise (the `go/`,
  `python/`, etc. subdirectory). One exercise may have several implementations;
  elf selects which to run via the `settings.language` / `ELF_LANGUAGE` setting,
  and `benchmark.json` records timings per implementation by `name` ("Go", …).

- **Runner** — an elf-supported language backend that executes an Implementation,
  driven by a template under `~/.config/elf/runners/`. As of v0.4.1 (`elf runners
  list`): Python (`py`), Go (`go`), Bash (`bash`), Rust (`rs`), Lua (`lua`),
  Fortran 77 (`f77`). This is the set elf can *run*; the Exercism `tracks/`
  directories cover more languages (e.g. nim, jq, vimscript) that are not elf
  runners.

- **Test case** — an `{input, expected}` pair stored in `info.json` under
  `data.testCases.one` / `.two`. Used to validate a part before running it
  against the real input.

- **Input** — the real puzzle input, stored in the file named by
  `data.inputFile` (typically `input.txt`) alongside the exercise.

- **Answer** — the verified final result for the real input, stored in
  `info.json` under `data.answers.a` / `.b`.

- **Benchmark** — timing data for an implementation (`mean`/`min`/`max` per part,
  over `numRuns`), produced by `elf benchmark` and stored in `benchmark.json`.
  `normalization` scales results for cross-machine comparison. `elf analyze`
  graphs them (box/line plots); `run-times.png` is such a graph.

- **Track** — a language-organized collection of Exercism exercises under
  `tracks/<language>/` (go, nim, python, rust, bash, jq, fortran, vimscript).
  Distinct from the AoC `exercises/` tree, and broader than the set of elf
  **Runners** (Exercism tracks are not run by elf).

- **Year** — AoC exercises are grouped by puzzle year (2015–2024) under
  `exercises/<year>/`.

## Code layout

- `exercises/<year>/<NN-slug>/` — AoC exercises and their per-language implementations.
- `tracks/<language>/` — Exercism solutions, grouped by language track.
- `internal/common/` — shared AoC types, notably `BaseExercise`.
- `pkg/` — reusable Go utilities (`set`, `pq`, `ring`, `permutation`, `utilities`).
- `lib/bash/aoc.bash`, `lib/aocpy/` — language-specific helper libraries for
  bash and python implementations.

## Conventions

- AoC exercises are managed with **elf**, not run directly. To run a solution by
  hand, follow the layout elf expects (above); to run via elf, use
  `elf solve <year>/<day>`. The exercises root is `settings.advent.dir` (default
  `exercises`).
- VCS is **jujutsu** (`jj`), not raw `git`.
- Go module: `github.com/asphaltbuffet/advent-of-code`.
