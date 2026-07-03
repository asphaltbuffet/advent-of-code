# [Day 23: A Long Walk](https://adventofcode.com/2023/day/23)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 23: A Long Walk][rm23]
[Go][go23]
[Python][py23]

[rm23]: 23-aLongWalk/README.md
[go23]: 23-aLongWalk/go
[py23]: 23-aLongWalk/py

-->

## Go

```text
< section intentionally left blank >
```

## Python

The trails compress to a weighted graph of junctions (cells with more than two
open neighbors) joined by corridor lengths, so the longest-path DFS runs over
~36 nodes rather than every cell. Part one honors slopes as one-way corridors;
part two makes them ordinary. The DFS is the plain recursive version — correct,
but the slowest of the three on the exponential part two.

```text
──────────────────────────────
      ADVENT OF CODE 2023
     Day 23: A Long Walk
──────────────────────────────

Solving (Python)…
1.0:  PASS              0.030s
2.0:  PASS             20.334s
```

## Rust

The same junction-graph compression, but nodes are indexed so the DFS "visited"
set is a `u64` bitmask — making the exponential longest-path search fast enough to
finish part two in a third of a second.

```text
──────────────────────────────
      ADVENT OF CODE 2023
     Day 23: A Long Walk
──────────────────────────────

Solving (Rust)…
1.0:  PASS             3.093ms
2.0:  PASS           332.700ms
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
