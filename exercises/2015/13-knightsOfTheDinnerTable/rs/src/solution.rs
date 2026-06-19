// Solution for Advent of Code 2015 day 13.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

use itertools::Itertools;
use std::collections::HashMap;

// Build happiness[(a, b)] = a's happiness change when seated next to b, plus the
// sorted list of distinct people.
fn parse(input: &str) -> (HashMap<(String, String), i32>, Vec<String>) {
    let mut happiness = HashMap::new();
    let mut people = Vec::new();

    for line in input.trim().lines() {
        let f: Vec<&str> = line.trim_end_matches('.').split_whitespace().collect();
        // <A> would <gain|lose> <N> happiness units by sitting next to <B>
        let (a, b) = (f[0].to_string(), f[10].to_string());
        let n: i32 = f[3].parse::<i32>().unwrap() * if f[2] == "gain" { 1 } else { -1 };
        if !people.contains(&a) {
            people.push(a.clone());
        }
        happiness.insert((a, b), n);
    }

    people.sort();
    (happiness, people)
}

// Max total happiness over all circular seatings. The first person is fixed to
// factor out rotations, so we permute only the rest.
fn best(happiness: &HashMap<(String, String), i32>, people: &[String]) -> i32 {
    let head = &people[0];
    let rest = &people[1..];

    rest.iter()
        .permutations(rest.len())
        .map(|order| {
            let seating: Vec<&String> = std::iter::once(head).chain(order).collect();
            let n = seating.len();
            (0..n)
                .map(|i| {
                    let a = seating[i];
                    let b = seating[(i + 1) % n];
                    let f = |x: &String, y: &String| {
                        *happiness.get(&(x.clone(), y.clone())).unwrap_or(&0)
                    };
                    f(a, b) + f(b, a)
                })
                .sum()
        })
        .max()
        .unwrap()
}

pub fn part_one(input: &str) -> String {
    let (happiness, people) = parse(input);
    best(&happiness, &people).to_string()
}

pub fn part_two(input: &str) -> String {
    let (mut happiness, mut people) = parse(input);

    for p in &people {
        happiness.insert(("me".to_string(), p.clone()), 0);
        happiness.insert((p.clone(), "me".to_string()), 0);
    }
    people.push("me".to_string());

    best(&happiness, &people).to_string()
}
