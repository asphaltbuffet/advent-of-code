# [Day 7: Laboratories](https://adventofcode.com/2025/day/7)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 7: Laboratories][rm7]
[Go][go7]
[Python][py7]

[rm7]: 07-laboratories/README.md
[go7]: 07-laboratories/go
[py7]: 07-laboratories/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025            
           Day 7: Laboratories            
──────────────────────────────────────────
Testing (Go)...
  1.1: PASS                  8µs
  2.1: PASS                  8µs
Solving (Go)...
  1.1: PASS                184µs
      ⤷ 1573
  2.1: PASS                190µs
      ⤷ 15093663987272
```

## Python

The beam wavefront is carried row by row as a `{column: count}` map. A beam on a
splitter (^) sends its count down-left and down-right. Part one counts split
events; part two sums the surviving timelines. Building a fresh map each row
avoids re-splitting children in the row that created them.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
           Day 7: Laboratories
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS             4.096ms
2.0:  PASS             3.770ms
```

## Rust

The same sweep with a dense `Vec<u64>` of per-column counts instead of a map —
columns are bounded by the grid width — advancing into a fresh vector each row.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
           Day 7: Laboratories
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS            68.050µs
2.0:  PASS            54.812µs
```

## 2025 Run Times

![2025 exercise run-time graphs](../run-times.png)
