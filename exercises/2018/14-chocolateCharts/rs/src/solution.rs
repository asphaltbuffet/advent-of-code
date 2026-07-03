// Solution for Advent of Code 2018 day 14.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

/// A step appends the digits of the two current recipes' sum (one or two
/// digits), then advances each elf forward by one plus its current recipe.
#[inline]
fn step(scores: &mut Vec<u8>, a: &mut usize, b: &mut usize) {
    let sum = scores[*a] + scores[*b];
    if sum >= 10 {
        scores.push(sum / 10);
    }
    scores.push(sum % 10);

    *a = (*a + 1 + scores[*a] as usize) % scores.len();
    *b = (*b + 1 + scores[*b] as usize) % scores.len();
}

pub fn part_one(input: &str) -> String {
    let n: usize = input.trim().parse().expect("input must be a number");

    let mut scores: Vec<u8> = Vec::with_capacity(n + 12);
    scores.extend_from_slice(&[3, 7]);
    let (mut a, mut b) = (0usize, 1usize);

    while scores.len() < n + 10 {
        step(&mut scores, &mut a, &mut b);
    }

    scores[n..n + 10].iter().map(|d| (b'0' + d) as char).collect()
}

pub fn part_two(input: &str) -> String {
    let input = input.trim();
    let target: Vec<u8> = input.bytes().map(|c| c - b'0').collect();
    let tlen = target.len();

    let mut scores: Vec<u8> = Vec::with_capacity(1 << 24);
    scores.extend_from_slice(&[3, 7]);
    let (mut a, mut b) = (0usize, 1usize);

    // A step appends one or two recipes, so check for the target after each
    // append by comparing the tail. `checked` tracks how far we have already
    // compared so no suffix is missed across a two-digit step.
    let mut checked = 0usize;
    loop {
        step(&mut scores, &mut a, &mut b);
        while checked + tlen <= scores.len() {
            if scores[checked..checked + tlen] == target[..] {
                return checked.to_string();
            }
            checked += 1;
        }
    }
}
