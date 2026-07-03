from collections import Counter
from itertools import combinations

from aocpy import BaseExercise


# Exercise for Advent of Code 2018 day 2.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        counts = [Counter(box_id).values() for box_id in instr.split()]
        twos = sum(1 for c in counts if 2 in c)
        threes = sum(1 for c in counts if 3 in c)
        return twos * threes

    @staticmethod
    def two(instr: str) -> str:
        for a, b in combinations(instr.split(), 2):
            if len(a) != len(b):
                continue
            common = [x for x, y in zip(a, b) if x == y]
            if len(common) == len(a) - 1:
                return "".join(common)
        return ""
