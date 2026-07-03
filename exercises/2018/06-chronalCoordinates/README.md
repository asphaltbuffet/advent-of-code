# [Day 6: Chronal Coordinates](https://adventofcode.com/2018/day/6)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 6: Chronal Coordinates][rm6]
[Go][go6]
[Rust][rs6]
[Python][py6]

[rm6]: 06-chronalCoordinates/README.md
[go6]: 06-chronalCoordinates/go
[rs6]: 06-chronalCoordinates/rs
[py6]: 06-chronalCoordinates/py

-->

## Notes

The coordinates seed a Manhattan-distance Voronoi diagram on the grid.

- **Part One** assigns each cell in the bounding box to its single nearest
  coordinate (ties leave the cell unowned) and reports the largest *finite* area.
  Any coordinate that owns a cell on the bounding-box edge owns an unbounded
  region — it keeps winning outward forever — so those coordinates are excluded
  from the contest.
- **Part Two** counts the cells whose *summed* distance to every coordinate is
  below a threshold (10000 for the puzzle, 32 for the worked example). Stepping
  one cell outside the box adds at least one distance per coordinate, so the
  region reaches at most `threshold / numCoords` cells beyond the box; scanning
  that padded box captures all of it.

The example's threshold differs from the real run, so the threshold is chosen
from the coordinate count rather than hardcoded.

## Go

```text
────────────────────────────────────────
─   2018 Day 6: Chronal Coordinates    ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS            19.321ms
2.0:  PASS            70.756ms
```

## Rust

Parsing splits each line on non-digit characters and takes the first two integers.
The two scans are plain nested-range loops; Part Two chains the coordinate ranges
with `flat_map` and finishes with a `filter(...).count()`. The threshold (32 for the
worked example, 10000 otherwise) is picked from the coordinate count, and the padded
bounding box is sized as `threshold / numPts + 1`.

```text
────────────────────────────────────────
─   2018 Day 6: Chronal Coordinates    ─
────────────────────────────────────────

Solving (Rust)…
1.0:  PASS             9.044ms
2.0:  PASS            29.913ms
```

## Python

Coordinates are pulled with a regex over each line, then the whole grid is solved
with numpy instead of Python loops. A single broadcast builds the `(H, W, N)`
Manhattan-distance tensor `|gx - cx| + |gy - cy|`. Part One takes `argmin` over the
coordinate axis for ownership, masks cells whose minimum is tied (more than one
distance equals the row min) as unowned, reads the four border edges to find the
coordinates with unbounded regions, and `bincount`s the finite owners. Part Two
reuses the tensor, sums over the coordinate axis, and counts cells below the
threshold with `(total < threshold).sum()`. The threshold (32 for the example,
10000 for the real input) is chosen by coordinate count, and the scan pads the box
by `threshold // numPts + 1`.

```text
────────────────────────────────────────
─   2018 Day 6: Chronal Coordinates    ─
────────────────────────────────────────

Solving (Python)…
1.0:  PASS           167.008ms
2.0:  PASS           576.305ms
```

## Visualization

The Manhattan-distance Voronoi map: every cell is tinted by which coordinate
owns it, with contested (tied) cells left dark. The largest finite territory —
the Part One answer — is brightened, and the safe region where total distance to
all coordinates stays under 10000 (the Part Two answer) is outlined as a bright
contour. Coordinate seeds are marked. Territories vary in brightness and the two
answers are carried by brightness and outline, so the map reads in grayscale.

![Voronoi territories with the safe region](chronal-coordinates.png)

## Run Times

![Day 6 run-time graphs](run-times.png)
