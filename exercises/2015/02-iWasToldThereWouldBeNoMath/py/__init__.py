from typing import *
from aocpy import BaseExercise


def parse(instr: str) -> List[Tuple[int, int, int]]:
    dims = []
    for line in instr.splitlines():
        dims.append(tuple(int(i) for i in line.split("x")))

    return dims


# Exercise for Advent of Code 2015 day 2.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        packages = parse(instr)

        total = 0
        for p in packages:
            a = p[0] * p[1]
            b = p[1] * p[2]
            c = p[2] * p[0]

            m = min(a, b, c)

            total += 2 * (a + b + c) + m

        return total

    @staticmethod
    def two(instr: str) -> int:
        raise NotImplementedError
