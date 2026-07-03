# [Day 4: Scratchcards](https://adventofcode.com/2023/day/4)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 4: Scratchcards][rm4]
[Go][go4]
[Rust][rs4]
[Python][py4]

[rm4]: exercises/2023/04-scratchcards/README.md
[go4]: exercises/2023/04-scratchcards/go
[rs4]: exercises/2023/04-scratchcards/rs
[py4]: exercises/2023/04-scratchcards/py

-->

## Go

```text
 ───────────────────
 ADVENT OF CODE 2023
 Day 4: Scratchcards
 ───────────────────
  Test 1.0: PASS in 26.9 µs
  Test 2.0: PASS in 17.7 µs
  Part 1: 23441 in 446.2 µs
  Part 2: 5923918 in 495.6 µs
```

## Python

Match counts come from a set intersection of the two number groups. Part one
sums `1 << (m - 1)`; part two keeps a `copies` list and, for each card, adds its
copy count forward onto the next `m` cards in one pass.

```text
 ───────────────────
 ADVENT OF CODE 2023
 Day 4: Scratchcards
 ───────────────────

Solving (Python)…
1.0:  PASS           838.279µs
2.0:  PASS           885.248µs
```

## Rust

`split_once` peels off the card label and splits the two number groups; a
`HashSet` of winners drives a `filter().count()` for the match count. Part two
propagates a `Vec<u64>` of copies forward with a bounded inner range — the same
single-pass cascade.

```text
 ───────────────────
 ADVENT OF CODE 2023
 Day 4: Scratchcards
 ───────────────────

Solving (Rust)…
1.0:  PASS           214.462µs
2.0:  PASS           224.236µs
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
