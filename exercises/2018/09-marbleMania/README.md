# [Day 9: Marble Mania](https://adventofcode.com/2018/day/9)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 9: Marble Mania][rm9]
[Go][go9]
[Rust][rs9]
[Python][py9]

[rm9]: 09-marbleMania/README.md
[go9]: 09-marbleMania/go
[rs9]: 09-marbleMania/rs
[py9]: 09-marbleMania/py

-->

## Notes

The marble ring is a **circular doubly-linked list**, so both game moves are O(1):
a normal marble is inserted between the two marbles clockwise of the current one,
and every 23rd marble scores — the current player takes it plus the marble seven
counter-clockwise, which is then removed.

- **Part One** plays the game as given and reports the winning player's score.
- **Part Two** is the same game with the last marble worth 100 times as much (~7.1
  million marbles). The linked-list ring keeps every insert and removal constant
  time, so the larger game still finishes in well under a second — a slice-based
  ring, with its O(n) shifts, would not.

## Go

```text
────────────────────────────────────────
─      2018 Day 9: Marble Mania        ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS             1.846ms
2.0:  PASS           563.857ms
```

## Rust

A `VecDeque` stands in for the ring, with the current marble kept at the back.
`rotate_left`/`rotate_right` walk the circle clockwise or counter-clockwise in O(1),
so the normal two-clockwise insertion and the seven-counter-clockwise scoring
removal both stay constant time — the part-two game of millions of marbles finishes
in a fraction of a second.

```text
────────────────────────────────────────
─      2018 Day 9: Marble Mania        ─
────────────────────────────────────────

Solving (Rust)…
1.0:  PASS             1.266ms
2.0:  PASS           118.104ms
```

## Python

`collections.deque` models the ring, and `deque.rotate` is the idiomatic fast tool:
rotating keeps the current marble at the right end so both moves are O(1) appends and
pops. A plain list would shift elements on every insert and removal, turning the
part-two game into an O(n²) crawl; the deque keeps it near a second.

```text
────────────────────────────────────────
─      2018 Day 9: Marble Mania        ─
────────────────────────────────────────

Solving (Python)…
1.0:  PASS              0.012s
2.0:  PASS              1.260s
```

## Visualization

The Part One game as a bar chart of every player's final score, in player order.
The winning player (107) is drawn bright and labeled; the rest share one dim tint,
so the winner reads by brightness alone and the chart survives grayscale. The
uneven profile comes from the every-23rd scoring marble landing on players in a
fixed cycle, so a player's take depends on their position modulo the scoring
period.

![Part One scores by player](marble-mania.png)

## Run Times

![Day 9 run-time graphs](run-times.png)
