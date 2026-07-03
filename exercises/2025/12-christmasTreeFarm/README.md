# [Day 12: Christmas Tree Farm](https://adventofcode.com/2025/day/12)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 12: Christmas Tree Farm][rm12]
[Go][go12]
[Python][py12]

[rm12]: 12-christmasTreeFarm/README.md
[go12]: 12-christmasTreeFarm/go
[py12]: 12-christmasTreeFarm/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025            
       Day 12: Christmas Tree Farm        
──────────────────────────────────────────
Testing (Go)...
  1.1: FAIL 
      ⤷ got "3", but expected "2"
Solving (Go)...
  1.1: PASS              4.241ms
      ⤷ 474
```

## Python

A true packing check rather than an area heuristic: presents (rotated and flipped)
are placed by covering the first empty cell, trying every orientation and anchor,
or leaving the cell as one of the allowed holes. A quick cell-area test rejects the
impossible regions up front. Unlike the area-only reference — which mis-counts the
small example — this reproduces the example's answer of 2 as well as the real 474.
The puzzle's example is deliberately tighter than the real regions, so its test is
far slower than the solve.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
       Day 12: Christmas Tree Farm
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS             4.409s
```

## Rust

The same real packing search over `(i32, i32)` cells, with orientations
de-duplicated into a `HashSet`. Native speed settles the full farm in a fraction of
a second.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
       Day 12: Christmas Tree Farm
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS           142.648ms
```

## 2025 Run Times

![2025 exercise run-time graphs](../run-times.png)
