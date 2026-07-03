use std::collections::{HashMap, HashSet};

/// Parse the initial planted pots (as a set of indices) and the set of five-pot
/// neighborhoods that grow a plant, each encoded as a 5-bit mask.
fn parse(input: &str) -> (HashSet<i64>, HashSet<u8>) {
    let mut lines = input.trim().lines();

    let initial = lines
        .next()
        .unwrap_or("")
        .trim_start_matches("initial state:")
        .trim();
    let pots: HashSet<i64> = initial
        .bytes()
        .enumerate()
        .filter(|&(_, c)| c == b'#')
        .map(|(i, _)| i as i64)
        .collect();

    let mut rules: HashSet<u8> = HashSet::new();
    for line in lines {
        let line = line.trim();
        if line.is_empty() {
            continue;
        }
        let (lhs, rhs) = line.split_once(" => ").expect("bad rule");
        if rhs == "#" {
            rules.insert(mask(lhs.bytes().map(|c| c == b'#')));
        }
    }

    (pots, rules)
}

/// Pack a five-pot window (leftmost first) into a bit mask.
fn mask(window: impl Iterator<Item = bool>) -> u8 {
    window.fold(0u8, |acc, live| (acc << 1) | u8::from(live))
}

fn bounds(pots: &HashSet<i64>) -> (i64, i64) {
    (*pots.iter().min().unwrap(), *pots.iter().max().unwrap())
}

/// Advance one generation. Only pots within two of a live pot can change.
fn step(pots: &HashSet<i64>, rules: &HashSet<u8>) -> HashSet<i64> {
    let (lo, hi) = bounds(pots);
    ((lo - 2)..=(hi + 2))
        .filter(|&i| rules.contains(&mask((i - 2..=i + 2).map(|j| pots.contains(&j)))))
        .collect()
}

fn sum_indices(pots: &HashSet<i64>) -> i64 {
    pots.iter().sum()
}

/// The planted pattern normalized to start at its leftmost pot, paired with that
/// offset, so two generations with the same shape differ only by a shift.
fn shape(pots: &HashSet<i64>) -> (String, i64) {
    let (lo, hi) = bounds(pots);
    let s = (lo..=hi)
        .map(|i| if pots.contains(&i) { '#' } else { '.' })
        .collect();
    (s, lo)
}

pub fn part_one(input: &str) -> String {
    let (mut pots, rules) = parse(input);
    for _ in 0..20 {
        pots = step(&pots, &rules);
    }
    sum_indices(&pots).to_string()
}

pub fn part_two(input: &str) -> String {
    let (mut pots, rules) = parse(input);
    const TARGET: i64 = 50_000_000_000;

    // Once a shape repeats, the planted count is fixed and the pattern just
    // drifts sideways by a constant offset per generation, so the index sum
    // grows linearly and we extrapolate to TARGET.
    let mut seen: HashMap<String, (i64, i64)> = HashMap::new();
    for gen in 0..TARGET {
        let (sh, lo) = shape(&pots);
        if let Some(&(prev_gen, prev_lo)) = seen.get(&sh) {
            let period = gen - prev_gen;
            let drift = lo - prev_lo;
            let remaining = TARGET - gen;
            let count = pots.len() as i64;
            return (sum_indices(&pots) + remaining * count * drift / period).to_string();
        }
        seen.insert(sh, (gen, lo));
        pots = step(&pots, &rules);
    }
    sum_indices(&pots).to_string()
}
