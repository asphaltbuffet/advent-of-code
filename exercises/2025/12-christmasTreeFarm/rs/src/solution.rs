// Solution for Advent of Code 2025 day 12.
//
// Each region must hold a given multiset of present shapes (rotated/reflected as
// needed) on a grid without overlap, though cells may be left empty. We count the
// regions where a packing exists. A quick necessary check — the presents' filled
// cells must not exceed the area — rejects the impossible regions outright; the
// rest are settled by backtracking that either covers the first empty cell with a
// present or marks it permanently unused. Day 12 is the finale, so part two is 0.

use std::collections::HashSet;

type Cell = (i32, i32);

fn normalize(cells: &HashSet<Cell>) -> Vec<Cell> {
    let r0 = cells.iter().map(|&(r, _)| r).min().unwrap();
    let c0 = cells.iter().map(|&(_, c)| c).min().unwrap();
    let mut v: Vec<Cell> = cells.iter().map(|&(r, c)| (r - r0, c - c0)).collect();
    v.sort_unstable();
    v
}

/// All eight rotations/reflections of a shape, de-duplicated.
fn orientations(shape: &[Cell]) -> Vec<Vec<Cell>> {
    let mut seen: HashSet<Vec<Cell>> = HashSet::new();
    let mut cur: HashSet<Cell> = shape.iter().copied().collect();
    for _ in 0..4 {
        cur = cur.iter().map(|&(r, c)| (c, -r)).collect();
        seen.insert(normalize(&cur));
        let reflected: HashSet<Cell> = cur.iter().map(|&(r, c)| (r, -c)).collect();
        seen.insert(normalize(&reflected));
    }
    seen.into_iter().collect()
}

fn parse(input: &str) -> (Vec<Vec<Cell>>, Vec<(usize, usize, Vec<usize>)>) {
    let mut blocks: Vec<&str> = input.split("\n\n").collect();
    let region_block = blocks.pop().unwrap();

    let shapes: Vec<Vec<Cell>> = blocks
        .iter()
        .map(|block| {
            block
                .lines()
                .skip(1)
                .enumerate()
                .flat_map(|(r, row)| {
                    row.bytes().enumerate().filter_map(move |(c, b)| {
                        (b == b'#').then_some((r as i32, c as i32))
                    })
                })
                .collect()
        })
        .collect();

    let regions = region_block
        .lines()
        .map(|line| {
            let mut it = line.split_whitespace();
            let head = it.next().unwrap().trim_end_matches(':');
            let (w, l) = head.split_once('x').unwrap();
            let counts = it.map(|c| c.parse().unwrap()).collect();
            (w.parse().unwrap(), l.parse().unwrap(), counts)
        })
        .collect();

    (shapes, regions)
}

struct Packer {
    w: i32,
    l: i32,
    grid: Vec<u8>,
    need: Vec<usize>,
    orients: Vec<Vec<Vec<Cell>>>,
}

impl Packer {
    fn first_empty(&self, start: usize) -> Option<usize> {
        (start..self.grid.len()).find(|&p| self.grid[p] == 0)
    }

    fn backtrack(&mut self, start: usize, remaining: usize) -> bool {
        if remaining == 0 {
            return true;
        }
        let fe = match self.first_empty(start) {
            Some(p) => p,
            None => return false,
        };
        let (r0, c0) = ((fe / self.w as usize) as i32, (fe % self.w as usize) as i32);

        // Cover the first empty cell with some present.
        for si in 0..self.need.len() {
            if self.need[si] == 0 {
                continue;
            }
            for oi in 0..self.orients[si].len() {
                for ai in 0..self.orients[si][oi].len() {
                    let (ar, ac) = self.orients[si][oi][ai];
                    let mut placed = [0usize; 9];
                    let mut n = 0;
                    let mut ok = true;
                    for &(r, c) in &self.orients[si][oi] {
                        let (x, y) = (r0 + (r - ar), c0 + (c - ac));
                        if x < 0 || x >= self.l || y < 0 || y >= self.w {
                            ok = false;
                            break;
                        }
                        let p = (x * self.w + y) as usize;
                        if self.grid[p] != 0 || p < fe {
                            ok = false;
                            break;
                        }
                        placed[n] = p;
                        n += 1;
                    }
                    if ok {
                        for &p in &placed[..n] {
                            self.grid[p] = 1;
                        }
                        self.need[si] -= 1;
                        if self.backtrack(fe + 1, remaining - 1) {
                            return true;
                        }
                        self.need[si] += 1;
                        for &p in &placed[..n] {
                            self.grid[p] = 0;
                        }
                    }
                }
            }
        }

        // Or leave the first empty cell permanently unused.
        self.grid[fe] = 2;
        if self.backtrack(fe + 1, remaining) {
            return true;
        }
        self.grid[fe] = 0;
        false
    }
}

fn can_pack(w: usize, l: usize, counts: &[usize], shapes: &[Vec<Cell>]) -> bool {
    let used: usize = counts.iter().enumerate().map(|(i, &n)| shapes[i].len() * n).sum();
    if used > w * l {
        return false;
    }
    let remaining: usize = counts.iter().sum();
    let orients = shapes.iter().map(|s| orientations(s)).collect();
    let mut packer = Packer {
        w: w as i32,
        l: l as i32,
        grid: vec![0u8; w * l],
        need: counts.to_vec(),
        orients,
    };
    packer.backtrack(0, remaining)
}

pub fn part_one(input: &str) -> String {
    let (shapes, regions) = parse(input);
    regions
        .iter()
        .filter(|(w, l, counts)| can_pack(*w, *l, counts, &shapes))
        .count()
        .to_string()
}

pub fn part_two(_input: &str) -> String {
    // Day 12 is the finale; the last star is granted for the other puzzles.
    "0".to_string()
}
