# [Day 11: Reactor](https://adventofcode.com/2025/day/11)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 11: Reactor][rm11]
[Go][go11]
[Python][py11]

[rm11]: 11-reactor/README.md
[go11]: 11-reactor/go
[py11]: 11-reactor/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025            
             Day 11: Reactor              
──────────────────────────────────────────
Testing (Go)...
  1.1: PASS                 18µs
  2.1: PASS                 22µs
Solving (Go)...
  1.1: PASS                180µs
      ⤷ 470
  2.1: PASS              1.836ms
      ⤷ 384151614084875
```

## Python

Both parts count paths through ordered waypoints. A path is valid when it visits
each required waypoint in sequence and ends at the last. `functools.cache` memoizes
on (node, remaining waypoints); part two picks the waypoint order by reachability.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
             Day 11: Reactor
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS             0.480ms
2.0:  PASS             2.176ms
```

## Rust

The same ordered-waypoint count over interned node indices. Since the remaining
waypoints are always a suffix, the memo keys on (node, next waypoint index) rather
than hashing a slice.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
             Day 11: Reactor
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS           193.552µs
2.0:  PASS           353.420µs
```

## 2025 Run Times

![2025 exercise run-time graphs](../run-times.png)
