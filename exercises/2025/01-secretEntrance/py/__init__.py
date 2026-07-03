from typing import *
from aocpy import BaseExercise


def _moves(instr: str) -> Iterator[tuple[str, int]]:
    # Each line is a direction (L/R) and a distance.
    for line in instr.splitlines():
        yield line[0], int(line[1:])


# Exercise for Advent of Code 2025 day 1.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        # The dial has 100 positions and starts at 50. Count how many moves end
        # with the pointer resting exactly on 0.
        pos = 50
        landings = 0
        for direction, dist in _moves(instr):
            pos = (pos + dist if direction == "R" else pos - dist) % 100
            landings += pos == 0
        return landings

    @staticmethod
    def two(instr: str) -> int:
        # Count every time the pointer sweeps onto 0 during a move. Working from
        # the reduced position each step, a rightward move of `d` from `pos`
        # passes zero (pos + d) // 100 times; a leftward move counts the
        # multiples of 100 in the half-open arc it sweeps into.
        pos = 50
        clicks = 0
        for direction, dist in _moves(instr):
            if direction == "R":
                clicks += (pos + dist) // 100
                pos = (pos + dist) % 100
            else:
                clicks += (pos - 1) // 100 - (pos - dist - 1) // 100
                pos = (pos - dist) % 100
        return clicks
