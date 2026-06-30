// Solution for Advent of Code 2015 day 6.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

// Parse a line into (op, x1, y1, x2, y2) where op is 0=off, 1=on, 2=toggle.
fn parse(line: &str) -> (u8, usize, usize, usize, usize) {
    let (op, rest) = if let Some(r) = line.strip_prefix("turn on ") {
        (1u8, r)
    } else if let Some(r) = line.strip_prefix("turn off ") {
        (0u8, r)
    } else {
        (2u8, line.strip_prefix("toggle ").unwrap())
    };

    let (a, b) = rest.split_once(" through ").unwrap();
    let (x1, y1) = a.split_once(',').unwrap();
    let (x2, y2) = b.split_once(',').unwrap();

    (
        op,
        x1.parse().unwrap(),
        y1.parse().unwrap(),
        x2.parse().unwrap(),
        y2.parse().unwrap(),
    )
}

pub fn part_one(input: &str) -> String {
    let mut grid = vec![false; 1000 * 1000];

    for line in input.lines() {
        let (op, x1, y1, x2, y2) = parse(line);
        for x in x1..=x2 {
            for cell in &mut grid[x * 1000 + y1..=x * 1000 + y2] {
                *cell = match op {
                    1 => true,
                    0 => false,
                    _ => !*cell,
                };
            }
        }
    }

    grid.iter().filter(|&&b| b).count().to_string()
}

pub fn part_two(input: &str) -> String {
    let mut grid = vec![0i32; 1000 * 1000];

    for line in input.lines() {
        let (op, x1, y1, x2, y2) = parse(line);
        for x in x1..=x2 {
            for cell in &mut grid[x * 1000 + y1..=x * 1000 + y2] {
                *cell = match op {
                    1 => *cell + 1,
                    0 => (*cell - 1).max(0),
                    _ => *cell + 2,
                };
            }
        }
    }

    grid.iter().sum::<i32>().to_string()
}
