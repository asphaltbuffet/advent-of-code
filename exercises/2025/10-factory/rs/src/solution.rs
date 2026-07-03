// Solution for Advent of Code 2025 day 10.
//
// Each machine has a target light pattern, a set of buttons (each toggling some
// light bits), and a required joltage per output. Part one is the fewest buttons
// whose XOR reproduces the light pattern. Part two reaches the joltage vector in
// base 2: the buttons pressed at a level fix the low bit of every component
// (parities must match), then the remainder is halved and solved recursively at
// double cost. Each button-subset's component-sum vector is precomputed so the
// recursion only subtracts.

use std::collections::HashMap;

struct Machine {
    target: u32,
    buttons: Vec<u32>,
    joltage: Vec<i64>,
}

fn parse(line: &str) -> Machine {
    let tokens: Vec<&str> = line.split_whitespace().collect();
    let lights = tokens[0].trim_matches(|c| c == '[' || c == ']');
    let target = lights
        .bytes()
        .enumerate()
        .filter(|&(_, c)| c == b'#')
        .fold(0u32, |acc, (i, _)| acc | (1 << i));

    let joltage = tokens[tokens.len() - 1]
        .trim_matches(|c| c == '{' || c == '}')
        .split(',')
        .map(|n| n.parse().unwrap())
        .collect();

    let buttons = tokens[1..tokens.len() - 1]
        .iter()
        .map(|tok| {
            tok.trim_matches(|c| c == '(' || c == ')')
                .split(',')
                .filter_map(|o| o.parse::<u32>().ok())
                .filter(|&idx| (idx as usize) < lights.len())
                .fold(0u32, |acc, idx| acc | (1 << idx))
        })
        .collect();

    Machine { target, buttons, joltage }
}

/// Fewest buttons whose XOR equals `target`, searching by increasing count.
fn min_button_presses(target: u32, buttons: &[u32]) -> u32 {
    let n = buttons.len();
    for k in 1..=n {
        // Iterate k-subsets in lexicographic order via an index combination.
        let mut idx: Vec<usize> = (0..k).collect();
        loop {
            if idx.iter().fold(0u32, |x, &j| x ^ buttons[j]) == target {
                return k as u32;
            }
            let mut i = k;
            while i > 0 && idx[i - 1] == n - k + (i - 1) {
                i -= 1;
            }
            if i == 0 {
                break;
            }
            idx[i - 1] += 1;
            for j in i..k {
                idx[j] = idx[j - 1] + 1;
            }
        }
    }
    0
}

fn min_joltage_presses(machine: &Machine) -> i64 {
    let rlen = machine.joltage.len();
    let limit = 1usize << machine.buttons.len();

    // button_vecs[b][i] = whether button b feeds output i.
    let button_vecs: Vec<Vec<i64>> = machine
        .buttons
        .iter()
        .map(|&b| (0..rlen).map(|i| ((b >> i) & 1) as i64).collect())
        .collect();

    // mask_sum[mask] = component sums for that button subset; popcount alongside.
    let mut mask_sum: Vec<Vec<i64>> = vec![vec![0; rlen]; limit];
    let mut popcount = vec![0i64; limit];
    for mask in 1..limit {
        let low = mask & mask.wrapping_neg();
        let b = low.trailing_zeros() as usize;
        let prev = mask ^ low;
        popcount[mask] = popcount[prev] + 1;
        for i in 0..rlen {
            mask_sum[mask][i] = mask_sum[prev][i] + button_vecs[b][i];
        }
    }

    let mut memo: HashMap<Vec<i64>, i64> = HashMap::new();
    solve(&machine.joltage, limit, &mask_sum, &popcount, rlen, &mut memo)
}

fn solve(
    target: &[i64],
    limit: usize,
    mask_sum: &[Vec<i64>],
    popcount: &[i64],
    rlen: usize,
    memo: &mut HashMap<Vec<i64>, i64>,
) -> i64 {
    if target.iter().all(|&t| t == 0) {
        return 0;
    }
    if let Some(&v) = memo.get(target) {
        return v;
    }

    let mut best = -1i64;
    for mask in 0..limit {
        let sums = &mask_sum[mask];
        let mut next = Vec::with_capacity(rlen);
        let mut ok = true;
        for i in 0..rlen {
            let rem = target[i] - sums[i];
            if rem < 0 || rem & 1 == 1 {
                ok = false;
                break;
            }
            next.push(rem >> 1);
        }
        if ok {
            let sub = solve(&next, limit, mask_sum, popcount, rlen, memo);
            if sub != -1 {
                let total = popcount[mask] + 2 * sub;
                if best == -1 || total < best {
                    best = total;
                }
            }
        }
    }

    memo.insert(target.to_vec(), best);
    best
}

pub fn part_one(input: &str) -> String {
    input
        .lines()
        .map(|line| {
            let m = parse(line);
            min_button_presses(m.target, &m.buttons)
        })
        .sum::<u32>()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    input
        .lines()
        .map(|line| min_joltage_presses(&parse(line)))
        .sum::<i64>()
        .to_string()
}
