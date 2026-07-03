// Solution for Advent of Code 2023 day 1.
//
// Each line's calibration value is its first digit times ten plus its last
// digit; part two also reads spelled-out words. Rather than searching the whole
// line once per token, we scan positions left to right and yield the digit
// starting at each one — an iterator whose `.next()` and `.last()` give the two
// digits we need in a single pass, with no intermediate allocation.

const WORDS: [&str; 9] = [
    "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
];

/// The digit value beginning at byte offset `i`, if any. ASCII digits are read
/// directly; when `words` is set, a spelled digit whose spelling starts at `i`
/// also counts (overlaps like "eightwo" resolve because we test every offset).
fn digit_at(line: &str, i: usize, words: bool) -> Option<u32> {
    let bytes = line.as_bytes();
    if bytes[i].is_ascii_digit() {
        return Some((bytes[i] - b'0') as u32);
    }
    if words {
        let tail = &line[i..];
        for (w, v) in WORDS.iter().zip(1..) {
            if tail.starts_with(w) {
                return Some(v);
            }
        }
    }
    None
}

fn calibration(line: &str, words: bool) -> u32 {
    let mut digits = (0..line.len()).filter_map(|i| digit_at(line, i, words));
    match digits.next() {
        Some(first) => {
            let last = digits.last().unwrap_or(first);
            first * 10 + last
        }
        None => 0,
    }
}

fn solve(input: &str, words: bool) -> u32 {
    input.lines().map(|line| calibration(line, words)).sum()
}

pub fn part_one(input: &str) -> String {
    solve(input, false).to_string()
}

pub fn part_two(input: &str) -> String {
    solve(input, true).to_string()
}
