from math import isqrt, prod
from typing import *
from aocpy import BaseExercise


def _wins(time: int, record: int) -> int:
    # Holding the button h ms travels h*(time-h); beat the record when
    # h**2 - time*h + record < 0. The roots bound the winning h; count the
    # integers strictly between them via the quadratic formula (integer-safe).
    disc = time * time - 4 * record
    if disc <= 0:
        return 0
    root = isqrt(disc)
    # Smallest h with h*(time-h) > record, and the symmetric largest.
    lo = (time - root) // 2
    while lo * (time - lo) <= record:
        lo += 1
    hi = (time + root) // 2
    while hi * (time - hi) <= record:
        hi -= 1
    return hi - lo + 1


# Exercise for Advent of Code 2023 day 6.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        times, dists = (
            [int(n) for n in line.split(":")[1].split()]
            for line in instr.splitlines()
        )
        return prod(_wins(t, d) for t, d in zip(times, dists))

    @staticmethod
    def two(instr: str) -> int:
        time, record = (
            int(line.split(":")[1].replace(" ", ""))
            for line in instr.splitlines()
        )
        return _wins(time, record)
