// Solution for Advent of Code 2018 day 20.
//
// The input is a route regex describing a facility. Walking it with a position
// stack — push on '(', reset to the group start on '|', pop on ')' — and recording
// the running door-distance (keeping the minimum on revisits) yields the shortest
// distance to every room directly. Part one is the furthest room; part two counts
// rooms at least 1000 doors away.

use std::collections::HashMap;

fn walk(input: &str) -> HashMap<(i32, i32), u32> {
    let route = input.trim().as_bytes();
    let mut dist: HashMap<(i32, i32), u32> = HashMap::new();
    dist.insert((0, 0), 0);

    let mut stack: Vec<(i32, i32)> = Vec::new();
    let mut pos = (0, 0);

    for &b in route {
        match b {
            b'(' => stack.push(pos),
            b'|' => pos = *stack.last().unwrap(),
            b')' => pos = stack.pop().unwrap(),
            b'N' | b'S' | b'E' | b'W' => {
                let d = dist[&pos] + 1;
                match b {
                    b'N' => pos.1 -= 1,
                    b'S' => pos.1 += 1,
                    b'E' => pos.0 += 1,
                    b'W' => pos.0 -= 1,
                    _ => unreachable!(),
                }
                let e = dist.entry(pos).or_insert(u32::MAX);
                if d < *e {
                    *e = d;
                }
            }
            _ => {} // '^' and '$'
        }
    }
    dist
}

pub fn part_one(input: &str) -> String {
    walk(input).values().max().unwrap().to_string()
}

pub fn part_two(input: &str) -> String {
    walk(input).values().filter(|&&d| d >= 1000).count().to_string()
}
