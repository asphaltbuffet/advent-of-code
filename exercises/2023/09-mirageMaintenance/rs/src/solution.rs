// Solution for Advent of Code 2023 day 9.
//
// Each history is extended by repeatedly taking differences until the row is all
// zeros, then summing the last value of every row. Part two extends backwards,
// which is the same computation on the reversed sequence. We recurse on the
// difference row via `windows(2)`.

fn extrapolate(seq: &[i64]) -> i64 {
    if seq.iter().all(|&n| n == 0) {
        return 0;
    }
    let diffs: Vec<i64> = seq.windows(2).map(|w| w[1] - w[0]).collect();
    seq.last().unwrap() + extrapolate(&diffs)
}

fn solve(input: &str, reverse: bool) -> i64 {
    input
        .lines()
        .map(|line| {
            let mut seq: Vec<i64> =
                line.split_whitespace().map(|n| n.parse().unwrap()).collect();
            if reverse {
                seq.reverse();
            }
            extrapolate(&seq)
        })
        .sum()
}

pub fn part_one(input: &str) -> String {
    solve(input, false).to_string()
}

pub fn part_two(input: &str) -> String {
    solve(input, true).to_string()
}
