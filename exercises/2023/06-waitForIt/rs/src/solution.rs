// Solution for Advent of Code 2023 day 6.
//
// Holding the button for h of a t-millisecond race travels h*(t-h); you beat the
// record d when h*(t-h) > d, i.e. h^2 - t*h + d < 0. The winning h lie strictly
// between the roots of that quadratic, so each race is answered in O(1) by the
// quadratic formula plus a small integer nudge to stay exact. Part one multiplies
// the counts across races; part two removes the spaces to form one large race.

/// Count the integer button-hold times that beat the record for one race.
fn wins(time: u64, record: u64) -> u64 {
    let disc = time * time - 4 * record;
    if disc == 0 {
        return 0;
    }
    let root = (disc as f64).sqrt();
    let t = time as f64;

    // Real roots bound the winning interval; refine to exact integer bounds.
    let mut lo = ((t - root) / 2.0).floor() as u64;
    while lo * (time - lo) <= record {
        lo += 1;
    }
    let mut hi = ((t + root) / 2.0).ceil() as u64;
    while hi * (time - hi) <= record {
        hi -= 1;
    }
    hi - lo + 1
}

/// Read the two label lines, applying `clean` to each value list before parsing.
fn parse(input: &str, join_digits: bool) -> Vec<(u64, u64)> {
    let nums: Vec<Vec<u64>> = input
        .lines()
        .map(|line| {
            let values = line.split_once(':').unwrap().1;
            if join_digits {
                vec![values.replace(' ', "").parse().unwrap()]
            } else {
                values.split_whitespace().map(|n| n.parse().unwrap()).collect()
            }
        })
        .collect();

    nums[0].iter().cloned().zip(nums[1].iter().cloned()).collect()
}

pub fn part_one(input: &str) -> String {
    parse(input, false)
        .into_iter()
        .map(|(t, d)| wins(t, d))
        .product::<u64>()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    let (t, d) = parse(input, true)[0];
    wins(t, d).to_string()
}
