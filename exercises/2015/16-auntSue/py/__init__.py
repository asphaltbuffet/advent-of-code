import operator
import re
from typing import *
from aocpy import BaseExercise

# MFCSAM ticker-tape reading from the problem statement (not the input).
TARGET = {
    "children": 3, "cats": 7, "samoyeds": 2, "pomeranians": 3, "akitas": 0,
    "vizslas": 0, "goldfish": 5, "trees": 3, "cars": 2, "perfumes": 1,
}

# Part 2: cats/trees are lower bounds, pomeranians/goldfish upper bounds; the
# rest still match exactly.
RANGED = {"cats": operator.gt, "trees": operator.gt,
          "pomeranians": operator.lt, "goldfish": operator.lt}


def parse(instr: str) -> List[Dict[str, int]]:
    """Each Sue becomes a dict of the compounds she lists -> count."""
    sues = []
    for line in instr.strip().splitlines():
        obs = {k: int(v) for k, v in re.findall(r"(\w+): (\d+)", line.split(": ", 1)[1])}
        sues.append(obs)
    return sues


def find(sues: List[Dict[str, int]], cmp: Callable[[str, int], bool]) -> int:
    for i, obs in enumerate(sues, start=1):
        if all(cmp(k, got) for k, got in obs.items()):
            return i
    return -1


# Exercise for Advent of Code 2015 day 16.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return find(parse(instr), lambda k, got: got == TARGET[k])

    @staticmethod
    def two(instr: str) -> int:
        return find(parse(instr),
                    lambda k, got: RANGED.get(k, operator.eq)(got, TARGET[k]))
