# [Day 13: Point of Incidence](https://adventofcode.com/2023/day/13)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 13: Point of Incidence][rm13]
[Go][go13]
[Python][py13]

[rm13]: 13-pointOfIncidence/README.md
[go13]: 13-pointOfIncidence/go
[py13]: 13-pointOfIncidence/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
        Day 13: Point of Incidence
──────────────────────────────────────────

Testing...
  1.1: PASS              0.01 ms
  2.1: PASS              0.01 ms

Solving...
    1: PASS              0.72 ms
      ⤷ 32723
    2: PASS              0.74 ms
      ⤷ 34536
```

## Python

Rows and columns are encoded as integer bitmasks, so a reflection test becomes
XOR: a mirror is valid when the bit differences across the folded pairs sum to
the target. Part one wants a perfect mirror (0); part two wants exactly one
smudge (1) — the only difference between the parts.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
        Day 13: Point of Incidence
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS             3.976ms
2.0:  PASS             4.014ms
```

## Rust

The same bitmask/XOR idea with `u32` masks and `count_ones()` summing each folded
pair's difference. Rows and columns are filled in a single grid pass, and both
parts share `reflection`, parameterized by the smudge count.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
        Day 13: Point of Incidence
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS           147.735µs
2.0:  PASS           123.411µs
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
