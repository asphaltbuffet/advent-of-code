# [Day 16: Chronal Classification](https://adventofcode.com/2018/day/16)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 16: Chronal Classification][rm16]
[Go][go16]
[Rust][rs16]
[Python][py16]

[rm16]: 16-chronalClassification/README.md
[go16]: 16-chronalClassification/go
[rs16]: 16-chronalClassification/rs
[py16]: 16-chronalClassification/py

-->

## Notes

The device has sixteen operations (register/immediate variants of add, multiply,
bitwise and/or, assignment, and greater-than/equality tests). The input gives
Before/After samples for numbered opcodes, then a program using those numbers.

- **Part One** counts the samples that behave like three or more of the sixteen
  operations.
- **Part Two** deduces which number is which operation: each opcode's candidate set
  is the operations consistent with every sample that uses it, then the mapping is
  resolved by elimination — any opcode left with a single candidate is fixed and
  that operation is struck from the others, repeating until all sixteen are pinned.
  The program is then run and register 0 read out.

## Go

```text
────────────────────────────────────────
─ 2018 Day 16: Chronal Classification  ─
────────────────────────────────────────

Solving (Go)…
1.0:  PASS             3.842ms
2.0:  PASS             4.605ms
```

## Rust

The sixteen operations are a `const` array of `fn` pointers indexed like the Go
table, so counting the ones that reproduce a sample's After registers is a filter
over `0..16`. Each sample's matches are folded into a `u16` bitmask, so a number's
candidate set is a single `&=` and elimination is `count_ones() == 1` plus a
`&= !bit` to strike the resolved operation from the others.

```text
────────────────────────────────────────
─ 2018 Day 16: Chronal Classification  ─
────────────────────────────────────────

Solving (Rust)…
1.0:  PASS             0.980ms
2.0:  PASS             1.605ms
```

## Python

The operations are lambdas in a name-keyed `dict`, and a sample's matching
operations come back as a `set` from a comprehension. Part One counts the samples
with three or more; Part Two intersects each opcode number's candidate `set` across
its samples, then repeatedly fixes any number down to a single candidate and
`discard`s that operation from the rest until all sixteen are pinned.

```text
────────────────────────────────────────
─ 2018 Day 16: Chronal Classification  ─
────────────────────────────────────────

Solving (Python)…
1.0:  PASS            19.438ms
2.0:  PASS            19.592ms
```

## Visualization

The opcode deduction as a 16×16 matrix: rows are opcode numbers, columns the
sixteen operations. A cell is bright where that operation is still a candidate for
the number (consistent with every sample) and dark where it is ruled out; the one
operation the solver pins down for each number is boxed. You can read the puzzle's
structure directly — some numbers start nearly unambiguous while opcode 6 is
consistent with almost every operation — and the boxes trace the one-to-one mapping
that elimination recovers. Candidacy is carried by brightness and the resolution by
the box, so it reads in grayscale.

![Opcode constraint matrix](chronal-classification.png)

## Run Times

![Day 16 run-time graphs](run-times.png)
