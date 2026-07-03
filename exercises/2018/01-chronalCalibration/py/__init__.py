from itertools import accumulate, cycle
from typing import *

from aocpy import BaseExercise


# Exercise for Advent of Code 2018 day 1.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return sum(int(n) for n in instr.split())

    @staticmethod
    def two(instr: str) -> int:
        changes = [int(n) for n in instr.split()]
        seen = {0}
        for freq in accumulate(cycle(changes)):
            if freq in seen:
                return freq
            seen.add(freq)
