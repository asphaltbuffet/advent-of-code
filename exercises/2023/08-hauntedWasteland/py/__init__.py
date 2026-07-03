import re
from itertools import cycle
from math import lcm
from typing import *
from aocpy import BaseExercise


def _parse(instr: str) -> tuple[str, dict[str, tuple[str, str]]]:
    moves, network = instr.split("\n\n")
    graph = {}
    for node, left, right in re.findall(r"(\w+) = \((\w+), (\w+)\)", network):
        graph[node] = (left, right)
    return moves, graph


def _steps(start: str, moves: str, graph: dict[str, tuple[str, str]],
           done: Callable[[str], bool]) -> int:
    node = start
    for count, move in enumerate(cycle(moves), start=1):
        node = graph[node][0] if move == "L" else graph[node][1]
        if done(node):
            return count


# Exercise for Advent of Code 2023 day 8.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        moves, graph = _parse(instr)
        return _steps("AAA", moves, graph, lambda n: n == "ZZZ")

    @staticmethod
    def two(instr: str) -> int:
        moves, graph = _parse(instr)
        # Each ghost cycles independently; they align at the LCM of their lengths.
        starts = [n for n in graph if n.endswith("A")]
        return lcm(*(
            _steps(start, moves, graph, lambda n: n.endswith("Z"))
            for start in starts
        ))
