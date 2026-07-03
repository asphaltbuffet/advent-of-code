import heapq
import re

from aocpy import BaseExercise

Bot = tuple[int, int, int, int]  # x, y, z, r


def _parse(instr: str) -> list[Bot]:
    bots: list[Bot] = []
    for line in instr.strip().splitlines():
        x, y, z, r = map(int, re.findall(r"-?\d+", line))
        bots.append((x, y, z, r))
    return bots


def _axis_dist(v: int, lo: int, size: int) -> int:
    """How far v lies outside the interval [lo, lo+size) (0 if inside)."""
    if v < lo:
        return lo - v
    if v > lo + size - 1:
        return v - (lo + size - 1)
    return 0


# Exercise for Advent of Code 2018 day 23.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        bots = _parse(instr)
        sx, sy, sz, sr = max(bots, key=lambda b: b[3])
        return sum(
            1 for x, y, z, _ in bots if abs(x - sx) + abs(y - sy) + abs(z - sz) <= sr
        )

    @staticmethod
    def two(instr: str) -> int:
        bots = _parse(instr)

        # Cube covering every bot, side a power of two.
        span = max(max(abs(x), abs(y), abs(z)) for x, y, z, _ in bots)
        size = 1
        while size < span:
            size *= 2

        def score(x: int, y: int, z: int, s: int) -> tuple:
            in_range = sum(
                1
                for bx, by, bz, br in bots
                if _axis_dist(bx, x, s) + _axis_dist(by, y, s) + _axis_dist(bz, z, s) <= br
            )
            dist = _axis_dist(0, x, s) + _axis_dist(0, y, s) + _axis_dist(0, z, s)
            # Min-heap key: most bots (negated), then closest, then smallest.
            return (-in_range, dist, s, x, y, z)

        heap = [score(-size, -size, -size, 2 * size)]
        while heap:
            _, dist, s, x, y, z = heapq.heappop(heap)
            if s == 1:
                return dist
            half = s // 2
            for dx, dy, dz in (
                (0, 0, 0), (half, 0, 0), (0, half, 0), (0, 0, half),
                (half, half, 0), (half, 0, half), (0, half, half), (half, half, half),
            ):
                heapq.heappush(heap, score(x + dx, y + dy, z + dz, half))
        raise ValueError("no point found")
