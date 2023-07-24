from typing import *
from aocpy import BaseExercise


def parse(instr: str) -> List[str]:
    cubes = []

    for line in instr.splitlines():
        c = tuple(map(int, line.split(",")))
        cubes.append(c)

    return cubes


# Exercise for Advent of Code 2022 day 18.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        cubes = parse(instr)

        cube_set = set(cubes)

        # start with all cube sides exposed
        total_exposure = len(cubes) * 6

        adjacent = [
            (0, 0, -1),
            (0, 0, 1),
            (0, -1, 0),
            (0, 1, 0),
            (-1, 0, 0),
            (1, 0, 0),
        ]

        for c in cubes:
            for a in adjacent:
                x = c[0] + a[0]
                y = c[1] + a[1]
                z = c[2] + a[2]

                if (x, y, z) in cube_set:
                    total_exposure -= 1

        return total_exposure

    @staticmethod
    def two(instr: str) -> int:
        raise NotImplementedError
