// Solution for Advent of Code 2025 day 9.
//
// The vertices trace a rectilinear polygon. Part one is the largest inclusive
// axis-aligned rectangle spanned by any vertex pair. Part two restricts to
// rectangles that lie wholly inside the polygon: every corner on an edge or in
// the interior, and no polygon edge crossing the rectangle's interior. The tests
// are the classic on-segment, ray-cast parity, and edge-vs-rectangle checks.

type Point = (i64, i64);

fn points(input: &str) -> Vec<Point> {
    input
        .lines()
        .map(|line| {
            let (x, y) = line.split_once(',').unwrap();
            (x.parse().unwrap(), y.parse().unwrap())
        })
        .collect()
}

/// Inclusive lattice-cell area of the rectangle spanning `a` and `b`.
fn area(a: Point, b: Point) -> i64 {
    ((a.0 - b.0).abs() + 1) * ((a.1 - b.1).abs() + 1)
}

/// Iterate the polygon edges (wrapping from last vertex to first).
fn edges(pts: &[Point]) -> impl Iterator<Item = (Point, Point)> + '_ {
    let last = *pts.last().unwrap();
    std::iter::once(last).chain(pts.iter().copied()).zip(pts.iter().copied())
}

fn on_segment(a: Point, b: Point, p: Point) -> bool {
    if a.1 == b.1 && p.1 == a.1 {
        a.0.min(b.0) <= p.0 && p.0 <= a.0.max(b.0)
    } else if a.0 == b.0 && p.0 == a.0 {
        a.1.min(b.1) <= p.1 && p.1 <= a.1.max(b.1)
    } else {
        false
    }
}

fn on_edge(pts: &[Point], p: Point) -> bool {
    edges(pts).any(|(a, b)| on_segment(a, b, p))
}

fn inside(pts: &[Point], p: Point) -> bool {
    // Ray-cast parity; the integer division truncates toward zero, matching the
    // reference (safe here because the crossing edges are axis-aligned).
    let (px, py) = p;
    let mut inside = false;
    for (a, b) in edges(pts) {
        if (a.1 > py) != (b.1 > py) {
            let x_cross = (b.0 - a.0) * (py - a.1) / (b.1 - a.1) + a.0;
            if px < x_cross {
                inside = !inside;
            }
        }
    }
    inside
}

fn edge_crosses_rect(a: Point, b: Point, lo: Point, hi: Point) -> bool {
    if a.1 == b.1 {
        lo.1 < a.1 && a.1 < hi.1 && a.0.max(b.0) > lo.0 && a.0.min(b.0) < hi.0
    } else {
        lo.0 < a.0 && a.0 < hi.0 && a.1.max(b.1) > lo.1 && a.1.min(b.1) < hi.1
    }
}

fn fits(pts: &[Point], p1: Point, p2: Point) -> bool {
    let corners = [p1, (p1.0, p2.1), p2, (p2.0, p1.1)];
    if corners.iter().any(|&c| !on_edge(pts, c) && !inside(pts, c)) {
        return false;
    }
    let lo = (p1.0.min(p2.0), p1.1.min(p2.1));
    let hi = (p1.0.max(p2.0), p1.1.max(p2.1));
    !edges(pts).any(|(a, b)| edge_crosses_rect(a, b, lo, hi))
}

pub fn part_one(input: &str) -> String {
    let pts = points(input);
    let mut best = 0;
    for i in 0..pts.len() {
        for j in 0..i {
            best = best.max(area(pts[i], pts[j]));
        }
    }
    best.to_string()
}

pub fn part_two(input: &str) -> String {
    let pts = points(input);
    let mut best = 0;
    for i in 0..pts.len() {
        for j in 0..i {
            let a = area(pts[i], pts[j]);
            if a > best && fits(&pts, pts[i], pts[j]) {
                best = a;
            }
        }
    }
    best.to_string()
}
