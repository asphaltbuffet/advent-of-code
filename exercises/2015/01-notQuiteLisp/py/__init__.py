from typing import *
from aocpy import BaseExercise


# Exercise for Advent of Code 2015 day 1.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        up = instr.count("(")
        down = instr.count(")")
        return up - down

    @staticmethod
    def two(instr: str) -> int:
        floor = 0
        for i, c in enumerate(instr):
            if c == "(":
                floor += 1
            elif c == ")":
                floor -= 1

            if floor == -1:
                return i + 1
