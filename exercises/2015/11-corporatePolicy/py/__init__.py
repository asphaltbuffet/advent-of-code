from typing import *
from aocpy import BaseExercise


def increment(p: List[int]) -> None:
    """Advance the password by one as a base-26 odometer over a..z, in place."""
    for i in range(len(p) - 1, -1, -1):
        if p[i] == 25:  # 'z'
            p[i] = 0
            continue
        p[i] += 1
        break


def valid(p: List[int]) -> bool:
    """Check the three corporate policy rules (p is a list of 0..25 offsets)."""
    # Rule 2: no i (8), o (14), or l (11).
    if any(c in (8, 14, 11) for c in p):
        return False

    # Rule 1: an increasing straight of at least three letters.
    if not any(
        p[i + 1] == p[i] + 1 and p[i + 2] == p[i] + 2 for i in range(len(p) - 2)
    ):
        return False

    # Rule 3: at least two different, non-overlapping pairs.
    pairs = set()
    i = 0
    while i < len(p) - 1:
        if p[i] == p[i + 1]:
            pairs.add(p[i])
            i += 2
        else:
            i += 1
    return len(pairs) >= 2


def next_password(s: str) -> str:
    """Return the next valid password strictly after s."""
    p = [ord(c) - ord("a") for c in s]
    increment(p)
    while not valid(p):
        increment(p)
    return "".join(chr(c + ord("a")) for c in p)


# Exercise for Advent of Code 2015 day 11.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> str:
        return next_password(instr.strip())

    @staticmethod
    def two(instr: str) -> str:
        return next_password(next_password(instr.strip()))
