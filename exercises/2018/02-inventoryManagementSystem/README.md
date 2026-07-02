# [Day 2: Inventory Management System](https://adventofcode.com/2018/day/2)

<!-- [Day 2: Inventory Management System](02-inventoryManagementSystem) -->

## Notes

The box IDs are strings of lowercase letters.

Part One is a checksum: count the IDs that contain some letter exactly twice, and
those that contain some letter exactly three times, then multiply the two counts.
A fixed `[26]int` tally per ID is enough — scan for a count of 2 and a count of 3.

Part Two finds the two correct boxes, whose IDs differ by a single character at
the same position, and returns their common letters. A pairwise scan compares
each pair character by character, bailing as soon as a second mismatch appears;
the one pair with exactly one difference yields the answer by dropping that index.

## Go

```text
────────────────────────────────────────
─   2018 Day 2: Inventory Management   ─
─                System                ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS            36.993µs
2.0:  PASS           109.755µs
```

## Python

```text
    < section intentionally left blank >
```

## Run Times

![Day 2 run-time graphs](run-times.png)
