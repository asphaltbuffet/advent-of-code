from typing import *
from aocpy import BaseExercise


def parse(inpt: str) -> List[int]:
    return [int(x) for x in inpt.splitlines() if x != ""]


def increase_count(data: List[int]) -> int:
    c = 0

    for i in range(1, len(data)):
        if data[i] > data[i - 1]:
            c += 1
    return c


class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        data = parse(instr)
        return increase_count(data)

    @staticmethod
    def two(instr: str) -> int:
        data = parse(instr)
        counts = [sum(data[i : i + 3]) for i in range(len(data) - 2)]
        return increase_count(counts)
