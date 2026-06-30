// Solution for Advent of Code 2015 day 24.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

fn parse(input: &str) -> Vec<i64> {
    input.split_whitespace().map(|x| x.parse().unwrap()).collect()
}

// min_qe finds the smallest first-group size summing to an equal share, then
// the lowest quantum entanglement (product) among groups of that size.
// Searching by increasing size yields the fewest-packages tier. QE needs i64.
fn min_qe(weights: &[i64], groups: i64) -> i64 {
    let target = weights.iter().sum::<i64>() / groups;

    for size in 1..=weights.len() {
        let mut best: Option<i64> = None;
        pick(weights, 0, target, size, 1, &mut best);
        if let Some(b) = best {
            return b;
        }
    }
    -1
}

// pick recursively chooses `count` more weights from index `start`, tracking
// the minimum product among selections that exactly hit `remaining`.
fn pick(weights: &[i64], start: usize, remaining: i64, count: usize, qe: i64, best: &mut Option<i64>) {
    if count == 0 {
        if remaining == 0 && best.map_or(true, |b| qe < b) {
            *best = Some(qe);
        }
        return;
    }
    let n = weights.len();
    for i in start..=n.saturating_sub(count) {
        if weights[i] <= remaining {
            pick(weights, i + 1, remaining - weights[i], count - 1, qe * weights[i], best);
        }
    }
}

pub fn part_one(input: &str) -> String {
    min_qe(&parse(input), 3).to_string()
}

pub fn part_two(input: &str) -> String {
    min_qe(&parse(input), 4).to_string()
}
