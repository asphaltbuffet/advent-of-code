# [Day 22: Sand Slabs](https://adventofcode.com/2023/day/22)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 22: Sand Slabs][rm22]
[Go][go22]

[rm22]: 22-sandSlabs/README.md
[go22]: 22-sandSlabs/go

-->

## Go

```text
────────────────────────────────────────
─       2023 Day 22: Sand Slabs        ─
────────────────────────────────────────
Solving (Go)…
1.0:  PASS           304.984ms
2.0:  PASS           296.814ms
```

## Python

Bricks are settled in z order using a `heights[(x,y)]` map of the topmost cell, so
each drop is a lookup rather than a collision scan; that yields the support graph
directly. Part one counts bricks that aren't a sole supporter; part two BFS-es the
chain of bricks that fall once all their supporters have.

```text
────────────────────────────────────────
─       2023 Day 22: Sand Slabs        ─
────────────────────────────────────────
Solving (Python)…
1.0:  PASS             6.600ms
2.0:  PASS            25.138ms
```

## Rust

The same height-map settle and support graph over `Vec<HashSet>`, with the chain
reaction using `is_subset` to test whether a brick has lost all its supporters.
The map-based settle keeps both parts in the low tens of milliseconds.

```text
────────────────────────────────────────
─       2023 Day 22: Sand Slabs        ─
────────────────────────────────────────
Solving (Rust)…
1.0:  PASS             1.205ms
2.0:  PASS            13.408ms
```

## Visualization

The settled brick pile in the two side elevations the puzzle uses: looking down
the y axis (x vs z, left) and down the x axis (y vs z, right). Each brick keeps a
stable colour across both views, so the packed tower structure after everything
falls is visible from either side.

![Sand Slabs settled pile](sand-slabs.png)

## Run Times

![Day 22 run-time graphs](run-times.png)
