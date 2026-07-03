# [Day 10: Pipe Maze](https://adventofcode.com/2023/day/10)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 10: Pipe Maze][rm10]
[Go][go10]

[rm10]: 10-pipeMaze/README.md
[go10]: 10-pipeMaze/go

-->

## Go

```text
────────────────────────────────────────
─        2023 Day 10: Pipe Maze        ─
────────────────────────────────────────
Solving (Go)…
1.0:  PASS            12.954ms
2.0:  PASS           228.745ms
```

## Python

The loop is traced once by inferring `S`'s shape from a neighbor that connects
back, then following each pipe's two links. Part one is half the loop length;
part two applies the shoelace formula for the enclosed area and Pick's theorem
(`i = A - b/2 + 1`) for the interior tile count — no flood fill or ray casting.

```text
────────────────────────────────────────
─        2023 Day 10: Pipe Maze        ─
────────────────────────────────────────
Solving (Python)…
1.0:  PASS             6.346ms
2.0:  PASS             6.452ms
```

## Rust

The same single trace over a byte grid, then shoelace + Pick for the interior.
Because both answers fall out of the loop's geometry rather than scanning every
enclosed cell, part two runs in well under a millisecond.

```text
────────────────────────────────────────
─        2023 Day 10: Pipe Maze        ─
────────────────────────────────────────
Solving (Rust)…
1.0:  PASS           241.760µs
2.0:  PASS           186.461µs
```

## Visualization

The main loop traced through the maze (teal), starting from `S` (red), with the
tiles it encloses — the Part Two answer, found by a winding-number test — filled
in orange. Everything outside the loop is left dark.

![Pipe Maze loop and enclosed area](pipe-maze.png)

## Run Times

![Day 10 run-time graphs](run-times.png)
