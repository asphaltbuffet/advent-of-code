use std::collections::HashSet;

/// Parse the input into signed frequency changes, tolerating blank lines and
/// surrounding whitespace (the runner may deliver input with or without a
/// trailing newline).
fn changes(input: &str) -> impl Iterator<Item = i64> + '_ {
    input.split_whitespace().map(|f| f.parse().unwrap())
}

pub fn part_one(input: &str) -> String {
    changes(input).sum::<i64>().to_string()
}

pub fn part_two(input: &str) -> String {
    let mut seen = HashSet::from([0]);
    let mut freq = 0;

    changes(input)
        .collect::<Vec<_>>()
        .into_iter()
        .cycle()
        .find_map(|c| {
            freq += c;
            (!seen.insert(freq)).then_some(freq)
        })
        .unwrap()
        .to_string()
}
