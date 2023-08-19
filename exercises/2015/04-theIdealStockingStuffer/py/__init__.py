from typing import *
from aocpy import BaseExercise
from hashlib import md5


# Exercise for Advent of Code 2015 day 4.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        i = 0
        while True:
            i += 1
            h = md5((instr + str(i)).encode()).hexdigest()
            if h[:5] == "00000":
                return i

    @staticmethod
    def two(instr: str) -> int:
        i = 0
        while True:
            i += 1
            h = md5((instr + str(i)).encode()).hexdigest()
            if h[:6] == "000000":
                return i
