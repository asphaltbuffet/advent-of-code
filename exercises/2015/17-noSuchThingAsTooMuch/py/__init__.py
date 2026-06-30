from itertools import combinations
from typing import *
from aocpy import BaseExercise


def parse(instr: str) -> Tuple[List[int], int]:
    """Container sizes plus the target volume. The target isn't in the input,
    so it's inferred: the 5-container AoC example fills 25 liters, the real 150.
    """
    sizes = [int(x) for x in instr.split()]
    target = 25 if len(sizes) <= 5 else 150
    return sizes, target


def size_counts(sizes: List[int], target: int) -> List[int]:
    """counts[k] = number of k-container subsets that sum to exactly target."""
    counts = [0] * (len(sizes) + 1)
    for k in range(1, len(sizes) + 1):
        for combo in combinations(sizes, k):
            if sum(combo) == target:
                counts[k] += 1
    return counts


# Exercise for Advent of Code 2015 day 17.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        sizes, target = parse(instr)
        return sum(size_counts(sizes, target))

    @staticmethod
    def two(instr: str) -> int:
        sizes, target = parse(instr)
        # The first non-zero size bucket is the minimum-container count.
        return next(c for c in size_counts(sizes, target) if c > 0)
