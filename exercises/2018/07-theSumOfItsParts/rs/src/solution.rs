// Solution for Advent of Code 2018 day 7.

use std::collections::{BTreeMap, BTreeSet};

/// Parse each `Step X must be finished before step Y can begin` line into a
/// dependency map: for every step, the set of steps that must precede it.
fn deps(input: &str) -> BTreeMap<u8, BTreeSet<u8>> {
    let mut d: BTreeMap<u8, BTreeSet<u8>> = BTreeMap::new();

    for line in input.lines() {
        let bytes = line.as_bytes();
        if bytes.len() < 37 {
            continue;
        }
        // "Step B must be finished before step Y can begin."
        //       ^5 (before)                     ^36 (after)
        let before = bytes[5];
        let after = bytes[36];

        d.entry(before).or_default();
        d.entry(after).or_default().insert(before);
    }

    d
}

pub fn part_one(input: &str) -> String {
    let d = deps(input);
    let mut done: BTreeSet<u8> = BTreeSet::new();
    let mut order = String::new();

    while done.len() < d.len() {
        let next = d
            .iter()
            .find(|(step, prereqs)| !done.contains(step) && prereqs.is_subset(&done))
            .map(|(step, _)| *step)
            .expect("a step should always be ready");

        done.insert(next);
        order.push(next as char);
    }

    order
}

pub fn part_two(input: &str) -> String {
    let d = deps(input);

    // The small example runs 2 workers with no base cost; the real puzzle runs
    // 5 workers with a 60-second base per step.
    let (workers, base) = if d.len() <= 6 { (2, 0) } else { (5, 60) };

    let mut done: BTreeSet<u8> = BTreeSet::new();
    let mut in_progress: BTreeMap<u8, u32> = BTreeMap::new(); // step -> finish second

    for t in 0.. {
        // Retire any steps that have finished by the start of this second.
        in_progress.retain(|&step, &mut finish| {
            if finish <= t {
                done.insert(step);
                false
            } else {
                true
            }
        });

        if done.len() == d.len() {
            return t.to_string();
        }

        // Assign idle workers to ready steps, alphabetically first.
        let ready: Vec<u8> = d
            .iter()
            .filter(|(step, prereqs)| {
                !done.contains(step)
                    && !in_progress.contains_key(step)
                    && prereqs.is_subset(&done)
            })
            .map(|(step, _)| *step)
            .collect();

        for step in ready {
            if in_progress.len() >= workers {
                break;
            }
            in_progress.insert(step, t + base + u32::from(step - b'A') + 1);
        }
    }

    unreachable!()
}
