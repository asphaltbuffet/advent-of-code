// Solution for Advent of Code 2023 day 23.
//
// We want the longest hike from top to bottom without revisiting a tile. The
// trails collapse to a small graph of junctions (cells with more than two open
// neighbors) joined by weighted corridors, so the longest-path search runs over
// ~36 nodes instead of thousands of cells. Part one honors the slopes as one-way
// edges; part two treats them as normal paths (edges become bidirectional).
// Nodes are indexed so the DFS "visited" set is a u64 bitmask.

use std::collections::HashMap;

const NEIGHBORS: [(i32, i32); 4] = [(-1, 0), (1, 0), (0, -1), (0, 1)];

fn slope_dir(tile: u8) -> Option<(i32, i32)> {
    match tile {
        b'^' => Some((-1, 0)),
        b'v' => Some((1, 0)),
        b'<' => Some((0, -1)),
        b'>' => Some((0, 1)),
        _ => None,
    }
}

/// Build the compressed adjacency (node index -> [(neighbor index, length)]),
/// returning it along with the start and end node indices.
fn compress(grid: &[&[u8]], slippery: bool) -> (Vec<Vec<(usize, u32)>>, usize, usize) {
    let (rows, cols) = (grid.len() as i32, grid[0].len() as i32);
    let open = |r: i32, c: i32| {
        r >= 0 && c >= 0 && r < rows && c < cols && grid[r as usize][c as usize] != b'#'
    };

    let start = (0i32, grid[0].iter().position(|&b| b == b'.').unwrap() as i32);
    let end = (
        rows - 1,
        grid[rows as usize - 1].iter().position(|&b| b == b'.').unwrap() as i32,
    );

    // Collect junctions: start, end, and any cell with >2 open neighbors.
    let mut nodes = vec![start, end];
    for r in 0..rows {
        for c in 0..cols {
            if open(r, c) && NEIGHBORS.iter().filter(|&&(dr, dc)| open(r + dr, c + dc)).count() > 2 {
                nodes.push((r, c));
            }
        }
    }
    let index: HashMap<(i32, i32), usize> =
        nodes.iter().enumerate().map(|(i, &p)| (p, i)).collect();

    // Walk each corridor from every junction to the next junction.
    let mut graph = vec![Vec::new(); nodes.len()];
    for (&(jr, jc), &ji) in &index {
        let mut stack = vec![(jr, jc, 0u32)];
        let mut seen = vec![(jr, jc)];
        while let Some((r, c, dist)) = stack.pop() {
            if dist > 0 {
                if let Some(&ni) = index.get(&(r, c)) {
                    graph[ji].push((ni, dist));
                    continue;
                }
            }
            let tile = grid[r as usize][c as usize];
            // On a slope, part one may leave only in the slope's direction.
            let forced = slippery.then(|| slope_dir(tile)).flatten();
            let dirs: &[(i32, i32)] = match &forced {
                Some(d) => std::slice::from_ref(d),
                None => &NEIGHBORS,
            };
            for &(dr, dc) in dirs {
                let (nr, nc) = (r + dr, c + dc);
                if open(nr, nc) && !seen.contains(&(nr, nc)) {
                    seen.push((nr, nc));
                    stack.push((nr, nc, dist + 1));
                }
            }
        }
    }

    (graph, index[&start], index[&end])
}

fn longest(graph: &[Vec<(usize, u32)>], node: usize, end: usize, visited: u64) -> Option<u32> {
    if node == end {
        return Some(0);
    }
    let mut best = None;
    for &(next, weight) in &graph[node] {
        if visited & (1 << next) == 0 {
            if let Some(rest) = longest(graph, next, end, visited | (1 << next)) {
                let total = weight + rest;
                best = Some(best.map_or(total, |b: u32| b.max(total)));
            }
        }
    }
    best
}

fn solve(input: &str, slippery: bool) -> u32 {
    let grid: Vec<&[u8]> = input.lines().map(str::as_bytes).collect();
    let (graph, start, end) = compress(&grid, slippery);
    longest(&graph, start, end, 1 << start).unwrap()
}

pub fn part_one(input: &str) -> String {
    solve(input, true).to_string()
}

pub fn part_two(input: &str) -> String {
    solve(input, false).to_string()
}
