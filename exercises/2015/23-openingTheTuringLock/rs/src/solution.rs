// Solution for Advent of Code 2015 day 23.
//
// Implement part_one and part_two. Each receives the puzzle input as a &str and
// returns the answer as a String. elf's generated harness calls these and
// handles the wire protocol, timing, and panic reporting — you only edit this
// file.

// run executes the assembly with register a = start_a and returns final b. A
// program counter outside the program halts. jio jumps if the register equals
// one (not if it is odd).
fn run(input: &str, start_a: i64) -> i64 {
    let prog: Vec<Vec<&str>> = input
        .trim()
        .lines()
        .map(|l| l.split([',', ' ']).filter(|s| !s.is_empty()).collect())
        .collect();

    let mut regs = [start_a, 0i64]; // [a, b]
    let ri = |s: &str| if s == "b" { 1 } else { 0 };
    let mut pc: i64 = 0;

    while pc >= 0 && (pc as usize) < prog.len() {
        let ins = &prog[pc as usize];
        match ins[0] {
            "hlf" => {
                regs[ri(ins[1])] /= 2;
                pc += 1;
            }
            "tpl" => {
                regs[ri(ins[1])] *= 3;
                pc += 1;
            }
            "inc" => {
                regs[ri(ins[1])] += 1;
                pc += 1;
            }
            "jmp" => pc += ins[1].parse::<i64>().unwrap(),
            "jie" => {
                if regs[ri(ins[1])] % 2 == 0 {
                    pc += ins[2].parse::<i64>().unwrap();
                } else {
                    pc += 1;
                }
            }
            "jio" => {
                if regs[ri(ins[1])] == 1 {
                    pc += ins[2].parse::<i64>().unwrap();
                } else {
                    pc += 1;
                }
            }
            _ => pc += 1,
        }
    }

    regs[1]
}

pub fn part_one(input: &str) -> String {
    run(input, 0).to_string()
}

pub fn part_two(input: &str) -> String {
    run(input, 1).to_string()
}
