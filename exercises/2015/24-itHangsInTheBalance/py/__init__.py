from itertools import combinations
from math import prod
from typing import *
from aocpy import BaseExercise


def min_qe(weights: List[int], groups: int) -> int:
    """Lowest quantum entanglement of the smallest first group summing to an
    equal share. Searching by increasing size yields the fewest-packages tier;
    the minimum product within it is the answer."""
    target = sum(weights) // groups
    for size in range(1, len(weights) + 1):
        qes = [prod(c) for c in combinations(weights, size) if sum(c) == target]
        if qes:
            return min(qes)
    return -1


# Exercise for Advent of Code 2015 day 24.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return min_qe([int(x) for x in instr.split()], 3)

    @staticmethod
    def two(instr: str) -> int:
        return min_qe([int(x) for x in instr.split()], 4)
