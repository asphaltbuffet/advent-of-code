import re
from collections import defaultdict
from math import prod
from typing import *
from aocpy import BaseExercise

NUMBER_RE = re.compile(r"\d+")


def _numbers(grid: list[str]) -> Iterator[tuple[int, set[tuple[int, int]]]]:
    # Yield each number with the set of cells bordering it (its 8-neighborhood).
    for y, row in enumerate(grid):
        for m in NUMBER_RE.finditer(row):
            border = {
                (x, ny)
                for ny in (y - 1, y, y + 1)
                for x in range(m.start() - 1, m.end() + 1)
            }
            yield int(m.group()), border


# Exercise for Advent of Code 2023 day 3.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        grid = instr.splitlines()
        symbols = {
            (x, y)
            for y, row in enumerate(grid)
            for x, ch in enumerate(row)
            if ch != "." and not ch.isdigit()
        }
        return sum(
            value
            for value, border in _numbers(grid)
            if border & symbols
        )

    @staticmethod
    def two(instr: str) -> int:
        grid = instr.splitlines()
        stars = {
            (x, y)
            for y, row in enumerate(grid)
            for x, ch in enumerate(row)
            if ch == "*"
        }
        # Group each gear's adjacent numbers, then multiply the pairs.
        adjacent = defaultdict(list)
        for value, border in _numbers(grid):
            for star in border & stars:
                adjacent[star].append(value)
        return sum(prod(nums) for nums in adjacent.values() if len(nums) == 2)
