// Solution for Advent of Code 2025 day 2.
//
// Each comma-separated token is an inclusive ID range. We never scan the range;
// instead we *generate* the qualifying IDs from their structure. An "invalid" ID
// (part one) is a `p`-digit pattern tiled exactly twice; a "repeated" ID (part
// two) is a pattern tiled any number of times >= 2. Building candidates directly
// keeps the work proportional to the count of qualifying IDs, not the range span.

use std::collections::HashSet;

/// Inclusive (lo, hi) ranges from the comma-separated input.
fn ranges(input: &str) -> impl Iterator<Item = (u64, u64)> + '_ {
    input.trim().split(',').map(|chunk| {
        let (lo, hi) = chunk.split_once('-').unwrap();
        (lo.parse().unwrap(), hi.parse().unwrap())
    })
}

/// Tile pattern `q` (which has `p` digits) `r` times into one integer.
fn tile(q: u64, p: u32, r: u32) -> u64 {
    let base = 10u64.pow(p);
    let mut n = 0;
    for _ in 0..r {
        n = n * base + q;
    }
    n
}

pub fn part_one(input: &str) -> String {
    // Invalid IDs: even width, both halves equal — a half repeated exactly twice.
    let mut sum: u64 = 0;
    for (lo, hi) in ranges(input) {
        let width = hi.to_string().len() as u32;
        for half in 1..=width / 2 {
            for q in 10u64.pow(half - 1)..10u64.pow(half) {
                let n = tile(q, half, 2);
                if (lo..=hi).contains(&n) {
                    sum += n;
                }
            }
        }
    }
    sum.to_string()
}

pub fn part_two(input: &str) -> String {
    // Repeated IDs: any p-digit pattern tiled r >= 2 times. A HashSet folds away
    // IDs reachable through more than one (pattern, repeat) factoring.
    let mut total: u64 = 0;
    for (lo, hi) in ranges(input) {
        let width = hi.to_string().len() as u32;
        let mut found: HashSet<u64> = HashSet::new();
        for p in 1..=width / 2 {
            for r in 2..=width / p {
                for q in 10u64.pow(p - 1)..10u64.pow(p) {
                    let n = tile(q, p, r);
                    if n > hi {
                        break;
                    }
                    if n >= lo {
                        found.insert(n);
                    }
                }
            }
        }
        total += found.iter().sum::<u64>();
    }
    total.to_string()
}
