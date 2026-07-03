# [Day 5: If You Give A Seed A Fertilizer](https://adventofcode.com/2023/day/5)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 5: If You Give A Seed A Fertilizer][rm05]
[Go][go05]
[Rust][rs05]
[Python][py05]

[rm05]: 05-ifYouGiveASeedAFertilizer/README.md
[go05]: 05-ifYouGiveASeedAFertilizer/go
[rs05]: 05-ifYouGiveASeedAFertilizer/rs
[py05]: 05-ifYouGiveASeedAFertilizer/py

-->

## Go

```text
────────────────────────────────────────
─   2023 Day 5: If You Give A Seed A   ─
─              Fertilizer              ─
────────────────────────────────────────
Solving (Go)…
1.0:  PASS             2.034ms
2.0:  PASS             1.230ms
```

## Python

Part one maps each seed through the chain with a linear rule scan. Part two keeps
a worklist of half-open intervals and splits each against the overlapping rule,
pushing the unmatched fragments back on — so entire seed ranges flow through the
maps at once.

```text
────────────────────────────────────────
─   2023 Day 5: If You Give A Seed A   ─
─              Fertilizer              ─
────────────────────────────────────────
Solving (Python)…
1.0:  PASS             0.363ms
2.0:  PASS             1.581ms
```

## Rust

Part one `fold`s each seed through the maps. Part two carries a `Vec<(u64,u64)>`
of intervals and, per map, pops each interval, clips it to the first overlapping
rule, emits the mapped slice, and re-queues the leftover ends — the same
interval-splitting strategy, allocation-light.

```text
────────────────────────────────────────
─   2023 Day 5: If You Give A Seed A   ─
─              Fertilizer              ─
────────────────────────────────────────
Solving (Rust)…
1.0:  PASS            44.992µs
2.0:  PASS            71.767µs
```

## Run Times

![Day 5 run-time graphs](run-times.png)
