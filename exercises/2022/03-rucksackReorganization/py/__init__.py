from typing import *
from aocpy import BaseExercise
import string

priority_value = {c: i + 1 for i, c in enumerate(string.ascii_lowercase)}
priority_value.update({c: i + 27 for i, c in enumerate(string.ascii_uppercase)})


def score_mispacked(line: str) -> int:
    compartment_one = set(line[: len(line) // 2])
    priority = 0

    for ch in line[len(line) // 2 :]:
        if ch in compartment_one:
            priority += priority_value[ch]
            compartment_one.remove(ch)

    return priority


def score_badges(a, b, c):
    shared_items = set(a)

    # list comprehension to filter the shared items
    shared_items = [item for item in shared_items if item in b and item in c]

    return priority_value[shared_items[0]] if shared_items else 0


# Exercise for Advent of Code 2022 day 3
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        scores = [score_mispacked(line) for line in instr.split("\n")]
        return sum(scores)

    @staticmethod
    def two(instr: str) -> int:
        data = instr.split("\n")

        # use a generator expression to calculate the scores for each triplet of lines
        score = sum(
            score_badges(data[i], data[i + 1], data[i + 2])
            for i in range(0, len(data), 3)
        )

        return score
