from typing import *
from aocpy import BaseExercise

from itertools import groupby


def look_and_say(s: str) -> str:
    """Apply one round: each run of identical digits becomes <count><digit>."""
    return "".join(f"{len(list(g))}{d}" for d, g in groupby(s))


def iterate(instr: str, n: int) -> int:
    """Apply look_and_say n times and return the resulting length."""
    s = instr.strip()
    for _ in range(n):
        s = look_and_say(s)
    return len(s)


# Exercise for Advent of Code 2015 day 10.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return iterate(instr, 40)

    @staticmethod
    def two(instr: str) -> int:
        return iterate(instr, 50)
