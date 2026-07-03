// Solution for Advent of Code 2018 day 23.
//
// Part one counts the bots in range of the strongest. Part two finds the point in
// range of the most bots (closest to the origin on ties) with an octree search: a
// max-heap of cubes ordered by how many bots could reach them, subdividing the best
// cube into eight until a single point wins.

use std::cmp::Ordering;
use std::collections::BinaryHeap;

struct Bot {
    x: i64,
    y: i64,
    z: i64,
    r: i64,
}

fn parse(input: &str) -> Vec<Bot> {
    input
        .trim()
        .lines()
        .map(|line| {
            let n: Vec<i64> = line
                .split(|c: char| !(c.is_ascii_digit() || c == '-'))
                .filter(|s| !s.is_empty() && *s != "-")
                .map(|s| s.parse().unwrap())
                .collect();
            Bot { x: n[0], y: n[1], z: n[2], r: n[3] }
        })
        .collect()
}

pub fn part_one(input: &str) -> String {
    let bots = parse(input);
    let s = bots.iter().max_by_key(|b| b.r).unwrap();
    bots.iter()
        .filter(|b| (b.x - s.x).abs() + (b.y - s.y).abs() + (b.z - s.z).abs() <= s.r)
        .count()
        .to_string()
}

// How far v lies outside the interval [lo, lo+size) on one axis (0 if inside).
fn axis_dist(v: i64, lo: i64, size: i64) -> i64 {
    if v < lo {
        lo - v
    } else if v > lo + size - 1 {
        v - (lo + size - 1)
    } else {
        0
    }
}

struct Cube {
    x: i64,
    y: i64,
    z: i64,
    size: i64,
    in_range: usize,
    dist: i64,
}

impl Cube {
    fn new(bots: &[Bot], x: i64, y: i64, z: i64, size: i64) -> Cube {
        let in_range = bots
            .iter()
            .filter(|b| {
                axis_dist(b.x, x, size) + axis_dist(b.y, y, size) + axis_dist(b.z, z, size) <= b.r
            })
            .count();
        let dist = axis_dist(0, x, size) + axis_dist(0, y, size) + axis_dist(0, z, size);
        Cube { x, y, z, size, in_range, dist }
    }
}

// Best-first order for the max-heap: more bots in range wins; then closer to the
// origin; then smaller. Distance and size are minimize criteria, so their
// comparisons are flipped.
impl Ord for Cube {
    fn cmp(&self, other: &Self) -> Ordering {
        self.in_range
            .cmp(&other.in_range)
            .then_with(|| other.dist.cmp(&self.dist))
            .then_with(|| other.size.cmp(&self.size))
    }
}
impl PartialOrd for Cube {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}
impl Eq for Cube {}
impl PartialEq for Cube {
    fn eq(&self, other: &Self) -> bool {
        self.cmp(other) == Ordering::Equal
    }
}

pub fn part_two(input: &str) -> String {
    let bots = parse(input);

    let span = bots
        .iter()
        .map(|b| b.x.abs().max(b.y.abs()).max(b.z.abs()))
        .max()
        .unwrap();
    let mut size = 1;
    while size < span {
        size *= 2;
    }

    let mut heap = BinaryHeap::new();
    heap.push(Cube::new(&bots, -size, -size, -size, 2 * size));

    while let Some(c) = heap.pop() {
        if c.size == 1 {
            return c.dist.to_string();
        }
        let half = c.size / 2;
        for (dx, dy, dz) in [
            (0, 0, 0),
            (half, 0, 0),
            (0, half, 0),
            (0, 0, half),
            (half, half, 0),
            (half, 0, half),
            (0, half, half),
            (half, half, half),
        ] {
            heap.push(Cube::new(&bots, c.x + dx, c.y + dy, c.z + dz, half));
        }
    }
    unreachable!()
}
