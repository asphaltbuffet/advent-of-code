// Solution for Advent of Code 2025 day 4.
//
// The floor is a grid of paper rolls (`@`). A roll is *accessible* when fewer than
// four of its eight neighbors are occupied. Part one counts the accessible rolls;
// part two peels them away in synchronous waves until none remain. Removing a roll
// only lowers its neighbors' counts, so the peeling is monotone and its total is
// independent of order. We store rolls as a coordinate set — the only query is
// occupancy.

use std::collections::HashSet;

const ADJ: [(i32, i32); 8] = [
    (-1, -1), (-1, 0), (-1, 1),
    (0, -1), (0, 1),
    (1, -1), (1, 0), (1, 1),
];

fn rolls(input: &str) -> HashSet<(i32, i32)> {
    input
        .lines()
        .enumerate()
        .flat_map(|(y, line)| {
            line.bytes().enumerate().filter_map(move |(x, b)| {
                (b == b'@').then_some((x as i32, y as i32))
            })
        })
        .collect()
}

fn accessible(rolls: &HashSet<(i32, i32)>, &(x, y): &(i32, i32)) -> bool {
    ADJ.iter()
        .filter(|(dx, dy)| rolls.contains(&(x + dx, y + dy)))
        .count()
        < 4
}

pub fn part_one(input: &str) -> String {
    let rolls = rolls(input);
    rolls.iter().filter(|p| accessible(&rolls, p)).count().to_string()
}

pub fn part_two(input: &str) -> String {
    let mut rolls = rolls(input);
    let mut removed = 0;
    loop {
        let wave: Vec<(i32, i32)> = rolls
            .iter()
            .copied()
            .filter(|p| accessible(&rolls, p))
            .collect();
        if wave.is_empty() {
            return removed.to_string();
        }
        removed += wave.len();
        for p in wave {
            rolls.remove(&p);
        }
    }
}
