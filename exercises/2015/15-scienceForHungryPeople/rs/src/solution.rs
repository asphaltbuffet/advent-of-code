// Solution for Advent of Code 2015 day 15.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

use regex::Regex;

const TEASPOONS: i64 = 100;

// Each ingredient: [capacity, durability, flavor, texture, calories].
fn parse(input: &str) -> Vec<[i64; 5]> {
    let re = Regex::new(r"-?\d+").unwrap();
    input
        .trim()
        .lines()
        .map(|line| {
            let mut props = [0i64; 5];
            for (i, m) in re.find_iter(line).enumerate() {
                props[i] = m.as_str().parse().unwrap();
            }
            props
        })
        .collect()
}

// score multiplies the per-property totals (negatives clamped to 0). If
// calorie_goal is Some and the recipe misses it, the recipe scores 0.
fn score(ingredients: &[[i64; 5]], amounts: &[i64], calorie_goal: Option<i64>) -> i64 {
    let mut totals = [0i64; 5];
    for (ing, &a) in ingredients.iter().zip(amounts) {
        for p in 0..5 {
            totals[p] += ing[p] * a;
        }
    }
    if let Some(goal) = calorie_goal {
        if totals[4] != goal {
            return 0;
        }
    }
    totals[..4].iter().map(|&t| t.max(0)).product()
}

// best walks every distribution of TEASPOONS across the ingredients, returning
// the highest score. Recursion enumerates compositions of 100 into len parts.
fn best(ingredients: &[[i64; 5]], calorie_goal: Option<i64>) -> i64 {
    let n = ingredients.len();
    let mut amounts = vec![0i64; n];
    let mut highest = 0;

    fn rec(
        i: usize,
        remaining: i64,
        ingredients: &[[i64; 5]],
        amounts: &mut [i64],
        calorie_goal: Option<i64>,
        highest: &mut i64,
    ) {
        if i == amounts.len() - 1 {
            amounts[i] = remaining;
            let s = score(ingredients, amounts, calorie_goal);
            if s > *highest {
                *highest = s;
            }
            return;
        }
        for a in 0..=remaining {
            amounts[i] = a;
            rec(i + 1, remaining - a, ingredients, amounts, calorie_goal, highest);
        }
    }

    rec(0, TEASPOONS, ingredients, &mut amounts, calorie_goal, &mut highest);
    highest
}

pub fn part_one(input: &str) -> String {
    best(&parse(input), None).to_string()
}

pub fn part_two(input: &str) -> String {
    best(&parse(input), Some(500)).to_string()
}
