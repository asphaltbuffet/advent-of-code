// Solution for Advent of Code 2023 day 14.
//
// Round rocks ('O') roll until they hit a wall or a cube rock ('#'). Part one
// tilts north once and reports the load (each rock scores rows-below-it). Part
// two runs a billion N->W->S->E spin cycles; the grid is eventually periodic, so
// we hash each state, detect the repeat, jump forward by whole periods, and
// finish the remainder. Each direction rolls the grid in place with a per-lane
// "next free slot" cursor — no rotation or reallocation per tilt.

use std::collections::HashMap;

type Grid = Vec<Vec<u8>>;

fn parse(input: &str) -> Grid {
    input.lines().map(|l| l.as_bytes().to_vec()).collect()
}

fn tilt_north(g: &mut Grid) {
    let (h, w) = (g.len(), g[0].len());
    for c in 0..w {
        let mut free = 0; // next row a rolling rock can occupy
        for r in 0..h {
            match g[r][c] {
                b'#' => free = r + 1,
                b'O' => {
                    g[r][c] = b'.';
                    g[free][c] = b'O';
                    free += 1;
                }
                _ => {}
            }
        }
    }
}

fn tilt_south(g: &mut Grid) {
    let (h, w) = (g.len(), g[0].len());
    for c in 0..w {
        let mut free = h as isize - 1;
        for r in (0..h).rev() {
            match g[r][c] {
                b'#' => free = r as isize - 1,
                b'O' => {
                    g[r][c] = b'.';
                    g[free as usize][c] = b'O';
                    free -= 1;
                }
                _ => {}
            }
        }
    }
}

fn tilt_west(g: &mut Grid) {
    for row in g.iter_mut() {
        let w = row.len();
        let mut free = 0;
        for c in 0..w {
            match row[c] {
                b'#' => free = c + 1,
                b'O' => {
                    row[c] = b'.';
                    row[free] = b'O';
                    free += 1;
                }
                _ => {}
            }
        }
    }
}

fn tilt_east(g: &mut Grid) {
    for row in g.iter_mut() {
        let w = row.len();
        let mut free = w as isize - 1;
        for c in (0..w).rev() {
            match row[c] {
                b'#' => free = c as isize - 1,
                b'O' => {
                    row[c] = b'.';
                    row[free as usize] = b'O';
                    free -= 1;
                }
                _ => {}
            }
        }
    }
}

fn spin(g: &mut Grid) {
    tilt_north(g);
    tilt_west(g);
    tilt_south(g);
    tilt_east(g);
}

fn load(g: &Grid) -> usize {
    let h = g.len();
    g.iter()
        .enumerate()
        .map(|(r, row)| (h - r) * row.iter().filter(|&&c| c == b'O').count())
        .sum()
}

pub fn part_one(input: &str) -> String {
    let mut g = parse(input);
    tilt_north(&mut g);
    load(&g).to_string()
}

pub fn part_two(input: &str) -> String {
    let mut g = parse(input);
    let target = 1_000_000_000usize;
    let mut seen: HashMap<Grid, usize> = HashMap::new();

    let mut i = 0;
    while i < target {
        spin(&mut g);
        i += 1;
        if let Some(&prev) = seen.get(&g) {
            let period = i - prev;
            i += ((target - i) / period) * period;
            seen.clear(); // don't jump again
        } else {
            seen.insert(g.clone(), i);
        }
    }

    load(&g).to_string()
}
