from typing import *
from aocpy import BaseExercise


class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        raise NotImplementedError

    @staticmethod
    def two(instr: str) -> int:
        raise NotImplementedError