// Solution for Advent of Code 2015 day 7.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

use std::collections::HashMap;

// Parse "<signal> -> <wire>" lines into a map of wire name -> signal expression.
fn parse(input: &str) -> HashMap<&str, &str> {
    input
        .lines()
        .map(|line| {
            let (signal, wire) = line.split_once(" -> ").unwrap();
            (wire, signal)
        })
        .collect()
}

// Resolve an operand that is either a literal number or a wire name.
fn value<'a>(
    token: &'a str,
    circuit: &HashMap<&'a str, &'a str>,
    cache: &mut HashMap<&'a str, u16>,
) -> u16 {
    match token.parse::<u16>() {
        Ok(n) => n,
        Err(_) => eval(token, circuit, cache),
    }
}

// Memoized evaluation of the signal feeding `wire`.
fn eval<'a>(
    wire: &'a str,
    circuit: &HashMap<&'a str, &'a str>,
    cache: &mut HashMap<&'a str, u16>,
) -> u16 {
    if let Some(&v) = cache.get(wire) {
        return v;
    }

    let expr = circuit[wire];
    let tokens: Vec<&str> = expr.split_whitespace().collect();

    let result = match tokens.as_slice() {
        [x] => value(x, circuit, cache),
        ["NOT", x] => !value(x, circuit, cache),
        [a, "AND", b] => value(a, circuit, cache) & value(b, circuit, cache),
        [a, "OR", b] => value(a, circuit, cache) | value(b, circuit, cache),
        [a, "LSHIFT", b] => value(a, circuit, cache) << value(b, circuit, cache),
        [a, "RSHIFT", b] => value(a, circuit, cache) >> value(b, circuit, cache),
        _ => panic!("unknown signal: {expr}"),
    };

    cache.insert(wire, result);
    result
}

pub fn part_one(input: &str) -> String {
    let circuit = parse(input);
    let mut cache = HashMap::new();
    eval("a", &circuit, &mut cache).to_string()
}

pub fn part_two(input: &str) -> String {
    let circuit = parse(input);
    if !circuit.contains_key("a") {
        return String::new();
    }

    // First pass: solve for "a" as in part one.
    let mut cache = HashMap::new();
    let a = eval("a", &circuit, &mut cache);

    // Second pass: override "b" with that value and re-solve from scratch.
    let mut cache = HashMap::new();
    cache.insert("b", a);
    eval("a", &circuit, &mut cache).to_string()
}
