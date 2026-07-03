import re

from aocpy import BaseExercise

INT = re.compile(r"-?\d+")


def _parse(instr: str) -> list[tuple[int, ...]]:
    pts = []
    for line in instr.strip().splitlines():
        nums = INT.findall(line)
        if len(nums) == 4:
            pts.append(tuple(int(n) for n in nums))
    return pts


def _manhattan(a: tuple[int, ...], b: tuple[int, ...]) -> int:
    return sum(abs(x - y) for x, y in zip(a, b))


# Exercise for Advent of Code 2018 day 25.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        pts = _parse(instr)
        parent = list(range(len(pts)))

        def find(x: int) -> int:
            while parent[x] != x:
                parent[x] = parent[parent[x]]  # path compression
                x = parent[x]
            return x

        for i in range(len(pts)):
            for j in range(i + 1, len(pts)):
                if _manhattan(pts[i], pts[j]) <= 3:
                    parent[find(i)] = find(j)

        return len({find(i) for i in range(len(pts))})

    @staticmethod
    def two(instr: str) -> str:
        # Day 25 finale: no second puzzle, only the closing message.
        return "Merry Christmas!"
