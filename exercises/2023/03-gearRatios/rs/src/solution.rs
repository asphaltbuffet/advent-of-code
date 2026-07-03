// Solution for Advent of Code 2023 day 3.
//
// The schematic is a grid. A run of digits is a "part number" if any of the
// cells in its 8-neighborhood holds a symbol (anything but a dot or digit).
// Part one sums such numbers; part two finds every '*' adjacent to exactly two
// numbers and sums the products (gear ratios). We scan the grid as raw bytes,
// walking each digit run once and inspecting the frame of cells around it, so a
// single pass over the grid handles both the symbol test and gear collection.

use std::collections::HashMap;

/// A number found in the grid, with its row and inclusive column span.
struct Number {
    value: u32,
    row: usize,
    c0: usize,
    c1: usize,
}

/// Parse the grid and every digit run within it.
fn scan(input: &str) -> (Vec<&[u8]>, Vec<Number>) {
    let grid: Vec<&[u8]> = input.lines().map(str::as_bytes).collect();
    let mut numbers = Vec::new();

    for (row, line) in grid.iter().enumerate() {
        let mut c = 0;
        while c < line.len() {
            if !line[c].is_ascii_digit() {
                c += 1;
                continue;
            }
            let c0 = c;
            let mut value = 0u32;
            while c < line.len() && line[c].is_ascii_digit() {
                value = value * 10 + (line[c] - b'0') as u32;
                c += 1;
            }
            numbers.push(Number { value, row, c0, c1: c - 1 });
        }
    }

    (grid, numbers)
}

/// Visit each cell in the frame surrounding a number, passing its byte and
/// position to `f`.
fn for_border<F: FnMut(u8, usize, usize)>(grid: &[&[u8]], n: &Number, mut f: F) {
    let r0 = n.row.saturating_sub(1);
    let r1 = (n.row + 1).min(grid.len() - 1);
    let x0 = n.c0.saturating_sub(1);
    for r in r0..=r1 {
        let row = grid[r];
        let x1 = (n.c1 + 1).min(row.len() - 1);
        for (x, &b) in row.iter().enumerate().take(x1 + 1).skip(x0) {
            f(b, x, r);
        }
    }
}

fn is_symbol(b: u8) -> bool {
    b != b'.' && !b.is_ascii_digit()
}

pub fn part_one(input: &str) -> String {
    let (grid, numbers) = scan(input);
    let mut sum = 0u32;

    for n in &numbers {
        let mut touches = false;
        for_border(&grid, n, |b, _, _| {
            if is_symbol(b) {
                touches = true;
            }
        });
        if touches {
            sum += n.value;
        }
    }

    sum.to_string()
}

pub fn part_two(input: &str) -> String {
    let (grid, numbers) = scan(input);
    let mut gears: HashMap<(usize, usize), Vec<u32>> = HashMap::new();

    for n in &numbers {
        for_border(&grid, n, |b, x, y| {
            if b == b'*' {
                gears.entry((x, y)).or_default().push(n.value);
            }
        });
    }

    gears
        .values()
        .filter(|parts| parts.len() == 2)
        .map(|parts| parts[0] * parts[1])
        .sum::<u32>()
        .to_string()
}
