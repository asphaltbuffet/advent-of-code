// Solution for Advent of Code 2025 day 7.
//
// A beam starts at `S` and travels down. On a splitter (^) it forks down-left and
// down-right, and the beam at that column stops. We carry the wavefront as a dense
// vector of per-column beam counts and advance it one row at a time, building a
// fresh vector each step so children created this row are not re-split until the
// next. Part one counts split events; part two sums the surviving timelines.

fn sweep(input: &str) -> (u64, u64) {
    let lines: Vec<&[u8]> = input.lines().map(str::as_bytes).collect();
    let width = lines.iter().map(|l| l.len()).max().unwrap();

    let mut beams = vec![0u64; width];
    let start = lines[0].iter().position(|&c| c == b'S').unwrap();
    beams[start] = 1;

    let mut splits = 0u64;
    for row in &lines[1..] {
        let mut next = vec![0u64; width];
        for (x, &count) in beams.iter().enumerate() {
            if count == 0 {
                continue;
            }
            if row.get(x) == Some(&b'^') {
                next[x - 1] += count;
                next[x + 1] += count;
                splits += 1;
            } else {
                next[x] += count;
            }
        }
        beams = next;
    }

    (splits, beams.iter().sum())
}

pub fn part_one(input: &str) -> String {
    sweep(input).0.to_string()
}

pub fn part_two(input: &str) -> String {
    sweep(input).1.to_string()
}
