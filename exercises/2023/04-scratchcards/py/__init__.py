from typing import *
from aocpy import BaseExercise


def _match_counts(instr: str) -> list[int]:
    # Number of winning numbers present on each card, via set intersection.
    counts = []
    for line in instr.splitlines():
        winning, have = line.split(":")[1].split("|")
        counts.append(len(set(winning.split()) & set(have.split())))
    return counts


# Exercise for Advent of Code 2023 day 4.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return sum(1 << (m - 1) for m in _match_counts(instr) if m)

    @staticmethod
    def two(instr: str) -> int:
        matches = _match_counts(instr)
        copies = [1] * len(matches)
        for i, m in enumerate(matches):
            for j in range(i + 1, i + 1 + m):
                copies[j] += copies[i]
        return sum(copies)
