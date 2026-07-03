import heapq
import re

from aocpy import BaseExercise

TORCH = 1  # start and finish equipped with the torch
# How far past the target to extend the grid; a shortest path never detours much
# further (a 7-minute tool switch bounds any useful detour).
MARGIN = 50


def _parse(instr: str) -> tuple[int, int, int]:
    depth, tx, ty = map(int, re.findall(r"\d+", instr))
    return depth, tx, ty


def _regions(depth: int, tx: int, ty: int, w: int, h: int) -> list[list[int]]:
    """Region type (erosion % 3) over the (w+1) x (h+1) grid."""
    erosion = [[0] * (w + 1) for _ in range(h + 1)]
    for y in range(h + 1):
        for x in range(w + 1):
            if (x, y) in ((0, 0), (tx, ty)):
                geo = 0
            elif y == 0:
                geo = x * 16807
            elif x == 0:
                geo = y * 48271
            else:
                geo = erosion[y][x - 1] * erosion[y - 1][x]
            erosion[y][x] = (geo + depth) % 20183
    return [[e % 3 for e in row] for row in erosion]


# Exercise for Advent of Code 2018 day 22.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        depth, tx, ty = _parse(instr)
        region = _regions(depth, tx, ty, tx, ty)
        return sum(region[y][x] for y in range(ty + 1) for x in range(tx + 1))

    @staticmethod
    def two(instr: str) -> int:
        depth, tx, ty = _parse(instr)
        w, h = tx + MARGIN, ty + MARGIN
        region = _regions(depth, tx, ty, w, h)

        # Dijkstra over (x, y, tool): move costs 1, tool switch costs 7, and a tool
        # is allowed in a region iff it differs from the region's type.
        dist: dict[tuple[int, int, int], int] = {(0, 0, TORCH): 0}
        pq = [(0, 0, 0, TORCH)]
        while pq:
            cost, x, y, tool = heapq.heappop(pq)
            if (x, y, tool) == (tx, ty, TORCH):
                return cost
            if cost > dist.get((x, y, tool), 1 << 30):
                continue
            # Switch to the other tool allowed here.
            for t in range(3):
                if t != region[y][x] and t != tool:
                    nd = cost + 7
                    if nd < dist.get((x, y, t), 1 << 30):
                        dist[(x, y, t)] = nd
                        heapq.heappush(pq, (nd, x, y, t))
            # Move to a neighbor keeping the current tool.
            for nx, ny in ((x + 1, y), (x - 1, y), (x, y + 1), (x, y - 1)):
                if 0 <= nx <= w and 0 <= ny <= h and region[ny][nx] != tool:
                    nd = cost + 1
                    if nd < dist.get((nx, ny, tool), 1 << 30):
                        dist[(nx, ny, tool)] = nd
                        heapq.heappush(pq, (nd, nx, ny, tool))
        raise ValueError("target unreachable")
