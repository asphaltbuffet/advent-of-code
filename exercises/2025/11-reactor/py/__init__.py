from functools import cache
from typing import *
from aocpy import BaseExercise


def _parse(instr: str) -> dict[str, list[str]]:
    graph: dict[str, list[str]] = {}
    for line in instr.splitlines():
        head, *tail = line.split()
        graph[head.rstrip(":")] = tail
    return graph


def _count_paths(graph: dict[str, list[str]], start: str, waypoints: tuple[str, ...]) -> int:
    # Count paths from `start` that visit `waypoints` in order, ending at the last
    # one. State is (current node, remaining waypoints); memoized.
    @cache
    def trace(cur: str, keys: tuple[str, ...]) -> int:
        rest = keys
        if cur in keys:
            if len(keys) == 1:  # reached the final destination
                return 1
            if cur != keys[0]:  # hit a later waypoint out of order
                return 0
            rest = keys[1:]  # cleared the next waypoint; require the remainder
        return sum(trace(nxt, rest) for nxt in graph.get(cur, ()))

    return trace(start, waypoints)


def _reaches(graph: dict[str, list[str]], src: str, dst: str) -> bool:
    stack, seen = [src], set()
    while stack:
        node = stack.pop()
        if node == dst:
            return True
        if node in seen:
            continue
        seen.add(node)
        stack.extend(graph.get(node, ()))
    return False


# Exercise for Advent of Code 2025 day 11.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _count_paths(_parse(instr), "you", ("out",))

    @staticmethod
    def two(instr: str) -> int:
        # Thread both waypoints in whichever order is actually reachable.
        graph = _parse(instr)
        order = ("dac", "fft") if _reaches(graph, "dac", "fft") else ("fft", "dac")
        return _count_paths(graph, "svr", (*order, "out"))
