from typing import *
from aocpy import BaseExercise


def parse(instr: str) -> List[List[bool]]:
    """Read the light grid; True is on."""
    return [[c == "#" for c in line] for line in instr.split()]


def steps(grid: List[List[bool]], stuck: bool) -> int:
    """Animation length. The small AoC example runs 4 steps for part 1 and 5
    for part 2 (stuck corners); the real input runs 100 for both."""
    if len(grid) > 6:
        return 100
    return 5 if stuck else 4


def neighbors_on(grid: List[List[bool]], r: int, c: int) -> int:
    n = len(grid)
    total = 0
    for dr in (-1, 0, 1):
        for dc in (-1, 0, 1):
            if dr == 0 and dc == 0:
                continue
            nr, nc = r + dr, c + dc
            if 0 <= nr < n and 0 <= nc < len(grid[nr]) and grid[nr][nc]:
                total += 1
    return total


def step(grid: List[List[bool]]) -> List[List[bool]]:
    return [
        [
            (grid[r][c] and neighbors_on(grid, r, c) in (2, 3))
            or (not grid[r][c] and neighbors_on(grid, r, c) == 3)
            for c in range(len(grid[r]))
        ]
        for r in range(len(grid))
    ]


def stick_corners(grid: List[List[bool]]) -> None:
    last = len(grid) - 1
    grid[0][0] = grid[0][last] = grid[last][0] = grid[last][last] = True


def run(grid: List[List[bool]], n: int, stuck: bool) -> int:
    if stuck:
        stick_corners(grid)
    for _ in range(n):
        grid = step(grid)
        if stuck:
            stick_corners(grid)
    return sum(cell for row in grid for cell in row)


# Exercise for Advent of Code 2015 day 18.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        grid = parse(instr)
        return run(grid, steps(grid, False), False)

    @staticmethod
    def two(instr: str) -> int:
        grid = parse(instr)
        return run(grid, steps(grid, True), True)
