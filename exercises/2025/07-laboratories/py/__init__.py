from collections import defaultdict
from typing import *
from aocpy import BaseExercise


def _sweep(instr: str) -> tuple[int, int]:
    # Propagate the beam wavefront row by row as a {column: count} map. A beam on
    # a splitter (^) sends its count down-left and down-right; every such event is
    # one split, and the final counts are the number of timelines.
    lines = instr.splitlines()
    beams: dict[int, int] = {lines[0].index("S"): 1}
    splits = 0
    for row in lines[1:]:
        nxt: dict[int, int] = defaultdict(int)
        for x, count in beams.items():
            if 0 <= x < len(row) and row[x] == "^":
                nxt[x - 1] += count
                nxt[x + 1] += count
                splits += 1
            else:
                nxt[x] += count
        beams = nxt
    return splits, sum(beams.values())


# Exercise for Advent of Code 2025 day 7.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _sweep(instr)[0]

    @staticmethod
    def two(instr: str) -> int:
        return _sweep(instr)[1]
