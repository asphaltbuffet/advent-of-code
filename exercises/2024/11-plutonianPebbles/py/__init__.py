from collections import Counter
from typing import *
from aocpy import BaseExercise


def _blink(n: int) -> Tuple[int, ...]:
    if n == 0:
        return (1,)
    s = str(n)
    if len(s) % 2 == 0:
        half = len(s) // 2
        return (int(s[:half]), int(s[half:]))
    return (n * 2024,)


def _count_after(instr: str, blinks: int) -> int:
    stones = Counter(int(x) for x in instr.split())
    for _ in range(blinks):
        nxt: Counter = Counter()
        for n, cnt in stones.items():
            for m in _blink(n):
                nxt[m] += cnt
        stones = nxt
    return sum(stones.values())


# Exercise for Advent of Code 2024 day 11.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _count_after(instr, 25)

    @staticmethod
    def two(instr: str) -> int:
        return _count_after(instr, 75)
