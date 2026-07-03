// Solution for Advent of Code 2023 day 17.
//
// The crucible seeks the minimum heat-loss path across the city grid, but may
// only travel between `lo` and `hi` cells in a straight line before turning.
// This is Dijkstra over states (row, col, axis-just-travelled): because a move
// must turn onto the perpendicular axis, expanding every lo..hi straight step in
// one transition captures the constraint without tracking a run length. Part one
// uses (1, 3); part two the ultra crucible's (4, 10).

use std::cmp::Reverse;
use std::collections::BinaryHeap;

const HORIZONTAL: usize = 0;
const VERTICAL: usize = 1;

fn min_heat(grid: &[Vec<u32>], lo: i32, hi: i32) -> u32 {
    let (h, w) = (grid.len() as i32, grid[0].len() as i32);
    let goal = (h - 1, w - 1);

    // best[axis][r][c]: cheapest arrival travelling along `axis`.
    let mut best = vec![vec![vec![u32::MAX; w as usize]; h as usize]; 2];
    let mut heap: BinaryHeap<Reverse<(u32, i32, i32, usize)>> = BinaryHeap::new();
    heap.push(Reverse((0, 0, 0, HORIZONTAL)));
    heap.push(Reverse((0, 0, 0, VERTICAL)));

    while let Some(Reverse((cost, r, c, axis))) = heap.pop() {
        if (r, c) == goal {
            return cost;
        }
        if cost > best[axis][r as usize][c as usize] {
            continue;
        }

        // Turn onto the other axis, stepping lo..hi cells.
        let (dr, dc) = if axis == VERTICAL { (0, 1) } else { (1, 0) };
        for sign in [1, -1] {
            let (dr, dc) = (dr * sign, dc * sign);
            let (mut nr, mut nc, mut acc) = (r, c, cost);
            for step in 1..=hi {
                nr += dr;
                nc += dc;
                if nr < 0 || nc < 0 || nr >= h || nc >= w {
                    break;
                }
                acc += grid[nr as usize][nc as usize];
                if step < lo {
                    continue;
                }
                let naxis = axis ^ 1;
                if acc < best[naxis][nr as usize][nc as usize] {
                    best[naxis][nr as usize][nc as usize] = acc;
                    heap.push(Reverse((acc, nr, nc, naxis)));
                }
            }
        }
    }

    u32::MAX
}

fn parse(input: &str) -> Vec<Vec<u32>> {
    input
        .lines()
        .map(|l| l.bytes().map(|b| (b - b'0') as u32).collect())
        .collect()
}

pub fn part_one(input: &str) -> String {
    min_heat(&parse(input), 1, 3).to_string()
}

pub fn part_two(input: &str) -> String {
    min_heat(&parse(input), 4, 10).to_string()
}
