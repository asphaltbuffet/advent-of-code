# [Day 7: Camel Cards](https://adventofcode.com/2023/day/7)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 7: Camel Cards][rm7]
[Go][go7]
[Python][py7]

[rm7]: 07-camelCards/README.md
[go7]: 07-camelCards/go
[py7]: 07-camelCards/py

-->

## Go

```text
 ───────────────────
 ADVENT OF CODE 2023
 Day 7: Camel Cards
 ───────────────────
  Test 1.1: PASS in 28.5 µs
  Test 2.1: PASS in 19.2 µs
    Part 1: 246912307 1.3 ms
    Part 2: 246894760 1.4 ms
```

## Python

Each hand becomes a sort key `(type, card_values)`: `Counter` gives the
count-multiset that ranks the type, and part two pops the jokers onto the most
common card before typing while ordering `J` weakest. Sorting the keyed hands and
weighting by rank gives the winnings.

```text
 ───────────────────
 ADVENT OF CODE 2023
 Day 7: Camel Cards
 ───────────────────

Solving (Python)…
1.0:  PASS             4.956ms
2.0:  PASS             4.925ms
```

## Rust

The same key `(type_rank, [card_rank; 5])`, but leaning on Rust's derived tuple
ordering — one `sort_by` ranks every hand. A small count array yields the type;
jokers are folded onto the largest group for typing and mapped to strength 0 for
tiebreaks.

```text
 ───────────────────
 ADVENT OF CODE 2023
 Day 7: Camel Cards
 ───────────────────

Solving (Rust)…
1.0:  PASS           193.682µs
2.0:  PASS           180.492µs
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
