from typing import *
from aocpy import BaseExercise

Grid = tuple[str, ...]


def _roll_north(grid: Grid) -> Grid:
    # Slide every 'O' as far up as it can within each column segment bounded by
    # '#'. Working on columns, sorting a segment descending puts 'O' (79) before
    # '.' (46), i.e. toward the top.
    cols = ["".join(col) for col in zip(*grid)]
    rolled = []
    for col in cols:
        packed = "#".join(
            "".join(sorted(seg, reverse=True)) for seg in col.split("#")
        )
        rolled.append(packed)
    # Transpose back to rows.
    return tuple("".join(row) for row in zip(*rolled))


def _rotate_cw(grid: Grid) -> Grid:
    # Rotate 90 degrees clockwise; the west edge becomes the new north edge.
    return tuple("".join(col) for col in zip(*grid[::-1]))


def _tilt_north(grid: Grid) -> Grid:
    return _roll_north(grid)


def _spin(grid: Grid) -> Grid:
    # A full N->W->S->E cycle: roll north, then rotate CW so the next roll_north
    # tilts what was west, and so on around all four directions.
    for _ in range(4):
        grid = _rotate_cw(_roll_north(grid))
    return grid


def _load(grid: Grid) -> int:
    height = len(grid)
    return sum(
        height - r
        for r, row in enumerate(grid)
        for ch in row
        if ch == "O"
    )


# Exercise for Advent of Code 2023 day 14.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _load(_tilt_north(tuple(instr.splitlines())))

    @staticmethod
    def two(instr: str) -> int:
        grid = tuple(instr.splitlines())
        target = 1_000_000_000
        seen: dict[Grid, int] = {}
        i = 0
        while i < target:
            grid = _spin(grid)
            i += 1
            if grid in seen:
                # Jump forward by whole cycles, then finish the remainder.
                period = i - seen[grid]
                i += ((target - i) // period) * period
                seen = {}  # avoid re-jumping
            else:
                seen[grid] = i
        return _load(grid)
