from functools import cache
from typing import *
from aocpy import BaseExercise


@cache
def _count(springs: str, groups: tuple[int, ...]) -> int:
    # Memoized DP over (remaining pattern, remaining group sizes).
    if not groups:
        # Valid only if no forced '#' remain.
        return 0 if "#" in springs else 1
    if len(springs) < sum(groups) + len(groups) - 1:
        return 0  # not enough room for the groups plus separators

    total = 0
    first = springs[0]
    # Option A: treat the first char as operational and skip it.
    if first in ".?":
        total += _count(springs[1:], groups)
    # Option B: place the next group of '#' starting here.
    n = groups[0]
    if (
        "." not in springs[:n]          # the run has no forced gap
        and (len(springs) == n or springs[n] != "#")  # boundary after the run
    ):
        total += _count(springs[n + 1:], groups[1:])
    return total


def _arrangements(instr: str, unfold: int) -> int:
    total = 0
    for line in instr.splitlines():
        pattern, group_str = line.split()
        groups = tuple(int(n) for n in group_str.split(","))
        pattern = "?".join([pattern] * unfold)
        groups = groups * unfold
        total += _count(pattern, groups)
    return total


# Exercise for Advent of Code 2023 day 12.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _arrangements(instr, 1)

    @staticmethod
    def two(instr: str) -> int:
        return _arrangements(instr, 5)
