from typing import *
from aocpy import BaseExercise


def _axis_sum(coords: list[int], size: int, factor: int) -> int:
    # Sum of pairwise distances along one axis, expanding every empty track by
    # `factor`. Remapping each coordinate to its expanded position first, the
    # sorted-order identity sum_{i<j}(x_j - x_i) = sum_i x_i*(2i - n + 1) makes
    # this linear rather than quadratic.
    occupied = set(coords)
    expanded, offset = [], 0
    for x in range(size):
        if x not in occupied:
            offset += factor - 1
        expanded.append(x + offset)

    positions = sorted(expanded[c] for c in coords)
    n = len(positions)
    return sum(p * (2 * i - n + 1) for i, p in enumerate(positions))


def _total(instr: str, factor: int) -> int:
    grid = instr.splitlines()
    rows = [r for r, line in enumerate(grid) for ch in line if ch == "#"]
    cols = [c for line in grid for c, ch in enumerate(line) if ch == "#"]
    return _axis_sum(rows, len(grid), factor) + _axis_sum(cols, len(grid[0]), factor)


# Exercise for Advent of Code 2023 day 11.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _total(instr, 2)

    @staticmethod
    def two(instr: str) -> int:
        return _total(instr, 1_000_000)
