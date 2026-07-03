// Solution for Advent of Code 2025 day 5.
//
// The inventory is a set of numeric ranges; the ingredients are IDs. Part one
// counts IDs that fall inside some range, part two counts the integers the ranges
// cover. Both reduce to the same interval merge — sort by low, coalesce overlaps.
// Only genuinely overlapping ranges merge (`high >= next.low`); a one-unit gap
// keeps them separate.

/// Parse the two blocks into merged ranges and the ingredient IDs.
fn parse(input: &str) -> (Vec<(i64, i64)>, Vec<i64>) {
    let (inventory, ingredients) = input.split_once("\n\n").unwrap();

    let mut ranges: Vec<(i64, i64)> = inventory
        .lines()
        .map(|line| {
            let (lo, hi) = line.split_once('-').unwrap();
            (lo.parse().unwrap(), hi.parse().unwrap())
        })
        .collect();
    ranges.sort_unstable();

    let mut merged: Vec<(i64, i64)> = Vec::with_capacity(ranges.len());
    for (lo, hi) in ranges {
        match merged.last_mut() {
            Some(last) if last.1 >= lo => last.1 = last.1.max(hi),
            _ => merged.push((lo, hi)),
        }
    }

    let ids = ingredients.lines().map(|l| l.parse().unwrap()).collect();
    (merged, ids)
}

pub fn part_one(input: &str) -> String {
    let (ranges, ids) = parse(input);
    // The merged ranges are sorted and disjoint, so a binary search on the low
    // bounds finds the only range that could contain each ID.
    let fresh = ids
        .iter()
        .filter(|&&id| {
            let idx = ranges.partition_point(|&(lo, _)| lo <= id);
            idx > 0 && ranges[idx - 1].1 >= id
        })
        .count();
    fresh.to_string()
}

pub fn part_two(input: &str) -> String {
    let (ranges, _) = parse(input);
    ranges
        .iter()
        .map(|(lo, hi)| hi - lo + 1)
        .sum::<i64>()
        .to_string()
}
