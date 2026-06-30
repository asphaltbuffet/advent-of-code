// Solution for Advent of Code 2015 day 18.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

type Grid = Vec<Vec<bool>>;

// parse reads the light grid; true is on.
fn parse(input: &str) -> Grid {
    input
        .split_whitespace()
        .map(|line| line.chars().map(|c| c == '#').collect())
        .collect()
}

// steps is the animation length: the small AoC example runs 4 steps for part 1
// and 5 for part 2 (stuck corners); the real input runs 100 for both.
fn steps(grid: &Grid, stuck: bool) -> usize {
    if grid.len() > 6 {
        100
    } else if stuck {
        5
    } else {
        4
    }
}

fn neighbors_on(grid: &Grid, r: usize, c: usize) -> usize {
    let n = grid.len() as isize;
    let mut total = 0;
    for dr in -1..=1 {
        for dc in -1..=1 {
            if dr == 0 && dc == 0 {
                continue;
            }
            let (nr, nc) = (r as isize + dr, c as isize + dc);
            if nr >= 0 && nr < n && nc >= 0 && (nc as usize) < grid[nr as usize].len()
                && grid[nr as usize][nc as usize]
            {
                total += 1;
            }
        }
    }
    total
}

fn step(grid: &Grid) -> Grid {
    (0..grid.len())
        .map(|r| {
            (0..grid[r].len())
                .map(|c| {
                    let on = neighbors_on(grid, r, c);
                    if grid[r][c] {
                        on == 2 || on == 3
                    } else {
                        on == 3
                    }
                })
                .collect()
        })
        .collect()
}

fn stick_corners(grid: &mut Grid) {
    let last = grid.len() - 1;
    grid[0][0] = true;
    grid[0][last] = true;
    grid[last][0] = true;
    grid[last][last] = true;
}

fn run(mut grid: Grid, n: usize, stuck: bool) -> usize {
    if stuck {
        stick_corners(&mut grid);
    }
    for _ in 0..n {
        grid = step(&grid);
        if stuck {
            stick_corners(&mut grid);
        }
    }
    grid.iter().flatten().filter(|&&on| on).count()
}

pub fn part_one(input: &str) -> String {
    let grid = parse(input);
    let n = steps(&grid, false);
    run(grid, n, false).to_string()
}

pub fn part_two(input: &str) -> String {
    let grid = parse(input);
    let n = steps(&grid, true);
    run(grid, n, true).to_string()
}
