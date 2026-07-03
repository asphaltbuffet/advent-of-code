# [Day 9: Mirage Maintenance](https://adventofcode.com/2023/day/9)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 9: Mirage Maintenance][rm9]
[Go][go9]
[Python][py9]

[rm9]: 09-mirageMaintenance/README.md
[go9]: 09-mirageMaintenance/go
[py9]: 09-mirageMaintenance/py

-->

## Go

```text
 ─────────────────────────
 ADVENT OF CODE 2023
 Day 9: Mirage Maintenance
 ─────────────────────────
  Test 1.1: PASS in 4.6 µs
  Test 2.1: PASS in 4.5 µs
    Part 1: 1921197370 620.2 µs
    Part 2: 1124 476.2 µs
```

## Python

`_extrapolate` recurses on the `pairwise` difference row until it is all zeros,
adding each row's last value. Part two reverses each history first, since
extending backwards is the forward computation on the reversed sequence.

```text
 ─────────────────────────
 ADVENT OF CODE 2023
 Day 9: Mirage Maintenance
 ─────────────────────────

Solving (Python)…
1.0:  PASS             2.428ms
2.0:  PASS             2.614ms
```

## Rust

The same recursion over `windows(2)` differences, summing each row's last value.
A shared `solve` takes a `reverse` flag so part two just flips each sequence
before extrapolating.

```text
 ─────────────────────────
 ADVENT OF CODE 2023
 Day 9: Mirage Maintenance
 ─────────────────────────

Solving (Rust)…
1.0:  PASS           156.550µs
2.0:  PASS           146.543µs
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
