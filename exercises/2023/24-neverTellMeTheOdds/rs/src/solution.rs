// Solution for Advent of Code 2023 day 24.
//
// Part one counts pairs of hailstones whose XY paths cross inside the test area
// at non-negative time. Part two finds a rock position R and velocity W that hits
// every stone: in the rock's frame each stone passes through the origin, so
// (P_i - R) is parallel to (V_i - W), giving (P_i - R) x (V_i - W) = 0.
// Subtracting stone 0's equation from stones 1 and 2 cancels the R x W term and
// leaves a 6x6 linear system in (R, W), solved here with exact i128 rational
// Gaussian elimination so the 15-digit answer stays precise.

/// A hailstone: position and velocity, each (x, y, z).
type Hail = ([i128; 3], [i128; 3]);

fn parse(input: &str) -> Vec<Hail> {
    input
        .lines()
        .map(|line| {
            let (pos, vel) = line.split_once(" @ ").unwrap();
            let read = |s: &str| -> [i128; 3] {
                let mut it = s.split(',').map(|n| n.trim().parse().unwrap());
                [it.next().unwrap(), it.next().unwrap(), it.next().unwrap()]
            };
            (read(pos), read(vel))
        })
        .collect()
}

fn crosses_xy(a: &Hail, b: &Hail, lo: f64, hi: f64) -> bool {
    let ([px, py, _], [vx, vy, _]) = (a.0, a.1);
    let ([qx, qy, _], [wx, wy, _]) = (b.0, b.1);
    let denom = (vx * wy - vy * wx) as f64;
    if denom == 0.0 {
        return false;
    }
    // Path parameters for a (t) and b (s); both must be non-negative (future).
    let t = (((qx - px) * wy - (qy - py) * wx) as f64) / denom;
    let s = (((qx - px) * vy - (qy - py) * vx) as f64) / denom;
    if t < 0.0 || s < 0.0 {
        return false;
    }
    let x = px as f64 + vx as f64 * t;
    let y = py as f64 + vy as f64 * t;
    (lo..=hi).contains(&x) && (lo..=hi).contains(&y)
}

/// Solve the 6x6 rational system A x = b, returning x as exact i128 values.
fn solve6(mut a: [[i128; 6]; 6], mut b: [i128; 6]) -> [i128; 6] {
    // Fraction-free (Bareiss-style) forward elimination keeps entries integral.
    for col in 0..6 {
        let pivot = (col..6).find(|&r| a[r][col] != 0).unwrap();
        a.swap(col, pivot);
        b.swap(col, pivot);
        for r in 0..6 {
            if r != col {
                let f = a[r][col];
                let g = a[col][col];
                for c in 0..6 {
                    a[r][c] = a[r][c] * g - a[col][c] * f;
                }
                b[r] = b[r] * g - b[col] * f;
                // Reduce the row by its gcd to keep magnitudes bounded.
                let mut divisor = b[r].abs();
                for c in 0..6 {
                    divisor = gcd(divisor, a[r][c].abs());
                }
                if divisor > 1 {
                    for c in 0..6 {
                        a[r][c] /= divisor;
                    }
                    b[r] /= divisor;
                }
            }
        }
    }
    let mut x = [0i128; 6];
    for i in 0..6 {
        x[i] = b[i] / a[i][i];
    }
    x
}

fn gcd(a: i128, b: i128) -> i128 {
    if b == 0 { a } else { gcd(b, a % b) }
}

fn solve_rock(hail: &[Hail]) -> i128 {
    let (p0, v0) = hail[0];
    let mut rows = [[0i128; 6]; 6];
    let mut rhs = [0i128; 6];
    let mut idx = 0;

    for other in 1..=2 {
        let (pb, vb) = hail[other];
        // Three component equations per stone pair (0 vs other).
        for &(i, j) in &[(0usize, 1usize), (1, 2), (2, 0)] {
            rows[idx][i] = vb[j] - v0[j];
            rows[idx][j] = v0[i] - vb[i];
            rows[idx][3 + i] = p0[j] - pb[j];
            rows[idx][3 + j] = pb[i] - p0[i];
            rhs[idx] = (pb[i] * vb[j] - pb[j] * vb[i]) - (p0[i] * v0[j] - p0[j] * v0[i]);
            idx += 1;
        }
    }

    let x = solve6(rows, rhs);
    x[0] + x[1] + x[2] // Rx + Ry + Rz
}

pub fn part_one(input: &str) -> String {
    let hail = parse(input);
    // The example uses a small area; the real puzzle a huge one — pick by scale.
    let big = hail.iter().any(|(p, _)| p.iter().any(|&c| c.abs() > 1_000_000));
    let (lo, hi) = if big {
        (2e14, 4e14)
    } else {
        (7.0, 27.0)
    };

    let mut count = 0;
    for i in 0..hail.len() {
        for j in (i + 1)..hail.len() {
            if crosses_xy(&hail[i], &hail[j], lo, hi) {
                count += 1;
            }
        }
    }
    count.to_string()
}

pub fn part_two(input: &str) -> String {
    solve_rock(&parse(input)).to_string()
}
