// Solution for Advent of Code 2019 day 1.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

pub fn part_one(input: &str) -> String {
    input
        .lines()
        .filter_map(|l| l.trim().parse::<i64>().ok())
        .map(|m| m / 3 - 2)
        .sum::<i64>()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    input
        .lines()
        .filter_map(|l| l.trim().parse::<i64>().ok())
        .map(fuel_cost)
        .sum::<i64>()
        .to_string()
}

fn fuel_cost(mass: i64) -> i64 {
    let f = mass / 3 - 2;
    if f <= 0 { 0 } else { f + fuel_cost(f) }
}
