// Solution for Advent of Code 2025 day 3.
//
// Each line is a run of digits; we want the largest number formed by an in-order
// subsequence of length k (k=2 for part one, k=12 for part two), summed over
// lines. This is the classic "delete n-k digits to maximize" greedy: a monotonic
// stack that pops smaller trailing digits while deletions remain.

/// Largest length-`k` in-order subsequence of `digits`, as an integer.
fn largest_subsequence(digits: &[u8], k: usize) -> u64 {
    let mut drop = digits.len() - k;
    let mut stack: Vec<u8> = Vec::with_capacity(digits.len());
    for &c in digits {
        while drop > 0 && stack.last().is_some_and(|&top| top < c) {
            stack.pop();
            drop -= 1;
        }
        stack.push(c);
    }
    stack[..k]
        .iter()
        .fold(0u64, |acc, &d| acc * 10 + (d - b'0') as u64)
}

fn sum_over_lines(input: &str, k: usize) -> u64 {
    input
        .lines()
        .map(|line| largest_subsequence(line.as_bytes(), k))
        .sum()
}

pub fn part_one(input: &str) -> String {
    sum_over_lines(input, 2).to_string()
}

pub fn part_two(input: &str) -> String {
    sum_over_lines(input, 12).to_string()
}
