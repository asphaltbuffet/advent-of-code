# [Day 5: Cafeteria](https://adventofcode.com/2025/day/5)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 5: Cafeteria][rm5]
[Go][go5]
[Python][py5]

[rm5]: 05-cafeteria/README.md
[go5]: 05-cafeteria/go
[py5]: 05-cafeteria/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025            
             Day 5: Cafeteria             
──────────────────────────────────────────
Testing (Go)...
  1.1: PASS                 50µs
  2.1: PASS                 10µs
Solving (Go)...
  1.1: PASS                691µs
      ⤷ 848
  2.1: PASS                243µs
      ⤷ 334714395325710
```

## Python

Sort the inventory ranges and coalesce overlaps (a one-unit gap keeps ranges
separate). Part one tests each ingredient ID for membership; part two sums the
merged widths.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
             Day 5: Cafeteria
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS             4.187ms
2.0:  PASS             0.338ms
```

## Rust

The same interval merge. Since the merged ranges are sorted and disjoint, part one
uses `partition_point` to binary-search the one range that could contain each ID,
rather than scanning every range.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
             Day 5: Cafeteria
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS           100.279µs
2.0:  PASS            53.753µs
```

## 2025 Run Times

![2025 exercise run-time graphs](../run-times.png)
