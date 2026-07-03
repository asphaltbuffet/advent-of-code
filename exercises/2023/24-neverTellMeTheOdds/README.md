# [Day 24: Never Tell Me The Odds](https://adventofcode.com/2023/day/24)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 24: Never Tell Me The Odds][rm24]
[Go][go24]
[Python][py24]

[rm24]: 24-neverTellMeTheOdds/README.md
[go24]: 24-neverTellMeTheOdds/go
[py24]: 24-neverTellMeTheOdds/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023            
      Day 24: Never Tell Me The Odds      
──────────────────────────────────────────
          
Testing...
  1.1: PASS              0.04 ms
  2.1: ERROR 
      ⤷ saying: no solution found
          
Solving...
    1: PASS              3.77 ms
      ⤷ 16665
    2: PASS             12.11 ms
      ⤷ 769840447420960
```

## Python

Part one tests each pair's XY crossing with exact `Fraction` arithmetic, keeping
only future intersections inside the test area (auto-selected small/large by
coordinate scale). Part two uses the rock-frame identity `(P_i - R) x (V_i - W) =
0`; differencing stone 0 against stones 1 and 2 cancels the `R x W` term, and the
resulting 6x6 system is solved by exact Gaussian elimination.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
      Day 24: Never Tell Me The Odds
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS           337.959ms
2.0:  PASS             1.058ms
```

## Rust

The same crossing test (in `f64`, which is ample for the inequality) and the same
6x6 linear system, solved with `i128` fraction-free Gaussian elimination — rows
reduced by their gcd each step so the 15-digit answer stays exact without a
bignum crate.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
      Day 24: Never Tell Me The Odds
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS             1.324ms
2.0:  PASS             0.137ms
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
