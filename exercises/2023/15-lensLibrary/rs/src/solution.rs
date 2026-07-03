// Solution for Advent of Code 2023 day 15.
//
// The HASH of a string folds each byte as ((acc + b) * 17) mod 256. Part one
// sums the HASH of every comma-separated step. Part two runs the HASHMAP: 256
// boxes, each an ordered list of (label, focal length). A '-' removes a label; a
// '=' updates it in place or appends it. Focusing power sums
// (box+1) * slot * focal over every lens.

fn hash(s: &str) -> usize {
    s.bytes().fold(0usize, |acc, b| (acc + b as usize) * 17 % 256)
}

pub fn part_one(input: &str) -> String {
    input
        .trim()
        .split(',')
        .map(hash)
        .sum::<usize>()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    // Order matters within a box, so each is a Vec of (label, focal).
    let mut boxes: Vec<Vec<(&str, u8)>> = vec![Vec::new(); 256];

    for step in input.trim().split(',') {
        if let Some(label) = step.strip_suffix('-') {
            let b = &mut boxes[hash(label)];
            if let Some(pos) = b.iter().position(|&(l, _)| l == label) {
                b.remove(pos);
            }
        } else {
            let (label, focal) = step.split_once('=').unwrap();
            let focal: u8 = focal.parse().unwrap();
            let b = &mut boxes[hash(label)];
            match b.iter_mut().find(|(l, _)| *l == label) {
                Some(slot) => slot.1 = focal, // update in place, keeping order
                None => b.push((label, focal)),
            }
        }
    }

    boxes
        .iter()
        .enumerate()
        .flat_map(|(bi, b)| {
            b.iter()
                .enumerate()
                .map(move |(si, &(_, focal))| (bi + 1) * (si + 1) * focal as usize)
        })
        .sum::<usize>()
        .to_string()
}
