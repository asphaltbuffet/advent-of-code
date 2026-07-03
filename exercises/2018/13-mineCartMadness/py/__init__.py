from typing import *
from aocpy import BaseExercise


# Turn deltas keyed by heading; curves reflect, intersections cycle L/S/R.
def _parse(instr: str) -> Tuple[List[str], List[List[int]]]:
    """Read the grid and lift carts off it, recording the track underneath."""
    grid = instr.rstrip("\n").split("\n")
    carts: List[List[int]] = []
    rows: List[str] = []
    headings = {"<": (-1, 0), ">": (1, 0), "^": (0, -1), "v": (0, 1)}
    under = {"<": "-", ">": "-", "^": "|", "v": "|"}
    for y, line in enumerate(grid):
        chars = list(line)
        for x, ch in enumerate(chars):
            if ch in headings:
                dx, dy = headings[ch]
                # [x, y, dx, dy, turns]
                carts.append([x, y, dx, dy, 0])
                chars[x] = under[ch]
        rows.append("".join(chars))
    return rows, carts


def _at(grid: List[str], x: int, y: int) -> str:
    if 0 <= y < len(grid) and 0 <= x < len(grid[y]):
        return grid[y][x]
    return " "


def _advance(grid: List[str], cart: List[int]) -> None:
    """Step one cell, then turn according to the track landed on."""
    cart[0] += cart[2]
    cart[1] += cart[3]
    dx, dy = cart[2], cart[3]
    cell = _at(grid, cart[0], cart[1])
    if cell == "/":
        cart[2], cart[3] = -dy, -dx
    elif cell == "\\":
        cart[2], cart[3] = dy, dx
    elif cell == "+":
        turn = cart[4] % 3
        if turn == 0:  # left
            cart[2], cart[3] = dy, -dx
        elif turn == 2:  # right
            cart[2], cart[3] = -dy, dx
        cart[4] += 1


def _simulate(grid: List[str], carts: List[List[int]], last_standing: bool) -> str:
    """Run until the first crash, or until one cart remains, as 'x,y'."""
    while True:
        carts.sort(key=lambda c: (c[1], c[0]))
        dead: Set[int] = set()
        for i, c in enumerate(carts):
            if i in dead:
                continue
            _advance(grid, c)
            for j, other in enumerate(carts):
                if i == j or j in dead or other[0] != c[0] or other[1] != c[1]:
                    continue
                if not last_standing:
                    return f"{c[0]},{c[1]}"
                dead.add(i)
                dead.add(j)
                break
        if last_standing:
            carts = [c for k, c in enumerate(carts) if k not in dead]
            if len(carts) == 1:
                return f"{carts[0][0]},{carts[0][1]}"


# Exercise for Advent of Code 2018 day 13.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> str:
        grid, carts = _parse(instr)
        return _simulate(grid, carts, False)

    @staticmethod
    def two(instr: str) -> str:
        grid, carts = _parse(instr)
        return _simulate(grid, carts, True)
