# [Day 3: No Matter How You Slice It](https://adventofcode.com/2018/day/3)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 3: No Matter How You Slice It][rm3]
[Go][go3]
[Rust][rs3]
[Python][py3]

[rm3]: 03-noMatterHowYouSliceIt/README.md
[go3]: 03-noMatterHowYouSliceIt/go
[rs3]: 03-noMatterHowYouSliceIt/rs
[py3]: 03-noMatterHowYouSliceIt/py

-->

## Notes

Each claim is a rectangle on a shared piece of fabric, given as
`#id @ left,top: widthxheight`. Marking every cell each claim covers in a
coverage map answers both parts:

- **Part One** counts the cells covered by two or more claims — the contested
  fabric.
- **Part Two** finds the single claim whose every cell is covered exactly once,
  and returns its id — the one intact claim.

The claims are parsed by scanning the five integers on each line, which is robust
to whitespace and delivery quirks in the input.

## Go

```text
────────────────────────────────────────
─2018 Day 3: No Matter How You Slice It─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS            65.918ms
2.0:  PASS            63.028ms
```

## Rust

The claim is a plain struct and coverage is a flat `vec![0u16; 1000 * 1000]`
grid indexed by `y * SIDE + x`, avoiding hashing. Both parts read as iterator
chains: Part One `filter(|&&n| n >= 2).count()`, and Part Two `find`s the claim
whose cells are all covered exactly once with nested `all`.

```text
────────────────────────────────────────
─2018 Day 3: No Matter How You Slice It─
────────────────────────────────────────

Solving (Rust)…
1.0:  PASS            54.099ms
2.0:  PASS             4.934ms
```

## Python

A regex scans the five integers per line and `collections.Counter.update` tallies
each claim's cells in one call. Part One counts values `>= 2`; Part Two returns the
first claim whose cells are all `== 1`, tested with a generator `all(...)`.

```text
────────────────────────────────────────
─2018 Day 3: No Matter How You Slice It─
────────────────────────────────────────

Solving (Python)…
1.0:  PASS              1.394s
2.0:  PASS              0.990s
```

## Visualization

The fabric as a heatmap: each cell is shaded by how many claims cover it, so
uncontested single-claim fabric is dim and the contested overlaps (the Part One
answer) stand out brightly. The one intact claim (the Part Two answer) is framed
and labeled. Coverage is carried by brightness and the answer by an outline and
label, so the image reads in grayscale.

![Fabric coverage heatmap](no-matter-how-you-slice-it.png)

## Run Times

![Day 3 run-time graphs](run-times.png)
