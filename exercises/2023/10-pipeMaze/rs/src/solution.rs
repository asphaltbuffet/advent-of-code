// Solution for Advent of Code 2023 day 10.
//
// The pipes form a single closed loop through the start tile S. We trace it once,
// stepping from tile to tile via each pipe's two connections. Part one's farthest
// point is half the loop length. Part two counts the enclosed tiles from the loop
// geometry: the shoelace formula gives the polygon area, and Pick's theorem
// (A = interior + boundary/2 - 1) recovers the interior count.

type Pos = (i64, i64);

/// The two direction offsets a pipe character connects, if it is a pipe.
fn connections(ch: u8) -> Option<[Pos; 2]> {
    Some(match ch {
        b'|' => [(-1, 0), (1, 0)],
        b'-' => [(0, -1), (0, 1)],
        b'L' => [(-1, 0), (0, 1)],
        b'J' => [(-1, 0), (0, -1)],
        b'7' => [(1, 0), (0, -1)],
        b'F' => [(1, 0), (0, 1)],
        _ => return None,
    })
}

fn trace(input: &str) -> Vec<Pos> {
    let grid: Vec<&[u8]> = input.lines().map(str::as_bytes).collect();
    let at = |r: i64, c: i64| -> u8 {
        if r >= 0 && (r as usize) < grid.len() && c >= 0 && (c as usize) < grid[r as usize].len() {
            grid[r as usize][c as usize]
        } else {
            b'.'
        }
    };

    let start = (0..grid.len() as i64)
        .flat_map(|r| (0..grid[r as usize].len() as i64).map(move |c| (r, c)))
        .find(|&(r, c)| at(r, c) == b'S')
        .unwrap();

    // Pick a first step: a neighbor whose pipe connects back toward S.
    let mut prev = start;
    let mut cur = (0..4)
        .map(|i| [(-1i64, 0i64), (1, 0), (0, -1), (0, 1)][i])
        .find_map(|(dr, dc)| {
            let n = (start.0 + dr, start.1 + dc);
            connections(at(n.0, n.1))
                .filter(|conns| conns.contains(&(-dr, -dc)))
                .map(|_| n)
        })
        .unwrap();

    let mut loop_tiles = vec![start];
    while cur != start {
        loop_tiles.push(cur);
        let conns = connections(at(cur.0, cur.1)).unwrap();
        let next = conns
            .iter()
            .map(|&(dr, dc)| (cur.0 + dr, cur.1 + dc))
            .find(|&n| n != prev)
            .unwrap();
        prev = cur;
        cur = next;
    }
    loop_tiles
}

pub fn part_one(input: &str) -> String {
    (trace(input).len() / 2).to_string()
}

pub fn part_two(input: &str) -> String {
    let loop_tiles = trace(input);
    let n = loop_tiles.len();

    // Shoelace over the loop vertices, then Pick's theorem for interior points.
    let twice_area: i64 = (0..n)
        .map(|i| {
            let (r0, c0) = loop_tiles[i];
            let (r1, c1) = loop_tiles[(i + 1) % n];
            r0 * c1 - r1 * c0
        })
        .sum();
    let area = twice_area.abs() / 2;
    (area - n as i64 / 2 + 1).to_string()
}
