// Solution for Advent of Code 2015 day 11.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

// Advance the password by one, as a base-26 odometer over b'a'..=b'z' with
// carry from the rightmost byte.
fn increment(p: &mut [u8]) {
    for c in p.iter_mut().rev() {
        if *c == b'z' {
            *c = b'a';
        } else {
            *c += 1;
            break;
        }
    }
}

// Check the three corporate policy rules.
fn valid(p: &[u8]) -> bool {
    // Rule 2: no i, o, or l.
    if p.iter().any(|&c| c == b'i' || c == b'o' || c == b'l') {
        return false;
    }

    // Rule 1: an increasing straight of at least three letters.
    let straight = p.windows(3).any(|w| w[1] == w[0] + 1 && w[2] == w[0] + 2);
    if !straight {
        return false;
    }

    // Rule 3: at least two different, non-overlapping pairs.
    let mut pairs = std::collections::HashSet::new();
    let mut i = 0;
    while i + 1 < p.len() {
        if p[i] == p[i + 1] {
            pairs.insert(p[i]);
            i += 2;
        } else {
            i += 1;
        }
    }
    pairs.len() >= 2
}

// Return the next valid password strictly after `s`.
fn next_password(s: &str) -> String {
    let mut p = s.trim().as_bytes().to_vec();
    increment(&mut p);
    while !valid(&p) {
        increment(&mut p);
    }
    String::from_utf8(p).unwrap()
}

pub fn part_one(input: &str) -> String {
    next_password(input)
}

pub fn part_two(input: &str) -> String {
    next_password(&next_password(input))
}
