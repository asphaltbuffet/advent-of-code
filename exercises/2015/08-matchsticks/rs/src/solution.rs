// Solution for Advent of Code 2015 day 8.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

// For one quoted string literal, return (code chars, in-memory chars).
fn decode_len(line: &str) -> (usize, usize) {
    let b = line.as_bytes();
    let code = b.len();

    // Walk the characters between the surrounding quotes, counting decoded chars.
    let mut mem = 0;
    let mut i = 1; // skip opening quote
    while i < b.len() - 1 {
        if b[i] == b'\\' {
            match b[i + 1] {
                b'x' => i += 4, // \xNN -> 1 char
                _ => i += 2,    // \\ or \" -> 1 char
            }
        } else {
            i += 1;
        }
        mem += 1;
    }

    (code, mem)
}

// For one line, return the length of its re-encoded (quoted) form.
fn encode_len(line: &str) -> usize {
    // 2 for the new surrounding quotes, +1 for each " or \ that must be escaped.
    line.len() + 2 + line.bytes().filter(|&c| c == b'"' || c == b'\\').count()
}

pub fn part_one(input: &str) -> String {
    input
        .lines()
        .map(|l| {
            let (code, mem) = decode_len(l);
            code - mem
        })
        .sum::<usize>()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    input
        .lines()
        .map(|l| encode_len(l) - l.len())
        .sum::<usize>()
        .to_string()
}
