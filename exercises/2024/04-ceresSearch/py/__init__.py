from typing import *
from aocpy import BaseExercise

_DIRS = [(0, 1), (0, -1), (1, 0), (-1, 0), (1, 1), (1, -1), (-1, 1), (-1, -1)]
_WORD = "XMAS"


def _grid(instr: str) -> List[str]:
    return instr.split()


# Exercise for Advent of Code 2024 day 4.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        g = _grid(instr)
        rows = len(g)
        count = 0
        for r in range(rows):
            for c in range(len(g[r])):
                if g[r][c] != "X":
                    continue
                for dr, dc in _DIRS:
                    for k in range(1, len(_WORD)):
                        nr, nc = r + dr * k, c + dc * k
                        if not (0 <= nr < rows and 0 <= nc < len(g[nr])) or g[nr][nc] != _WORD[k]:
                            break
                    else:
                        count += 1
        return count

    @staticmethod
    def two(instr: str) -> int:
        g = _grid(instr)
        rows = len(g)

        def is_mas(a: str, b: str) -> bool:
            return {a, b} == {"M", "S"}

        count = 0
        for r in range(1, rows - 1):
            for c in range(1, len(g[r]) - 1):
                if g[r][c] != "A":
                    continue
                if is_mas(g[r - 1][c - 1], g[r + 1][c + 1]) and is_mas(g[r - 1][c + 1], g[r + 1][c - 1]):
                    count += 1
        return count
