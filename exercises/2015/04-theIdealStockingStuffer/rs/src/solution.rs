// Solution for Advent of Code 2015 day 4.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

use md5::{Digest, Md5};

// Find the lowest positive integer that, appended to `key`, yields an MD5 hash
// whose hex digest begins with `zeros` leading '0' characters.
fn mine(key: &str, zeros: usize) -> u64 {
    // Each leading hex zero is a zero nibble; check the raw digest bytes so we
    // avoid formatting a hex string for every candidate. `full` is the number of
    // whole zero bytes; `half` is true when there's an extra zero nibble.
    let full = zeros / 2;
    let half = zeros % 2 == 1;

    let key = key.as_bytes();

    for i in 1u64.. {
        let mut hasher = Md5::new();
        hasher.update(key);
        hasher.update(i.to_string().as_bytes());
        let digest = hasher.finalize();

        if digest[..full].iter().all(|&b| b == 0) && (!half || digest[full] < 0x10) {
            return i;
        }
    }

    unreachable!()
}

pub fn part_one(input: &str) -> String {
    mine(input.trim(), 5).to_string()
}

pub fn part_two(input: &str) -> String {
    mine(input.trim(), 6).to_string()
}
