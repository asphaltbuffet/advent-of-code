// Solution for Advent of Code 2023 day 13.
//
// Each pattern is summarized by its mirror line. Encoding every row and every
// column as an integer bitmask ('#' = 1) turns a reflection test into XOR: a
// mirror after line i is valid when the bit differences across all folded pairs
// sum to the target — 0 for a perfect mirror (part one) or exactly 1 for the
// single smudge (part two). `count_ones` on the XOR gives each pair's difference.

/// Bitmask rows for a block, and the same for its transpose (columns).
fn rows_and_cols(block: &str) -> (Vec<u32>, Vec<u32>) {
    let grid: Vec<&[u8]> = block.lines().map(str::as_bytes).collect();
    let h = grid.len();
    let w = grid[0].len();

    let mut rows = vec![0u32; h];
    let mut cols = vec![0u32; w];
    for (r, line) in grid.iter().enumerate() {
        for (c, &ch) in line.iter().enumerate() {
            if ch == b'#' {
                rows[r] |= 1 << c;
                cols[c] |= 1 << r;
            }
        }
    }
    (rows, cols)
}

/// The mirror index (lines above the reflection) whose folded pairs differ by
/// exactly `smudges` bits, or 0 if none.
fn reflection(lines: &[u32], smudges: u32) -> u32 {
    for i in 1..lines.len() {
        let diff: u32 = lines[..i]
            .iter()
            .rev()
            .zip(&lines[i..])
            .map(|(a, b)| (a ^ b).count_ones())
            .sum();
        if diff == smudges {
            return i as u32;
        }
    }
    0
}

fn summary(input: &str, smudges: u32) -> u32 {
    input
        .split("\n\n")
        .map(|block| {
            let (rows, cols) = rows_and_cols(block);
            100 * reflection(&rows, smudges) + reflection(&cols, smudges)
        })
        .sum()
}

pub fn part_one(input: &str) -> String {
    summary(input, 0).to_string()
}

pub fn part_two(input: &str) -> String {
    summary(input, 1).to_string()
}
