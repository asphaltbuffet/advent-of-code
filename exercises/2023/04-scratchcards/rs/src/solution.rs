// Solution for Advent of Code 2023 day 4.
//
// Each card lists winning numbers and the numbers you have; the match count is
// how many held numbers are winners. Part one scores a card with matches as
// 2^(matches-1); part two plays out the cascade where a card with m matches
// wins one copy of each of the next m cards. We compute the match counts once
// with a HashSet intersection, then a single forward pass propagates copies.

use std::collections::HashSet;

fn match_counts(input: &str) -> Vec<usize> {
    input
        .lines()
        .map(|line| {
            let (winning, have) = line
                .split_once(':')
                .and_then(|(_, rest)| rest.split_once('|'))
                .expect("card format");
            let winners: HashSet<&str> = winning.split_whitespace().collect();
            have.split_whitespace()
                .filter(|n| winners.contains(n))
                .count()
        })
        .collect()
}

pub fn part_one(input: &str) -> String {
    match_counts(input)
        .into_iter()
        .filter(|&m| m > 0)
        .map(|m| 1u32 << (m - 1))
        .sum::<u32>()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    let matches = match_counts(input);
    let mut copies = vec![1u64; matches.len()];

    for i in 0..matches.len() {
        let won = copies[i];
        for j in (i + 1)..(i + 1 + matches[i]).min(matches.len()) {
            copies[j] += won;
        }
    }

    copies.iter().sum::<u64>().to_string()
}
