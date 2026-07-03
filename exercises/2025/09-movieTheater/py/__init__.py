from itertools import combinations
from typing import *
from aocpy import BaseExercise

Point = tuple[int, int]


def _points(instr: str) -> list[Point]:
    return [(int(x), int(y)) for x, y in (l.split(",") for l in instr.splitlines())]


def _area(a: Point, b: Point) -> int:
    # Inclusive lattice-cell area of the axis-aligned rectangle spanning a and b.
    return (abs(a[0] - b[0]) + 1) * (abs(a[1] - b[1]) + 1)


def _edges(pts: list[Point]) -> Iterator[tuple[Point, Point]]:
    a = pts[-1]
    for b in pts:
        yield a, b
        a = b


def _on_segment(a: Point, b: Point, p: Point) -> bool:
    # a-b is axis-aligned (rectilinear polygon); is p on it?
    if a[1] == b[1] == p[1]:
        return min(a[0], b[0]) <= p[0] <= max(a[0], b[0])
    if a[0] == b[0] == p[0]:
        return min(a[1], b[1]) <= p[1] <= max(a[1], b[1])
    return False


def _on_edge(pts: list[Point], p: Point) -> bool:
    return any(_on_segment(a, b, p) for a, b in _edges(pts))


def _inside(pts: list[Point], p: Point) -> bool:
    # Ray-cast parity test.
    px, py = p
    inside = False
    for a, b in _edges(pts):
        if (a[1] > py) != (b[1] > py):
            x_cross = (b[0] - a[0]) * (py - a[1]) // (b[1] - a[1]) + a[0]
            if px < x_cross:
                inside = not inside
    return inside


def _edge_crosses_rect(a: Point, b: Point, lo: Point, hi: Point) -> bool:
    # Does a polygon edge pass through the open interior of the rectangle?
    if a[1] == b[1]:  # horizontal edge
        return (
            lo[1] < a[1] < hi[1]
            and max(a[0], b[0]) > lo[0]
            and min(a[0], b[0]) < hi[0]
        )
    return (  # vertical edge
        lo[0] < a[0] < hi[0]
        and max(a[1], b[1]) > lo[1]
        and min(a[1], b[1]) < hi[1]
    )


def _fits(pts: list[Point], p1: Point, p2: Point) -> bool:
    corners = [p1, (p1[0], p2[1]), p2, (p2[0], p1[1])]
    if any(not _on_edge(pts, c) and not _inside(pts, c) for c in corners):
        return False
    lo = (min(p1[0], p2[0]), min(p1[1], p2[1]))
    hi = (max(p1[0], p2[0]), max(p1[1], p2[1]))
    return not any(_edge_crosses_rect(a, b, lo, hi) for a, b in _edges(pts))


# Exercise for Advent of Code 2025 day 9.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        pts = _points(instr)
        return max(_area(a, b) for a, b in combinations(pts, 2))

    @staticmethod
    def two(instr: str) -> int:
        # Largest rectangle spanned by a vertex pair that lies wholly inside the
        # rectilinear polygon: every corner on-edge or interior, and no polygon
        # edge crossing its interior.
        pts = _points(instr)
        best = 0
        for a, b in combinations(pts, 2):
            area = _area(a, b)
            if area > best and _fits(pts, a, b):
                best = area
        return best
