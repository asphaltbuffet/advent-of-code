// Solution for Advent of Code 2015 day 9.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

use std::collections::HashMap;

// Parse the edge list into (n cities, distance matrix indexed by city id).
fn parse(input: &str) -> (usize, Vec<Vec<Option<u32>>>) {
    let mut edges: Vec<(usize, usize, u32)> = Vec::new();

    let id = |name: &str, ids: &mut HashMap<String, usize>| -> usize {
        let next = ids.len();
        *ids.entry(name.to_string()).or_insert(next)
    };

    let mut ids: HashMap<String, usize> = HashMap::new();
    for line in input.lines() {
        let t: Vec<&str> = line.split_whitespace().collect();
        // <c1> to <c2> = <dist>
        let a = id(t[0], &mut ids);
        let b = id(t[2], &mut ids);
        let d: u32 = t[4].parse().unwrap();
        edges.push((a, b, d));
    }

    let n = ids.len();
    // None marks a missing edge: routes traversing it are not valid.
    let mut dist = vec![vec![None; n]; n];
    for (a, b, d) in edges {
        dist[a][b] = Some(d);
        dist[b][a] = Some(d);
    }

    (n, dist)
}

// Enumerate all permutations of [0..n), returning (min route, max route).
fn routes(n: usize, dist: &[Vec<Option<u32>>]) -> (u32, u32) {
    let mut perm: Vec<usize> = (0..n).collect();
    let mut min = u32::MAX;
    let mut max = 0;
    permute(&mut perm, 0, dist, &mut min, &mut max);
    (min, max)
}

fn permute(
    perm: &mut [usize],
    k: usize,
    dist: &[Vec<Option<u32>>],
    min: &mut u32,
    max: &mut u32,
) {
    if k == perm.len() {
        // Sum the legs; skip the route entirely if any edge is missing.
        let mut total = 0;
        for w in perm.windows(2) {
            match dist[w[0]][w[1]] {
                Some(d) => total += d,
                None => return,
            }
        }
        *min = (*min).min(total);
        *max = (*max).max(total);
        return;
    }
    for i in k..perm.len() {
        perm.swap(k, i);
        permute(perm, k + 1, dist, min, max);
        perm.swap(k, i);
    }
}

pub fn part_one(input: &str) -> String {
    let (n, dist) = parse(input);
    routes(n, &dist).0.to_string()
}

pub fn part_two(input: &str) -> String {
    let (n, dist) = parse(input);
    routes(n, &dist).1.to_string()
}
