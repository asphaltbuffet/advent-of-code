from collections import Counter

from aocpy import BaseExercise

OPEN, TREES, LUMBER = ".", "|", "#"


def _parse(instr: str) -> list[list[str]]:
    return [list(line) for line in instr.strip().splitlines()]


def _step(grid: list[list[str]]) -> list[list[str]]:
    h, w = len(grid), len(grid[0])
    nxt = [row[:] for row in grid]
    for y in range(h):
        row = grid[y]
        y0, y1 = max(y - 1, 0), min(y + 2, h)
        for x in range(w):
            x0, x1 = max(x - 1, 0), min(x + 2, w)
            trees = lumber = 0
            for ny in range(y0, y1):
                nrow = grid[ny]
                for nx in range(x0, x1):
                    if ny == y and nx == x:
                        continue
                    c = nrow[nx]
                    if c == TREES:
                        trees += 1
                    elif c == LUMBER:
                        lumber += 1
            cell = row[x]
            if cell == OPEN and trees >= 3:
                nxt[y][x] = TREES
            elif cell == TREES and lumber >= 3:
                nxt[y][x] = LUMBER
            elif cell == LUMBER and not (lumber >= 1 and trees >= 1):
                nxt[y][x] = OPEN
    return nxt


def _resource_value(grid: list[list[str]]) -> int:
    flat = Counter(c for row in grid for c in row)
    return flat[TREES] * flat[LUMBER]


def _key(grid: list[list[str]]) -> str:
    return "".join("".join(row) for row in grid)


# Exercise for Advent of Code 2018 day 18.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        grid = _parse(instr)
        for _ in range(10):
            grid = _step(grid)
        return _resource_value(grid)

    @staticmethod
    def two(instr: str) -> int:
        target = 1_000_000_000
        grid = _parse(instr)
        # The state settles into a cycle; fast-forward once a state repeats.
        seen: dict[str, int] = {}
        minute = 0
        while minute < target:
            key = _key(grid)
            if key in seen:
                period = minute - seen[key]
                for _ in range((target - minute) % period):
                    grid = _step(grid)
                break
            seen[key] = minute
            grid = _step(grid)
            minute += 1
        return _resource_value(grid)
