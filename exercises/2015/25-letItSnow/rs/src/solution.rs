// Solution for Advent of Code 2015 day 25.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

const FIRST: i64 = 20151125;
const MULT: i64 = 252533;
const MOD: i64 = 33554393;

// mod_pow computes base^exp mod m by fast exponentiation.
fn mod_pow(mut base: i64, mut exp: i64, m: i64) -> i64 {
    let mut result = 1i64;
    base %= m;
    while exp > 0 {
        if exp & 1 == 1 {
            result = result * base % m;
        }
        base = base * base % m;
        exp >>= 1;
    }
    result
}

pub fn part_one(input: &str) -> String {
    // Extract the two integers (row, col) from the prose.
    let nums: Vec<i64> = input
        .split(|c: char| !c.is_ascii_digit())
        .filter(|s| !s.is_empty())
        .map(|s| s.parse().unwrap())
        .collect();
    let (row, col) = (nums[0], nums[1]);

    // Cell (row, col) is the n-th code in diagonal fill order; the n-th code is
    // FIRST * MULT^(n-1) mod MOD.
    let diag = row + col - 2;
    let n = diag * (diag + 1) / 2 + col;
    (FIRST * mod_pow(MULT, n - 1, MOD) % MOD).to_string()
}

pub fn part_two(_input: &str) -> String {
    // Day 25 has no part 2 — the final star comes from finishing the rest.
    "Merry Christmas!".to_string()
}
