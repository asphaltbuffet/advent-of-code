// Solution for Advent of Code 2018 day 15 (Beverage Bandits).
//
// A reading-order combat simulation. Each round, live units act in reading order
// (row-major of their position at round start). A unit moves toward the nearest
// in-range square via BFS (ties broken by reading order twice: destination square,
// then first step) and attacks the weakest adjacent enemy.

use std::collections::{HashMap, VecDeque};

// Neighbor offsets in reading order (up, left, right, down) so BFS naturally
// discovers the reading-order-first shortest path.
const READING_DIRS: [(i32, i32); 4] = [(0, -1), (-1, 0), (1, 0), (0, 1)];

#[derive(Clone)]
struct Unit {
    x: i32,
    y: i32,
    kind: u8, // b'E' or b'G'
    hp: i32,
    alive: bool,
}

type Point = (i32, i32);

fn parse_cave(instr: &str) -> (Vec<Vec<u8>>, Vec<Unit>) {
    let mut grid: Vec<Vec<u8>> = Vec::new();
    let mut units: Vec<Unit> = Vec::new();

    for (y, line) in instr.trim_end_matches('\n').lines().enumerate() {
        let mut row: Vec<u8> = line.bytes().collect();
        for x in 0..row.len() {
            let c = row[x];
            if c == b'E' || c == b'G' {
                units.push(Unit {
                    x: x as i32,
                    y: y as i32,
                    kind: c,
                    hp: 200,
                    alive: true,
                });
                row[x] = b'.';
            }
        }
        grid.push(row);
    }

    (grid, units)
}

fn is_open(grid: &[Vec<u8>], occ: &HashMap<Point, usize>, x: i32, y: i32) -> bool {
    if y < 0 || y as usize >= grid.len() || x < 0 || x as usize >= grid[y as usize].len() {
        return false;
    }
    if grid[y as usize][x as usize] != b'.' {
        return false;
    }
    !occ.contains_key(&(x, y))
}

// Returns the index of the weakest enemy adjacent to (x, y), preferring lower HP
// then reading order, or None.
fn adjacent_enemy(units: &[Unit], x: i32, y: i32, kind: u8) -> Option<usize> {
    let mut best: Option<usize> = None;
    for (dx, dy) in READING_DIRS {
        let (nx, ny) = (x + dx, y + dy);
        for (i, e) in units.iter().enumerate() {
            if !e.alive || e.kind == kind || e.x != nx || e.y != ny {
                continue;
            }
            match best {
                None => best = Some(i),
                Some(b) if e.hp < units[b].hp => best = Some(i),
                _ => {}
            }
        }
    }
    best
}

// Finds the first step of the shortest path to the nearest in-range square, with
// all ties broken by reading order.
fn step_toward(
    grid: &[Vec<u8>],
    occ: &HashMap<Point, usize>,
    units: &[Unit],
    u: &Unit,
) -> Option<Point> {
    // In-range squares: open floor adjacent to a live enemy.
    let mut in_range: HashMap<Point, ()> = HashMap::new();
    for e in units {
        if !e.alive || e.kind == u.kind {
            continue;
        }
        for (dx, dy) in READING_DIRS {
            let (nx, ny) = (e.x + dx, e.y + dy);
            if is_open(grid, occ, nx, ny) {
                in_range.insert((nx, ny), ());
            }
        }
    }
    if in_range.is_empty() {
        return None;
    }

    // BFS from the unit, recording distance and the reading-order-first initial
    // step to reach each square. Expanding neighbors in reading order guarantees
    // the first step recorded is reading-order-minimal.
    let start = (u.x, u.y);
    let mut dist: HashMap<Point, i32> = HashMap::new();
    dist.insert(start, 0);
    let mut first_step: HashMap<Point, Point> = HashMap::new();
    let mut queue: VecDeque<Point> = VecDeque::new();
    queue.push_back(start);

    let mut chosen: Point = (0, 0);
    let mut found = false;
    let mut best_dist = 0;

    while let Some(cur) = queue.pop_front() {
        let cur_dist = dist[&cur];
        if found && cur_dist > best_dist {
            break;
        }
        if cur != start && in_range.contains_key(&cur) && !found {
            found = true;
            best_dist = cur_dist;
            chosen = cur;
        }
        for (dx, dy) in READING_DIRS {
            let np = (cur.0 + dx, cur.1 + dy);
            if !is_open(grid, occ, np.0, np.1) {
                continue;
            }
            if dist.contains_key(&np) {
                continue;
            }
            dist.insert(np, cur_dist + 1);
            let fs = if cur == start { np } else { first_step[&cur] };
            first_step.insert(np, fs);
            queue.push_back(np);
        }
    }

    if !found {
        return None;
    }
    Some(first_step[&chosen])
}

fn combat(instr: &str, elf_ap: i32, stop_on_elf_death: bool) -> i32 {
    let (grid, mut units) = parse_cave(instr);

    let mut rounds = 0;
    loop {
        units.sort_by(|a, b| (a.y, a.x).cmp(&(b.y, b.x)));

        for i in 0..units.len() {
            if !units[i].alive {
                continue;
            }
            let (ux, uy, ukind) = (units[i].x, units[i].y, units[i].kind);

            // Any enemies left? If not, combat ends mid-round.
            if !units.iter().any(|e| e.alive && e.kind != ukind) {
                let total: i32 = units.iter().filter(|e| e.alive).map(|e| e.hp).sum();
                return rounds * total;
            }

            let mut occ: HashMap<Point, usize> = HashMap::new();
            for (j, e) in units.iter().enumerate() {
                if e.alive {
                    occ.insert((e.x, e.y), j);
                }
            }

            // Move unless already adjacent to an enemy.
            if adjacent_enemy(&units, ux, uy, ukind).is_none() {
                let u = units[i].clone();
                if let Some((fx, fy)) = step_toward(&grid, &occ, &units, &u) {
                    units[i].x = fx;
                    units[i].y = fy;
                }
            }

            // Attack the weakest adjacent enemy (reading order breaks HP ties).
            let (ax, ay) = (units[i].x, units[i].y);
            if let Some(t) = adjacent_enemy(&units, ax, ay, ukind) {
                let ap = if ukind == b'E' { elf_ap } else { 3 };
                units[t].hp -= ap;
                if units[t].hp <= 0 {
                    units[t].alive = false;
                    if stop_on_elf_death && units[t].kind == b'E' {
                        return -1;
                    }
                }
            }
        }

        rounds += 1;
    }
}

pub fn part_one(input: &str) -> String {
    combat(input, 3, false).to_string()
}

pub fn part_two(input: &str) -> String {
    let mut ap = 4;
    loop {
        let outcome = combat(input, ap, true);
        if outcome != -1 {
            return outcome.to_string();
        }
        ap += 1;
    }
}
