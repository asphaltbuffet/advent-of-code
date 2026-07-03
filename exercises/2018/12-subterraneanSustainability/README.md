# [Day 12: Subterranean Sustainability](https://adventofcode.com/2018/day/12)

<!-- [Day 12: Subterranean Sustainability](12-subterraneanSustainability) -->

## Notes

A one-dimensional cellular automaton: each generation, a pot's next state comes
from a rule keyed on its five-pot neighborhood. The pots live in a `map[int]bool`
so the pattern can grow left or right without bounds.

- **Part One** runs 20 generations and sums the indices of the planted pots.
- **Part Two** runs 50 billion generations, which is only feasible because the
  pattern settles into a fixed shape that drifts sideways by a constant offset
  each generation. Once a normalized shape repeats, the planted count is fixed, so
  the index sum grows linearly and the remaining generations are extrapolated
  rather than simulated.

## Go

```text
────────────────────────────────────────
─ 2018 Day 12: Subterranean Sustainab… ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS             0.324ms
2.0:  PASS             5.795ms
```

## Python

```text
    < section intentionally left blank >
```

## Visualization

The automaton as a space-time diagram: time flows downward, one row per
generation, planted pots drawn bright. The early generations churn, then the
pattern resolves into clean parallel diagonal streaks — the fixed shape drifting
sideways at constant speed. That steady drift is exactly what lets Part Two skip
50 billion generations: once the diagonals are straight, the index sum just grows
linearly. The generation where the shape locks in is marked.

![Space-time diagram of the automaton](subterranean-sustainability.png)

## Run Times

![Day 12 run-time graphs](run-times.png)
