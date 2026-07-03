from typing import *
from aocpy import BaseExercise


def _parse(instr: str) -> tuple[set[tuple[int, int]], tuple[int, int], int]:
    rocks = set()
    start = (0, 0)
    grid = instr.splitlines()
    size = len(grid)
    for r, row in enumerate(grid):
        for c, ch in enumerate(row):
            if ch == "#":
                rocks.add((r, c))
            elif ch == "S":
                start = (r, c)
    return rocks, start, size


def _reachable(rocks, start, size, steps: int, infinite: bool) -> int:
    # BFS by whole steps; the frontier is the plots reachable in exactly `s`
    # steps. On the infinite grid, positions wrap into the tile via modulo.
    frontier = {start}
    for _ in range(steps):
        nxt = set()
        for r, c in frontier:
            for dr, dc in ((-1, 0), (1, 0), (0, -1), (0, 1)):
                nr, nc = r + dr, c + dc
                if infinite:
                    if (nr % size, nc % size) not in rocks:
                        nxt.add((nr, nc))
                elif 0 <= nr < size and 0 <= nc < size and (nr, nc) not in rocks:
                    nxt.add((nr, nc))
        frontier = nxt
    return len(frontier)


# Exercise for Advent of Code 2023 day 21.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        rocks, start, size = _parse(instr)
        return _reachable(rocks, start, size, 64, infinite=False)

    @staticmethod
    def two(instr: str) -> int:
        rocks, start, size = _parse(instr)
        target = 26501365
        # The grid is a clear-bordered 131x131 square, so reachable counts grow
        # quadratically once per tile period. Sample at the three offsets
        # target % size, +size, +2*size, then fit and evaluate a quadratic.
        offset = target % size
        samples = [
            _reachable(rocks, start, size, offset + i * size, infinite=True)
            for i in range(3)
        ]
        # Lagrange over x = 0,1,2 with y = samples; evaluate at n = target // size.
        y0, y1, y2 = samples
        a = (y2 - 2 * y1 + y0) // 2
        b = y1 - y0 - a
        c = y0
        n = target // size
        return a * n * n + b * n + c
