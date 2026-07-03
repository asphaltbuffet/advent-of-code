# [Day 2: Gift Shop](https://adventofcode.com/2025/day/2)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 2: Gift Shop][rm2]
[Go][go2]
[Python][py2]

[rm2]: 02-giftShop/README.md
[go2]: 02-giftShop/go
[py2]: 02-giftShop/py

-->

## Go

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025            
             Day 2: Gift Shop             
──────────────────────────────────────────
Solving (Go)...
  1.1: PASS                 42µs
      ⤷ 12599655151
  2.1: PASS              1.834ms
      ⤷ 20942028255
```

## Python

Rather than scan each range, generate the qualifying IDs from their structure. An
"invalid" ID (part one) is a half repeated exactly twice, built with `str(h) * 2`;
a "repeated" ID (part two) is a `p`-digit pattern tiled `r >= 2` times, `int(str(q)
* r)`. A `set` folds away IDs reachable through more than one factoring.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
             Day 2: Gift Shop
──────────────────────────────────────────

Solving (Python)…
1.0:  PASS           118.211ms
2.0:  PASS            63.646ms
```

## Rust

The same structural generation, tiling patterns arithmetically (`n = n * 10^p + q`)
to avoid string allocation in the hot loop, with a `HashSet<u64>` for the part-two
dedup.

```text
──────────────────────────────────────────
           ADVENT OF CODE 2025
             Day 2: Gift Shop
──────────────────────────────────────────

Solving (Rust)…
1.0:  PASS           786.718µs
2.0:  PASS           664.114µs
```

## 2025 Run Times

![2025 exercise run-time graphs](../run-times.png)
