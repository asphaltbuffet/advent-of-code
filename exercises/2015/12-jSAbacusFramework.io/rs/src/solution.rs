// Solution for Advent of Code 2015 day 12.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

use serde_json::Value;

// Part 1: sum every signed integer appearing in the text, ignoring structure.
pub fn part_one(input: &str) -> String {
    let b = input.as_bytes();
    let mut sum: i64 = 0;
    let mut i = 0;

    while i < b.len() {
        if b[i].is_ascii_digit() {
            // A leading '-' immediately before the digits makes it negative.
            let neg = i > 0 && b[i - 1] == b'-';
            let mut n: i64 = 0;
            while i < b.len() && b[i].is_ascii_digit() {
                n = n * 10 + (b[i] - b'0') as i64;
                i += 1;
            }
            sum += if neg { -n } else { n };
        } else {
            i += 1;
        }
    }

    sum.to_string()
}

// Recursively sum numbers, skipping any object (and its descendants) that has a
// value of exactly "red". Arrays are never skipped.
fn sum_value(v: &Value) -> i64 {
    match v {
        Value::Number(n) => n.as_i64().unwrap_or(0),
        Value::Array(items) => items.iter().map(sum_value).sum(),
        Value::Object(map) => {
            if map.values().any(|e| e == "red") {
                0
            } else {
                map.values().map(sum_value).sum()
            }
        }
        _ => 0,
    }
}

// Part 2: sum numbers, ignoring any object containing the value "red".
pub fn part_two(input: &str) -> String {
    let doc: Value = serde_json::from_str(input.trim()).unwrap();
    sum_value(&doc).to_string()
}
