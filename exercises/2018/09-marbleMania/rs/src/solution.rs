// Solution for Advent of Code 2018 day 9.

use std::collections::VecDeque;

/// Parse the player count and the last marble's worth from the puzzle input.
fn parse(input: &str) -> (usize, usize) {
    let nums: Vec<usize> = input
        .split(|c: char| !c.is_ascii_digit())
        .filter(|s| !s.is_empty())
        .map(|s| s.parse().unwrap())
        .collect();

    (nums[0], nums[1])
}

/// Play the full game and return the winning player's score.
///
/// A `VecDeque` models the ring: the current marble sits at the back and
/// `rotate_left`/`rotate_right` walk clockwise or counter-clockwise in O(1), so
/// both the normal insertion and the scoring removal stay cheap even across the
/// millions of marbles in part two.
fn high_score(players: usize, last: usize) -> usize {
    let mut scores = vec![0usize; players];
    let mut ring: VecDeque<usize> = VecDeque::with_capacity(last + 1);
    ring.push_back(0);

    for m in 1..=last {
        if m % 23 == 0 {
            // Scoring move: step seven marbles counter-clockwise and remove one.
            ring.rotate_right(7);
            scores[m % players] += m + ring.pop_back().unwrap();
            // The marble one clockwise becomes current.
            ring.rotate_left(1);
        } else {
            // Normal move: place the marble two positions clockwise.
            ring.rotate_left(1);
            ring.push_back(m);
        }
    }

    scores.into_iter().max().unwrap_or(0)
}

pub fn part_one(input: &str) -> String {
    let (players, last) = parse(input);
    high_score(players, last).to_string()
}

pub fn part_two(input: &str) -> String {
    let (players, last) = parse(input);
    high_score(players, last * 100).to_string()
}
