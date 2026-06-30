from typing import *
from aocpy import BaseExercise

_NEIGHBORS = [(-1, 0), (1, 0), (0, -1), (0, 1)]
_CORNERS = [(-1, -1), (-1, 1), (1, -1), (1, 1)]


def _regions(instr: str) -> List[Tuple[int, int, int]]:
    """Return (area, perimeter, sides) for each connected same-letter region."""
    grid = instr.split()
    rows = len(grid)

    def at(r: int, c: int) -> Optional[str]:
        if 0 <= r < rows and 0 <= c < len(grid[r]):
            return grid[r][c]
        return None

    seen: Set[Tuple[int, int]] = set()
    out: List[Tuple[int, int, int]] = []
    for sr in range(rows):
        for sc in range(len(grid[sr])):
            if (sr, sc) in seen:
                continue
            letter = grid[sr][sc]
            area = perimeter = sides = 0
            stack = [(sr, sc)]
            seen.add((sr, sc))
            while stack:
                r, c = stack.pop()
                area += 1

                def same(nr: int, nc: int) -> bool:
                    return at(nr, nc) == letter

                for dr, dc in _NEIGHBORS:
                    nr, nc = r + dr, c + dc
                    if not same(nr, nc):
                        perimeter += 1
                    elif (nr, nc) not in seen:
                        seen.add((nr, nc))
                        stack.append((nr, nc))

                for dr, dc in _CORNERS:
                    vert = same(r + dr, c)
                    horiz = same(r, c + dc)
                    diag = same(r + dr, c + dc)
                    if not vert and not horiz:
                        sides += 1
                    elif vert and horiz and not diag:
                        sides += 1

            out.append((area, perimeter, sides))
    return out


# Exercise for Advent of Code 2024 day 12.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return sum(area * perimeter for area, perimeter, _ in _regions(instr))

    @staticmethod
    def two(instr: str) -> int:
        return sum(area * sides for area, _, sides in _regions(instr))
