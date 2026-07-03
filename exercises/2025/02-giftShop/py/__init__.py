from typing import *
from aocpy import BaseExercise


def _ranges(instr: str) -> Iterator[tuple[int, int]]:
    for chunk in instr.strip().split(","):
        lo, hi = chunk.split("-")
        yield int(lo), int(hi)


def _invalid_ids(lo: int, hi: int) -> Iterator[int]:
    # An "invalid" ID has an even digit count whose two halves are identical,
    # i.e. some half repeated exactly twice. Generate those directly rather than
    # scanning the range.
    width = len(str(hi))
    for half in range(1, width // 2 + 1):
        for h in range(10 ** (half - 1), 10**half):
            n = int(str(h) * 2)
            if lo <= n <= hi:
                yield n


def _repeated_ids(lo: int, hi: int) -> set[int]:
    # A "repeated" ID is any pattern of `p` digits tiled `r >= 2` times. Build
    # each candidate straight from its pattern; a set folds away the numbers that
    # are reachable through more than one (pattern, repeat) factoring.
    width = len(str(hi))
    found: set[int] = set()
    for p in range(1, width // 2 + 1):
        for r in range(2, width // p + 1):
            for q in range(10 ** (p - 1), 10**p):
                n = int(str(q) * r)
                if n > hi:
                    break
                if n >= lo:
                    found.add(n)
    return found


# Exercise for Advent of Code 2025 day 2.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return sum(n for lo, hi in _ranges(instr) for n in _invalid_ids(lo, hi))

    @staticmethod
    def two(instr: str) -> int:
        return sum(sum(_repeated_ids(lo, hi)) for lo, hi in _ranges(instr))
