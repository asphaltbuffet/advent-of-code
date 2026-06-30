from collections import defaultdict
from itertools import combinations
from typing import *
from aocpy import BaseExercise


def _parse(instr: str) -> Tuple[int, int, Dict[str, List[Tuple[int, int]]]]:
    grid = instr.split()
    rows = len(grid)
    cols = len(grid[0]) if grid else 0
    antennas: Dict[str, List[Tuple[int, int]]] = defaultdict(list)
    for r, row in enumerate(grid):
        for c, ch in enumerate(row):
            if ch != ".":
                antennas[ch].append((r, c))
    return rows, cols, antennas


# Exercise for Advent of Code 2024 day 8.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        rows, cols, antennas = _parse(instr)

        def in_bounds(r: int, c: int) -> bool:
            return 0 <= r < rows and 0 <= c < cols

        nodes: Set[Tuple[int, int]] = set()
        for pts in antennas.values():
            for (ar, ac), (br, bc) in combinations(pts, 2):
                dr, dc = br - ar, bc - ac
                for nr, nc in ((ar - dr, ac - dc), (br + dr, bc + dc)):
                    if in_bounds(nr, nc):
                        nodes.add((nr, nc))
        return len(nodes)

    @staticmethod
    def two(instr: str) -> int:
        rows, cols, antennas = _parse(instr)

        def in_bounds(r: int, c: int) -> bool:
            return 0 <= r < rows and 0 <= c < cols

        nodes: Set[Tuple[int, int]] = set()
        for pts in antennas.values():
            for (ar, ac), (br, bc) in combinations(pts, 2):
                dr, dc = br - ar, bc - ac
                for step in (-1, 1):
                    nr, nc = ar, ac
                    while in_bounds(nr, nc):
                        nodes.add((nr, nc))
                        nr, nc = nr + step * dr, nc + step * dc
        return len(nodes)
