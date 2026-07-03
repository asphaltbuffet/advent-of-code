// Solution for Advent of Code 2018 day 19.
//
// The device from day 16 with one register bound to the instruction pointer, so
// the program can jump and loop. Part one runs it. Part two's program is a naive
// O(n^2) sum-of-divisors over a large number; running just long enough to build
// that number, then summing its divisors directly, is instant.

struct Instr {
    op: String,
    a: usize,
    b: usize,
    c: usize,
}

fn parse(input: &str) -> (usize, Vec<Instr>) {
    let mut ip_reg = 0;
    let mut prog = Vec::new();
    for line in input.trim().lines() {
        let f: Vec<&str> = line.split_whitespace().collect();
        if f[0] == "#ip" {
            ip_reg = f[1].parse().unwrap();
            continue;
        }
        prog.push(Instr {
            op: f[0].to_string(),
            a: f[1].parse().unwrap(),
            b: f[2].parse().unwrap(),
            c: f[3].parse().unwrap(),
        });
    }
    (ip_reg, prog)
}

fn apply(i: &Instr, r: &mut [i64; 6]) {
    let (a, b, c) = (i.a, i.b, i.c);
    r[c] = match i.op.as_str() {
        "addr" => r[a] + r[b],
        "addi" => r[a] + b as i64,
        "mulr" => r[a] * r[b],
        "muli" => r[a] * b as i64,
        "banr" => r[a] & r[b],
        "bani" => r[a] & b as i64,
        "borr" => r[a] | r[b],
        "bori" => r[a] | b as i64,
        "setr" => r[a],
        "seti" => a as i64,
        "gtir" => (a as i64 > r[b]) as i64,
        "gtri" => (r[a] > b as i64) as i64,
        "gtrr" => (r[a] > r[b]) as i64,
        "eqir" => (a as i64 == r[b]) as i64,
        "eqri" => (r[a] == b as i64) as i64,
        "eqrr" => (r[a] == r[b]) as i64,
        _ => unreachable!(),
    };
}

// run executes from the given registers, stopping after max_steps if it is Some.
fn run(ip_reg: usize, prog: &[Instr], mut regs: [i64; 6], max_steps: Option<usize>) -> [i64; 6] {
    let mut ip = 0i64;
    let mut steps = 0;
    while ip >= 0 && (ip as usize) < prog.len() {
        if let Some(m) = max_steps {
            if steps >= m {
                break;
            }
        }
        regs[ip_reg] = ip;
        apply(&prog[ip as usize], &mut regs);
        ip = regs[ip_reg] + 1;
        steps += 1;
    }
    regs
}

fn sum_divisors(n: i64) -> i64 {
    let mut sum = 0;
    let mut d = 1;
    while d * d <= n {
        if n % d == 0 {
            sum += d;
            if d != n / d {
                sum += n / d;
            }
        }
        d += 1;
    }
    sum
}

pub fn part_one(input: &str) -> String {
    let (ip_reg, prog) = parse(input);
    run(ip_reg, &prog, [0; 6], None)[0].to_string()
}

pub fn part_two(input: &str) -> String {
    let (ip_reg, prog) = parse(input);
    // Build the large number, then compute its divisor sum directly.
    let regs = run(ip_reg, &prog, [1, 0, 0, 0, 0, 0], Some(1000));
    let n = *regs.iter().max().unwrap();
    sum_divisors(n).to_string()
}
