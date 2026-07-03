from typing import *
from aocpy import BaseExercise

# Direction deltas as (dr, dc).
UP, DOWN, LEFT, RIGHT = (-1, 0), (1, 0), (0, -1), (0, 1)

# How each tile transforms an incoming direction into outgoing directions.
TILES = {
    "/": {RIGHT: [UP], LEFT: [DOWN], UP: [RIGHT], DOWN: [LEFT]},
    "\\": {RIGHT: [DOWN], LEFT: [UP], UP: [LEFT], DOWN: [RIGHT]},
    "|": {LEFT: [UP, DOWN], RIGHT: [UP, DOWN]},
    "-": {UP: [LEFT, RIGHT], DOWN: [LEFT, RIGHT]},
}


def _energized(grid: list[str], start: tuple[tuple[int, int], tuple[int, int]]) -> int:
    height, width = len(grid), len(grid[0])
    seen: set[tuple[tuple[int, int], tuple[int, int]]] = set()
    stack = [start]

    while stack:
        (r, c), d = stack.pop()
        if not (0 <= r < height and 0 <= c < width):
            continue
        if ((r, c), d) in seen:
            continue
        seen.add(((r, c), d))

        # A tile either passes the beam through or redirects/splits it.
        for nd in TILES.get(grid[r][c], {}).get(d, [d]):
            stack.append(((r + nd[0], c + nd[1]), nd))

    return len({pos for pos, _ in seen})


# Exercise for Advent of Code 2023 day 16.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        grid = instr.splitlines()
        return _energized(grid, ((0, 0), RIGHT))

    @staticmethod
    def two(instr: str) -> int:
        grid = instr.splitlines()
        height, width = len(grid), len(grid[0])
        starts = (
            [((0, c), DOWN) for c in range(width)]
            + [((height - 1, c), UP) for c in range(width)]
            + [((r, 0), RIGHT) for r in range(height)]
            + [((r, width - 1), LEFT) for r in range(height)]
        )
        return max(_energized(grid, start) for start in starts)
