# [Day 6: Trash Compactor](https://adventofcode.com/2025/day/6)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 6: Trash Compactor][rm6]
[Go][go6]
[Python][py6]

[rm6]: 06-trashCompactor/README.md
[go6]: 06-trashCompactor/go
[py6]: 06-trashCompactor/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025            
          Day 6: Trash Compactor          
──────────────────────────────────────────
Testing (Go)...
  1.1: PASS                 52µs
  2.1: PASS                 18µs
Solving (Go)...
  1.1: PASS              2.176ms
      ⤷ 6378679666679
  2.1: PASS              2.017ms
      ⤷ 11494432585168
```

## Python

The same grid is read two ways. Part one treats whitespace-delimited tokens as a
column table, combining each column's numbers by its operator. Part two reads
digits stacked vertically in each character column, right to left, with operator
columns delimiting problems.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
          Day 6: Trash Compactor
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS             1.002ms
2.0:  PASS             2.677ms
```

## Rust

The same two readings over byte slices, treating right-trimmed short rows as
space-padded and building each vertical number arithmetically.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
          Day 6: Trash Compactor
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS           128.445µs
2.0:  PASS            41.227µs
```

## 2025 Run Times

![2025 exercise run-time graphs](../run-times.png)
