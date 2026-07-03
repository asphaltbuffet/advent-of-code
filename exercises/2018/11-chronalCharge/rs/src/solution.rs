// Solution for Advent of Code 2018 day 11.

const GRID_SIZE: usize = 300;

/// Power level of the fuel cell at (x, y) for the given serial.
fn cell_power(x: i32, y: i32, serial: i32) -> i32 {
    let rack_id = x + 10;
    let power = (rack_id * y + serial) * rack_id;
    (power / 100) % 10 - 5
}

/// Build a summed-area table so any square's total power is an O(1) lookup.
/// `sat[y * stride + x]` holds the sum of all cells from (1,1) to (x,y) inclusive.
fn summed_area_table(serial: i32) -> Vec<i32> {
    let stride = GRID_SIZE + 1;
    let mut sat = vec![0i32; stride * stride];
    for y in 1..=GRID_SIZE {
        for x in 1..=GRID_SIZE {
            let p = cell_power(x as i32, y as i32, serial);
            sat[y * stride + x] = p
                + sat[(y - 1) * stride + x]
                + sat[y * stride + (x - 1)]
                - sat[(y - 1) * stride + (x - 1)];
        }
    }
    sat
}

/// Find the top-left corner and total of the highest-power square of `size`.
fn best(sat: &[i32], size: usize) -> (usize, usize, i32) {
    let stride = GRID_SIZE + 1;
    let (mut bx, mut by, mut best_sum) = (0, 0, i32::MIN);
    for y in 1..=GRID_SIZE - size + 1 {
        for x in 1..=GRID_SIZE - size + 1 {
            let (x2, y2) = (x + size - 1, y + size - 1);
            let s = sat[y2 * stride + x2]
                - sat[(y - 1) * stride + x2]
                - sat[y2 * stride + (x - 1)]
                + sat[(y - 1) * stride + (x - 1)];
            if s > best_sum {
                best_sum = s;
                bx = x;
                by = y;
            }
        }
    }
    (bx, by, best_sum)
}

pub fn part_one(input: &str) -> String {
    let serial: i32 = input.trim().parse().unwrap();
    let sat = summed_area_table(serial);
    let (x, y, _) = best(&sat, 3);
    format!("{x},{y}")
}

pub fn part_two(input: &str) -> String {
    let serial: i32 = input.trim().parse().unwrap();
    let sat = summed_area_table(serial);

    // Scan every square size; the summed-area table keeps each total O(1),
    // so the whole search is O(GRID_SIZE^3).
    let (mut bx, mut by, mut bsize, mut bsum) = (0, 0, 0, i32::MIN);
    for size in 1..=GRID_SIZE {
        let (x, y, sum) = best(&sat, size);
        if sum > bsum {
            bsum = sum;
            bx = x;
            by = y;
            bsize = size;
        }
    }
    format!("{bx},{by},{bsize}")
}
