import re

from typing import *
from aocpy import BaseExercise

# The sixteen operations, each producing the value stored in register C from the
# current registers r and the two inputs a and b.
OPS: dict[str, Callable[[list[int], int, int], int]] = {
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


def _apply(fn: Callable[[list[int], int, int], int], r: list[int], ins: list[int]) -> list[int]:
    out = r[:]
    out[ins[3]] = fn(r, ins[1], ins[2])
    return out


def _matches(before: list[int], ins: list[int], after: list[int]) -> set[str]:
    """Names of the operations whose effect reproduces ``after``."""
    return {name for name, fn in OPS.items() if _apply(fn, before, ins) == after}


def _parse(instr: str) -> tuple[list[tuple[list[int], list[int], list[int]]], list[list[int]]]:
    """Split the input into observed samples and the test program."""
    head, _, tail = instr.replace("\r\n", "\n").partition("\n\n\n")

    samples = []
    for blk in head.split("\n\n"):
        lines = blk.strip().splitlines()
        if len(lines) < 3:
            continue
        before, ins, after = (list(map(int, re.findall(r"-?\d+", ln))) for ln in lines[:3])
        samples.append((before, ins, after))

    program = [list(map(int, re.findall(r"-?\d+", ln))) for ln in tail.splitlines() if ln.strip()]
    return samples, program


# Exercise for Advent of Code 2018 day 16.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        samples, _ = _parse(instr)
        return sum(len(_matches(*s)) >= 3 for s in samples)

    @staticmethod
    def two(instr: str) -> int:
        samples, program = _parse(instr)

        # Each opcode number's candidates are the operations consistent with every
        # sample that uses that number.
        candidates: list[set[str]] = [set(OPS) for _ in range(16)]
        for before, ins, after in samples:
            candidates[ins[0]] &= _matches(before, ins, after)

        # Resolve by elimination: an opcode with a single candidate is fixed and
        # that operation is struck from the others, until all sixteen are pinned.
        op_for: dict[int, str] = {}
        while len(op_for) < 16:
            for opcode, cand in enumerate(candidates):
                if opcode not in op_for and len(cand) == 1:
                    (name,) = cand
                    op_for[opcode] = name
                    for other in range(16):
                        if other != opcode:
                            candidates[other].discard(name)

        r = [0, 0, 0, 0]
        for ins in program:
            r = _apply(OPS[op_for[ins[0]]], r, ins)
        return r[0]
