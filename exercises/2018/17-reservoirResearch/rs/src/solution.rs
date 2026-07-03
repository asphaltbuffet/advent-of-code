// Solution for Advent of Code 2018 day 17.
//
// Water flows from the spring at x=500 through clay veins. A recursive flood-fill
// drops water down, spreads it across floors, and settles rows that are walled on
// both sides. Both parts share one simulation; part one counts every wet tile,
// part two only the settled water.

const SAND: u8 = 0;
const CLAY: u8 = 1;
const FLOWING: u8 = 2;
const SETTLED: u8 = 3;

struct Grid {
    tiles: Vec<u8>,
    width: usize,
    min_x: i32,
    min_y: i32,
    max_y: i32,
}

impl Grid {
    fn idx(&self, x: i32, y: i32) -> usize {
        (y as usize) * self.width + (x - self.min_x) as usize
    }
    fn at(&self, x: i32, y: i32) -> u8 {
        self.tiles[self.idx(x, y)]
    }
    fn set(&mut self, x: i32, y: i32, v: u8) {
        let i = self.idx(x, y);
        self.tiles[i] = v;
    }
    fn is_floor(&self, x: i32, y: i32) -> bool {
        matches!(self.at(x, y), CLAY | SETTLED)
    }

    // Drop water from (x, y); return true if it comes to rest as settled water,
    // making it a floor for the row above.
    fn flow(&mut self, x: i32, y: i32) -> bool {
        self.set(x, y, FLOWING);
        let below = y + 1;
        if below > self.max_y {
            return false; // falls off the bottom of the scan
        }
        if !self.is_floor(x, below) {
            match self.at(x, below) {
                SAND if self.flow(x, below) => {}
                SAND => return false,      // fell through and kept flowing
                _ => return false,         // already flowing here: not a floor
            }
        }
        let (lx, left_wall) = self.spread(x, y, -1);
        let (rx, right_wall) = self.spread(x, y, 1);
        if left_wall && right_wall {
            for fx in lx..=rx {
                self.set(fx, y, SETTLED);
            }
            return true;
        }
        false
    }

    // Walk from x in direction dir, marking flowing water and spilling down open
    // edges. If a spill settles into a new floor, keep extending the row. Returns
    // the furthest reachable x and whether it ended against a clay wall.
    fn spread(&mut self, mut x: i32, y: i32, dir: i32) -> (i32, bool) {
        loop {
            if self.at(x + dir, y) == CLAY {
                return (x, true);
            }
            x += dir;
            self.set(x, y, FLOWING);
            if !self.is_floor(x, y + 1) {
                if self.at(x, y + 1) == SAND && self.flow(x, y + 1) {
                    continue; // spill settled: extend the row past this edge
                }
                return (x, false); // water escapes here
            }
        }
    }
}

fn simulate(input: &str) -> Grid {
    // Parse each vein as its three integers; track the clay bounding box.
    let mut veins: Vec<(i32, i32, i32, i32)> = Vec::new();
    let (mut min_x, mut max_x) = (i32::MAX, i32::MIN);
    let (mut min_y, mut max_y) = (i32::MAX, i32::MIN);

    for line in input.trim().lines() {
        let mut nums = line
            .split(|c: char| !c.is_ascii_digit())
            .filter(|s| !s.is_empty())
            .map(|s| s.parse::<i32>().unwrap());
        let (a, b, c) = (nums.next().unwrap(), nums.next().unwrap(), nums.next().unwrap());
        let v = if line.starts_with('x') {
            (a, a, b, c) // x=a, y=b..c
        } else {
            (b, c, a, a) // y=a, x=b..c
        };
        veins.push(v);
        min_x = min_x.min(v.0);
        max_x = max_x.max(v.1);
        min_y = min_y.min(v.2);
        max_y = max_y.max(v.3);
    }

    // Pad x by one each side so water can spill past the outermost clay.
    min_x -= 1;
    max_x += 1;
    let width = (max_x - min_x + 1) as usize;

    let mut grid = Grid {
        tiles: vec![SAND; width * (max_y as usize + 1)],
        width,
        min_x,
        min_y,
        max_y,
    };
    for (x1, x2, y1, y2) in veins {
        for y in y1..=y2 {
            for x in x1..=x2 {
                grid.set(x, y, CLAY);
            }
        }
    }
    grid.flow(500, 0);
    grid
}

fn count(grid: &Grid, include_flowing: bool) -> usize {
    let start = grid.idx(grid.min_x, grid.min_y);
    grid.tiles[start..]
        .iter()
        .filter(|&&t| t == SETTLED || (include_flowing && t == FLOWING))
        .count()
}

pub fn part_one(input: &str) -> String {
    count(&simulate(input), true).to_string()
}

pub fn part_two(input: &str) -> String {
    count(&simulate(input), false).to_string()
}
