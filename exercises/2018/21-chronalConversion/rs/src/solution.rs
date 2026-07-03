// Solution for Advent of Code 2018 day 21.
//
// The bound-ip opcode machine from day 19. The program only halts at its single
// `eqrr … 0` instruction, comparing register 0 against a value that evolves in a
// loop. The first such value halts fastest (part one); the last distinct value
// before the sequence repeats halts slowest (part two).
//
// The inner loop is an O(n) digit-by-digit division by 256 — tens of seconds to
// interpret. Instead the two per-input constants (the seed and the multiplier) are
// read out of the program and the recurrence is evaluated directly.

use std::collections::HashSet;

const MASK: u64 = 0xFF_FFFF;

// Read the seed the compared register is reset to and the multiplier applied to it.
fn constants(input: &str) -> (u64, u64) {
    let prog: Vec<Vec<&str>> = input
        .trim()
        .lines()
        .filter(|l| !l.starts_with('#'))
        .map(|l| l.split_whitespace().collect())
        .collect();

    // The value register is the operand of `eqrr … 0` that isn't register 0.
    let val_reg = prog
        .iter()
        .find(|f| f[0] == "eqrr" && (f[1] == "0" || f[2] == "0"))
        .map(|f| if f[1] == "0" { f[2] } else { f[1] })
        .unwrap();

    let seed = prog
        .iter()
        .filter(|f| f[0] == "seti" && f[3] == val_reg)
        .map(|f| f[1].parse::<u64>().unwrap())
        .max()
        .unwrap();
    let mult = prog
        .iter()
        .find(|f| f[0] == "muli" && f[3] == val_reg)
        .map(|f| f[2].parse::<u64>().unwrap())
        .unwrap();
    (seed, mult)
}

// A lazy iterator over the distinct halt values, in order, ending when the
// sequence would repeat. Part one takes the first value; part two drains it to the
// last, so neither computes more of the cycle than it needs.
fn halt_values(input: &str) -> impl Iterator<Item = u64> {
    let (seed, mult) = constants(input);
    let mut seen = HashSet::new();
    let mut acc: u64 = 0;
    std::iter::from_fn(move || {
        let mut hi = acc | 0x1_0000;
        acc = seed;
        loop {
            acc = ((acc + (hi & 0xFF)) & MASK) * mult & MASK;
            if hi < 256 {
                break;
            }
            hi /= 256;
        }
        seen.insert(acc).then_some(acc) // None once acc repeats, ending the iterator
    })
}

pub fn part_one(input: &str) -> String {
    halt_values(input).next().unwrap().to_string()
}

pub fn part_two(input: &str) -> String {
    halt_values(input).last().unwrap().to_string()
}
