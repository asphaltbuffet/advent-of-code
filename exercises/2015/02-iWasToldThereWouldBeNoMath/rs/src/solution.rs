// Solution for Advent of Code 2015 day 2.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

fn parse(input: &str) -> impl Iterator<Item = (u32, u32, u32)> + '_ {
    input.lines().map(|line| {
        let mut parts = line.split('x').map(|n| n.parse::<u32>().unwrap());
        (
            parts.next().unwrap(),
            parts.next().unwrap(),
            parts.next().unwrap(),
        )
    })
}

pub fn part_one(input: &str) -> String {
    parse(input)
        .map(|(l, w, h)| {
            let (a, b, c) = (l * w, w * h, h * l);
            2 * (a + b + c) + a.min(b).min(c)
        })
        .sum::<u32>()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    parse(input)
        .map(|(l, w, h)| {
            let max = l.max(w).max(h);
            let perimeter = 2 * (l + w + h - max);
            let volume = l * w * h;
            perimeter + volume
        })
        .sum::<u32>()
        .to_string()
}
