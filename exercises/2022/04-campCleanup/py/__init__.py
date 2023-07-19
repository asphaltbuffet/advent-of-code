from typing import *
from aocpy import BaseExercise

import re


# Exercise for Advent of Code 2022 day 4
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        count = 0
        for line in instr.split("\n"):
            elf_one, elf_two = parse_pair(line)
            if elf_one and elf_two and is_fully_overlapping(elf_one, elf_two):
                count += 1
        return count

    @staticmethod
    def two(instr: str) -> int:
        count = 0
        for line in instr.split("\n"):
            elf_one, elf_two = parse_pair(line)
            if elf_one and elf_two and is_any_overlapping(elf_one, elf_two):
                count += 1
        return count


class Elf:
    def __init__(self, low, high) -> None:
        self.low = low
        self.high = high


def is_any_overlapping(a, b) -> bool:
    return (a.low <= b.low and a.high >= b.low) or (b.low <= a.low and b.high >= a.low)


def is_fully_overlapping(a, b) -> bool:
    return (a.low <= b.low and a.high >= b.high) or (
        b.low <= a.low and b.high >= a.high
    )


def parse_pair(line) -> Tuple[Elf, Elf]:
    try:
        a_low, a_high, b_low, b_high = map(int, re.findall(r"\d+", line))
        a = Elf(a_low, a_high)
        b = Elf(b_low, b_high)
    except ValueError:
        print(f"error parsing line: {line}")
        a = b = None
    return a, b
