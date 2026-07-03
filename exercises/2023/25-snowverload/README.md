# [Day 25: Snowverload](https://adventofcode.com/2023/day/25)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 25: Snowverload][rm25]
[Go][go25]
[Python][py25]

[rm25]: 25-snowverload/README.md
[go25]: 25-snowverload/go
[py25]: 25-snowverload/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023            
           Day 25: Snowverload            
──────────────────────────────────────────
          
Testing...
  1.1: PASS              0.54 ms
          
Solving...
    1: PASS            179.09 ms
      ⤷ 562912
```

## Python

The two clusters are joined by exactly three wires, so the min cut between nodes
on opposite sides is 3. Fixing a source and running unit-capacity max flow (BFS
augmenting paths) to each other node, the first that saturates at flow 3 splits
the graph; the residual-reachable set is one cluster, and the sizes multiply. Day
25 has no second puzzle, so part two returns the empty finale answer.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
           Day 25: Snowverload
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS            17.416ms
```

## Rust

The same max-flow min-cut over interned node indices, with residual capacities in
a `HashMap`. Because the answer is a genuine min cut rather than a betweenness
heuristic, it is exact on both the small example and the full graph.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
           Day 25: Snowverload
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS             3.394ms
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
