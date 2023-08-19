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
        raise NotImplementedError
