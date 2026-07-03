// Solution for Advent of Code 2025 day 1.
//
// The dial has 100 positions and starts at 50. Part one counts moves that end on
// zero; part two counts every sweep across zero. We reduce the pointer modulo 100
// each step (mirroring the physical dial) and count the multiples of 100 the move
// arc covers — a half-open convention so a landing is credited once, on arrival.

/// (direction byte, distance) for each instruction line.
fn moves(input: &str) -> impl Iterator<Item = (u8, i32)> + '_ {
    input
        .lines()
        .map(|line| (line.as_bytes()[0], line[1..].parse().unwrap()))
}

pub fn part_one(input: &str) -> String {
    let mut pos: i32 = 50;
    let landings = moves(input)
        .filter(|&(dir, dist)| {
            pos = (if dir == b'R' { pos + dist } else { pos - dist }).rem_euclid(100);
            pos == 0
        })
        .count();
    landings.to_string()
}

pub fn part_two(input: &str) -> String {
    // `div_euclid` gives true floor division so the leftward arc is counted
    // correctly even as the pointer dips below zero before reduction.
    let mut pos: i32 = 50;
    let mut clicks: i64 = 0;
    for (dir, dist) in moves(input) {
        if dir == b'R' {
            clicks += ((pos + dist) / 100) as i64;
            pos = (pos + dist).rem_euclid(100);
        } else {
            clicks += ((pos - 1).div_euclid(100) - (pos - dist - 1).div_euclid(100)) as i64;
            pos = (pos - dist).rem_euclid(100);
        }
    }
    clicks.to_string()
}
