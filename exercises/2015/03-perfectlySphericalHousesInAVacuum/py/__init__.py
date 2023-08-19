from typing import *
from aocpy import BaseExercise


def move(x: int, y: int, c: str) -> Tuple[int, int]:
    if c == "^":
        y -= 1
    elif c == "v":
        y += 1
    elif c == ">":
        x += 1
    elif c == "<":
        x -= 1

    return x, y


# Exercise for Advent of Code 2015 day 3.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        x, y = 0, 0
        visited = set()
        visited.add((0, 0))

        for c in instr:
            x, y = move(x, y, c)
            visited.add((x, y))

        return len(visited)

    @staticmethod
    def two(instr: str) -> int:
        x, y = 0, 0
        a, b = 0, 0
        visited = set()
        visited.add((0, 0))

        for i, c in enumerate(instr):
            if i % 2 == 0:
                x, y = move(x, y, c)
                visited.add((x, y))
            else:
                a, b = move(a, b, c)
                visited.add((a, b))

        return len(visited)
