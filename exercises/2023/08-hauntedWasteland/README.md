# [Day 8: Haunted Wasteland](https://adventofcode.com/2023/day/8)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 8: Haunted Wasteland][rm8]
[Go][go8]
[Python][py8]

[rm8]: 08-hauntedWasteland/README.md
[go8]: 08-hauntedWasteland/go
[py8]: 08-hauntedWasteland/py

-->

## Go

```text
 ────────────────────────
 ADVENT OF CODE 2023
 Day 8: Haunted Wasteland
 ────────────────────────
  Test 1.1: PASS in 89.7 µs
  Test 1.2: PASS in 54.2 µs
  Test 2.1: PASS in 188.9 µs
    Part 1: 19199 8.3 ms
    Part 2: 13663968099527 7.3 ms
```

## Python

A regex builds the `node -> (left, right)` graph and `itertools.cycle` repeats
the instructions. Part one walks AAA to ZZZ; part two walks every `..A` start to
its `..Z`, and since each ghost loops on a fixed period they meet at `math.lcm`
of the periods.

```text
 ────────────────────────
 ADVENT OF CODE 2023
 Day 8: Haunted Wasteland
 ────────────────────────

Solving (Python)…
1.0:  PASS             3.019ms
2.0:  PASS            14.319ms
```

## Rust

The graph is a `HashMap<&str, (&str, &str)>` borrowing straight from the input via
fixed-offset slicing, and `moves.iter().cycle()` drives the walk. A shared `steps`
helper takes a `done` closure for both parts; part two folds the per-ghost cycle
lengths with a gcd-based `lcm`.

```text
 ────────────────────────
 ADVENT OF CODE 2023
 Day 8: Haunted Wasteland
 ────────────────────────

Solving (Rust)…
1.0:  PASS             0.798ms
2.0:  PASS             2.469ms
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
