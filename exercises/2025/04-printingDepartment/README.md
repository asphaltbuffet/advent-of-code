# [Day 4: Printing Department](https://adventofcode.com/2025/day/4)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 4: Printing Department][rm4]
[Go][go4]
[Python][py4]

[rm4]: 04-printingDepartment/README.md
[go4]: 04-printingDepartment/go
[py4]: 04-printingDepartment/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025            
        Day 4: Printing Department        
──────────────────────────────────────────
Solving (Go)...
  1.1: PASS              6.154ms
      ⤷ 1537
  2.1: PASS              62.68ms
      ⤷ 8707
```

## Python

Rolls are stored as a coordinate `set`; a roll is accessible when fewer than four
of its eight neighbors are occupied. Part two peels the floor in synchronous waves
— each round removes every currently-accessible roll — until a fixed point.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
        Day 4: Printing Department
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS            28.282ms
2.0:  PASS           443.514ms
```

## Rust

The same coordinate `HashSet` and wave peeling. Because removal is monotone
(it only lowers neighbor counts), the total is order-independent, so collecting
each wave and then deleting it is safe.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
        Day 4: Printing Department
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS             2.716ms
2.0:  PASS            54.546ms
```

## 2025 Run Times

![2025 exercise run-time graphs](../run-times.png)
