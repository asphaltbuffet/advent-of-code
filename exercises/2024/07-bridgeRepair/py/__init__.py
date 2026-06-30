from typing import *
from aocpy import BaseExercise


def _parse(instr: str) -> List[Tuple[int, List[int]]]:
    eqs = []
    for line in instr.strip().splitlines():
        target, rest = line.split(":", 1)
        eqs.append((int(target), [int(n) for n in rest.split()]))
    return eqs


def _solvable(target: int, acc: int, nums: List[int], concat: bool) -> bool:
    if not nums:
        return acc == target
    if acc > target:
        return False
    n, rest = nums[0], nums[1:]
    if _solvable(target, acc + n, rest, concat):
        return True
    if _solvable(target, acc * n, rest, concat):
        return True
    if concat and _solvable(target, int(f"{acc}{n}"), rest, concat):
        return True
    return False


def _calibrate(instr: str, concat: bool) -> int:
    return sum(
        target
        for target, nums in _parse(instr)
        if _solvable(target, nums[0], nums[1:], concat)
    )


# Exercise for Advent of Code 2024 day 7.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _calibrate(instr, False)

    @staticmethod
    def two(instr: str) -> int:
        return _calibrate(instr, True)
