from functools import lru_cache
from typing import *
from aocpy import BaseExercise

_STEPS = [(-1, 0), (1, 0), (0, -1), (0, 1)]


def _parse(instr: str) -> List[List[int]]:
    return [[int(ch) if ch.isdigit() else -1 for ch in line] for line in instr.split()]


def _at(grid: List[List[int]], r: int, c: int) -> int:
    if 0 <= r < len(grid) and 0 <= c < len(grid[r]):
        return grid[r][c]
    return -1


def _reachable_9s(grid: List[List[int]], r: int, c: int, ends: Set[Tuple[int, int]]) -> None:
    if grid[r][c] == 9:
        ends.add((r, c))
        return
    for dr, dc in _STEPS:
        if _at(grid, r + dr, c + dc) == grid[r][c] + 1:
            _reachable_9s(grid, r + dr, c + dc, ends)


# Exercise for Advent of Code 2024 day 10.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        grid = _parse(instr)
        score = 0
        for r in range(len(grid)):
            for c in range(len(grid[r])):
                if grid[r][c] == 0:
                    ends: Set[Tuple[int, int]] = set()
                    _reachable_9s(grid, r, c, ends)
                    score += len(ends)
        return score

    @staticmethod
    def two(instr: str) -> int:
        grid = _parse(instr)

        @lru_cache(maxsize=None)
        def ratings(r: int, c: int) -> int:
            if grid[r][c] == 9:
                return 1
            total = 0
            for dr, dc in _STEPS:
                if _at(grid, r + dr, c + dc) == grid[r][c] + 1:
                    total += ratings(r + dr, c + dc)
            return total

        return sum(
            ratings(r, c)
            for r in range(len(grid))
            for c in range(len(grid[r]))
            if grid[r][c] == 0
        )
