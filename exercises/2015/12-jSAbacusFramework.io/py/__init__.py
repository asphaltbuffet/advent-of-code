from typing import *
from aocpy import BaseExercise

import json
import re


def sum_value(v: Any) -> int:
    """Recursively sum numbers, skipping objects with a "red" value."""
    if isinstance(v, int):
        return v
    if isinstance(v, list):
        return sum(sum_value(e) for e in v)
    if isinstance(v, dict):
        if "red" in v.values():
            return 0
        return sum(sum_value(e) for e in v.values())
    return 0


# Exercise for Advent of Code 2015 day 12.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return sum(int(n) for n in re.findall(r"-?\d+", instr))

    @staticmethod
    def two(instr: str) -> int:
        return sum_value(json.loads(instr.strip()))
