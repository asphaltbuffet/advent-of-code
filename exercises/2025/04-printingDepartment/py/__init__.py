from typing import *
from aocpy import BaseExercise

_ADJ = [(-1, -1), (-1, 0), (-1, 1), (0, -1), (0, 1), (1, -1), (1, 0), (1, 1)]


def _rolls(instr: str) -> set[tuple[int, int]]:
    return {
        (x, y)
        for y, line in enumerate(instr.splitlines())
        for x, c in enumerate(line)
        if c == "@"
    }


def _accessible(rolls: set[tuple[int, int]], p: tuple[int, int]) -> bool:
    # A roll can be reached iff fewer than four of its eight neighbors are filled.
    x, y = p
    return sum((x + dx, y + dy) in rolls for dx, dy in _ADJ) < 4


# Exercise for Advent of Code 2025 day 4.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        rolls = _rolls(instr)
        return sum(_accessible(rolls, p) for p in rolls)

    @staticmethod
    def two(instr: str) -> int:
        # Peel the floor in synchronous waves: each round removes every roll
        # accessible in the current state. Removing a roll only lowers neighbor
        # counts, so the process is monotone and converges to a fixed point.
        rolls = _rolls(instr)
        removed = 0
        while True:
            wave = {p for p in rolls if _accessible(rolls, p)}
            if not wave:
                return removed
            removed += len(wave)
            rolls -= wave
