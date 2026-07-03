# [Day 8: Memory Maneuver](https://adventofcode.com/2018/day/8)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 8: Memory Maneuver][rm8]
[Go][go8]
[Rust][rs8]
[Python][py8]

[rm8]: 08-memoryManeuver/README.md
[go8]: 08-memoryManeuver/go
[rs8]: 08-memoryManeuver/rs
[py8]: 08-memoryManeuver/py

-->

## Notes

The input is a flat list of numbers encoding a tree; each node is `childCount
metaCount`, then its children, then its metadata. A single recursive reader with
an index cursor walks the whole structure in one pass.

- **Part One** sums every metadata entry in the tree.
- **Part Two** computes the root's value: a leaf's value is the sum of its
  metadata, while an internal node's metadata are 1-based indices into its
  children, and its value is the sum of the referenced children's values
  (out-of-range indices skipped, repeats counted).

## Go

```text
────────────────────────────────────────
─    2018 Day 8: Memory Maneuver       ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS             5.263ms
2.0:  PASS             5.821ms
```

## Rust

The number stream is a plain `Iterator<Item = usize>` threaded through the
recursion by `&mut`, so each node just pulls its header, recurses for its
children, then consumes its metadata — no index bookkeeping. Part Two collects
child values into a small `Vec` and uses `get(r - 1)` (via `wrapping_sub` so the
out-of-range check is one lookup) to sum the referenced children.

```text
────────────────────────────────────────
─     2018 Day 8: Memory Maneuver      ─
────────────────────────────────────────

Solving (Rust)…
1.0:  PASS           296.911µs
2.0:  PASS           346.133µs
```

## Python

Mirrors the Rust: a single `iter(...)` over the split tokens is passed down the
recursion, and each node calls `next()` for its header, children, and metadata.
Part Two builds a list of child values and sums the ones its 1-based metadata
index into, falling back to summing the metadata directly for leaves.

```text
────────────────────────────────────────
─     2018 Day 8: Memory Maneuver      ─
────────────────────────────────────────

Solving (Python)…
1.0:  PASS             8.043ms
2.0:  PASS             7.699ms
```

## Run Times

![Day 8 run-time graphs](run-times.png)
