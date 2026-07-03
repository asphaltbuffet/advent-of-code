# [Day 14: Chocolate Charts](https://adventofcode.com/2018/day/14)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 14: Chocolate Charts][rm14]
[Go][go14]
[Rust][rs14]
[Python][py14]

[rm14]: 14-chocolateCharts/README.md
[go14]: 14-chocolateCharts/go
[rs14]: 14-chocolateCharts/rs
[py14]: 14-chocolateCharts/py

-->

## Notes

Two elves build a scoreboard of single-digit recipes. Each step appends the digits
of their two current recipes' sum, then each elf steps forward by one plus its own
recipe. The board is a `[]byte` of digits that only ever grows.

- **Part One** treats the input as a count: build until there are ten recipes past
  it, then read those ten digits.
- **Part Two** treats the input as a digit pattern and reports how many recipes
  precede its first appearance. Since a step can append one or two digits, the tail
  is checked after every append, tracking how far it has already compared so no
  suffix is skipped across a two-digit step.

## Go

```text
────────────────────────────────────────
─    2018 Day 14: Chocolate Charts     ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS             8.264ms
2.0:  PASS           357.976ms
```

## Rust

A growable `Vec<u8>` of single digits with two `usize` elf indices. Part Two keeps
a `checked` cursor and compares the target slice against the board after every
append, so a pattern landing on either digit of a two-digit step is caught without
re-scanning the whole board.

```text
────────────────────────────────────────
─    2018 Day 14: Chocolate Charts     ─
────────────────────────────────────────

Solving (Rust)…
1.0:  PASS             9.012ms
2.0:  PASS           419.665ms
```

## Python

A `bytearray` of single digits gives cheap appends and low memory. Part Two grows
the board in fixed-size batches, then scans each new region (overlapping the prior
tail) with `bytearray.find` — a C-level search — so the interpreter never slices a
prefix per step.

```text
────────────────────────────────────────
─    2018 Day 14: Chocolate Charts     ─
────────────────────────────────────────

Solving (Python)…
1.0:  PASS              0.206s
2.0:  PASS              9.193s
```

## Run Times

![Day 14 run-time graphs](run-times.png)
