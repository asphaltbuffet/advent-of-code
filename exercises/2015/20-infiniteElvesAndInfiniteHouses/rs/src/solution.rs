// Solution for Advent of Code 2015 day 20.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

fn target(input: &str) -> usize {
    input.trim().parse().unwrap()
}

pub fn part_one(input: &str) -> String {
    let t = target(input);
    let limit = t / 10 + 1;
    // House h receives 10 * (sum of divisors of h). Sieve each elf's gift onto
    // its multiples; bound target/10 suffices (elf h delivers 10*h to house h).
    let mut houses = vec![0usize; limit + 1];
    for n in 1..=limit {
        let mut h = n;
        while h <= limit {
            houses[h] += 10 * n;
            h += n;
        }
    }
    houses
        .iter()
        .position(|&g| g >= t)
        .unwrap()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    let t = target(input);
    let limit = t / 10 + 1;
    // Each elf gives 11*n but visits only its first 50 multiples.
    let mut houses = vec![0usize; limit + 1];
    for n in 1..=limit {
        let mut h = n;
        let mut c = 0;
        while h <= limit && c < 50 {
            houses[h] += 11 * n;
            h += n;
            c += 1;
        }
    }
    houses
        .iter()
        .position(|&g| g >= t)
        .unwrap()
        .to_string()
}
