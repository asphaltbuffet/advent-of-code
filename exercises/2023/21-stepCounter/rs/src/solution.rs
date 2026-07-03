// Solution for Advent of Code 2023 day 21.
//
// The elf reaches a plot in exactly N steps iff its BFS distance has the same
// parity as N and is <= N. Part one runs a bounded BFS to 64 steps. Part two's
// 26,501,365 steps on the infinitely tiled garden is intractable directly, but
// the input's clear center row/column and border make the reachable count grow
// as a quadratic in the number of tile widths. Sampling at target%size,
// +size, +2*size lets us fit that quadratic and evaluate it at the real n.

use std::collections::HashSet;

fn parse(input: &str) -> (HashSet<(i64, i64)>, (i64, i64), i64) {
    let grid: Vec<&[u8]> = input.lines().map(str::as_bytes).collect();
    let size = grid.len() as i64;
    let mut rocks = HashSet::new();
    let mut start = (0, 0);
    for (r, row) in grid.iter().enumerate() {
        for (c, &ch) in row.iter().enumerate() {
            match ch {
                b'#' => {
                    rocks.insert((r as i64, c as i64));
                }
                b'S' => start = (r as i64, c as i64),
                _ => {}
            }
        }
    }
    (rocks, start, size)
}

/// Plots reachable in exactly `steps` steps, wrapping into the tile on the
/// infinite grid.
fn reachable(rocks: &HashSet<(i64, i64)>, start: (i64, i64), size: i64, steps: i64, infinite: bool) -> u64 {
    let mut frontier: HashSet<(i64, i64)> = HashSet::new();
    frontier.insert(start);

    for _ in 0..steps {
        let mut next = HashSet::with_capacity(frontier.len() * 2);
        for &(r, c) in &frontier {
            for (dr, dc) in [(-1, 0), (1, 0), (0, -1), (0, 1)] {
                let (nr, nc) = (r + dr, c + dc);
                let blocked = if infinite {
                    rocks.contains(&(nr.rem_euclid(size), nc.rem_euclid(size)))
                } else {
                    nr < 0 || nc < 0 || nr >= size || nc >= size || rocks.contains(&(nr, nc))
                };
                if !blocked {
                    next.insert((nr, nc));
                }
            }
        }
        frontier = next;
    }

    frontier.len() as u64
}

pub fn part_one(input: &str) -> String {
    let (rocks, start, size) = parse(input);
    reachable(&rocks, start, size, 64, false).to_string()
}

pub fn part_two(input: &str) -> String {
    let (rocks, start, size) = parse(input);
    let target: i64 = 26_501_365;
    let offset = target % size;

    // Three samples one tile-width apart, then a quadratic fit.
    let y: Vec<i64> = (0..3)
        .map(|i| reachable(&rocks, start, size, offset + i * size, true) as i64)
        .collect();
    let a = (y[2] - 2 * y[1] + y[0]) / 2;
    let b = y[1] - y[0] - a;
    let c = y[0];
    let n = target / size;

    (a * n * n + b * n + c).to_string()
}
