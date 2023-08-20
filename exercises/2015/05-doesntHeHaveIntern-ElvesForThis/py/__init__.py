from typing import *
from aocpy import BaseExercise


# Exercise for Advent of Code 2015 day 5.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        has_vowels = lambda s: sum([s.count(v) for v in "aeiou"]) >= 3
        has_double = lambda s: any([s[i] == s[i + 1] for i in range(len(s) - 1)])
        has_bad = lambda s: any([s.count(b) > 0 for b in ["ab", "cd", "pq", "xy"]])

        return sum(
            [
                has_vowels(s) and has_double(s) and not has_bad(s)
                for s in instr.split("\n")
            ]
        )

    @staticmethod
    def two(instr: str) -> int:
        has_pair = lambda s: any([s.count(s[i : i + 2]) > 1 for i in range(len(s) - 1)])
        has_repeat = lambda s: any([s[i] == s[i + 2] for i in range(len(s) - 2)])

        return sum([has_pair(s) and has_repeat(s) for s in instr.split("\n")])
