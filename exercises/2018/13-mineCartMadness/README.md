# [Day 13: Mine Cart Madness](https://adventofcode.com/2018/day/13)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 13: Mine Cart Madness][rm13]
[Go][go13]
[Rust][rs13]
[Python][py13]

[rm13]: 13-mineCartMadness/README.md
[go13]: 13-mineCartMadness/go
[rs13]: 13-mineCartMadness/rs
[py13]: 13-mineCartMadness/py

-->

## Notes

Carts ride a track of straights, curves (`/` `\`), and intersections (`+`). The
carts are lifted off the grid at parse time, each remembering the straight track
it was sitting on, so the map underneath stays intact. Every tick the carts move
in reading order; a curve reflects the heading and an intersection cycles the cart
through left, straight, right.

- **Part One** reports the location of the first collision.
- **Part Two** instead removes both carts in every collision and keeps going until
  a single cart remains, reporting its position at the end of that tick.

## Go

```text
────────────────────────────────────────
─    2018 Day 13: Mine Cart Madness    ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS             0.222ms
2.0:  PASS             1.795ms
```

## Rust

Carts are a small struct in a `Vec`, re-sorted by `(y, x)` each tick so they move
in reading order. The track lives as `Vec<Vec<u8>>` with carts lifted off onto the
straight beneath them. Collisions are checked the instant a cart lands — mid-tick,
matching the immediate-crash rule — and part two `retain`s the survivors after each
tick until one remains.

```text
────────────────────────────────────────
─    2018 Day 13: Mine Cart Madness    ─
────────────────────────────────────────

Solving (Rust)…
1.0:  PASS             0.321ms
2.0:  PASS             1.464ms
```

## Python

Carts are `[x, y, dx, dy, turns]` lists sorted by `(y, x)` each tick; the grid is a
list of strings with the underlying track preserved. Each cart is checked for a
collision immediately after it steps, so a cart moved later in a tick can still hit
one that already moved. Part two collects dead indices per tick and rebuilds the
cart list until a single survivor is left.

```text
────────────────────────────────────────
─    2018 Day 13: Mine Cart Madness    ─
────────────────────────────────────────

Solving (Python)…
1.0:  PASS            10.312ms
2.0:  PASS            73.621ms
```

## Visualization

An animation of the carts running the track. The track is drawn dim; live carts
are bright blue and every crash leaves a lingering orange marker, each haloed so it
stays legible over the busy network. There are more than ten thousand ticks, so
frames are sampled — but any tick with a crash is always captured, so every wreck
and the final lone survivor appear. Events read by brightness and halo rather than
hue, so it works in grayscale too.

![Carts running the track until one survives](mine-cart-madness.gif)

## Run Times

![Day 13 run-time graphs](run-times.png)
