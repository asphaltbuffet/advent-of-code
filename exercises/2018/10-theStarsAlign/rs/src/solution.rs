// Solution for Advent of Code 2018 day 10.
//
// Points of light drift at constant velocity. They spell a message when their
// bounding box is tightest; that minimum time is also the answer to part two.

/// A point of light with a position and constant velocity.
struct Star {
    x: i64,
    y: i64,
    vx: i64,
    vy: i64,
}

/// Parse every signed integer on each line as `x, y, vx, vy`.
fn parse(input: &str) -> Vec<Star> {
    input
        .lines()
        .filter(|line| !line.trim().is_empty())
        .map(|line| {
            let nums: Vec<i64> = line
                .split(|c: char| !c.is_ascii_digit() && c != '-')
                .filter(|s| !s.is_empty() && *s != "-")
                .map(|s| s.parse().unwrap())
                .collect();
            Star { x: nums[0], y: nums[1], vx: nums[2], vy: nums[3] }
        })
        .collect()
}

/// Combined width and height of the bounding box at time `t`.
fn extent(stars: &[Star], t: i64) -> i64 {
    let (mut min_x, mut min_y) = (i64::MAX, i64::MAX);
    let (mut max_x, mut max_y) = (i64::MIN, i64::MIN);
    for s in stars {
        let (x, y) = (s.x + s.vx * t, s.y + s.vy * t);
        min_x = min_x.min(x);
        max_x = max_x.max(x);
        min_y = min_y.min(y);
        max_y = max_y.max(y);
    }
    (max_x - min_x) + (max_y - min_y)
}

/// The second at which the message appears: the time minimizing the extent.
/// The extent shrinks to that minimum then grows, so we step until it stops.
fn converge(stars: &[Star]) -> i64 {
    let mut t = 0;
    while extent(stars, t + 1) < extent(stars, t) {
        t += 1;
    }
    t
}

/// Draw the stars at time `t` as a grid of '█' and ' ', rows joined by '\n'.
fn render(stars: &[Star], t: i64) -> String {
    let (mut min_x, mut min_y) = (i64::MAX, i64::MAX);
    let (mut max_x, mut max_y) = (i64::MIN, i64::MIN);
    for s in stars {
        let (x, y) = (s.x + s.vx * t, s.y + s.vy * t);
        min_x = min_x.min(x);
        max_x = max_x.max(x);
        min_y = min_y.min(y);
        max_y = max_y.max(y);
    }

    let (w, h) = ((max_x - min_x + 1) as usize, (max_y - min_y + 1) as usize);
    let mut grid = vec![vec![' '; w]; h];
    for s in stars {
        let x = (s.x + s.vx * t - min_x) as usize;
        let y = (s.y + s.vy * t - min_y) as usize;
        grid[y][x] = '\u{2588}';
    }

    grid.iter()
        .map(|row| row.iter().collect::<String>())
        .collect::<Vec<_>>()
        .join("\n")
}

pub fn part_one(input: &str) -> String {
    let stars = parse(input);
    render(&stars, converge(&stars))
}

pub fn part_two(input: &str) -> String {
    let stars = parse(input);
    converge(&stars).to_string()
}
