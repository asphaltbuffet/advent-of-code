// Solution for Advent of Code 2015 day 22.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

#[derive(Clone, Copy)]
struct State {
    php: i32,
    mana: i32,
    bhp: i32,
    bdmg: i32,
    shield: i32,
    poison: i32,
    recharge: i32,
    spent: i32,
}

impl State {
    // Apply active effects once and decrement their timers.
    fn tick(&mut self) {
        if self.poison > 0 {
            self.bhp -= 3;
            self.poison -= 1;
        }
        if self.recharge > 0 {
            self.mana += 101;
            self.recharge -= 1;
        }
        if self.shield > 0 {
            self.shield -= 1;
        }
    }
}

// (cost, instant boss damage, instant heal, shield, poison, recharge); the last
// three set an effect timer when nonzero.
const SPELLS: [(i32, i32, i32, i32, i32, i32); 5] = [
    (53, 4, 0, 0, 0, 0),
    (73, 2, 2, 0, 0, 0),
    (113, 0, 0, 6, 0, 0),
    (173, 0, 0, 0, 6, 0),
    (229, 0, 0, 0, 0, 5),
];

// Branch-and-bound DFS for the least mana spent to win. hard adds the part-two
// rule (player loses 1 HP at the start of each of their turns).
fn search(mut s: State, hard: bool, best: &mut i32) {
    if s.spent >= *best {
        return;
    }

    // Player turn.
    if hard {
        s.php -= 1;
        if s.php <= 0 {
            return;
        }
    }
    s.tick();
    if s.bhp <= 0 {
        *best = (*best).min(s.spent);
        return;
    }

    for &(cost, dmg, heal, sh, po, re) in &SPELLS {
        if cost > s.mana {
            continue;
        }
        if (sh > 0 && s.shield > 0) || (po > 0 && s.poison > 0) || (re > 0 && s.recharge > 0) {
            continue;
        }

        let mut n = s;
        n.mana -= cost;
        n.spent += cost;
        n.bhp -= dmg;
        n.php += heal;
        if sh > 0 {
            n.shield = sh;
        }
        if po > 0 {
            n.poison = po;
        }
        if re > 0 {
            n.recharge = re;
        }

        if n.bhp <= 0 {
            *best = (*best).min(n.spent);
            continue;
        }

        // Boss turn.
        n.tick();
        if n.bhp <= 0 {
            *best = (*best).min(n.spent);
            continue;
        }
        let armor = if n.shield > 0 { 7 } else { 0 };
        n.php -= (n.bdmg - armor).max(1);
        if n.php <= 0 {
            continue;
        }

        search(n, hard, best);
    }
}

fn parse_boss(input: &str) -> (i32, i32) {
    let v: Vec<i32> = input
        .trim()
        .lines()
        .map(|l| l.split(": ").nth(1).unwrap().trim().parse().unwrap())
        .collect();
    (v[0], v[1])
}

fn solve(input: &str, hard: bool) -> i32 {
    let (bhp, bdmg) = parse_boss(input);
    let mut best = 1 << 30;
    search(
        State { php: 50, mana: 500, bhp, bdmg, shield: 0, poison: 0, recharge: 0, spent: 0 },
        hard,
        &mut best,
    );
    best
}

pub fn part_one(input: &str) -> String {
    solve(input, false).to_string()
}

pub fn part_two(input: &str) -> String {
    solve(input, true).to_string()
}
