use std::collections::HashMap;

pub fn part_one(input: &str) -> String {
    let (twos, threes) = input.split_whitespace().fold((0, 0), |(twos, threes), id| {
        let mut counts: HashMap<u8, usize> = HashMap::new();
        for &b in id.as_bytes() {
            *counts.entry(b).or_insert(0) += 1;
        }
        let any = |n| counts.values().any(|&c| c == n);
        (twos + any(2) as i64, threes + any(3) as i64)
    });

    (twos * threes).to_string()
}

pub fn part_two(input: &str) -> String {
    let ids: Vec<&[u8]> = input.split_whitespace().map(str::as_bytes).collect();

    for (i, &a) in ids.iter().enumerate() {
        for &b in &ids[i + 1..] {
            if a.len() != b.len() {
                continue;
            }
            let mut diffs = a.iter().zip(b).enumerate().filter(|(_, (x, y))| x != y);
            if let Some((at, _)) = diffs.next() {
                if diffs.next().is_none() {
                    return a
                        .iter()
                        .enumerate()
                        .filter(|&(k, _)| k != at)
                        .map(|(_, &c)| c as char)
                        .collect();
                }
            }
        }
    }

    String::new()
}
