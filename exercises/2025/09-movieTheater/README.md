# [Day 9: Movie Theater](https://adventofcode.com/2025/day/9)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 9: Movie Theater][rm9]
[Go][go9]
[Python][py9]

[rm9]: 09-movieTheater/README.md
[go9]: 09-movieTheater/go
[py9]: 09-movieTheater/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025            
           Day 9: Movie Theater           
──────────────────────────────────────────
Testing (Go)...
  1.1: PASS                  3µs
  2.1: PASS                  6µs
Solving (Go)...
  1.1: PASS                378µs
      ⤷ 4781546175
  2.1: PASS            203.018ms
      ⤷ 1573359081
```

## Python

The vertices trace a rectilinear polygon. Part one takes the largest inclusive
rectangle spanned by any vertex pair; part two keeps only rectangles that lie
wholly inside — every corner on an edge or interior (on-segment and ray-cast
tests), and no polygon edge crossing the interior.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
           Day 9: Movie Theater
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS             0.034s
2.0:  PASS            20.485s
```

## Rust

The same geometry over `i64` points, with native truncating division matching the
reference's ray cast. The O(n³) validity search that costs Python ~20s finishes in
a fraction of a second.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
           Day 9: Movie Theater
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS             0.159ms
2.0:  PASS           224.796ms
```

## 2025 Run Times

![2025 exercise run-time graphs](../run-times.png)
