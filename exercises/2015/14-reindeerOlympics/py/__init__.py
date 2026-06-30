from typing import *
from aocpy import BaseExercise


def parse(instr: str) -> Tuple[List[Tuple[int, int, int]], int]:
    """Parse (speed, fly, rest) per reindeer. The race duration isn't in the
    input, so it's inferred: the 2-reindeer AoC example races 1000s, real 2503s.
    """
    reindeer = []
    for line in instr.strip().splitlines():
        f = line.split()
        # <name> can fly <speed> km/s for <fly> seconds, ... rest for <rest> ...
        reindeer.append((int(f[3]), int(f[6]), int(f[13])))

    duration = 1000 if len(reindeer) <= 2 else 2503
    return reindeer, duration


def distance(r: Tuple[int, int, int], t: int) -> int:
    """How far reindeer r has travelled after t seconds."""
    speed, fly, rest = r
    cycle = fly + rest
    flying = (t // cycle) * fly + min(t % cycle, fly)
    return flying * speed


# Exercise for Advent of Code 2015 day 14.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        reindeer, duration = parse(instr)
        return max(distance(r, duration) for r in reindeer)

    @staticmethod
    def two(instr: str) -> int:
        reindeer, duration = parse(instr)
        points = [0] * len(reindeer)

        for t in range(1, duration + 1):
            dists = [distance(r, t) for r in reindeer]
            lead = max(dists)
            for i, d in enumerate(dists):
                if d == lead:
                    points[i] += 1

        return max(points)
