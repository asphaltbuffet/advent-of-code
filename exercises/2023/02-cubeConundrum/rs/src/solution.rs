// Solution for Advent of Code 2023 day 2.
//
// Each game lists several draws of colored cubes. We reduce a game to the
// maximum count needed of each color — the fewest cubes that make every draw
// possible. Part one keeps games within (12 red, 13 green, 14 blue) and sums
// their ids; part two sums the product of the three maxima. The parse leans on
// stdlib split iterators rather than a regex crate, keeping the build dependency
// free.

/// The minimum (red, green, blue) cubes required for one game line.
fn maxima(line: &str) -> (u32, u32, u32) {
    // Skip the "Game N:" prefix, then read every "<count> <color>" token.
    let draws = line.split(':').nth(1).unwrap_or("");
    let mut max = (0, 0, 0);

    for cube in draws.split([',', ';']) {
        let mut it = cube.split_whitespace();
        let count: u32 = match it.next().and_then(|n| n.parse().ok()) {
            Some(n) => n,
            None => continue,
        };
        match it.next() {
            Some("red") => max.0 = max.0.max(count),
            Some("green") => max.1 = max.1.max(count),
            Some("blue") => max.2 = max.2.max(count),
            _ => {}
        }
    }

    max
}

pub fn part_one(input: &str) -> String {
    input
        .lines()
        .enumerate()
        .filter(|(_, line)| {
            let (r, g, b) = maxima(line);
            r <= 12 && g <= 13 && b <= 14
        })
        .map(|(i, _)| i as u32 + 1)
        .sum::<u32>()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    input
        .lines()
        .map(|line| {
            let (r, g, b) = maxima(line);
            r * g * b
        })
        .sum::<u32>()
        .to_string()
}
