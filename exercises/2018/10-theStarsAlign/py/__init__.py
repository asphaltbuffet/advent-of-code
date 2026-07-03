import re

import numpy as np
from aocpy import BaseExercise

INT_RE = re.compile(r"-?\d+")


def parse(instr):
    """Parse each line's four integers into position and velocity arrays."""
    nums = np.array(
        [[int(n) for n in INT_RE.findall(line)] for line in instr.strip().splitlines()]
    )
    return nums[:, :2], nums[:, 2:]


def extent(pos, vel, t):
    """Combined width and height of the bounding box at time t."""
    p = pos + vel * t
    lo = p.min(axis=0)
    hi = p.max(axis=0)
    return (hi - lo).sum()


def converge(pos, vel):
    """Second at which the message appears (minimum bounding-box extent).

    The extent shrinks to that minimum then grows, so step until it stops.
    """
    t = 0
    while extent(pos, vel, t + 1) < extent(pos, vel, t):
        t += 1
    return t


def render(pos, vel, t):
    """Draw the points at time t as '█'/' ' rows joined by newlines."""
    p = pos + vel * t
    lo = p.min(axis=0)
    w, h = p.max(axis=0) - lo + 1
    grid = np.full((h, w), " ")
    xs, ys = (p - lo).T
    grid[ys, xs] = "█"
    return "\n".join("".join(row) for row in grid)


# Exercise for Advent of Code 2018 day 10.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> str:
        pos, vel = parse(instr)
        return render(pos, vel, converge(pos, vel))

    @staticmethod
    def two(instr: str) -> int:
        pos, vel = parse(instr)
        return int(converge(pos, vel))
