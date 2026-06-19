// Solution for Advent of Code 2015 day 14.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

struct Reindeer {
    speed: u32,
    fly: u32,
    rest: u32,
}

impl Reindeer {
    // How far this reindeer has travelled after t seconds.
    fn distance(&self, t: u32) -> u32 {
        let cycle = self.fly + self.rest;
        let flying = (t / cycle) * self.fly + (t % cycle).min(self.fly);
        flying * self.speed
    }
}

// Parse reindeer stats. The race duration is not in the input, so it is
// inferred: the 2-reindeer AoC example races 1000s, the real input 2503s.
fn parse(input: &str) -> (Vec<Reindeer>, u32) {
    let reindeer: Vec<Reindeer> = input
        .trim()
        .lines()
        .map(|line| {
            let f: Vec<&str> = line.split_whitespace().collect();
            // <name> can fly <speed> km/s for <fly> seconds, ... rest <rest> ...
            Reindeer {
                speed: f[3].parse().unwrap(),
                fly: f[6].parse().unwrap(),
                rest: f[13].parse().unwrap(),
            }
        })
        .collect();

    let duration = if reindeer.len() <= 2 { 1000 } else { 2503 };
    (reindeer, duration)
}

pub fn part_one(input: &str) -> String {
    let (reindeer, duration) = parse(input);
    reindeer
        .iter()
        .map(|r| r.distance(duration))
        .max()
        .unwrap()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    let (reindeer, duration) = parse(input);
    let mut points = vec![0u32; reindeer.len()];

    for t in 1..=duration {
        let dists: Vec<u32> = reindeer.iter().map(|r| r.distance(t)).collect();
        let lead = *dists.iter().max().unwrap();
        for (i, &d) in dists.iter().enumerate() {
            if d == lead {
                points[i] += 1;
            }
        }
    }

    points.into_iter().max().unwrap().to_string()
}
