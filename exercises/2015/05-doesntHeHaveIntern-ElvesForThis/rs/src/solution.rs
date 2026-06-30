// Solution for Advent of Code 2015 day 5.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

fn is_nice_one(s: &str) -> bool {
    let b = s.as_bytes();

    let vowels = s.chars().filter(|c| "aeiou".contains(*c)).count() >= 3;
    let doubled = b.windows(2).any(|w| w[0] == w[1]);
    let bad = ["ab", "cd", "pq", "xy"].iter().any(|p| s.contains(p));

    vowels && doubled && !bad
}

fn is_nice_two(s: &str) -> bool {
    let b = s.as_bytes();

    // A non-overlapping pair of two letters appearing twice.
    let pair = (0..b.len().saturating_sub(1)).any(|i| {
        let p = &b[i..i + 2];
        b[i + 2..].windows(2).any(|w| w == p)
    });
    // A letter repeating with exactly one letter between them.
    let repeat = b.windows(3).any(|w| w[0] == w[2]);

    pair && repeat
}

pub fn part_one(input: &str) -> String {
    input.lines().filter(|l| is_nice_one(l)).count().to_string()
}

pub fn part_two(input: &str) -> String {
    input.lines().filter(|l| is_nice_two(l)).count().to_string()
}
