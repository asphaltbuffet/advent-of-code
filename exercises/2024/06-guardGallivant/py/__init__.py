from typing import *
from aocpy import BaseExercise

# Facing directions, ordered so (dir + 1) % 4 is a right turn.
_DIRS = [(-1, 0), (0, 1), (1, 0), (0, -1)]  # up, right, down, left


def _find_start(grid: List[str]) -> Tuple[int, int]:
    for r, row in enumerate(grid):
        c = row.find("^")
        if c != -1:
            return r, c
    return -1, -1


def _walk(
    grid: List[str],
    start: Tuple[int, int],
    extra: Optional[Tuple[int, int]] = None,
) -> Tuple[Set[Tuple[int, int]], bool]:
    """Simulate the guard. Returns (visited cells, looped?)."""
    rows = len(grid)
    visited: Set[Tuple[int, int]] = set()
    seen: Set[Tuple[int, int, int]] = set()
    r, c = start
    d = 0
    while True:
        visited.add((r, c))
        if (r, c, d) in seen:
            return visited, True
        seen.add((r, c, d))

        dr, dc = _DIRS[d]
        nr, nc = r + dr, c + dc
        if not (0 <= nr < rows and 0 <= nc < len(grid[nr])):
            return visited, False
        if grid[nr][nc] == "#" or (nr, nc) == extra:
            d = (d + 1) % 4
            continue
        r, c = nr, nc


# Exercise for Advent of Code 2024 day 6.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        grid = instr.split()
        visited, _ = _walk(grid, _find_start(grid))
        return len(visited)

    @staticmethod
    def two(instr: str) -> int:
        grid = instr.split()
        start = _find_start(grid)
        path, _ = _walk(grid, start)
        path.discard(start)
        return sum(1 for cell in path if _walk(grid, start, cell)[1])
