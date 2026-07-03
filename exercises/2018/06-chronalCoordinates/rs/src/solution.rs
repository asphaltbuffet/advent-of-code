// Solution for Advent of Code 2018 day 6.
//
// Manhattan-distance Voronoi over the coordinates' bounding box.
//
// Part One assigns each cell to its single nearest coordinate (ties leave the
// cell unowned) and reports the largest finite area. A coordinate owning a cell
// on the bounding-box edge owns an unbounded region, so it is disqualified.
//
// Part Two counts the cells whose summed distance to every coordinate is below
// a threshold (10000 for the real puzzle, 32 for the worked example, chosen by
// coordinate count). Stepping one cell outside the box adds at least one
// distance per coordinate, so the region reaches at most threshold/numPts cells
// beyond the box; scanning that padded box captures all of it.

type Point = (i64, i64);

/// Parse each coordinate by scanning its two integers, tolerant of spacing.
fn parse(input: &str) -> Vec<Point> {
    input
        .lines()
        .filter_map(|line| {
            let mut nums = line
                .split(|c: char| !c.is_ascii_digit() && c != '-')
                .filter(|s| !s.is_empty())
                .map(|s| s.parse::<i64>().unwrap());
            Some((nums.next()?, nums.next()?))
        })
        .collect()
}

fn bounds(pts: &[Point]) -> (i64, i64, i64, i64) {
    let min_x = pts.iter().map(|p| p.0).min().unwrap();
    let max_x = pts.iter().map(|p| p.0).max().unwrap();
    let min_y = pts.iter().map(|p| p.1).min().unwrap();
    let max_y = pts.iter().map(|p| p.1).max().unwrap();
    (min_x, min_y, max_x, max_y)
}

fn dist(a: Point, b: Point) -> i64 {
    (a.0 - b.0).abs() + (a.1 - b.1).abs()
}

pub fn part_one(input: &str) -> String {
    let pts = parse(input);
    let (min_x, min_y, max_x, max_y) = bounds(&pts);

    let mut area = vec![0usize; pts.len()];
    let mut infinite = vec![false; pts.len()];

    for y in min_y..=max_y {
        for x in min_x..=max_x {
            // Find the unique nearest coordinate; a tie leaves the cell unowned.
            let mut best = i64::MAX;
            let mut best_idx = None;
            let mut tie = false;
            for (i, &p) in pts.iter().enumerate() {
                let d = dist(p, (x, y));
                if d < best {
                    best = d;
                    best_idx = Some(i);
                    tie = false;
                } else if d == best {
                    tie = true;
                }
            }
            if tie {
                continue;
            }
            if let Some(i) = best_idx {
                area[i] += 1;
                if x == min_x || x == max_x || y == min_y || y == max_y {
                    infinite[i] = true;
                }
            }
        }
    }

    area.iter()
        .zip(&infinite)
        .filter(|(_, &inf)| !inf)
        .map(|(&a, _)| a)
        .max()
        .unwrap()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    let pts = parse(input);
    let threshold: i64 = if pts.len() <= 10 { 32 } else { 10000 };

    let (min_x, min_y, max_x, max_y) = bounds(&pts);
    let pad = threshold / pts.len() as i64 + 1;

    let size = (min_y - pad..=max_y + pad)
        .flat_map(|y| (min_x - pad..=max_x + pad).map(move |x| (x, y)))
        .filter(|&cell| pts.iter().map(|&p| dist(p, cell)).sum::<i64>() < threshold)
        .count();

    size.to_string()
}
