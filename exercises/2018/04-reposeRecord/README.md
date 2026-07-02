# [Day 4: Repose Record](https://adventofcode.com/2018/day/4)

<!-- [Day 4: Repose Record](04-reposeRecord) -->

## Notes

The log lines arrive shuffled, but the `[YYYY-MM-DD HH:MM]` timestamp sorts
correctly as plain text, so sorting the lines lexically restores chronological
order. Walking the sorted log, a `Guard #N begins shift` line sets the active
guard and each `falls asleep` / `wakes up` pair brackets a run of asleep minutes.
Only the minute field (0–59) matters, so the whole record collapses to a tally
`asleep[guard][minute]` = number of days that guard was asleep at that minute.

- **Part One (Strategy 1):** pick the guard with the most total asleep minutes,
  then the minute they were asleep most often; answer is guard × minute.
- **Part Two (Strategy 2):** pick the single (guard, minute) cell with the
  highest count across every guard; answer is guard × minute.

## Go

```text
────────────────────────────────────────
─      2018 Day 4: Repose Record       ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS           193.115µs
2.0:  PASS           167.753µs
```

## Python

```text
    < section intentionally left blank >
```

## Visualization

A heatmap of the sleep tally: one row per guard, columns for minutes 0–59, each
cell shaded by how many days that guard was asleep at that minute. Strategy 1's
guard (brightest row overall) and Strategy 2's cell (the single brightest cell)
are outlined and labeled. Counts read as brightness and the two answers by
outline and label, so the heatmap reads in grayscale.

![Guard sleep heatmap](repose-record.png)

## Run Times

![Day 4 run-time graphs](run-times.png)
