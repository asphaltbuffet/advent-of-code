import re
from math import prod
from typing import *
from aocpy import BaseExercise

LIMITS = {"red": 12, "green": 13, "blue": 14}
DRAW_RE = re.compile(r"(\d+) (red|green|blue)")


def _maxima(line: str) -> dict[str, int]:
    # Collapse every "<n> <color>" draw in a game into the max seen per color.
    maxima = {"red": 0, "green": 0, "blue": 0}
    for count, color in DRAW_RE.findall(line):
        maxima[color] = max(maxima[color], int(count))
    return maxima


# Exercise for Advent of Code 2023 day 2.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        total = 0
        for game_id, line in enumerate(instr.splitlines(), start=1):
            need = _maxima(line)
            if all(need[c] <= limit for c, limit in LIMITS.items()):
                total += game_id
        return total

    @staticmethod
    def two(instr: str) -> int:
        return sum(prod(_maxima(line).values()) for line in instr.splitlines())
