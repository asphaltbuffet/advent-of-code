from functools import cmp_to_key
from typing import *
from aocpy import BaseExercise


def _parse(instr: str) -> Tuple[Set[Tuple[int, int]], List[List[int]]]:
    rules_block, updates_block = instr.strip("\n").split("\n\n", 1)
    rules = set()
    for line in rules_block.splitlines():
        a, b = line.split("|")
        rules.add((int(a), int(b)))
    updates = [
        [int(p) for p in line.split(",")]
        for line in updates_block.splitlines()
        if line
    ]
    return rules, updates


def _ordered(rules: Set[Tuple[int, int]], pages: List[int]) -> bool:
    for i in range(len(pages)):
        for j in range(i + 1, len(pages)):
            if (pages[j], pages[i]) in rules:
                return False
    return True


# Exercise for Advent of Code 2024 day 5.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        rules, updates = _parse(instr)
        return sum(p[len(p) // 2] for p in updates if _ordered(rules, p))

    @staticmethod
    def two(instr: str) -> int:
        rules, updates = _parse(instr)

        def cmp(a: int, b: int) -> int:
            if (a, b) in rules:
                return -1
            if (b, a) in rules:
                return 1
            return 0

        total = 0
        for pages in updates:
            if _ordered(rules, pages):
                continue
            fixed = sorted(pages, key=cmp_to_key(cmp))
            total += fixed[len(fixed) // 2]
        return total
