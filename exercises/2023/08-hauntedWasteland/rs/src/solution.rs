// Solution for Advent of Code 2023 day 8.
//
// The input is a cyclic L/R instruction string plus a graph of node -> (left,
// right). Part one counts steps from AAA to ZZZ. Part two starts a ghost at every
// node ending in 'A'; each reaches a '..Z' node on its own fixed cycle, so all
// ghosts coincide after the least common multiple of those cycle lengths.

use std::collections::HashMap;

fn parse(input: &str) -> (&[u8], HashMap<&str, (&str, &str)>) {
    let (moves, network) = input.split_once("\n\n").unwrap();
    let graph = network
        .lines()
        .map(|line| {
            // "AAA = (BBB, CCC)"
            let node = &line[0..3];
            let left = &line[7..10];
            let right = &line[12..15];
            (node, (left, right))
        })
        .collect();
    (moves.as_bytes(), graph)
}

/// Steps from `start` until `done` holds, following the cyclic instructions.
fn steps(
    start: &str,
    moves: &[u8],
    graph: &HashMap<&str, (&str, &str)>,
    done: impl Fn(&str) -> bool,
) -> u64 {
    let mut node = start;
    for (count, &m) in moves.iter().cycle().enumerate() {
        let (l, r) = graph[node];
        node = if m == b'L' { l } else { r };
        if done(node) {
            return count as u64 + 1;
        }
    }
    unreachable!()
}

fn gcd(a: u64, b: u64) -> u64 {
    if b == 0 { a } else { gcd(b, a % b) }
}

fn lcm(a: u64, b: u64) -> u64 {
    a / gcd(a, b) * b
}

pub fn part_one(input: &str) -> String {
    let (moves, graph) = parse(input);
    steps("AAA", moves, &graph, |n| n == "ZZZ").to_string()
}

pub fn part_two(input: &str) -> String {
    let (moves, graph) = parse(input);
    graph
        .keys()
        .filter(|n| n.ends_with('A'))
        .map(|start| steps(start, moves, &graph, |n| n.ends_with('Z')))
        .fold(1u64, lcm)
        .to_string()
}
