import re
from collections import Counter, defaultdict
from typing import *

from aocpy import BaseExercise


def _sleep_by_guard(instr: str) -> dict[int, Counter]:
    """Sort the shuffled log (timestamps sort lexically) and tally, per guard,
    how many days they were asleep at each minute 0..59."""
    asleep: dict[int, Counter] = defaultdict(Counter)
    guard = 0
    start = 0

    for line in sorted(instr.strip().splitlines()):
        minute = int(line[15:17])
        if "Guard" in line:
            guard = int(re.search(r"#(\d+)", line).group(1))
        elif "falls asleep" in line:
            start = minute
        elif "wakes up" in line:
            asleep[guard].update(range(start, minute))

    return asleep


# Exercise for Advent of Code 2018 day 4.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        asleep = _sleep_by_guard(instr)
        # Strategy 1: guard with the most total minutes asleep.
        guard = max(asleep, key=lambda g: sum(asleep[g].values()))
        minute = asleep[guard].most_common(1)[0][0]
        return guard * minute

    @staticmethod
    def two(instr: str) -> int:
        asleep = _sleep_by_guard(instr)
        # Strategy 2: the single (guard, minute) with the highest sleep count.
        guard, minute = max(
            (
                (g, mins.most_common(1)[0][0])
                for g, mins in asleep.items()
            ),
            key=lambda gm: asleep[gm[0]][gm[1]],
        )
        return guard * minute
