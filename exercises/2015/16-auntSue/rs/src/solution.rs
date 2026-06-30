// Solution for Advent of Code 2015 day 16.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

// MFCSAM ticker-tape reading from the problem statement (not the input).
fn target(compound: &str) -> i32 {
    match compound {
        "children" => 3,
        "cats" => 7,
        "samoyeds" => 2,
        "pomeranians" => 3,
        "akitas" => 0,
        "vizslas" => 0,
        "goldfish" => 5,
        "trees" => 3,
        "cars" => 2,
        "perfumes" => 1,
        _ => unreachable!(),
    }
}

// parse yields one Sue per line as a list of (compound, count) observations.
fn parse(input: &str) -> Vec<Vec<(&str, i32)>> {
    input
        .trim()
        .lines()
        .map(|line| {
            // Drop the "Sue N: " prefix, then split the remaining "k: v" pairs.
            let rest = line.split_once(": ").unwrap().1;
            rest.split(", ")
                .map(|pair| {
                    let (k, v) = pair.split_once(": ").unwrap();
                    (k, v.parse().unwrap())
                })
                .collect()
        })
        .collect()
}

// find returns the 1-based number of the first Sue all of whose observations
// satisfy `ok(compound, observed)`.
fn find(sues: &[Vec<(&str, i32)>], ok: impl Fn(&str, i32) -> bool) -> usize {
    sues.iter()
        .position(|obs| obs.iter().all(|&(k, got)| ok(k, got)))
        .map(|i| i + 1)
        .unwrap_or(0)
}

pub fn part_one(input: &str) -> String {
    let sues = parse(input);
    find(&sues, |k, got| got == target(k)).to_string()
}

pub fn part_two(input: &str) -> String {
    // cats/trees are lower bounds, pomeranians/goldfish upper bounds; rest exact.
    let sues = parse(input);
    find(&sues, |k, got| match k {
        "cats" | "trees" => got > target(k),
        "pomeranians" | "goldfish" => got < target(k),
        _ => got == target(k),
    })
    .to_string()
}
