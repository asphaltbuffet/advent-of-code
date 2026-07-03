from typing import *
from aocpy import BaseExercise

# Which directions each pipe connects, as (dr, dc) offsets.
PIPES = {
    "|": {(-1, 0), (1, 0)},
    "-": {(0, -1), (0, 1)},
    "L": {(-1, 0), (0, 1)},
    "J": {(-1, 0), (0, -1)},
    "7": {(1, 0), (0, -1)},
    "F": {(1, 0), (0, 1)},
}


def _trace(instr: str) -> list[tuple[int, int]]:
    grid = instr.splitlines()
    start = next(
        (r, c)
        for r, row in enumerate(grid)
        for c, ch in enumerate(row)
        if ch == "S"
    )

    # Infer S's shape by finding a neighbor that connects back to S.
    def connects(r, c):
        ch = grid[r][c] if 0 <= r < len(grid) and 0 <= c < len(grid[r]) else "."
        return PIPES.get(ch, set())

    for dr, dc in ((-1, 0), (1, 0), (0, -1), (0, 1)):
        nr, nc = start[0] + dr, start[1] + dc
        if (-dr, -dc) in connects(nr, nc):
            prev, cur = start, (nr, nc)
            break

    loop = [start]
    while cur != start:
        loop.append(cur)
        for dr, dc in PIPES[grid[cur[0]][cur[1]]]:
            nxt = (cur[0] + dr, cur[1] + dc)
            if nxt != prev:
                prev, cur = cur, nxt
                break
    return loop


def _interior(loop: list[tuple[int, int]]) -> int:
    # Shoelace gives the polygon area; Pick's theorem recovers interior points:
    # A = i + b/2 - 1  =>  i = A - b/2 + 1, with b the loop length.
    area = abs(
        sum(
            r0 * c1 - r1 * c0
            for (r0, c0), (r1, c1) in zip(loop, loop[1:] + loop[:1])
        )
    ) // 2
    return area - len(loop) // 2 + 1


# Exercise for Advent of Code 2023 day 10.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return len(_trace(instr)) // 2

    @staticmethod
    def two(instr: str) -> int:
        return _interior(_trace(instr))
