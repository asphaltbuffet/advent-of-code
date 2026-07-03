// Solution for Advent of Code 2018 day 22.
//
// The cave's region types come from a geologic-index recurrence turned into an
// erosion level, mod 3. Part one sums the risk over the target rectangle. Part two
// is a Dijkstra over (x, y, tool): moving costs 1 minute, switching tools costs 7,
// and each region forbids the one tool equal to its type. Start and finish are
// equipped with the torch.

use std::cmp::Reverse;
use std::collections::BinaryHeap;

const ROCKY_MOD: u64 = 20183;
const TORCH: usize = 1;
// How far past the target to extend the search grid; a shortest path never detours
// much further than this (a 7-minute tool switch bounds any useful detour).
const MARGIN: usize = 50;

fn parse(input: &str) -> (u64, usize, usize) {
    let mut depth = 0;
    let (mut tx, mut ty) = (0, 0);
    for line in input.trim().lines() {
        if let Some(d) = line.strip_prefix("depth: ") {
            depth = d.trim().parse().unwrap();
        } else if let Some(t) = line.strip_prefix("target: ") {
            let (a, b) = t.trim().split_once(',').unwrap();
            tx = a.parse().unwrap();
            ty = b.parse().unwrap();
        }
    }
    (depth, tx, ty)
}

// Build the region-type grid (erosion % 3) over a (w+1) x (h+1) area.
fn regions(depth: u64, tx: usize, ty: usize, w: usize, h: usize) -> Vec<u8> {
    let stride = w + 1;
    let mut erosion = vec![0u64; stride * (h + 1)];
    for y in 0..=h {
        for x in 0..=w {
            let geo = if (x == 0 && y == 0) || (x == tx && y == ty) {
                0
            } else if y == 0 {
                x as u64 * 16807
            } else if x == 0 {
                y as u64 * 48271
            } else {
                erosion[y * stride + (x - 1)] * erosion[(y - 1) * stride + x]
            };
            erosion[y * stride + x] = (geo + depth) % ROCKY_MOD;
        }
    }
    erosion.iter().map(|&e| (e % 3) as u8).collect()
}

pub fn part_one(input: &str) -> String {
    let (depth, tx, ty) = parse(input);
    let region = regions(depth, tx, ty, tx, ty);
    let stride = tx + 1;
    let risk: u64 = (0..=ty)
        .flat_map(|y| (0..=tx).map(move |x| (x, y)))
        .map(|(x, y)| region[y * stride + x] as u64)
        .sum();
    risk.to_string()
}

pub fn part_two(input: &str) -> String {
    let (depth, tx, ty) = parse(input);
    let (w, h) = (tx + MARGIN, ty + MARGIN);
    let stride = w + 1;
    let region = regions(depth, tx, ty, w, h);

    // dist indexed by (y*stride + x)*3 + tool.
    let idx = |x: usize, y: usize, tool: usize| (y * stride + x) * 3 + tool;
    let mut dist = vec![u32::MAX; stride * (h + 1) * 3];
    dist[idx(0, 0, TORCH)] = 0;

    let mut pq = BinaryHeap::new();
    pq.push(Reverse((0u32, 0usize, 0usize, TORCH)));

    while let Some(Reverse((cost, x, y, tool))) = pq.pop() {
        if cost > dist[idx(x, y, tool)] {
            continue;
        }
        if (x, y, tool) == (tx, ty, TORCH) {
            return cost.to_string();
        }
        let mut relax = |nx: usize, ny: usize, nt: usize, nd: u32, pq: &mut BinaryHeap<_>, dist: &mut [u32]| {
            if nd < dist[idx(nx, ny, nt)] {
                dist[idx(nx, ny, nt)] = nd;
                pq.push(Reverse((nd, nx, ny, nt)));
            }
        };

        // Switch to the other tool allowed here (7 minutes).
        let here = region[y * stride + x] as usize;
        for t in 0..3 {
            if t != here && t != tool {
                relax(x, y, t, cost + 7, &mut pq, &mut dist);
            }
        }
        // Move to a neighbor keeping the current tool (1 minute).
        for (dx, dy) in [(1i32, 0i32), (-1, 0), (0, 1), (0, -1)] {
            let (nx, ny) = (x as i32 + dx, y as i32 + dy);
            if nx < 0 || ny < 0 || nx as usize > w || ny as usize > h {
                continue;
            }
            let (nx, ny) = (nx as usize, ny as usize);
            if region[ny * stride + nx] as usize != tool {
                relax(nx, ny, tool, cost + 1, &mut pq, &mut dist);
            }
        }
    }
    unreachable!("target unreachable")
}
