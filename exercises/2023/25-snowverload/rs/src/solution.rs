// Solution for Advent of Code 2023 day 25.
//
// The wiring diagram splits into two clusters joined by exactly three wires. We
// find those wires by edge betweenness: a BFS from many sources marks each edge
// on its shortest-path tree, and summed over sources the three bridge wires carry
// by far the most traffic. Removing them and measuring one side's size (the other
// is the remainder) gives the product the puzzle asks for. Day 25 has no second
// puzzle, so part two returns the empty finale answer.

use std::collections::{HashMap, VecDeque};

/// Parse into an adjacency list over interned node indices.
fn parse<'a>(input: &'a str) -> Vec<Vec<usize>> {
    let mut ids: HashMap<&'a str, usize> = HashMap::new();
    let mut adj: Vec<Vec<usize>> = Vec::new();

    // Intern a name to its index, extending the adjacency list on first sight.
    fn intern<'a>(name: &'a str, ids: &mut HashMap<&'a str, usize>, adj: &mut Vec<Vec<usize>>) -> usize {
        if let Some(&i) = ids.get(name) {
            return i;
        }
        let i = adj.len();
        adj.push(Vec::new());
        ids.insert(name, i);
        i
    }

    for line in input.lines() {
        let (left, rights) = line.split_once(": ").unwrap();
        let a = intern(left, &mut ids, &mut adj);
        for right in rights.split_whitespace() {
            let b = intern(right, &mut ids, &mut adj);
            adj[a].push(b);
            adj[b].push(a);
        }
    }
    adj
}

/// Number of edge-disjoint paths from `s` to `t` (max flow with unit capacities),
/// plus the set of nodes still reachable from `s` in the residual graph.
fn max_flow(adj: &[Vec<usize>], s: usize, t: usize) -> (u32, Vec<bool>) {
    let n = adj.len();
    // Residual capacity per directed edge, keyed (from, to).
    let mut cap: HashMap<(usize, usize), i32> = HashMap::new();
    for (u, nbrs) in adj.iter().enumerate() {
        for &v in nbrs {
            *cap.entry((u, v)).or_insert(0) += 1;
        }
    }

    let mut flow = 0;
    loop {
        // BFS for an augmenting path in the residual graph.
        let mut prev = vec![usize::MAX; n];
        prev[s] = s;
        let mut queue = VecDeque::from([s]);
        while let Some(u) = queue.pop_front() {
            if u == t {
                break;
            }
            for &v in &adj[u] {
                if prev[v] == usize::MAX && cap.get(&(u, v)).copied().unwrap_or(0) > 0 {
                    prev[v] = u;
                    queue.push_back(v);
                }
            }
        }
        if prev[t] == usize::MAX {
            // No more augmenting paths: `prev` marks the source side of the cut.
            let reachable: Vec<bool> = prev.iter().map(|&p| p != usize::MAX).collect();
            return (flow, reachable);
        }
        // Push one unit along the found path.
        let mut v = t;
        while v != s {
            let u = prev[v];
            *cap.get_mut(&(u, v)).unwrap() -= 1;
            *cap.entry((v, u)).or_insert(0) += 1;
            v = u;
        }
        flow += 1;
    }
}

pub fn part_one(input: &str) -> String {
    let adj = parse(input);
    let n = adj.len();

    // The two clusters are joined by exactly three wires, so the max flow between
    // any two nodes on opposite sides is 3. Fix the source at 0 and find a node
    // whose min cut to it is 3; the residual-reachable set is one cluster.
    for t in 1..n {
        let (flow, reachable) = max_flow(&adj, 0, t);
        if flow == 3 {
            let size = reachable.iter().filter(|&&r| r).count();
            return (size * (n - size)).to_string();
        }
    }

    "0".to_string()
}

pub fn part_two(_input: &str) -> String {
    // Day 25 grants the final star for completing the other days.
    String::new()
}
