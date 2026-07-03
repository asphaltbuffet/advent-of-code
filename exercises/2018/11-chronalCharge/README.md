# [Day 11: Chronal Charge](https://adventofcode.com/2018/day/11)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 11: Chronal Charge][rm11]
[Go][go11]
[Rust][rs11]
[Python][py11]

[rm11]: 11-chronalCharge/README.md
[go11]: 11-chronalCharge/go
[rs11]: 11-chronalCharge/rs
[py11]: 11-chronalCharge/py

-->

## Notes

Each cell's power comes from the coordinate-and-serial formula, over a 300×300
grid. The key structure is a **summed-area table** — a prefix sum where each entry
holds the total power of the rectangle from `(1,1)` to that cell — so any square's
total is a four-corner O(1) lookup.

- **Part One** scans every 3×3 square for the largest total and reports its
  top-left corner.
- **Part Two** does the same across *all* square sizes 1–300. With the table making
  each square O(1), the whole search is O(gridSize³) and finishes in a fraction of
  a second.

## Go

```text
────────────────────────────────────────
─     2018 Day 11: Chronal Charge      ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS             0.609ms
2.0:  PASS            21.233ms
```

## Rust

A flat `Vec<i32>` holds the summed-area table with a zero-filled first row and
column, so every square total is the same four-corner subtraction with no
edge-casing. Part Two loops sizes 1–300, calling the shared `best` helper, each
square lookup staying O(1).

```text
────────────────────────────────────────
─     2018 Day 11: Chronal Charge      ─
────────────────────────────────────────

Solving (Rust)…
1.0:  PASS             0.766ms
2.0:  PASS            25.511ms
```

## Python

NumPy builds the summed-area table in two `cumsum` passes over a vectorized power
grid, then each size's square-sums are computed for *all* top-left corners at once
by slicing the padded table and taking `argmax`. Part Two just loops the size and
keeps the running best, so the heavy work stays in vectorized array ops.

```text
────────────────────────────────────────
─     2018 Day 11: Chronal Charge      ─
────────────────────────────────────────

Solving (Python)…
1.0:  PASS             4.267ms
2.0:  PASS            36.860ms
```

## Visualization

The 300×300 grid as a power heatmap — brighter cells carry more power. The
`rackID` term gives the field its striking vertical-band moiré texture. The Part
One 3×3 square (thin blue) and the Part Two 14×14 square (thick orange) are
outlined with a dark halo so they stand out over any brightness; both land in the
bright high-power band on the right, which is exactly why they win. The boxes
differ in line weight as well as hue, so the picture reads in grayscale.

![Power heatmap with the winning squares](chronal-charge.png)

## Run Times

![Day 11 run-time graphs](run-times.png)
