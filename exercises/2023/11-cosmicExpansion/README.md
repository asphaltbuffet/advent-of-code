# [Day 11: Cosmic Expansion](https://adventofcode.com/2023/day/11)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 11: Cosmic Expansion][rm11]
[Go][go11]
[Python][py11]

[rm11]: 11-cosmicExpansion/README.md
[go11]: 11-cosmicExpansion/go
[py11]: 11-cosmicExpansion/py

-->

## Go

```text
──────────────────────────
   ADVENT OF CODE 2023
 Day 11: Cosmic Expansion
──────────────────────────

Testing...
  1.1: PASS         19.6 µs
  2.1: PASS          7.7 µs

Solving...
    1: PASS        507.2 µs
      ⤷ 10490062
    2: PASS        508.9 µs
      ⤷ 382979724122
```

## Python

The distance sum splits into independent row and column axes; along each, empty
tracks are expanded by remapping coordinates, and the sorted-order identity
`sum_i x_i*(2i - n + 1)` yields the pairwise sum in one linear pass. Only the
expansion factor (2 vs. 1,000,000) differs between parts.

```text
──────────────────────────
   ADVENT OF CODE 2023
 Day 11: Cosmic Expansion
──────────────────────────

Solving (Python)…
1.0:  PASS             1.074ms
2.0:  PASS             1.286ms
```

## Rust

The same per-axis linear sum, accumulated in signed arithmetic so the negative
lower-half weights don't underflow. Galaxy rows and columns are gathered in one
grid pass; both parts share the code with only the factor changing.

```text
──────────────────────────
   ADVENT OF CODE 2023
 Day 11: Cosmic Expansion
──────────────────────────

Solving (Rust)…
1.0:  PASS           128.271µs
2.0:  PASS            25.319µs
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
