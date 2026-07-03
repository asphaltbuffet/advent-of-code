import re
import sys

from aocpy import BaseExercise

# The water column can be ~2000 tiles deep, so the recursive fill needs headroom
# beyond CPython's default recursion limit.
sys.setrecursionlimit(10_000)

CLAY, FLOWING, SETTLED = "#", "|", "~"


def _simulate(instr: str) -> tuple[dict[tuple[int, int], str], int, int]:
    """Flow water from the spring at x=500 and return (grid, min_y, max_y)."""
    grid: dict[tuple[int, int], str] = {}
    for line in instr.strip().splitlines():
        a, b, c = map(int, re.findall(r"-?\d+", line))
        if line[0] == "x":  # x=a, y=b..c
            for y in range(b, c + 1):
                grid[(a, y)] = CLAY
        else:  # y=a, x=b..c
            for x in range(b, c + 1):
                grid[(x, a)] = CLAY
    ys = [y for _, y in grid]
    min_y, max_y = min(ys), max(ys)

    def is_floor(x: int, y: int) -> bool:
        return grid.get((x, y)) in (CLAY, SETTLED)

    def spread(x: int, y: int, dx: int) -> tuple[int, bool]:
        """Walk along row y in direction dx, spilling down open edges. Return the
        furthest reachable x and whether it ended against a wall."""
        while True:
            if grid.get((x + dx, y)) == CLAY:
                return x, True
            x += dx
            grid[(x, y)] = FLOWING
            if not is_floor(x, y + 1):
                if grid.get((x, y + 1)) is None and flow(x, y + 1):
                    continue  # spill settled into a new floor: extend the row
                return x, False  # water escapes here

    def flow(x: int, y: int) -> bool:
        """Drop water from (x, y). Return True if it comes to rest as settled water."""
        grid[(x, y)] = FLOWING
        below = y + 1
        if below > max_y:
            return False
        if not is_floor(x, below):
            if grid.get((x, below)) is None:
                if not flow(x, below):
                    return False
            else:
                return False  # already flowing through: not a floor

        lx, left_wall = spread(x, y, -1)
        rx, right_wall = spread(x, y, +1)
        if left_wall and right_wall:
            for fx in range(lx, rx + 1):
                grid[(fx, y)] = SETTLED
            return True
        return False

    flow(500, 0)
    return grid, min_y, max_y


def _counts(instr: str) -> tuple[int, int]:
    grid, min_y, max_y = _simulate(instr)
    reached = sum(1 for (_, y), t in grid.items() if min_y <= y <= max_y and t != CLAY)
    settled = sum(1 for (_, y), t in grid.items() if min_y <= y <= max_y and t == SETTLED)
    return reached, settled


# Exercise for Advent of Code 2018 day 17.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _counts(instr)[0]

    @staticmethod
    def two(instr: str) -> int:
        return _counts(instr)[1]
