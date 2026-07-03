from typing import *
from aocpy import BaseExercise


def _parse(instr: str) -> tuple[list[tuple[int, int]], list[int]]:
    inventory, ingredients = instr.split("\n\n")
    ranges = [
        (int(lo), int(hi))
        for lo, hi in (line.split("-") for line in inventory.splitlines())
    ]
    ids = [int(line) for line in ingredients.splitlines()]
    return ranges, ids


def _merge(ranges: list[tuple[int, int]]) -> list[tuple[int, int]]:
    # Coalesce overlapping ranges. Touching-but-disjoint ranges (a gap of one)
    # are left separate, matching the reference's `high >= next.low` test.
    merged: list[list[int]] = []
    for lo, hi in sorted(ranges):
        if merged and merged[-1][1] >= lo:
            merged[-1][1] = max(merged[-1][1], hi)
        else:
            merged.append([lo, hi])
    return [(lo, hi) for lo, hi in merged]


# Exercise for Advent of Code 2025 day 5.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        ranges, ids = _parse(instr)
        merged = _merge(ranges)
        return sum(any(lo <= i <= hi for lo, hi in merged) for i in ids)

    @staticmethod
    def two(instr: str) -> int:
        ranges, _ = _parse(instr)
        return sum(hi - lo + 1 for lo, hi in _merge(ranges))
