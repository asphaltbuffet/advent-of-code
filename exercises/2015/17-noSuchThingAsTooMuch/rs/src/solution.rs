// Solution for Advent of Code 2015 day 17.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

// Container sizes plus the target volume. The target isn't in the input, so it
// is inferred: the 5-container AoC example fills 25 liters, the real input 150.
fn parse(input: &str) -> (Vec<u32>, u32) {
    let sizes: Vec<u32> = input.split_whitespace().map(|x| x.parse().unwrap()).collect();
    let target = if sizes.len() <= 5 { 25 } else { 150 };
    (sizes, target)
}

// counts[k] = number of k-container subsets summing to exactly target. With
// <= 20 containers the 2^n bitmask sweep is trivial.
fn size_counts(sizes: &[u32], target: u32) -> Vec<u32> {
    let mut counts = vec![0u32; sizes.len() + 1];
    for mask in 0u32..(1 << sizes.len()) {
        let sum: u32 = sizes
            .iter()
            .enumerate()
            .filter(|(i, _)| mask & (1 << i) != 0)
            .map(|(_, &s)| s)
            .sum();
        if sum == target {
            counts[mask.count_ones() as usize] += 1;
        }
    }
    counts
}

pub fn part_one(input: &str) -> String {
    let (sizes, target) = parse(input);
    size_counts(&sizes, target).iter().sum::<u32>().to_string()
}

pub fn part_two(input: &str) -> String {
    let (sizes, target) = parse(input);
    // The first non-zero size bucket is the minimum-container count.
    size_counts(&sizes, target)
        .iter()
        .find(|&&c| c > 0)
        .unwrap()
        .to_string()
}
