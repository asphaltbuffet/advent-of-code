// Solution for Advent of Code 2023 day 20.
//
// Modules pass low/high pulses: flip-flops toggle on a low pulse; conjunctions
// remember each input's last pulse and send low only when all are high. Each
// button press drains a FIFO of pulses. Part one multiplies the low and high
// counts over 1000 presses. Part two: rx is fed by one conjunction, which emits
// low only when all its inputs go high; those inputs each cycle, so the answer is
// the LCM of the presses at which each first sends a high pulse.

use std::collections::{HashMap, VecDeque};

#[derive(Clone, Copy, PartialEq)]
enum Kind {
    Broadcaster,
    FlipFlop,
    Conjunction,
}

struct Module<'a> {
    kind: Kind,
    outputs: Vec<&'a str>,
    flip_on: bool,
    memory: HashMap<&'a str, bool>, // conjunction: last pulse per input
}

fn parse(input: &str) -> HashMap<&str, Module> {
    let mut modules: HashMap<&str, Module> = HashMap::new();

    for line in input.lines() {
        let (left, right) = line.split_once(" -> ").unwrap();
        let (kind, name) = match left.as_bytes()[0] {
            b'%' => (Kind::FlipFlop, &left[1..]),
            b'&' => (Kind::Conjunction, &left[1..]),
            _ => (Kind::Broadcaster, left),
        };
        modules.insert(
            name,
            Module {
                kind,
                outputs: right.split(", ").collect(),
                flip_on: false,
                memory: HashMap::new(),
            },
        );
    }

    // Seed each conjunction's memory with its input modules.
    let edges: Vec<(&str, Vec<&str>)> =
        modules.iter().map(|(&n, m)| (n, m.outputs.clone())).collect();
    for (src, outs) in edges {
        for dst in outs {
            if let Some(m) = modules.get_mut(dst) {
                if m.kind == Kind::Conjunction {
                    m.memory.insert(src, false);
                }
            }
        }
    }

    modules
}

/// Run one button press. Returns (low, high) counts and, for each source that
/// sent a high pulse into `watch`, records its name in `fired`.
fn press<'a>(
    modules: &mut HashMap<&'a str, Module<'a>>,
    watch: &str,
    fired: &mut Vec<&'a str>,
) -> (u64, u64) {
    let (mut low, mut high) = (0u64, 0u64);
    let mut queue: VecDeque<(&str, &str, bool)> = VecDeque::new();
    queue.push_back(("button", "broadcaster", false));

    while let Some((src, name, pulse)) = queue.pop_front() {
        if pulse {
            high += 1;
            if name == watch {
                fired.push(src);
            }
        } else {
            low += 1;
        }

        let Some(module) = modules.get_mut(name) else {
            continue; // untyped sink
        };

        let out = match module.kind {
            Kind::FlipFlop => {
                if pulse {
                    continue;
                }
                module.flip_on = !module.flip_on;
                module.flip_on
            }
            Kind::Conjunction => {
                module.memory.insert(src, pulse);
                !module.memory.values().all(|&v| v)
            }
            Kind::Broadcaster => pulse,
        };

        // Collect targets first to release the mutable borrow.
        let targets: Vec<&str> = module.outputs.clone();
        for dst in targets {
            queue.push_back((name, dst, out));
        }
    }

    (low, high)
}

fn gcd(a: u64, b: u64) -> u64 {
    if b == 0 { a } else { gcd(b, a % b) }
}

pub fn part_one(input: &str) -> String {
    let mut modules = parse(input);
    let (mut low, mut high) = (0u64, 0u64);
    let mut sink = Vec::new();
    for _ in 0..1000 {
        let (l, h) = press(&mut modules, "", &mut sink);
        low += l;
        high += h;
    }
    (low * high).to_string()
}

pub fn part_two(input: &str) -> String {
    let mut modules = parse(input);

    // The single conjunction feeding rx, and its input modules.
    let feeder = modules
        .iter()
        .find(|(_, m)| m.outputs.contains(&"rx"))
        .map(|(&n, _)| n)
        .unwrap();
    let feeder_inputs: Vec<&str> = modules[feeder].memory.keys().copied().collect();

    let mut periods: HashMap<&str, u64> = HashMap::new();
    let mut presses = 0u64;
    while periods.len() < feeder_inputs.len() {
        presses += 1;
        let mut fired = Vec::new();
        press(&mut modules, feeder, &mut fired);
        for src in fired {
            periods.entry(src).or_insert(presses);
        }
    }

    periods
        .values()
        .fold(1u64, |acc, &p| acc / gcd(acc, p) * p)
        .to_string()
}
