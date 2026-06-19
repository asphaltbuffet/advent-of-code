import re
from typing import *
from aocpy import BaseExercise

FIRST = 20151125
MULT = 252533
MOD = 33554393


# Exercise for Advent of Code 2015 day 25.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        row, col = (int(x) for x in re.findall(r"\d+", instr))
        # Cell (row, col) is the n-th code filled in diagonal order.
        diag = row + col - 2
        n = diag * (diag + 1) // 2 + col
        # n-th code = FIRST * MULT^(n-1) mod MOD (3-arg pow does modular exp).
        return FIRST * pow(MULT, n - 1, MOD) % MOD

    @staticmethod
    def two(instr: str) -> str:
        # Day 25 has no part 2 — the final star comes from finishing the rest.
        return "Merry Christmas!"
