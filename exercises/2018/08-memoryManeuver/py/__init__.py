from typing import *
from aocpy import BaseExercise


def _meta_sum(nums: Iterator[int]) -> int:
    """Read one node from the stream, returning its total metadata sum."""
    children, meta = next(nums), next(nums)
    return sum(_meta_sum(nums) for _ in range(children)) + sum(
        next(nums) for _ in range(meta)
    )


def _node_value(nums: Iterator[int]) -> int:
    """Read one node, returning its value.

    A leaf's value is the sum of its metadata; an internal node's metadata are
    1-based indices into its children, summing the referenced children's values.
    """
    children, meta = next(nums), next(nums)
    values = [_node_value(nums) for _ in range(children)]
    refs = [next(nums) for _ in range(meta)]
    if not values:
        return sum(refs)
    return sum(values[r - 1] for r in refs if 1 <= r <= len(values))


# Exercise for Advent of Code 2018 day 8.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _meta_sum(iter(int(t) for t in instr.split()))

    @staticmethod
    def two(instr: str) -> int:
        return _node_value(iter(int(t) for t in instr.split()))
