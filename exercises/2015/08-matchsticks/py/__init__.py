from typing import *
from aocpy import BaseExercise

import re


def strip(s: str) -> int:
    r = re.sub(r'\\"', r'"', s)
    r = re.sub(r"\\\\", r"\\", r)
    r = re.sub(r"\\x[0-9a-f]{2}", r"x", r)

    return len(r)


def escape(s: str) -> int:
    r = re.sub(r"\\", r"\\\\", s)
    r = re.sub(r'"', r'\\"', r)

    return len(r) + 2


# Exercise for Advent of Code 2015 day 8.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        lines = instr.splitlines()

        count = 0
        for line in lines:
            count += len(line) - (strip(line) - 2)

        return count

    @staticmethod
    def two(instr: str) -> int:
        lines = instr.splitlines()

        count = 0
        for line in lines:
            count += escape(line) - len(line)

        return count
