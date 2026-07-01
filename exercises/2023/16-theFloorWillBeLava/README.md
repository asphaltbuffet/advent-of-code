# [Day 16: The Floor Will Be Lava](https://adventofcode.com/2023/day/16)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 16: The Floor Will Be Lava][rm16]
[Go][go16]

[rm16]: 16-theFloorWillBeLava/README.md
[go16]: 16-theFloorWillBeLava/go

-->

## Notes

Each of the ~440 edge entry points is an independent beam simulation, so Part
Two fans them out across `GOMAXPROCS` worker goroutines and takes the maximum
energized count.

## Go

```text
────────────────────────────────────────
─ 2023 Day 16: The Floor Will Be Lava  ─
────────────────────────────────────────
Solving (Go)…
1.0:  PASS              0.016s
2.0:  PASS              1.717s
```

## Visualization

The contraption energized by the strongest Part Two start beam. Tiles the beam
reaches glow gold; the dark cells are floor it never touches, and the blue-grey
specks are mirrors and splitters (brightened where the beam passes through
them). The dense grid of light traces how the `|` and `-` splitters fan the beam
across the whole floor.

![The Floor Will Be Lava energized map](floor-will-be-lava.png)

## Run Times

![Day 16 run-time graphs](run-times.png)
