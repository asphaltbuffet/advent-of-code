/// Fully collapse the polymer with a stack: push each unit, but when it is the
/// same letter as the current top in the opposite case, pop instead — one O(n)
/// pass captures the cascade. Bytes whose lowercase form equals `skip` are
/// ignored; `skip == 0` keeps every unit.
fn react(polymer: &[u8], skip: u8) -> usize {
    let mut stack: Vec<u8> = Vec::with_capacity(polymer.len());

    for &unit in polymer {
        if skip != 0 && unit.to_ascii_lowercase() == skip {
            continue;
        }
        match stack.last() {
            Some(&top) if top != unit && top.eq_ignore_ascii_case(&unit) => {
                stack.pop();
            }
            _ => stack.push(unit),
        }
    }

    stack.len()
}

pub fn part_one(input: &str) -> String {
    react(input.trim().as_bytes(), 0).to_string()
}

pub fn part_two(input: &str) -> String {
    let polymer = input.trim().as_bytes();

    (b'a'..=b'z')
        .map(|unit| react(polymer, unit))
        .min()
        .unwrap_or(polymer.len())
        .to_string()
}
