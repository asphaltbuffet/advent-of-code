// Solution for Advent of Code 2015 day 10.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

// Apply one round of look-and-say: each run of identical digits becomes
// <count><digit>. Operates on/returns ASCII digit bytes.
fn look_and_say(s: &[u8]) -> Vec<u8> {
    let mut out = Vec::with_capacity(s.len() * 2);
    let mut i = 0;

    while i < s.len() {
        let mut j = i;
        while j < s.len() && s[j] == s[i] {
            j += 1;
        }
        // Run length is small (1..=3 in look-and-say), but format generally.
        out.extend_from_slice((j - i).to_string().as_bytes());
        out.push(s[i]);
        i = j;
    }

    out
}

fn iterate(input: &str, n: usize) -> usize {
    let mut s = input.trim().as_bytes().to_vec();
    for _ in 0..n {
        s = look_and_say(&s);
    }
    s.len()
}

pub fn part_one(input: &str) -> String {
    iterate(input, 40).to_string()
}

pub fn part_two(input: &str) -> String {
    iterate(input, 50).to_string()
}
