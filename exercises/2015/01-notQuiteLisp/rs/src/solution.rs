// Solution for Advent of Code 2015 day 1.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

pub fn part_one(input: &str) -> String {
    let floor: i64 = input
        .chars()
        .map(|c| match c {
            '(' => 1,
            ')' => -1,
            _ => 0,
        })
        .sum();
    floor.to_string()
}

pub fn part_two(input: &str) -> String {
    let mut floor = 0i64;
    for (i, c) in input.chars().enumerate() {
        floor += match c {
            '(' => 1,
            ')' => -1,
            _ => 0,
        };
        if floor == -1 {
            return (i + 1).to_string();
        }
    }
    "never reaches basement".to_string()
}
