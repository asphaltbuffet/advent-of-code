# [Day 15: Lens Library](https://adventofcode.com/2023/day/15)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 15: Lens Library][rm15]
[Go][go15]
[Python][py15]

[rm15]: 15-lensLibrary/README.md
[go15]: 15-lensLibrary/go
[py15]: 15-lensLibrary/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
           Day 15: Lens Library
──────────────────────────────────────────

Testing...
  1.1: PASS              0.00 ms
  2.1: PASS              0.01 ms

Solving...
    1: PASS              0.07 ms
      ⤷ 511343
    2: PASS              0.56 ms
      ⤷ 294474
```

## Python

HASH is a `functools.reduce` fold. Part two uses a list of 256 dicts, one per
box: because a `dict` preserves insertion order, `=` updates a lens in place and
`-` pops it, matching the required slot semantics without a manual linked list.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
           Day 15: Lens Library
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS             3.695ms
2.0:  PASS             4.080ms
```

## Rust

HASH is a byte `fold`. Each box is a `Vec<(&str, u8)>` borrowing labels from the
input; `position`/`find` locate a lens to remove or update in place, preserving
order, and the focusing power falls out of a `flat_map` over the boxes.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
           Day 15: Lens Library
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS            99.030µs
2.0:  PASS           232.362µs
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
