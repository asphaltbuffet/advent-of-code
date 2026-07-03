from typing import *
from aocpy import BaseExercise

# Direction deltas as (dr, dc).
DIRS = {"R": (0, 1), "L": (0, -1), "U": (-1, 0), "D": (1, 0)}
# Part two decodes the last hex digit as a direction.
HEX_DIRS = {"0": "R", "1": "D", "2": "L", "3": "U"}


def _lagoon(steps: Iterable[tuple[str, int]]) -> int:
    # Trace the trench, accumulating the shoelace area and the perimeter length.
    # The lagoon holds interior + boundary cubes; Pick's theorem gives interior
    # i = A - b/2 + 1, so the total is A + b/2 + 1.
    r = c = area = perimeter = 0
    for direction, dist in steps:
        dr, dc = DIRS[direction]
        nr, nc = r + dr * dist, c + dc * dist
        area += r * nc - nr * c  # shoelace cross term
        perimeter += dist
        r, c = nr, nc
    return abs(area) // 2 + perimeter // 2 + 1


# Exercise for Advent of Code 2023 day 18.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        steps = (
            (parts[0], int(parts[1]))
            for line in instr.splitlines()
            for parts in [line.split()]
        )
        return _lagoon(steps)

    @staticmethod
    def two(instr: str) -> int:
        def decode(line: str) -> tuple[str, int]:
            code = line.split()[2].strip("(#)")
            return HEX_DIRS[code[5]], int(code[:5], 16)

        return _lagoon(decode(line) for line in instr.splitlines())
