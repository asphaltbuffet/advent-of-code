# [Day 7: The Sum of Its Parts](https://adventofcode.com/2018/day/7)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 7: The Sum of Its Parts][rm7]
[Go][go7]
[Rust][rs7]
[Python][py7]

[rm7]: 07-theSumOfItsParts/README.md
[go7]: 07-theSumOfItsParts/go
[rs7]: 07-theSumOfItsParts/rs
[py7]: 07-theSumOfItsParts/py

-->

## Notes

The instructions form a dependency graph; each `Step X must be finished before
step Y can begin` line is an edge, and we track for every step the set of steps
that must precede it.

- **Part One** is a lexicographic topological sort: repeatedly take the
  alphabetically first step whose prerequisites are all complete, append it to the
  order, and mark it done.
- **Part Two** simulates the work second by second with a pool of workers. Each
  second, finished steps are retired, then idle workers pick up ready steps
  (alphabetically first), each taking `base + letterIndex` seconds. The elapsed
  time when the last step finishes is the answer.

The example runs 2 workers with no base cost while the real puzzle runs 5 workers
with a 60-second base, so those parameters are chosen from the step count rather
than hardcoded.

## Go

```text
────────────────────────────────────────
─   2018 Day 7: The Sum of Its Parts   ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS             0.210ms
2.0:  PASS             2.042ms
```

## Rust

Parses lines by fixed byte offsets and stores prerequisites in a
`BTreeMap<u8, BTreeSet<u8>>`, so both parts iterate steps already in alphabetical
order and readiness is a single `prereqs.is_subset(&done)` check — the first ready
step in iteration order is the lexicographic pick with no extra sort. Part Two
drives the clock with `for t in 0..`, retiring finished steps via
`in_progress.retain`, and selects the worker/base regime (2 workers / base 0 for
the six-step example, 5 workers / base 60 for the real input) from the step count.

```text
────────────────────────────────────────
─   2018 Day 7: The Sum of Its Parts   ─
────────────────────────────────────────

Solving (Rust)…
1.0:  PASS            59.117µs
2.0:  PASS           867.429µs
```

## Python

A plain `dict[str, set[str]]` of prerequisites feeds a one-line
`sorted(... if pre <= done)` to pick the next ready step, so the lexicographic
topological sort reads directly from the set-subset test. Part Two runs the same
second-by-second simulation, retiring finished steps then filling idle workers, and
chooses 2 workers / base 0 for the example versus 5 workers / base 60 for the real
input by the number of distinct steps.

```text
────────────────────────────────────────
─   2018 Day 7: The Sum of Its Parts   ─
────────────────────────────────────────

Solving (Python)…
1.0:  PASS             0.282ms
2.0:  PASS             3.982ms
```

## Visualization

The Part Two worker schedule as a Gantt chart: one lane per worker, time running
left to right, each step a labeled bar from its start to its finish. Bars alternate
between Okabe-Ito blue and orange — distinct in hue *and* brightness — while idle
gaps stay dark. The picture shows why the total is 914 seconds: workers 4 and 5 go
idle early once their steps clear, and the long tail on worker 1 (`… V R H T`) is
the critical path that no amount of extra workers can shorten.

![Part Two worker schedule](the-sum-of-its-parts.png)

## Run Times

![Day 7 run-time graphs](run-times.png)
