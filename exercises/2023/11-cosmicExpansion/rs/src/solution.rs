// Solution for Advent of Code 2023 day 11.
//
// Empty rows and columns each expand to `factor` tracks. The total pairwise
// Manhattan distance separates into independent row and column sums, and along
// one axis the sorted-order identity sum_{i<j}(x_j - x_i) = sum_i x_i*(2i-n+1)
// turns each axis into a single linear pass — so the million-fold expansion in
// part two costs the same as part one.

/// Pairwise distance sum along one axis after expanding empty tracks.
fn axis_sum(coords: &[usize], size: usize, factor: u64) -> u64 {
    let occupied: Vec<bool> = {
        let mut v = vec![false; size];
        for &c in coords {
            v[c] = true;
        }
        v
    };

    // Cumulative expanded position of each track.
    let mut expanded = vec![0u64; size];
    let mut pos = 0u64;
    for x in 0..size {
        expanded[x] = pos;
        pos += if occupied[x] { 1 } else { factor };
    }

    let mut positions: Vec<u64> = coords.iter().map(|&c| expanded[c]).collect();
    positions.sort_unstable();

    // Accumulate in signed arithmetic: the per-term weight (2i - n + 1) is
    // negative for the lower half, so a u64 multiply would underflow.
    let n = positions.len() as i64;
    let sum: i64 = positions
        .iter()
        .enumerate()
        .map(|(i, &p)| p as i64 * (2 * i as i64 - n + 1))
        .sum();
    sum as u64
}

fn total(input: &str, factor: u64) -> u64 {
    let grid: Vec<&[u8]> = input.lines().map(str::as_bytes).collect();
    let height = grid.len();
    let width = grid[0].len();

    let mut rows = Vec::new();
    let mut cols = Vec::new();
    for (r, line) in grid.iter().enumerate() {
        for (c, &ch) in line.iter().enumerate() {
            if ch == b'#' {
                rows.push(r);
                cols.push(c);
            }
        }
    }

    axis_sum(&rows, height, factor) + axis_sum(&cols, width, factor)
}

pub fn part_one(input: &str) -> String {
    total(input, 2).to_string()
}

pub fn part_two(input: &str) -> String {
    total(input, 1_000_000).to_string()
}
