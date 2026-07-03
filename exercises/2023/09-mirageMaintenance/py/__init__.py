from itertools import pairwise
from typing import *
from aocpy import BaseExercise


def _extrapolate(seq: list[int]) -> int:
    # The next value is this row's last plus the extrapolation of its difference
    # row; a constant (all-zero difference) row contributes nothing.
    if any(seq):
        diffs = [b - a for a, b in pairwise(seq)]
        return seq[-1] + _extrapolate(diffs)
    return 0


# Exercise for Advent of Code 2023 day 9.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return sum(
            _extrapolate([int(n) for n in line.split()])
            for line in instr.splitlines()
        )

    @staticmethod
    def two(instr: str) -> int:
        # Extrapolating backwards is just extrapolating the reversed sequence.
        return sum(
            _extrapolate([int(n) for n in line.split()][::-1])
            for line in instr.splitlines()
        )
