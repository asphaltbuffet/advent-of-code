# [Day 19: Aplenty](https://adventofcode.com/2023/day/19)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 19: Aplenty][rm19]
[Go][go19]
[Python][py19]

[rm19]: 19-aplenty/README.md
[go19]: 19-aplenty/go
[py19]: 19-aplenty/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023            
             Day 19: Aplenty              
──────────────────────────────────────────
          
Testing...
  1.1: PASS              0.06 ms
  2.1: PASS              0.02 ms
          
Solving...
    1: PASS              1.19 ms
      ⤷ 325952
    2: PASS              1.44 ms
      ⤷ 125744206494820
```

## Python

Workflows parse into rule lists. Part one runs each part through them from `in`
and sums the accepted ratings. Part two recurses over the 4D range `xmas ∈
[1,4000]`, splitting each category at a rule's threshold — the matching sub-range
follows the target, the remainder falls through — and sums the accepted volumes.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
             Day 19: Aplenty
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS             2.070ms
2.0:  PASS             2.024ms
```

## Rust

The same evaluation and range-splitting, with ranges as `[(u64,u64); 4]` indexed
by category and workflows borrowing their rule targets from the input. Part two's
recursion multiplies the surviving span widths at each accepted leaf.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
             Day 19: Aplenty
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS           282.241µs
2.0:  PASS           187.141µs
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
