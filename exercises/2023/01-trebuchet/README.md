# [Day 1: Trebuchet?!](https://adventofcode.com/2023/day/1)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 1: Trebuchet?!][rm1]
[Go][go1]
[Rust][rs1]
[Python][py1]

[rm1]: exercises/2023/01-trebuchet/README.md
[go1]: exercises/2023/01-trebuchet/go
[rs1]: exercises/2023/01-trebuchet/rs
[py1]: exercises/2023/01-trebuchet/py

-->

## Go

```text
 ───────────────────
 ADVENT OF CODE 2023
 Day 1: Trebuchet?!
 ───────────────────
  Test 1.0: PASS in 3.1 µs
  Test 2.0: PASS in 7.9 µs
  Part 1: 53334 in 314.1 µs
  Part 2: 52834 in 757.7 µs
```

## Python

An overlapping-match regex does the work: `(?=(...))` lookahead captures every
digit (and, in part two, every spelled word) even where spellings overlap like
`eightwo`, and a dict maps the token to its value.

```text
 ───────────────────
 ADVENT OF CODE 2023
 Day 1: Trebuchet?!
 ───────────────────

Solving (Python)…
1.0:  PASS             1.214ms
2.0:  PASS             1.863ms
```

## Rust

A single left-to-right position scan yields the digit starting at each offset
via `filter_map`; the iterator's `.next()` and `.last()` give the first and last
digits in one pass with no intermediate allocation. Spelled words are matched by
`starts_with` at each offset, so overlaps resolve naturally.

```text
 ───────────────────
 ADVENT OF CODE 2023
 Day 1: Trebuchet?!
 ───────────────────

Solving (Rust)…
1.0:  PASS            89.784µs
2.0:  PASS           169.499µs
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
