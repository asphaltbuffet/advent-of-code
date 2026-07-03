# [Day 10: The Stars Align](https://adventofcode.com/2018/day/10)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 10: The Stars Align][rm10]
[Go][go10]
[Rust][rs10]
[Python][py10]

[rm10]: 10-theStarsAlign/README.md
[go10]: 10-theStarsAlign/go
[rs10]: 10-theStarsAlign/rs
[py10]: 10-theStarsAlign/py

-->

## Notes

Each star drifts at a constant velocity. The message appears at the one instant the
stars are tightest, which is the time that minimizes their bounding-box extent
(width + height). The extent shrinks to that minimum and then grows, so we step
forward while it keeps shrinking and stop at the turn.

- **Part One** renders that frame as solid blocks; the letters read **BFFZCNXE**.
- **Part Two** is the convergence time itself — **10391** seconds.

The rendered grid is returned as the answer (solid `█` blocks are far easier to
read than `#`).

## Go

`converge` steps the extent forward one second at a time until it stops
shrinking; `render` then paints that frame into a `[][]rune` grid and joins the
rows with `\n` (solid `█` blocks read far more legibly than `#`).

```text
────────────────────────────────────────
─     2018 Day 10: The Stars Align     ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS            16.234ms
2.0:  PASS            16.117ms
```

## Rust

Points are a plain `Star` struct in a `Vec`; parsing splits each line on any
non-digit (keeping `-`) to pull the four signed integers, so no regex crate is
needed. `converge` and `render` share the same single-pass bounding-box scan,
and the grid is a `Vec<Vec<char>>` whose rows are collected into `String`s and
joined with `\n`.

```text
────────────────────────────────────────
─     2018 Day 10: The Stars Align     ─
────────────────────────────────────────

Solving (Rust)…
1.0:  PASS            21.900ms
2.0:  PASS            21.101ms
```

## Python

positions and velocities are two NumPy arrays, so advancing every point to time
`t` is a single `pos + vel * t` and the bounding box is `min`/`max` along the
axes. `converge` walks the extent to its minimum; `render` stamps the lit cells
into a `full((h, w), " ")` char grid by fancy-indexing and joins the rows.

```text
────────────────────────────────────────
─     2018 Day 10: The Stars Align     ─
────────────────────────────────────────

Solving (Python)…
1.0:  PASS           845.396ms
2.0:  PASS           809.949ms
```

## Visualization

The stars caught at the moment they align into the message. Each star is a bright
point on a dark field; the frame is the tightest-bounding-box instant found by the
solver. Brightness alone carries the letters, so the picture reads in grayscale.

![The message the stars spell](the-stars-align.png)

## Run Times

![Day 10 run-time graphs](run-times.png)
