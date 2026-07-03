// Solution for Advent of Code 2023 day 22.
//
// The bricks fall and settle. Processing them in ascending z, each brick drops
// onto the highest occupied cell beneath its footprint; whatever tops out at that
// level supports it. From the resulting support graph — who rests on whom — part
// one counts bricks that aren't the sole support of anything, and part two sums,
// per brick, how many others would fall if it were removed (a brick falls once
// all its supporters have fallen).

use std::collections::{HashSet, VecDeque};

struct Brick {
    x1: i32,
    y1: i32,
    z1: i32,
    x2: i32,
    y2: i32,
    z2: i32,
}

/// Returns, per brick index, the supporters it rests on and the bricks it holds.
fn settle(input: &str) -> (Vec<HashSet<usize>>, Vec<HashSet<usize>>) {
    let mut bricks: Vec<Brick> = input
        .lines()
        .map(|line| {
            let n: Vec<i32> = line
                .split(['~', ','])
                .map(|v| v.parse().unwrap())
                .collect();
            Brick { x1: n[0], y1: n[1], z1: n[2], x2: n[3], y2: n[4], z2: n[5] }
        })
        .collect();
    bricks.sort_by_key(|b| b.z1.min(b.z2));

    // heights[(x,y)] = (top z, brick index) of the topmost settled cell.
    let mut heights: std::collections::HashMap<(i32, i32), (i32, usize)> =
        std::collections::HashMap::new();
    let mut supports = vec![HashSet::new(); bricks.len()];
    let mut held = vec![HashSet::new(); bricks.len()];

    for (i, b) in bricks.iter().enumerate() {
        let cells: Vec<(i32, i32)> = (b.x1.min(b.x2)..=b.x1.max(b.x2))
            .flat_map(|x| (b.y1.min(b.y2)..=b.y1.max(b.y2)).map(move |y| (x, y)))
            .collect();

        let rest_z = cells
            .iter()
            .filter_map(|c| heights.get(c).map(|&(z, _)| z))
            .max()
            .unwrap_or(0)
            + 1;

        for c in &cells {
            if let Some(&(z, j)) = heights.get(c) {
                if z == rest_z - 1 {
                    supports[i].insert(j);
                    held[j].insert(i);
                }
            }
        }

        let top = rest_z + (b.z1.max(b.z2) - b.z1.min(b.z2));
        for c in cells {
            heights.insert(c, (top, i));
        }
    }

    (supports, held)
}

pub fn part_one(input: &str) -> String {
    let (supports, held) = settle(input);
    (0..supports.len())
        .filter(|&i| held[i].iter().all(|&j| supports[j].len() > 1))
        .count()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    let (supports, held) = settle(input);
    let mut total = 0usize;

    for start in 0..supports.len() {
        let mut fallen = HashSet::from([start]);
        let mut queue = VecDeque::from([start]);
        while let Some(cur) = queue.pop_front() {
            for &j in &held[cur] {
                if !fallen.contains(&j) && supports[j].is_subset(&fallen) {
                    fallen.insert(j);
                    queue.push_back(j);
                }
            }
        }
        total += fallen.len() - 1;
    }

    total.to_string()
}
