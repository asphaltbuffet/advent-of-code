# [Day 15: Beverage Bandits](https://adventofcode.com/2018/day/15)

<!-- [Day 15: Beverage Bandits](15-beverageBandits) -->

## Notes

A full combat simulation, and the day's difficulty is entirely in its tie-breaking
rules. Units act in reading order each round. A unit that is not already next to an
enemy finds the open squares in range of some enemy and takes one step along the
shortest path to the nearest of them — with *every* tie (which target square, which
first step) broken by reading order. It then attacks the adjacent enemy with the
fewest hit points, ties again by reading order.

The reading-order tie-breaks fall out of one trick: a breadth-first search that
expands neighbors up, left, right, down. Because those are the reading-order
directions, the first shortest path the search reaches a square by is already the
reading-order-minimal one, so the recorded first step is correct with no extra
sorting.

- **Part One** runs the fight to the end and reports full rounds × remaining hit
  points.
- **Part Two** raises the elves' attack power from 4 upward until they win with no
  elf deaths (the fight aborts the instant an elf falls), and reports that
  outcome.

## Go

```text
────────────────────────────────────────
─    2018 Day 15: Beverage Bandits     ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS            50.118ms
2.0:  PASS           331.686ms
```

## Python

```text
    < section intentionally left blank >
```

## Visualization

An animation of the part-two winning battle — the lowest elf attack power at which
no elf dies — one frame per round. Walls are dim; elves are bright filled blocks
and goblins are hollow rings, so the two sides differ by shape as well as color and
stay distinct in grayscale. Units dim as their hit points fall and vanish on death.
By the final frame only elves remain, which is exactly the condition part two
searches for.

![The part-two battle, elves winning without a loss](beverage-bandits.gif)

## Run Times

![Day 15 run-time graphs](run-times.png)
