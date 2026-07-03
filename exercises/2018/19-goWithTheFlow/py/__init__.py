from math import isqrt

from aocpy import BaseExercise

# Each opcode as a single expression over registers r and operands a, b.
OPS = {
    "addr": lambda r, a, b: r[a] + r[b],
    "addi": lambda r, a, b: r[a] + b,
    "mulr": lambda r, a, b: r[a] * r[b],
    "muli": lambda r, a, b: r[a] * b,
    "banr": lambda r, a, b: r[a] & r[b],
    "bani": lambda r, a, b: r[a] & b,
    "borr": lambda r, a, b: r[a] | r[b],
    "bori": lambda r, a, b: r[a] | b,
    "setr": lambda r, a, b: r[a],
    "seti": lambda r, a, b: a,
    "gtir": lambda r, a, b: int(a > r[b]),
    "gtri": lambda r, a, b: int(r[a] > b),
    "gtrr": lambda r, a, b: int(r[a] > r[b]),
    "eqir": lambda r, a, b: int(a == r[b]),
    "eqri": lambda r, a, b: int(r[a] == b),
    "eqrr": lambda r, a, b: int(r[a] == r[b]),
}


def _parse(instr: str) -> tuple[int, list[tuple[str, int, int, int]]]:
    ip_reg = 0
    prog: list[tuple[str, int, int, int]] = []
    for line in instr.strip().splitlines():
        f = line.split()
        if f[0] == "#ip":
            ip_reg = int(f[1])
        else:
            prog.append((f[0], int(f[1]), int(f[2]), int(f[3])))
    return ip_reg, prog


def _run(ip_reg, prog, regs, max_steps=None) -> list[int]:
    ip = 0
    step = 0
    while 0 <= ip < len(prog):
        if max_steps is not None and step >= max_steps:
            break
        regs[ip_reg] = ip
        op, a, b, c = prog[ip]
        regs[c] = OPS[op](regs, a, b)
        ip = regs[ip_reg] + 1
        step += 1
    return regs


def _sum_divisors(n: int) -> int:
    total = 0
    for d in range(1, isqrt(n) + 1):
        if n % d == 0:
            total += d
            if d != n // d:
                total += n // d
    return total


# Exercise for Advent of Code 2018 day 19.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        ip_reg, prog = _parse(instr)
        return _run(ip_reg, prog, [0] * 6)[0]

    @staticmethod
    def two(instr: str) -> int:
        ip_reg, prog = _parse(instr)
        # The program builds a large number, then sums its divisors with an O(n^2)
        # loop. Run just long enough to assemble that number, then compute the
        # divisor sum directly.
        regs = _run(ip_reg, prog, [1, 0, 0, 0, 0, 0], max_steps=1000)
        return _sum_divisors(max(regs))
