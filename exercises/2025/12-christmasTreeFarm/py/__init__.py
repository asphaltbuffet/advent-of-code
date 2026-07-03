from typing import *
from aocpy import BaseExercise

Cell = tuple[int, int]


def _normalize(cells: Iterable[Cell]) -> frozenset[Cell]:
    cells = set(cells)
    r0 = min(r for r, _ in cells)
    c0 = min(c for _, c in cells)
    return frozenset((r - r0, c - c0) for r, c in cells)


def _orientations(shape: frozenset[Cell]) -> list[list[Cell]]:
    # All eight rotations/reflections a present may be placed in.
    seen: set[frozenset[Cell]] = set()
    cur = set(shape)
    for _ in range(4):
        cur = {(c, -r) for r, c in cur}
        seen.add(_normalize(cur))
        seen.add(_normalize((r, -c) for r, c in cur))
    return [sorted(o) for o in seen]


def _parse(instr: str) -> tuple[list[frozenset[Cell]], list[tuple[int, int, list[int]]]]:
    *shape_blocks, region_block = instr.split("\n\n")
    shapes = [
        frozenset(
            (r, c)
            for r, row in enumerate(block.splitlines()[1:])
            for c, ch in enumerate(row)
            if ch == "#"
        )
        for block in shape_blocks
    ]
    regions = []
    for line in region_block.splitlines():
        head, *counts = line.split()
        w, l = (int(x) for x in head.rstrip(":").split("x"))
        regions.append((w, l, [int(c) for c in counts]))
    return shapes, regions


def _can_pack(w: int, l: int, counts: list[int], shapes: list[frozenset[Cell]]) -> bool:
    total_cells = sum(len(shapes[i]) * n for i, n in enumerate(counts))
    # Necessary condition: the presents' filled cells must not exceed the area.
    if total_cells > w * l:
        return False

    area = w * l
    need = counts[:]
    orients = [_orientations(shapes[i]) if counts[i] else [] for i in range(len(shapes))]
    grid = bytearray(area)

    def backtrack(start: int, remaining: int, holes: int) -> bool:
        if remaining == 0:
            return True
        fe = -1
        for p in range(start, area):
            if not grid[p]:
                fe = p
                break
        if fe == -1:
            return False
        r0, c0 = divmod(fe, w)
        # Cover the first empty cell with some present, trying every orientation
        # and every cell of that present as the one landing on it.
        for si in range(len(shapes)):
            if need[si] == 0:
                continue
            for orient in orients[si]:
                for ar, ac in orient:
                    placed = []
                    ok = True
                    for r, c in orient:
                        x, y = r0 + (r - ar), c0 + (c - ac)
                        if not (0 <= x < l and 0 <= y < w):
                            ok = False
                            break
                        p = x * w + y
                        if p < fe or grid[p]:
                            ok = False
                            break
                        placed.append(p)
                    if ok:
                        for p in placed:
                            grid[p] = 1
                        need[si] -= 1
                        if backtrack(fe + 1, remaining - 1, holes):
                            return True
                        need[si] += 1
                        for p in placed:
                            grid[p] = 0
        # Or leave the first empty cell permanently unused, spending a hole.
        if holes > 0:
            grid[fe] = 2
            if backtrack(fe + 1, remaining, holes - 1):
                return True
            grid[fe] = 0
        return False

    return backtrack(0, sum(counts), area - total_cells)


# Exercise for Advent of Code 2025 day 12.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        shapes, regions = _parse(instr)
        return sum(_can_pack(w, l, counts, shapes) for w, l, counts in regions)

    @staticmethod
    def two(instr: str) -> int:
        # Day 12 is the finale — the last star is granted for the other puzzles.
        return 0
