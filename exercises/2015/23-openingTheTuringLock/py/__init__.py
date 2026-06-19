from typing import *
from aocpy import BaseExercise


def run(program: str, start_a: int) -> int:
    """Execute the assembly with register a = start_a; return final b. A program
    counter outside the program halts. Note jio jumps if the register is *one*,
    not if it is odd."""
    prog = [line.replace(",", "").split() for line in program.strip().splitlines()]
    regs = {"a": start_a, "b": 0}
    pc = 0
    while 0 <= pc < len(prog):
        op, *args = prog[pc]
        if op == "hlf":
            regs[args[0]] //= 2
            pc += 1
        elif op == "tpl":
            regs[args[0]] *= 3
            pc += 1
        elif op == "inc":
            regs[args[0]] += 1
            pc += 1
        elif op == "jmp":
            pc += int(args[0])
        elif op == "jie":
            pc += int(args[1]) if regs[args[0]] % 2 == 0 else 1
        elif op == "jio":
            pc += int(args[1]) if regs[args[0]] == 1 else 1
    return regs["b"]


# Exercise for Advent of Code 2015 day 23.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return run(instr, 0)

    @staticmethod
    def two(instr: str) -> int:
        return run(instr, 1)
