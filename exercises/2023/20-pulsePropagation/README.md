# [Day 20: Pulse Propagation](https://adventofcode.com/2023/day/20)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 20: Pulse Propagation][rm20]
[Go][go20]
[Python][py20]

[rm20]: 20-pulsePropagation/README.md
[go20]: 20-pulsePropagation/go
[py20]: 20-pulsePropagation/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
        Day 20: Pulse Propagation
──────────────────────────────────────────

Testing...
  1.1: PASS              0.87 ms
  1.2: PASS              0.72 ms
  2.1: PASS             18.75 ms

Solving...
    1: PASS              4.03 ms
      ⤷ 684125385
    2: PASS             14.63 ms
      ⤷ 225872806380073
```

## Python

Modules are dicts of kind, outputs, flip state, and conjunction memory; each
button press drains a `deque` of pulses. Part one multiplies the low/high counts
over 1000 presses. Part two finds the single conjunction feeding `rx` and LCMs
the presses at which each of its inputs first sends a high pulse.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
        Day 20: Pulse Propagation
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS            14.869ms
2.0:  PASS            62.647ms
```

## Rust

The same simulation over a `HashMap` of borrowed module names and a `VecDeque`
pulse queue; `press` reports which sources sent a high pulse into a watched
module, which part two uses to gather the per-input periods before folding them
with a gcd-based LCM.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2023
        Day 20: Pulse Propagation
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS             2.252ms
2.0:  PASS             8.446ms
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
