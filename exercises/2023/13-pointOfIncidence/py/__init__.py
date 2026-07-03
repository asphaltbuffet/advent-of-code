from typing import *
from aocpy import BaseExercise


def _rows_as_bits(block: list[str]) -> list[int]:
    # Each row becomes an integer bitmask of its '#' cells.
    return [int(row.replace(".", "0").replace("#", "1"), 2) for row in block]


def _reflection(lines: list[int], smudges: int) -> int:
    # A mirror after index i is valid when the total bit differences across all
    # folded-out pairs equals `smudges` (0 for a perfect mirror, 1 for one smudge).
    for i in range(1, len(lines)):
        diff = sum(
            bin(a ^ b).count("1")
            for a, b in zip(reversed(lines[:i]), lines[i:])
        )
        if diff == smudges:
            return i
    return 0


def _summary(instr: str, smudges: int) -> int:
    total = 0
    for chunk in instr.split("\n\n"):
        block = chunk.splitlines()
        rows = _rows_as_bits(block)
        cols = _rows_as_bits(["".join(col) for col in zip(*block)])
        total += 100 * _reflection(rows, smudges) + _reflection(cols, smudges)
    return total


# Exercise for Advent of Code 2023 day 13.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _summary(instr, 0)

    @staticmethod
    def two(instr: str) -> int:
        return _summary(instr, 1)
