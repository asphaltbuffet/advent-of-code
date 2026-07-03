# [Day 10: Factory](https://adventofcode.com/2025/day/10)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 10: Factory][rm10]
[Go][go10]
[Python][py10]

[rm10]: 10-factory/README.md
[go10]: 10-factory/go
[py10]: 10-factory/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
              Day 10: Factory
──────────────────────────────────────────
Solving (Go)...
  1.1: PASS              0.001s
      ⤷ 375
  2.1: PASS              6.138s
      ⤷ 15377
```

## Python

Part one finds the fewest buttons whose XOR reproduces the light pattern. Part two
reaches the joltage vector in base 2 — the buttons pressed at a level fix each
component's low bit, then the remainder halves and recurses at double cost.
Precomputing every button-subset's component-sum vector keeps the recursion cheap.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
              Day 10: Factory
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS             0.005s
2.0:  PASS            28.741s
```

## Rust

The same base-2 recursion with precomputed subset sums, over plain vectors and a
`HashMap` memo. Native speed collapses the multi-second part-two search to under a
second.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
              Day 10: Factory
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS             0.247ms
2.0:  PASS           506.937ms
```

## 2025 Run Times

![2025 exercise run-time graphs](../run-times.png)
