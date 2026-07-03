// Solution for Advent of Code 2018 day 18.
//
// A cellular automaton on the acre map. Each minute every acre updates from its
// eight neighbors. Part one runs 10 minutes; part two runs one billion by finding
// the cycle the state settles into and fast-forwarding.

use std::collections::HashMap;

const OPEN: u8 = b'.';
const TREES: u8 = b'|';
const LUMBER: u8 = b'#';

struct Grid {
    cells: Vec<u8>,
    w: usize,
    h: usize,
}

impl Grid {
    fn parse(input: &str) -> Grid {
        let lines: Vec<&str> = input.trim().lines().collect();
        let h = lines.len();
        let w = lines[0].len();
        let mut cells = Vec::with_capacity(w * h);
        for line in lines {
            cells.extend_from_slice(line.as_bytes());
        }
        Grid { cells, w, h }
    }

    // Count adjacent trees and lumberyards around (x, y).
    fn neighbors(&self, x: usize, y: usize) -> (u32, u32) {
        let (mut trees, mut lumber) = (0, 0);
        for dy in -1i32..=1 {
            for dx in -1i32..=1 {
                if dx == 0 && dy == 0 {
                    continue;
                }
                let (nx, ny) = (x as i32 + dx, y as i32 + dy);
                if nx < 0 || ny < 0 || nx as usize >= self.w || ny as usize >= self.h {
                    continue;
                }
                match self.cells[ny as usize * self.w + nx as usize] {
                    TREES => trees += 1,
                    LUMBER => lumber += 1,
                    _ => {}
                }
            }
        }
        (trees, lumber)
    }

    fn step(&self) -> Grid {
        let mut cells = vec![0u8; self.w * self.h];
        for y in 0..self.h {
            for x in 0..self.w {
                let (trees, lumber) = self.neighbors(x, y);
                let cur = self.cells[y * self.w + x];
                cells[y * self.w + x] = match cur {
                    OPEN if trees >= 3 => TREES,
                    TREES if lumber >= 3 => LUMBER,
                    LUMBER if !(lumber >= 1 && trees >= 1) => OPEN,
                    other => other,
                };
            }
        }
        Grid { cells, w: self.w, h: self.h }
    }

    fn resource_value(&self) -> usize {
        let trees = self.cells.iter().filter(|&&c| c == TREES).count();
        let lumber = self.cells.iter().filter(|&&c| c == LUMBER).count();
        trees * lumber
    }
}

pub fn part_one(input: &str) -> String {
    let mut grid = Grid::parse(input);
    for _ in 0..10 {
        grid = grid.step();
    }
    grid.resource_value().to_string()
}

pub fn part_two(input: &str) -> String {
    const TARGET: usize = 1_000_000_000;
    let mut grid = Grid::parse(input);

    // Record the minute each state was first seen; once one repeats, jump ahead
    // by the remaining minutes modulo the cycle length.
    let mut seen: HashMap<Vec<u8>, usize> = HashMap::new();
    let mut minute = 0;
    while minute < TARGET {
        if let Some(&first) = seen.get(&grid.cells) {
            let period = minute - first;
            for _ in 0..(TARGET - minute) % period {
                grid = grid.step();
            }
            return grid.resource_value().to_string();
        }
        seen.insert(grid.cells.clone(), minute);
        grid = grid.step();
        minute += 1;
    }
    grid.resource_value().to_string()
}
