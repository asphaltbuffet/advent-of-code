// Solution for Advent of Code 2023 day 12.
//
// Each row is a spring pattern ('.', '#', '?') plus the run lengths of damaged
// springs. We count the arrangements consistent with both by a memoized DP over
// (position in the pattern, index into the group list): at each position we may
// treat a '?' as operational, or try to place the next group of '#' starting
// here if it fits and is bounded by a non-'#'. Part two unfolds each row fivefold
// ('?'-joined pattern, repeated groups), which the memoization keeps tractable.

use std::collections::HashMap;

/// Count arrangements from byte offset `s` of `springs` using groups[g..].
fn count(
    springs: &[u8],
    s: usize,
    groups: &[usize],
    g: usize,
    memo: &mut HashMap<(usize, usize), u64>,
) -> u64 {
    if g == groups.len() {
        // No groups left: valid iff no forced '#' remain.
        return if springs[s..].contains(&b'#') { 0 } else { 1 };
    }
    // Not enough room for the remaining groups plus their separators.
    let need: usize = groups[g..].iter().sum::<usize>() + (groups.len() - g - 1);
    if springs.len() < s + need {
        return 0;
    }
    if let Some(&v) = memo.get(&(s, g)) {
        return v;
    }

    let mut total = 0;
    // Option A: current char operational, skip it.
    if springs[s] != b'#' {
        total += count(springs, s + 1, groups, g, memo);
    }
    // Option B: place a run of groups[g] damaged springs here.
    let n = groups[g];
    let fits = !springs[s..s + n].contains(&b'.');
    let bounded = s + n == springs.len() || springs[s + n] != b'#';
    if fits && bounded {
        // Skip the run and the separator that must follow it.
        let next = (s + n + 1).min(springs.len());
        total += count(springs, next, groups, g + 1, memo);
    }

    memo.insert((s, g), total);
    total
}

fn arrangements(input: &str, unfold: usize) -> u64 {
    input
        .lines()
        .map(|line| {
            let (pattern, group_str) = line.split_once(' ').unwrap();
            let groups: Vec<usize> =
                group_str.split(',').map(|n| n.parse().unwrap()).collect();

            let springs = vec![pattern; unfold].join("?").into_bytes();
            let groups: Vec<usize> = groups.repeat(unfold);

            let mut memo = HashMap::new();
            count(&springs, 0, &groups, 0, &mut memo)
        })
        .sum()
}

pub fn part_one(input: &str) -> String {
    arrangements(input, 1).to_string()
}

pub fn part_two(input: &str) -> String {
    arrangements(input, 5).to_string()
}
