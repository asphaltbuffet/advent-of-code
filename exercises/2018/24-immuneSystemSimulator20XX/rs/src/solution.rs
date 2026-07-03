// Solution for Advent of Code 2018 day 24.
//
// Two armies of groups fight in rounds: target selection (by effective power, then
// initiative) then attack (by initiative), with damage doubled against a weakness
// and nullified by an immunity. Part one is the winning army's surviving units.
// Part two is the smallest immune-system attack boost that wins, watching for
// stalemate rounds that kill nobody.

use std::cmp::Reverse;
use std::collections::HashSet;

#[derive(Clone)]
struct Group {
    army: u8, // 0 = immune system, 1 = infection
    units: i64,
    hp: i64,
    damage: i64,
    attack: String,
    initiative: i64,
    weak: HashSet<String>,
    immune: HashSet<String>,
}

impl Group {
    fn power(&self) -> i64 {
        self.units * self.damage
    }
    fn damage_to(&self, other: &Group) -> i64 {
        if other.immune.contains(&self.attack) {
            0
        } else if other.weak.contains(&self.attack) {
            self.power() * 2
        } else {
            self.power()
        }
    }
}

fn parse(input: &str) -> Vec<Group> {
    let mut groups = Vec::new();
    let mut army = 0u8;
    for line in input.trim().lines() {
        let line = line.trim();
        if line.starts_with("Immune System") {
            army = 0;
        } else if line.starts_with("Infection") {
            army = 1;
        } else if !line.is_empty() {
            groups.push(parse_group(line, army));
        }
    }
    groups
}

fn parse_group(line: &str, army: u8) -> Group {
    // Optional "(...)" modifier clause.
    let (mut weak, mut immune) = (HashSet::new(), HashSet::new());
    if let (Some(o), Some(c)) = (line.find('('), line.find(')')) {
        for clause in line[o + 1..c].split("; ") {
            if let Some(t) = clause.strip_prefix("weak to ") {
                weak = t.split(", ").map(String::from).collect();
            } else if let Some(t) = clause.strip_prefix("immune to ") {
                immune = t.split(", ").map(String::from).collect();
            }
        }
    }
    // Numbers and the attack word from the fixed sentence structure.
    let nums: Vec<i64> = line
        .split(|c: char| !c.is_ascii_digit())
        .filter(|s| !s.is_empty())
        .map(|s| s.parse().unwrap())
        .collect();
    let attack = line[..line.find(" damage").unwrap()]
        .rsplit(' ')
        .next()
        .unwrap()
        .to_string();
    Group {
        army,
        units: nums[0],
        hp: nums[1],
        damage: nums[nums.len() - 2],
        attack,
        initiative: nums[nums.len() - 1],
        weak,
        immune,
    }
}

// fight runs to completion, returning (winning army, its surviving units). A round
// that kills nobody is a stalemate, reported as an infection win.
fn fight(groups: &mut [Group]) -> (u8, i64) {
    loop {
        let alive: Vec<usize> = (0..groups.len()).filter(|&i| groups[i].units > 0).collect();
        let immune_units: i64 = alive.iter().filter(|&&i| groups[i].army == 0).map(|&i| groups[i].units).sum();
        let infect_units: i64 = alive.iter().filter(|&&i| groups[i].army == 1).map(|&i| groups[i].units).sum();
        if immune_units == 0 || infect_units == 0 {
            return if immune_units > 0 { (0, immune_units) } else { (1, infect_units) };
        }

        // Target selection.
        let mut order = alive.clone();
        order.sort_by_key(|&i| (Reverse(groups[i].power()), Reverse(groups[i].initiative)));
        let mut targets: Vec<Option<usize>> = vec![None; groups.len()];
        let mut taken: HashSet<usize> = HashSet::new();
        for &a in &order {
            let best = alive
                .iter()
                .filter(|&&d| groups[d].army != groups[a].army && !taken.contains(&d) && groups[a].damage_to(&groups[d]) > 0)
                .max_by_key(|&&d| (groups[a].damage_to(&groups[d]), groups[d].power(), groups[d].initiative))
                .copied();
            if let Some(d) = best {
                targets[a] = Some(d);
                taken.insert(d);
            }
        }

        // Attack in decreasing initiative order.
        let mut attackers = alive.clone();
        attackers.sort_by_key(|&i| Reverse(groups[i].initiative));
        let mut killed = 0i64;
        for a in attackers {
            if groups[a].units <= 0 {
                continue;
            }
            if let Some(d) = targets[a] {
                let dead = (groups[a].damage_to(&groups[d]) / groups[d].hp).min(groups[d].units);
                groups[d].units -= dead;
                killed += dead;
            }
        }
        if killed == 0 {
            return (1, infect_units); // stalemate
        }
    }
}

pub fn part_one(input: &str) -> String {
    let mut groups = parse(input);
    fight(&mut groups).1.to_string()
}

pub fn part_two(input: &str) -> String {
    let base = parse(input);
    for boost in 0.. {
        let mut groups = base.clone();
        for g in groups.iter_mut().filter(|g| g.army == 0) {
            g.damage += boost;
        }
        let (winner, units) = fight(&mut groups);
        if winner == 0 {
            return units.to_string();
        }
    }
    unreachable!()
}
