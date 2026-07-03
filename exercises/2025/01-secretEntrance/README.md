# [Day 1: Secret Entrance](https://adventofcode.com/2025/day/1)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 1: Secret Entrance][rm1]
[Go][go1]
[Python][py1]

[rm1]: 01-secretEntrance/README.md
[go1]: 01-secretEntrance/go
[py1]: 01-secretEntrance/py

-->
e0x434C49434B` translates to 'CLICK'.

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025            
          Day 1: Secret Entrance          
──────────────────────────────────────────
Solving (Go)...
  1.1: PASS                128µs
      ⤷ 1135
  2.1: PASS                135µs
      ⤷ 6558
```

## Python

The pointer is reduced modulo 100 each step, mirroring the physical dial. Part one
counts moves that end resting on 0; part two counts every sweep across 0 by tallying
the multiples of 100 the move arc covers, credited on arrival (a half-open
convention) so a landing is counted exactly once.

```text
────────────────────────────────────────
─     2025 Day 1: Secret Entrance      ─
────────────────────────────────────────

Solving (Python)…
1.0:  PASS             1.391ms
2.0:  PASS             1.595ms
```

## Rust

The same reduced-position walk, using `div_euclid`/`rem_euclid` for true floor
division so the leftward arc is counted correctly as the pointer dips below zero
before wrapping.

```text
────────────────────────────────────────
─     2025 Day 1: Secret Entrance      ─
────────────────────────────────────────

Solving (Rust)…
1.0:  PASS            71.217µs
2.0:  PASS           109.603µs
```

## 2025 Run Times

![2025 exercise run-time graphs](../run-times.png)
