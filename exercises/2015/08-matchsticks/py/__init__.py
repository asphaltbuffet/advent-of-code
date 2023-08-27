from typing import *
from aocpy import BaseExercise

import re


# Exercise for Advent of Code 2015 day 8.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        lines = instr.splitlines()

        count = 0
        for line in lines:
            orig = len(line)
            line = re.sub(r'\\"', r'"', line)
            line = re.sub(r"\\\\", r"\\", line)
            line = re.sub(r"\\x[0-9a-f]{2}", r"x", line)
            # print(f"{count} += {orig} - {len(line) - 2}")

            count += orig - (len(line) - 2)

        return count

    @staticmethod
    def two(instr: str) -> int:
        raise NotImplementedError("part 2 not implemented")
