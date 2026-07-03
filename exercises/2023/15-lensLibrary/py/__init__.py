from functools import reduce
from typing import *
from aocpy import BaseExercise


def _hash(s: str) -> int:
    # Fold each character: ((acc + ord) * 17) mod 256.
    return reduce(lambda acc, ch: (acc + ord(ch)) * 17 % 256, s, 0)


# Exercise for Advent of Code 2023 day 15.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return sum(_hash(step) for step in instr.strip().split(","))

    @staticmethod
    def two(instr: str) -> int:
        # 256 boxes; a dict per box preserves insertion order, so '=' updates in
        # place and '-' removes, exactly matching the linked-list semantics.
        boxes: list[dict[str, int]] = [{} for _ in range(256)]
        for step in instr.strip().split(","):
            if step.endswith("-"):
                label = step[:-1]
                boxes[_hash(label)].pop(label, None)
            else:
                label, focal = step.split("=")
                boxes[_hash(label)][label] = int(focal)

        return sum(
            (box_i + 1) * slot * focal
            for box_i, box in enumerate(boxes)
            for slot, focal in enumerate(box.values(), start=1)
        )
