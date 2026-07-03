import re
from collections import Counter
from typing import *

from aocpy import BaseExercise

CLAIM_RE = re.compile(r"-?\d+")


def parse(instr: str) -> list[tuple[int, int, int, int, int]]:
    """Parse each claim by scanning its five integers, which is robust to
    whitespace and delivery quirks in the input."""
    claims = []
    for line in instr.splitlines():
        if not line.strip():
            continue
        nums = [int(n) for n in CLAIM_RE.findall(line)]
        assert len(nums) == 5, f"expected 5 numbers in {line!r}"
        claims.append(tuple(nums))
    return claims


def coverage(claims) -> Counter:
    cover = Counter()
    for _id, left, top, width, height in claims:
        cover.update(
            (x, y)
            for x in range(left, left + width)
            for y in range(top, top + height)
        )
    return cover


# Exercise for Advent of Code 2018 day 3.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        cover = coverage(parse(instr))
        return sum(1 for n in cover.values() if n >= 2)

    @staticmethod
    def two(instr: str) -> int:
        claims = parse(instr)
        cover = coverage(claims)
        for _id, left, top, width, height in claims:
            if all(
                cover[(x, y)] == 1
                for x in range(left, left + width)
                for y in range(top, top + height)
            ):
                return _id
        raise ValueError("no non-overlapping claim found")
