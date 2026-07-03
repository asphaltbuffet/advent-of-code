// Solution for Advent of Code 2018 day 13.
//
// Carts ride a track of `-|/\+`. Each tick every cart, taken in reading order
// (top-to-bottom, left-to-right by current position), steps one cell. Curves
// reflect the heading; intersections cycle left / straight / right. A crash is
// detected the instant a cart lands on an occupied cell — mid-tick, not after.

/// A mine cart: position, heading, and how many intersections it has crossed
/// (which selects left / straight / right in a repeating cycle).
struct Cart {
    x: i32,
    y: i32,
    dx: i32,
    dy: i32,
    turns: usize,
    dead: bool,
}

/// Read the track grid and lift the carts off it, replacing each cart glyph with
/// the straight track it sits on.
fn parse(input: &str) -> (Vec<Vec<u8>>, Vec<Cart>) {
    let mut grid: Vec<Vec<u8>> = Vec::new();
    let mut carts = Vec::new();

    for (y, line) in input.trim_end_matches('\n').split('\n').enumerate() {
        let mut row = line.as_bytes().to_vec();
        for (x, cell) in row.iter_mut().enumerate() {
            let (dx, dy, under) = match *cell {
                b'<' => (-1, 0, b'-'),
                b'>' => (1, 0, b'-'),
                b'^' => (0, -1, b'|'),
                b'v' => (0, 1, b'|'),
                _ => continue,
            };
            carts.push(Cart {
                x: x as i32,
                y: y as i32,
                dx,
                dy,
                turns: 0,
                dead: false,
            });
            *cell = under;
        }
        grid.push(row);
    }

    (grid, carts)
}

/// Track glyph at (x, y), treating out-of-range as empty space.
fn at(grid: &[Vec<u8>], x: i32, y: i32) -> u8 {
    if y < 0 || x < 0 {
        return b' ';
    }
    grid.get(y as usize)
        .and_then(|row| row.get(x as usize))
        .copied()
        .unwrap_or(b' ')
}

impl Cart {
    /// Move one step, then turn according to the track landed on.
    fn advance(&mut self, grid: &[Vec<u8>]) {
        self.x += self.dx;
        self.y += self.dy;

        match at(grid, self.x, self.y) {
            b'/' => {
                let (dx, dy) = (self.dx, self.dy);
                self.dx = -dy;
                self.dy = -dx;
            }
            b'\\' => {
                let (dx, dy) = (self.dx, self.dy);
                self.dx = dy;
                self.dy = dx;
            }
            b'+' => {
                match self.turns % 3 {
                    0 => {
                        // left
                        let (dx, dy) = (self.dx, self.dy);
                        self.dx = dy;
                        self.dy = -dx;
                    }
                    2 => {
                        // right
                        let (dx, dy) = (self.dx, self.dy);
                        self.dx = -dy;
                        self.dy = dx;
                    }
                    _ => {}
                }
                self.turns += 1;
            }
            _ => {}
        }
    }
}

/// Run the carts until the first crash (`last_standing` false) or until one cart
/// remains (`last_standing` true), returning that location as "x,y".
fn simulate(grid: &[Vec<u8>], mut carts: Vec<Cart>, last_standing: bool) -> String {
    loop {
        carts.sort_by_key(|c| (c.y, c.x));

        for i in 0..carts.len() {
            if carts[i].dead {
                continue;
            }

            carts[i].advance(grid);
            let (cx, cy) = (carts[i].x, carts[i].y);

            // Collision with any other live cart at the new position.
            for j in 0..carts.len() {
                if i == j || carts[j].dead || carts[j].x != cx || carts[j].y != cy {
                    continue;
                }
                if !last_standing {
                    return format!("{cx},{cy}");
                }
                carts[i].dead = true;
                carts[j].dead = true;
                break;
            }
        }

        if last_standing {
            carts.retain(|c| !c.dead);
            if carts.len() == 1 {
                return format!("{},{}", carts[0].x, carts[0].y);
            }
        }
    }
}

pub fn part_one(input: &str) -> String {
    let (grid, carts) = parse(input);
    simulate(&grid, carts, false)
}

pub fn part_two(input: &str) -> String {
    let (grid, carts) = parse(input);
    simulate(&grid, carts, true)
}
