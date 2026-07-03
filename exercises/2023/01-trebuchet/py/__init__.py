import re
from typing import *
from aocpy import BaseExercise

WORDS = {
    "one": 1, "two": 2, "three": 3, "four": 4, "five": 5,
    "six": 6, "seven": 7, "eight": 8, "nine": 9,
}

# A single overlapping-scan regex per part. The lookahead `(?=(...))` lets
# adjacent spellings such as "eightwo" both match, which a consuming scan would
# miss. Part one uses only the digit alternative.
DIGITS_RE = re.compile(r"(?=([1-9]))")
WORDS_RE = re.compile(r"(?=([1-9]|" + "|".join(WORDS) + r"))")


def _value(digit: str) -> int:
    return WORDS[digit] if digit in WORDS else int(digit)


def _calibrate(instr: str, pattern: re.Pattern) -> int:
    total = 0
    for line in instr.splitlines():
        found = pattern.findall(line)
        if found:
            total += _value(found[0]) * 10 + _value(found[-1])
    return total


# Exercise for Advent of Code 2023 day 1.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _calibrate(instr, DIGITS_RE)

    @staticmethod
    def two(instr: str) -> int:
        return _calibrate(instr, WORDS_RE)
