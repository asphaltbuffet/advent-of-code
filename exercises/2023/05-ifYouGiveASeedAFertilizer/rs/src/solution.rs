// Solution for Advent of Code 2023 day 5.
//
// The almanac is a chain of maps, each a list of (dst, src, len) rules that
// relocate a contiguous source range onto a destination. Part one walks each
// seed value through the chain and takes the minimum. Part two treats the seeds
// as half-open intervals and pushes whole intervals through each map, splitting
// them against overlapping rules — so billion-wide ranges resolve without
// enumerating a single seed.

struct Rule {
    dst: u64,
    src: u64,
    len: u64,
}

struct Almanac {
    seeds: Vec<u64>,
    maps: Vec<Vec<Rule>>,
}

fn parse(input: &str) -> Almanac {
    let mut blocks = input.split("\n\n");
    let seeds = blocks
        .next()
        .unwrap()
        .split_once(':')
        .unwrap()
        .1
        .split_whitespace()
        .map(|n| n.parse().unwrap())
        .collect();

    let maps = blocks
        .map(|block| {
            block
                .lines()
                .skip(1)
                .map(|line| {
                    let mut it = line.split_whitespace().map(|n| n.parse().unwrap());
                    Rule {
                        dst: it.next().unwrap(),
                        src: it.next().unwrap(),
                        len: it.next().unwrap(),
                    }
                })
                .collect()
        })
        .collect();

    Almanac { seeds, maps }
}

fn map_value(value: u64, rules: &[Rule]) -> u64 {
    rules
        .iter()
        .find(|r| r.src <= value && value < r.src + r.len)
        .map(|r| r.dst + (value - r.src))
        .unwrap_or(value)
}

/// Push half-open intervals [start, end) through one map, splitting on overlap.
fn map_ranges(ranges: Vec<(u64, u64)>, rules: &[Rule]) -> Vec<(u64, u64)> {
    let mut out = Vec::new();
    let mut work = ranges;

    while let Some((start, end)) = work.pop() {
        let mut mapped = false;
        for r in rules {
            let lo = start.max(r.src);
            let hi = end.min(r.src + r.len);
            if lo < hi {
                out.push((lo - r.src + r.dst, hi - r.src + r.dst));
                if start < lo {
                    work.push((start, lo));
                }
                if hi < end {
                    work.push((hi, end));
                }
                mapped = true;
                break;
            }
        }
        if !mapped {
            out.push((start, end)); // identity for the untouched interval
        }
    }

    out
}

pub fn part_one(input: &str) -> String {
    let almanac = parse(input);
    almanac
        .seeds
        .iter()
        .map(|&seed| {
            almanac
                .maps
                .iter()
                .fold(seed, |value, rules| map_value(value, rules))
        })
        .min()
        .unwrap()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    let almanac = parse(input);
    let mut ranges: Vec<(u64, u64)> = almanac
        .seeds
        .chunks_exact(2)
        .map(|c| (c[0], c[0] + c[1]))
        .collect();

    for rules in &almanac.maps {
        ranges = map_ranges(ranges, rules);
    }

    ranges.iter().map(|&(start, _)| start).min().unwrap().to_string()
}
