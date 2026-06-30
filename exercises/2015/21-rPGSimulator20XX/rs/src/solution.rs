// Solution for Advent of Code 2015 day 21.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

const PLAYER_HP: i32 = 100;

// Shop tables (cost, damage, armor) from the problem statement, not the input.
const WEAPONS: [(i32, i32, i32); 5] =
    [(8, 4, 0), (10, 5, 0), (25, 6, 0), (40, 7, 0), (74, 8, 0)];
const ARMORS: [(i32, i32, i32); 6] =
    [(0, 0, 0), (13, 0, 1), (31, 0, 2), (53, 0, 3), (75, 0, 4), (102, 0, 5)];
const RINGS: [(i32, i32, i32); 6] =
    [(25, 1, 0), (50, 2, 0), (100, 3, 0), (20, 0, 1), (40, 0, 2), (80, 0, 3)];

fn parse_boss(input: &str) -> (i32, i32, i32) {
    let v: Vec<i32> = input
        .trim()
        .lines()
        .map(|l| l.split(": ").nth(1).unwrap().trim().parse().unwrap())
        .collect();
    (v[0], v[1], v[2]) // hp, damage, armor
}

// loadouts builds every legal (cost, damage, armor): one weapon, zero/one
// armor, and zero/one/two distinct rings.
fn loadouts() -> Vec<(i32, i32, i32)> {
    // Ring index sets: none, one, or two distinct.
    let mut ring_sets: Vec<Vec<usize>> = vec![vec![]];
    for i in 0..RINGS.len() {
        ring_sets.push(vec![i]);
        for j in i + 1..RINGS.len() {
            ring_sets.push(vec![i, j]);
        }
    }

    let mut out = Vec::new();
    for w in WEAPONS {
        for a in ARMORS {
            for rs in &ring_sets {
                let (mut c, mut d, mut ar) = (w.0 + a.0, w.1 + a.1, w.2 + a.2);
                for &ri in rs {
                    c += RINGS[ri].0;
                    d += RINGS[ri].1;
                    ar += RINGS[ri].2;
                }
                out.push((c, d, ar));
            }
        }
    }
    out
}

// player_wins: each side needs ceil(targetHP / effectiveDamage) hits; the
// player strikes first, so ties go to the player.
fn player_wins(p_dmg: i32, p_arm: i32, b_hp: i32, b_dmg: i32, b_arm: i32) -> bool {
    let boss_turns = ceil_div(b_hp, (p_dmg - b_arm).max(1));
    let player_turns = ceil_div(PLAYER_HP, (b_dmg - p_arm).max(1));
    boss_turns <= player_turns
}

fn ceil_div(a: i32, b: i32) -> i32 {
    (a + b - 1) / b
}

pub fn part_one(input: &str) -> String {
    let (bhp, bdmg, barm) = parse_boss(input);
    loadouts()
        .into_iter()
        .filter(|&(_, d, ar)| player_wins(d, ar, bhp, bdmg, barm))
        .map(|(c, _, _)| c)
        .min()
        .unwrap()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    let (bhp, bdmg, barm) = parse_boss(input);
    loadouts()
        .into_iter()
        .filter(|&(_, d, ar)| !player_wins(d, ar, bhp, bdmg, barm))
        .map(|(c, _, _)| c)
        .max()
        .unwrap()
        .to_string()
}
