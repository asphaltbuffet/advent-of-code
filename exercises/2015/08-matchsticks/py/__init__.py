from typing import *
from aocpy import BaseExercise

import re


def strip(s: str) -> str:
    s = re.sub(r'\\"', r'"', s)
    s = re.sub(r"\\\\", r"\\", s)
    s = re.sub(r"\\x[0-9a-f]{2}", r"x", s)
    # print(f"{count} += {orig} - {len(line) - 2}")

    return s


def escape(s: str) -> str:
    s = re.sub(r"\\", r"\\\\", s)
    s = re.sub(r'"', r'\\"', s)
    # s = re.sub(r"\\x[0-9a-f]{2}", r"x", s)

    return s


# Exercise for Advent of Code 2015 day 8.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        lines = instr.splitlines()

        count = 0
        for line in lines:
            orig = len(line)
            stripped = strip(line)
            count += orig - (len(stripped) - 2)

        return count

    @staticmethod
    def two(instr: str) -> int:
        lines = instr.splitlines()

        count = 0
        for line in lines:
            orig = len(line)
            esc = f'"{escape(line)}"'
            # print(f"{line} -> {esc}")
            # print(f"{count} += {len(esc)} - {orig}")
            count += len(esc) - orig

        return count
