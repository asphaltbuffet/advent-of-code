from collections import Counter
from math import prod
from typing import *
from aocpy import BaseExercise


def _junctions(instr: str) -> list[tuple[int, int, int]]:
    return [
        (int(x), int(y), int(z))
        for x, y, z in (line.split(",") for line in instr.splitlines())
    ]


def _edges(pts: list[tuple[int, int, int]]) -> list[tuple[int, int, int]]:
    # (squared distance, i, j) for every pair, sorted shortest first.
    edges = [
        (sum((pts[i][k] - pts[j][k]) ** 2 for k in range(3)), i, j)
        for i in range(len(pts))
        for j in range(i)
    ]
    edges.sort()
    return edges


class _DSU:
    def __init__(self, n: int) -> None:
        self.parent = list(range(n))

    def find(self, x: int) -> int:
        while self.parent[x] != x:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, a: int, b: int) -> bool:
        ra, rb = self.find(a), self.find(b)
        if ra == rb:
            return False
        self.parent[ra] = rb
        return True


# Exercise for Advent of Code 2025 day 8.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        pts = _junctions(instr)
        wires = 10 if len(pts) < 100 else 1000
        dsu = _DSU(len(pts))
        for _, a, b in _edges(pts)[:wires]:
            dsu.union(a, b)
        sizes = Counter(dsu.find(i) for i in range(len(pts))).values()
        return prod(sorted(sizes, reverse=True)[:3])

    @staticmethod
    def two(instr: str) -> int:
        # Kruskal's MST: connect the shortest edges that merge components until the
        # graph is a single tree. The last merging edge's junctions give the
        # answer as the product of their x-coordinates.
        pts = _junctions(instr)
        dsu = _DSU(len(pts))
        connections = 0
        for _, a, b in _edges(pts):
            if dsu.union(a, b):
                connections += 1
                if connections == len(pts) - 1:
                    return pts[a][0] * pts[b][0]
        return -1
