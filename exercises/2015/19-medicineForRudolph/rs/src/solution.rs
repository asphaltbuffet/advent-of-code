// Solution for Advent of Code 2015 day 19.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

use std::collections::HashSet;

// parse splits the input into replacement rules and the target molecule.
fn parse(input: &str) -> (Vec<(String, String)>, String) {
    let norm = input.replace("\r\n", "\n");
    let (block, mol) = norm.split_once("\n\n").unwrap();
    let rules = block
        .trim()
        .lines()
        .map(|line| {
            let (f, t) = line.split_once(" => ").unwrap();
            (f.to_string(), t.to_string())
        })
        .collect();
    (rules, mol.trim().to_string())
}

pub fn part_one(input: &str) -> String {
    let (rules, mol) = parse(input);
    let mut seen: HashSet<String> = HashSet::new();

    for (frm, to) in &rules {
        // Replace each occurrence of frm independently.
        let mut start = 0;
        while let Some(off) = mol[start..].find(frm) {
            let i = start + off;
            seen.insert(format!("{}{}{}", &mol[..i], to, &mol[i + frm.len()..]));
            start = i + 1;
        }
    }

    seen.len().to_string()
}

pub fn part_two(input: &str) -> String {
    let (mut rules, mol) = parse(input);

    // Work backwards greedily: collapse any production back to its source until
    // only "e" remains; reshuffle (small deterministic LCG) and retry on a dead
    // end. Input-agnostic for both the example and real Rn/Ar/Y grammars.
    let mut state: u64 = 1;
    let mut rng = || {
        state = state.wrapping_mul(6364136223846793005).wrapping_add(1);
        (state >> 33) as usize
    };

    loop {
        let mut cur = mol.clone();
        let mut steps = 0u32;
        let mut stuck = false;

        while cur != "e" {
            let mut applied = false;
            for (frm, to) in &rules {
                if let Some(i) = cur.find(to.as_str()) {
                    cur = format!("{}{}{}", &cur[..i], frm, &cur[i + to.len()..]);
                    steps += 1;
                    applied = true;
                    break;
                }
            }
            if !applied {
                stuck = true;
                break;
            }
        }

        if !stuck {
            return steps.to_string();
        }

        // Fisher-Yates shuffle with the LCG.
        for i in (1..rules.len()).rev() {
            let j = rng() % (i + 1);
            rules.swap(i, j);
        }
    }
}
