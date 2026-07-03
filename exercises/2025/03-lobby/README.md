# [Day 3: Lobby](https://adventofcode.com/2025/day/3)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 3: Lobby][rm3]
[Go][go3]
[Python][py3]

[rm3]: 03-lobby/README.md
[go3]: 03-lobby/go
[py3]: 03-lobby/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
                Day 3: Lobby
──────────────────────────────────────────
Solving (Go)...
  1.1: PASS              20.006µs
      ⤷ 17087
  2.1: PASS             160.063µs
      ⤷ 169019504359949
```

## Python

Each line asks for the largest length-k number formed by an in-order subsequence
(k=2, then 12). A monotonic stack sweeps left to right, popping a smaller trailing
digit whenever deletions remain; the first k survivors are the answer.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
                Day 3: Lobby
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS             3.237ms
2.0:  PASS             3.129ms
```

## Rust

The same greedy on `&[u8]`, accumulating the result as a `u64` (`acc * 10 + d`)
rather than building and parsing a string.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
                Day 3: Lobby
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS           146.276µs
2.0:  PASS           147.141µs
```

## 2025 Run Times

![2025 exercise run-time graphs](../run-times.png)
