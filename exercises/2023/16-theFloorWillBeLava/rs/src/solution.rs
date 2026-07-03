// Solution for Advent of Code 2023 day 16.
//
// A beam of light enters the grid and is bent by mirrors ('/', '\\') and split
// by splitters ('|', '-'). We trace it with an explicit stack, marking each
// (row, col, direction) state visited in a flat boolean grid so cycles stop and
// splitters don't loop forever; the energized count is the distinct cells that
// held any beam. Part one starts at the top-left heading right; part two takes
// the best of every edge entry point.

/// Directions indexed 0..4: up, down, left, right.
const DELTAS: [(i32, i32); 4] = [(-1, 0), (1, 0), (0, -1), (0, 1)];
const UP: usize = 0;
const DOWN: usize = 1;
const LEFT: usize = 2;
const RIGHT: usize = 3;

/// Outgoing directions for a tile given the incoming direction.
fn bounce(tile: u8, dir: usize) -> [Option<usize>; 2] {
    match tile {
        b'/' => [Some([RIGHT, LEFT, DOWN, UP][dir]), None],
        b'\\' => [Some([LEFT, RIGHT, UP, DOWN][dir]), None],
        b'|' if dir == LEFT || dir == RIGHT => [Some(UP), Some(DOWN)],
        b'-' if dir == UP || dir == DOWN => [Some(LEFT), Some(RIGHT)],
        _ => [Some(dir), None], // '.', or a splitter hit end-on: pass through
    }
}

fn energized(grid: &[&[u8]], start: (i32, i32, usize)) -> usize {
    let (h, w) = (grid.len() as i32, grid[0].len() as i32);
    let mut seen = vec![false; (h * w) as usize * 4];
    let mut lit = vec![false; (h * w) as usize];
    let mut stack = vec![start];

    while let Some((r, c, d)) = stack.pop() {
        if r < 0 || c < 0 || r >= h || c >= w {
            continue;
        }
        let state = ((r * w + c) as usize) * 4 + d;
        if seen[state] {
            continue;
        }
        seen[state] = true;
        lit[(r * w + c) as usize] = true;

        for nd in bounce(grid[r as usize][c as usize], d).into_iter().flatten() {
            let (dr, dc) = DELTAS[nd];
            stack.push((r + dr, c + dc, nd));
        }
    }

    lit.iter().filter(|&&b| b).count()
}

pub fn part_one(input: &str) -> String {
    let grid: Vec<&[u8]> = input.lines().map(str::as_bytes).collect();
    energized(&grid, (0, 0, RIGHT)).to_string()
}

pub fn part_two(input: &str) -> String {
    let grid: Vec<&[u8]> = input.lines().map(str::as_bytes).collect();
    let (h, w) = (grid.len() as i32, grid[0].len() as i32);

    let mut best = 0;
    for c in 0..w {
        best = best.max(energized(&grid, (0, c, DOWN)));
        best = best.max(energized(&grid, (h - 1, c, UP)));
    }
    for r in 0..h {
        best = best.max(energized(&grid, (r, 0, RIGHT)));
        best = best.max(energized(&grid, (r, w - 1, LEFT)));
    }

    best.to_string()
}
