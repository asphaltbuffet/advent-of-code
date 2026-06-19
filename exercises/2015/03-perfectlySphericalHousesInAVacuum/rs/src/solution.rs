// Solution for Advent of Code 2015 day 3.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

use std::collections::HashSet;

fn step(pos: &mut (i32, i32), c: char) {
    match c {
        '^' => pos.1 -= 1,
        'v' => pos.1 += 1,
        '>' => pos.0 += 1,
        '<' => pos.0 -= 1,
        _ => {}
    }
}

pub fn part_one(input: &str) -> String {
    let mut pos = (0, 0);
    let mut visited = HashSet::new();
    visited.insert(pos);

    for c in input.chars() {
        step(&mut pos, c);
        visited.insert(pos);
    }

    visited.len().to_string()
}

pub fn part_two(input: &str) -> String {
    let mut pos = [(0, 0), (0, 0)];
    let mut visited = HashSet::new();
    visited.insert((0, 0));

    for (i, c) in input.chars().filter(|c| "^v<>".contains(*c)).enumerate() {
        let p = &mut pos[i & 1];
        step(p, c);
        visited.insert(*p);
    }

    visited.len().to_string()
}
