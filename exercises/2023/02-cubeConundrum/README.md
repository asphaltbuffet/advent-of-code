# [Day 2: Cube Conundrum](https://adventofcode.com/2023/day/2)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 2: Cube Conundrum][rm2]
[Go][go2]
[Rust][rs2]
[Python][py2]

[rm2]: exercises/2023/02-cubeConundrum/README.md
[go2]: exercises/2023/02-cubeConundrum/go
[rs2]: exercises/2023/02-cubeConundrum/rs
[py2]: exercises/2023/02-cubeConundrum/py

-->

## Go

```text
 ─────────────────────
 ADVENT OF CODE 2023
 Day 2: Cube Conundrum
 ─────────────────────
  Test 1.0: PASS in 99.1 µs
  Test 2.0: PASS in 112.7 µs
  Part 1: 2439 in 2 ms
  Part 2: 63711 in 1.2 ms
```

## Python

Each game collapses to a per-color maximum via a regex over `<n> <color>`
draws and a dict; part one keeps games under the limits with `all(...)` and
part two sums `math.prod` of the three maxima.

```text
 ─────────────────────
 ADVENT OF CODE 2023
 Day 2: Cube Conundrum
 ─────────────────────

Solving (Python)…
1.0:  PASS           629.186µs
2.0:  PASS           600.576µs
```

## Rust

No regex crate: the parse splits each line on `:`, `,`, and `;` and folds the
`<count> <color>` tokens into a `(red, green, blue)` maxima tuple, then part one
filters by the limits and part two sums the products — all through iterator
chains over the standard library.

```text
 ─────────────────────
 ADVENT OF CODE 2023
 Day 2: Cube Conundrum
 ─────────────────────

Solving (Rust)…
1.0:  PASS            53.681µs
2.0:  PASS            47.513µs
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
