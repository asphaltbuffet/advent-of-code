from fractions import Fraction
from itertools import combinations
from typing import *
from aocpy import BaseExercise

Hail = tuple[tuple[int, int, int], tuple[int, int, int]]


def _parse(instr: str) -> list[Hail]:
    hail = []
    for line in instr.splitlines():
        pos, vel = line.split(" @ ")
        p = tuple(int(n) for n in pos.split(","))
        v = tuple(int(n) for n in vel.split(","))
        hail.append((p, v))
    return hail


def _crosses_xy(a: Hail, b: Hail, lo: int, hi: int) -> bool:
    (px, py, _), (vx, vy, _) = a
    (qx, qy, _), (wx, wy, _) = b
    denom = vx * wy - vy * wx
    if denom == 0:
        return False  # parallel
    # Solve for the two path parameters t (for a) and s (for b).
    t = Fraction((qx - px) * wy - (qy - py) * wx, denom)
    s = Fraction((qx - px) * vy - (qy - py) * vx, denom)
    if t < 0 or s < 0:
        return False  # crossing in the past for one of them
    x = px + vx * t
    y = py + vy * t
    return lo <= x <= hi and lo <= y <= hi


def _solve_rock(hail: list[Hail]) -> tuple[int, int, int]:
    # Relative to the rock, every hailstone passes through the origin, so
    # (P_i - R) is parallel to (V_i - W): (P_i - R) x (V_i - W) = 0. Subtracting
    # the equation for stone 0 from stones i cancels the R x W term, leaving a
    # linear system in (R, W). Three stones give six equations for six unknowns.
    (p0, v0), (p1, v1), (p2, v2) = hail[0], hail[1], hail[2]

    # Build the 6x6 linear system A x = b with unknowns [Rx,Ry,Rz,Wx,Wy,Wz].
    rows = []
    rhs = []

    def add_pair(pa, va, pb, vb):
        # From (Pa - R) x (Va - W) = (Pb - R) x (Vb - W), expand each axis pair.
        for i, j in ((0, 1), (1, 2), (2, 0)):
            # Coefficients for the R and W unknowns in the (i, j) component of
            # (Pa - R) x (Va - W) = (Pb - R) x (Vb - W), after the R x W terms
            # cancel between the two stones.
            row = [0] * 6
            row[i] += vb[j] - va[j]
            row[j] += va[i] - vb[i]
            row[3 + i] += pa[j] - pb[j]
            row[3 + j] += pb[i] - pa[i]
            b = (
                pb[i] * vb[j] - pb[j] * vb[i]
                - (pa[i] * va[j] - pa[j] * va[i])
            )
            rows.append([Fraction(c) for c in row])
            rhs.append(Fraction(b))

    add_pair(p0, v0, p1, v1)
    add_pair(p0, v0, p2, v2)

    # Gaussian elimination over exact fractions.
    n = 6
    for col in range(n):
        pivot = next(r for r in range(col, n) if rows[r][col] != 0)
        rows[col], rows[pivot] = rows[pivot], rows[col]
        rhs[col], rhs[pivot] = rhs[pivot], rhs[col]
        inv = rows[col][col]
        rows[col] = [c / inv for c in rows[col]]
        rhs[col] /= inv
        for r in range(n):
            if r != col and rows[r][col] != 0:
                f = rows[r][col]
                rows[r] = [a - f * b for a, b in zip(rows[r], rows[col])]
                rhs[r] -= f * rhs[col]

    rx, ry, rz = rhs[0], rhs[1], rhs[2]
    return int(rx), int(ry), int(rz)


# Exercise for Advent of Code 2023 day 24.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        hail = _parse(instr)
        # The example uses a small test area; the real puzzle a huge one. Pick by
        # coordinate magnitude so both pass.
        big = any(abs(c) > 1_000_000 for (p, _) in hail for c in p)
        lo, hi = (200000000000000, 400000000000000) if big else (7, 27)
        return sum(
            _crosses_xy(a, b, lo, hi) for a, b in combinations(hail, 2)
        )

    @staticmethod
    def two(instr: str) -> int:
        rx, ry, rz = _solve_rock(_parse(instr))
        return rx + ry + rz
