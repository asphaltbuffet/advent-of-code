// Solution for Advent of Code 2023 day 19.
//
// Parts with categories x,m,a,s flow through named workflows of conditional
// rules until Accepted or Rejected. Part one runs each concrete part and sums the
// ratings of the accepted ones. Part two counts accepted combinations over the
// whole space [1,4000]^4 by recursively splitting the 4D range at each rule's
// threshold: the matching sub-range follows the rule's target, the remainder
// falls through to the next rule.

use std::collections::HashMap;

/// One rule: an optional condition (category index, '<' or '>', threshold) plus
/// the target to jump to. A fallthrough rule has no condition.
struct Rule<'a> {
    cond: Option<(usize, u8, u64)>,
    target: &'a str,
}

fn cat_index(c: u8) -> usize {
    match c {
        b'x' => 0,
        b'm' => 1,
        b'a' => 2,
        b's' => 3,
        _ => unreachable!(),
    }
}

fn parse_workflows(block: &str) -> HashMap<&str, Vec<Rule>> {
    block
        .lines()
        .map(|line| {
            let (name, rest) = line.split_once('{').unwrap();
            let body = &rest[..rest.len() - 1]; // strip trailing '}'
            let rules = body
                .split(',')
                .map(|tok| match tok.split_once(':') {
                    Some((cond, target)) => {
                        let bytes = cond.as_bytes();
                        Rule {
                            cond: Some((
                                cat_index(bytes[0]),
                                bytes[1],
                                cond[2..].parse().unwrap(),
                            )),
                            target,
                        }
                    }
                    None => Rule { cond: None, target: tok },
                })
                .collect();
            (name, rules)
        })
        .collect()
}

fn accepts(part: &[u64; 4], workflows: &HashMap<&str, Vec<Rule>>) -> bool {
    let mut name = "in";
    while name != "A" && name != "R" {
        for rule in &workflows[name] {
            let matched = match rule.cond {
                None => true,
                Some((cat, b'<', t)) => part[cat] < t,
                Some((cat, _, t)) => part[cat] > t,
            };
            if matched {
                name = rule.target;
                break;
            }
        }
    }
    name == "A"
}

/// Accepted combinations within `ranges` (inclusive per category) at `name`.
fn count(mut ranges: [(u64, u64); 4], name: &str, workflows: &HashMap<&str, Vec<Rule>>) -> u64 {
    match name {
        "R" => return 0,
        "A" => return ranges.iter().map(|&(lo, hi)| hi - lo + 1).product(),
        _ => {}
    }

    let mut total = 0;
    for rule in &workflows[name] {
        let (cat, op, t) = match rule.cond {
            None => {
                total += count(ranges, rule.target, workflows);
                break;
            }
            Some(c) => c,
        };

        let (lo, hi) = ranges[cat];
        // Split [lo, hi] into the matching sub-range and the remainder.
        let (matched, rest) = if op == b'<' {
            ((lo, hi.min(t - 1)), (lo.max(t), hi))
        } else {
            ((lo.max(t + 1), hi), (lo, hi.min(t)))
        };

        if matched.0 <= matched.1 {
            let mut sub = ranges;
            sub[cat] = matched;
            total += count(sub, rule.target, workflows);
        }
        if rest.0 > rest.1 {
            break; // nothing falls through
        }
        ranges[cat] = rest;
    }
    total
}

pub fn part_one(input: &str) -> String {
    let (wf_block, part_block) = input.split_once("\n\n").unwrap();
    let workflows = parse_workflows(wf_block);

    part_block
        .lines()
        .filter_map(|line| {
            // "{x=787,m=2655,a=1222,s=2876}"
            let mut part = [0u64; 4];
            for (i, field) in line.trim_matches(['{', '}']).split(',').enumerate() {
                part[i] = field.split_once('=').unwrap().1.parse().unwrap();
            }
            accepts(&part, &workflows).then(|| part.iter().sum::<u64>())
        })
        .sum::<u64>()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    let (wf_block, _) = input.split_once("\n\n").unwrap();
    let workflows = parse_workflows(wf_block);
    count([(1, 4000); 4], "in", &workflows).to_string()
}
