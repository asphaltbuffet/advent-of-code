//! Advent of Code 2018 day 16: Chronal Classification.

type Regs = [i64; 4];

/// The sixteen operations, each producing the value to store in register C from
/// the current registers and the two inputs A and B.
const OPS: [fn(&Regs, i64, i64) -> i64; 16] = [
    |r, a, b| r[a as usize] + r[b as usize],           // addr
    |r, a, b| r[a as usize] + b,                       // addi
    |r, a, b| r[a as usize] * r[b as usize],           // mulr
    |r, a, b| r[a as usize] * b,                       // muli
    |r, a, b| r[a as usize] & r[b as usize],           // banr
    |r, a, b| r[a as usize] & b,                       // bani
    |r, a, b| r[a as usize] | r[b as usize],           // borr
    |r, a, b| r[a as usize] | b,                       // bori
    |r, a, _| r[a as usize],                           // setr
    |_, a, _| a,                                       // seti
    |r, a, b| (a > r[b as usize]) as i64,              // gtir
    |r, a, b| (r[a as usize] > b) as i64,              // gtri
    |r, a, b| (r[a as usize] > r[b as usize]) as i64,  // gtrr
    |r, a, b| (a == r[b as usize]) as i64,             // eqir
    |r, a, b| (r[a as usize] == b) as i64,             // eqri
    |r, a, b| (r[a as usize] == r[b as usize]) as i64, // eqrr
];

/// Apply one operation to a copy of the registers, returning the result.
fn apply(op: fn(&Regs, i64, i64) -> i64, mut r: Regs, ins: &[i64; 4]) -> Regs {
    r[ins[3] as usize] = op(&r, ins[1], ins[2]);
    r
}

/// A Before/instruction/After observation.
struct Sample {
    before: Regs,
    ins: [i64; 4],
    after: Regs,
}

impl Sample {
    /// Bitmask of the operation indices whose effect reproduces `after`.
    fn matches(&self) -> u16 {
        (0..16)
            .filter(|&i| apply(OPS[i], self.before, &self.ins) == self.after)
            .fold(0u16, |m, i| m | (1 << i))
    }
}

fn ints(s: &str) -> Vec<i64> {
    s.split(|c: char| !c.is_ascii_digit() && c != '-')
        .filter_map(|t| t.parse().ok())
        .collect()
}

fn four(v: &[i64]) -> [i64; 4] {
    [v[0], v[1], v[2], v[3]]
}

/// Split the input into observed samples and the test program.
fn parse(input: &str) -> (Vec<Sample>, Vec<[i64; 4]>) {
    let input = input.replace("\r\n", "\n");
    let (head, tail) = input.split_once("\n\n\n").unwrap_or((&input, ""));

    let samples = head
        .split("\n\n")
        .filter_map(|blk| {
            let lines: Vec<&str> = blk.trim().lines().collect();
            if lines.len() < 3 {
                return None;
            }
            Some(Sample {
                before: four(&ints(lines[0])),
                ins: four(&ints(lines[1])),
                after: four(&ints(lines[2])),
            })
        })
        .collect();

    let program = tail
        .lines()
        .filter(|l| !l.trim().is_empty())
        .map(|l| four(&ints(l)))
        .collect();

    (samples, program)
}

pub fn part_one(input: &str) -> String {
    let (samples, _) = parse(input);
    samples
        .iter()
        .filter(|s| s.matches().count_ones() >= 3)
        .count()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    let (samples, program) = parse(input);

    // Each opcode number's candidate operations start as all sixteen, then are
    // narrowed to those consistent with every sample that uses the number.
    let mut candidates = [0xFFFFu16; 16];
    for s in &samples {
        candidates[s.ins[0] as usize] &= s.matches();
    }

    // Resolve by elimination: an opcode with a single candidate is fixed and that
    // operation is struck from every other opcode, repeating until all are pinned.
    let mut op_for = [0usize; 16];
    let mut assigned = 0u16;
    while assigned.count_ones() < 16 {
        for opcode in 0..16 {
            if assigned & (1 << opcode) != 0 || candidates[opcode].count_ones() != 1 {
                continue;
            }
            let op = candidates[opcode].trailing_zeros() as usize;
            op_for[opcode] = op;
            assigned |= 1 << opcode;
            for other in 0..16 {
                if other != opcode {
                    candidates[other] &= !(1 << op);
                }
            }
        }
    }

    let mut r: Regs = [0; 4];
    for ins in &program {
        r = apply(OPS[op_for[ins[0] as usize]], r, ins);
    }
    r[0].to_string()
}
