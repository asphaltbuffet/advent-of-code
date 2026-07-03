from typing import *
from aocpy import BaseExercise

# Slope tiles and the single direction they permit, as (dr, dc).
SLOPES = {"^": (-1, 0), "v": (1, 0), "<": (0, -1), ">": (0, 1)}
NEIGHBORS = [(-1, 0), (1, 0), (0, -1), (0, 1)]


def _compress(grid: list[str], slippery: bool):
    # Collapse corridors into a weighted graph between junctions (cells with more
    # than two open neighbors), plus start and end.
    rows, cols = len(grid), len(grid[0])
    start = (0, grid[0].index("."))
    end = (rows - 1, grid[-1].index("."))

    def open_cell(r, c):
        return 0 <= r < rows and 0 <= c < cols and grid[r][c] != "#"

    junctions = {start, end}
    for r in range(rows):
        for c in range(cols):
            if not open_cell(r, c):
                continue
            if sum(open_cell(r + dr, c + dc) for dr, dc in NEIGHBORS) > 2:
                junctions.add((r, c))

    # For each junction, walk each corridor to the next junction, tracking length.
    graph: dict[tuple, dict[tuple, int]] = {j: {} for j in junctions}
    for jr, jc in junctions:
        stack = [(jr, jc, 0)]
        seen = {(jr, jc)}
        while stack:
            r, c, dist = stack.pop()
            if dist and (r, c) in junctions:
                graph[(jr, jc)][(r, c)] = dist
                continue
            tile = grid[r][c]
            dirs = [SLOPES[tile]] if slippery and tile in SLOPES else NEIGHBORS
            for dr, dc in dirs:
                nr, nc = r + dr, c + dc
                if open_cell(nr, nc) and (nr, nc) not in seen:
                    seen.add((nr, nc))
                    stack.append((nr, nc, dist + 1))

    return graph, start, end


def _longest(graph, start, end) -> int:
    best = -1
    seen: set[tuple] = set()

    def dfs(node, dist):
        nonlocal best
        if node == end:
            best = max(best, dist)
            return
        seen.add(node)
        for nxt, weight in graph[node].items():
            if nxt not in seen:
                dfs(nxt, dist + weight)
        seen.remove(node)

    dfs(start, 0)
    return best


# Exercise for Advent of Code 2023 day 23.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        graph, start, end = _compress(instr.splitlines(), slippery=True)
        return _longest(graph, start, end)

    @staticmethod
    def two(instr: str) -> int:
        graph, start, end = _compress(instr.splitlines(), slippery=False)
        return _longest(graph, start, end)
