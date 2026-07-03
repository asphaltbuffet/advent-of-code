// Solution for Advent of Code 2023 day 18.
//
// The dig plan traces a closed trench; we want the number of cubic metres it can
// hold, which is the polygon's interior plus its boundary. Walking the perimeter
// accumulates the shoelace area and the boundary length in one pass, then Pick's
// theorem (interior = A - b/2 + 1) gives the total as A + b/2 + 1. Part one reads
// the plain direction/distance; part two decodes them from the hex color.

/// (dr, dc) for a direction letter.
fn delta(dir: u8) -> (i64, i64) {
    match dir {
        b'R' => (0, 1),
        b'L' => (0, -1),
        b'U' => (-1, 0),
        b'D' => (1, 0),
        _ => unreachable!(),
    }
}

fn lagoon(steps: impl Iterator<Item = (u8, i64)>) -> i64 {
    let (mut r, mut c) = (0i64, 0i64);
    let mut area = 0i64;
    let mut perimeter = 0i64;

    for (dir, dist) in steps {
        let (dr, dc) = delta(dir);
        let (nr, nc) = (r + dr * dist, c + dc * dist);
        area += r * nc - nr * c; // shoelace cross term
        perimeter += dist;
        r = nr;
        c = nc;
    }

    area.abs() / 2 + perimeter / 2 + 1
}

pub fn part_one(input: &str) -> String {
    lagoon(input.lines().map(|line| {
        let mut it = line.split_whitespace();
        let dir = it.next().unwrap().as_bytes()[0];
        let dist = it.next().unwrap().parse().unwrap();
        (dir, dist)
    }))
    .to_string()
}

pub fn part_two(input: &str) -> String {
    lagoon(input.lines().map(|line| {
        // "... (#70c710)": first five hex digits are distance, last is direction.
        let code = line.split_whitespace().nth(2).unwrap();
        let hex = &code[2..code.len() - 1];
        let dist = i64::from_str_radix(&hex[..5], 16).unwrap();
        let dir = match hex.as_bytes()[5] {
            b'0' => b'R',
            b'1' => b'D',
            b'2' => b'L',
            b'3' => b'U',
            _ => unreachable!(),
        };
        (dir, dist)
    }))
    .to_string()
}
