from collections import deque
from typing import *
from aocpy import BaseExercise

Brick = tuple[int, int, int, int, int, int]  # x1,y1,z1,x2,y2,z2


def _settle(instr: str) -> tuple[list[list[int]], list[set[int]]]:
    # Parse bricks and let them fall in z order; return, per brick index, the
    # set of bricks it rests on (supports) and the set it holds up (held).
    bricks: list[list[int]] = [
        [int(n) for n in line.replace("~", ",").split(",")]
        for line in instr.splitlines()
    ]
    bricks.sort(key=lambda b: b[2])  # lowest first

    # heights[(x, y)] = (top z, brick index) of the topmost settled cell.
    heights: dict[tuple[int, int], tuple[int, int]] = {}
    supports: list[set[int]] = [set() for _ in bricks]  # i rests on these
    held: list[set[int]] = [set() for _ in bricks]       # i holds these up

    for i, (x1, y1, z1, x2, y2, z2) in enumerate(bricks):
        cells = [
            (x, y)
            for x in range(min(x1, x2), max(x1, x2) + 1)
            for y in range(min(y1, y2), max(y1, y2) + 1)
        ]
        rest_z = max((heights[c][0] for c in cells if c in heights), default=0) + 1
        # Whatever tops out at rest_z - 1 becomes a supporter.
        for c in cells:
            if c in heights and heights[c][0] == rest_z - 1:
                j = heights[c][1]
                supports[i].add(j)
                held[j].add(i)

        top = rest_z + (max(z1, z2) - min(z1, z2))
        for c in cells:
            heights[c] = (top, i)

    return supports, held


# Exercise for Advent of Code 2023 day 22.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        supports, held = _settle(instr)
        # A brick is safe to disintegrate unless it is the sole support of some
        # brick it holds up.
        return sum(
            all(len(supports[j]) > 1 for j in held[i])
            for i in range(len(supports))
        )

    @staticmethod
    def two(instr: str) -> int:
        supports, held = _settle(instr)
        total = 0
        for start in range(len(supports)):
            # Chain reaction: a brick falls once all of its supporters have fallen.
            fallen = {start}
            queue = deque([start])
            while queue:
                for j in held[queue.popleft()]:
                    if j not in fallen and supports[j] <= fallen:
                        fallen.add(j)
                        queue.append(j)
            total += len(fallen) - 1  # exclude the disintegrated brick itself
        return total
