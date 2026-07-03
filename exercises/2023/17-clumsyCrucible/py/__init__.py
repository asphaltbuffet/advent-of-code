import heapq
from typing import *
from aocpy import BaseExercise

# Two axes of travel; from a given axis the crucible must turn onto the other.
HORIZONTAL, VERTICAL = 0, 1


def _min_heat(grid: list[list[int]], lo: int, hi: int) -> int:
    height, width = len(grid), len(grid[0])
    goal = (height - 1, width - 1)

    # State = (r, c, axis just travelled). A move goes lo..hi cells along the
    # perpendicular axis, so the turn constraint is baked into the transitions and
    # no run-length needs tracking.
    best: dict[tuple[int, int, int], int] = {}
    # Seed both axes at the start so the first move can go either way.
    pq = [(0, 0, 0, HORIZONTAL), (0, 0, 0, VERTICAL)]

    while pq:
        cost, r, c, axis = heapq.heappop(pq)
        if (r, c) == goal:
            return cost
        if best.get((r, c, axis), 1 << 60) < cost:
            continue

        # Turn onto the other axis and step lo..hi cells, summing heat as we go.
        ndir = (0, 1) if axis == VERTICAL else (1, 0)
        for sign in (1, -1):
            dr, dc = ndir[0] * sign, ndir[1] * sign
            nr, nc, added = r, c, cost
            for step in range(1, hi + 1):
                nr, nc = nr + dr, nc + dc
                if not (0 <= nr < height and 0 <= nc < width):
                    break
                added += grid[nr][nc]
                if step < lo:
                    continue
                naxis = axis ^ 1
                key = (nr, nc, naxis)
                if added < best.get(key, 1 << 60):
                    best[key] = added
                    heapq.heappush(pq, (added, nr, nc, naxis))

    return -1


# Exercise for Advent of Code 2023 day 17.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        grid = [[int(ch) for ch in line] for line in instr.splitlines()]
        return _min_heat(grid, 1, 3)

    @staticmethod
    def two(instr: str) -> int:
        grid = [[int(ch) for ch in line] for line in instr.splitlines()]
        return _min_heat(grid, 4, 10)
