// Solution for Advent of Code 2025 day 11.
//
// The input is a directed graph of devices. Part one counts distinct paths from
// `you` to `out`. Part two counts paths from `svr` to `out` that pass through two
// ordered waypoints (`fft` and `dac`, ordered by reachability). Both are the same
// ordered-waypoint path count: a path is valid when it visits each required
// waypoint in sequence and ends at the last. We memoize on (node, next waypoint
// index) — the remaining waypoints are always a suffix, so one index suffices.

use std::collections::HashMap;

/// Intern device names to indices; return (adjacency, name→index map).
fn parse(input: &str) -> (Vec<Vec<usize>>, HashMap<String, usize>) {
    let mut ids: HashMap<String, usize> = HashMap::new();
    let mut adj: Vec<Vec<usize>> = Vec::new();

    fn intern(name: &str, ids: &mut HashMap<String, usize>, adj: &mut Vec<Vec<usize>>) -> usize {
        if let Some(&i) = ids.get(name) {
            return i;
        }
        let i = adj.len();
        adj.push(Vec::new());
        ids.insert(name.to_string(), i);
        i
    }

    for line in input.lines() {
        let mut parts = line.split_whitespace();
        let head = parts.next().unwrap().trim_end_matches(':');
        let h = intern(head, &mut ids, &mut adj);
        for name in parts {
            let n = intern(name, &mut ids, &mut adj);
            adj[h].push(n);
        }
    }
    (adj, ids)
}

/// Count paths from `cur` visiting `waypoints[wp..]` in order, ending at the last.
fn trace(
    adj: &[Vec<usize>],
    waypoints: &[usize],
    memo: &mut HashMap<(usize, usize), u64>,
    cur: usize,
    wp: usize,
) -> u64 {
    if cur == waypoints[wp] {
        if wp + 1 == waypoints.len() {
            return 1; // reached the final destination
        }
        // Cleared this waypoint; require the remainder from here on.
        let key = (cur, wp + 1);
        if let Some(&v) = memo.get(&key) {
            return v;
        }
        let v = adj[cur]
            .iter()
            .map(|&nxt| trace(adj, waypoints, memo, nxt, wp + 1))
            .sum();
        memo.insert(key, v);
        return v;
    }
    // Reaching a later waypoint before the next required one invalidates the path.
    if waypoints[wp + 1..].contains(&cur) {
        return 0;
    }

    let key = (cur, wp);
    if let Some(&v) = memo.get(&key) {
        return v;
    }
    let v = adj[cur]
        .iter()
        .map(|&nxt| trace(adj, waypoints, memo, nxt, wp))
        .sum();
    memo.insert(key, v);
    v
}

fn reaches(adj: &[Vec<usize>], src: usize, dst: usize) -> bool {
    let mut stack = vec![src];
    let mut seen = vec![false; adj.len()];
    while let Some(node) = stack.pop() {
        if node == dst {
            return true;
        }
        if seen[node] {
            continue;
        }
        seen[node] = true;
        stack.extend(&adj[node]);
    }
    false
}

pub fn part_one(input: &str) -> String {
    let (adj, ids) = parse(input);
    let waypoints = [ids["you"], ids["out"]];
    let mut memo = HashMap::new();
    trace(&adj, &waypoints, &mut memo, ids["you"], 0).to_string()
}

pub fn part_two(input: &str) -> String {
    let (adj, ids) = parse(input);
    let (dac, fft, out, svr) = (ids["dac"], ids["fft"], ids["out"], ids["svr"]);
    let waypoints = if reaches(&adj, dac, fft) {
        [svr, dac, fft, out]
    } else {
        [svr, fft, dac, out]
    };
    let mut memo = HashMap::new();
    trace(&adj, &waypoints, &mut memo, svr, 0).to_string()
}
