import re
from typing import *
from aocpy import BaseExercise

_MUL = re.compile(r"mul\((\d{1,3}),(\d{1,3})\)")
_INSTR = re.compile(r"mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)")


# Exercise for Advent of Code 2024 day 3.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return sum(int(a) * int(b) for a, b in _MUL.findall(instr))

    @staticmethod
    def two(instr: str) -> int:
        total = 0
        enabled = True
        for m in _INSTR.finditer(instr):
            tok = m.group(0)
            if tok == "do()":
                enabled = True
            elif tok == "don't()":
                enabled = False
            elif enabled:
                total += int(m.group(1)) * int(m.group(2))
        return total
