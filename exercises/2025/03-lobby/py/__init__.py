from typing import *
from aocpy import BaseExercise


def _largest_subsequence(digits: str, k: int) -> int:
    # Greedily keep the largest length-k subsequence in order: sweep left to
    # right on a monotonic stack, popping a smaller trailing digit whenever we
    # still have deletions to spend. The first k survivors form the answer.
    drop = len(digits) - k
    stack: list[str] = []
    for c in digits:
        while drop and stack and stack[-1] < c:
            stack.pop()
            drop -= 1
        stack.append(c)
    return int("".join(stack[:k]))


def _sum_over_lines(instr: str, k: int) -> int:
    return sum(_largest_subsequence(line, k) for line in instr.splitlines())


# Exercise for Advent of Code 2025 day 3.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _sum_over_lines(instr, 2)

    @staticmethod
    def two(instr: str) -> int:
        return _sum_over_lines(instr, 12)
