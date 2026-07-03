from collections import deque
from typing import *

from aocpy import BaseExercise

# Neighbor offsets in reading order (up, left, right, down) so BFS naturally
# discovers the reading-order-first shortest path.
READING_DIRS = ((0, -1), (-1, 0), (1, 0), (0, 1))


class Unit:
    __slots__ = ("x", "y", "kind", "hp", "alive")

    def __init__(self, x: int, y: int, kind: str):
        self.x = x
        self.y = y
        self.kind = kind
        self.hp = 200
        self.alive = True


def parse_cave(instr: str) -> Tuple[List[List[str]], List[Unit]]:
    lines = instr.rstrip("\n").split("\n")
    grid = [list(line) for line in lines]
    units: List[Unit] = []
    for y, row in enumerate(grid):
        for x, c in enumerate(row):
            if c in ("E", "G"):
                units.append(Unit(x, y, c))
                row[x] = "."
    return grid, units


def is_open(grid, occ, x: int, y: int) -> bool:
    if y < 0 or y >= len(grid) or x < 0 or x >= len(grid[y]):
        return False
    if grid[y][x] != ".":
        return False
    return (x, y) not in occ


def adjacent_enemy(units: List[Unit], x: int, y: int, kind: str) -> Optional[Unit]:
    best: Optional[Unit] = None
    for dx, dy in READING_DIRS:
        nx, ny = x + dx, y + dy
        for e in units:
            if not e.alive or e.kind == kind or e.x != nx or e.y != ny:
                continue
            if best is None or e.hp < best.hp:
                best = e
    return best


def step_toward(grid, occ, units: List[Unit], u: Unit) -> Optional[Tuple[int, int]]:
    # In-range squares: open floor adjacent to a live enemy.
    in_range = set()
    for e in units:
        if not e.alive or e.kind == u.kind:
            continue
        for dx, dy in READING_DIRS:
            nx, ny = e.x + dx, e.y + dy
            if is_open(grid, occ, nx, ny):
                in_range.add((nx, ny))
    if not in_range:
        return None

    # BFS from the unit, recording distance and the reading-order-first initial
    # step to reach each square. Expanding neighbors in reading order guarantees
    # the first step recorded is reading-order-minimal.
    start = (u.x, u.y)
    dist = {start: 0}
    first_step: Dict[Tuple[int, int], Tuple[int, int]] = {}
    queue = deque([start])

    chosen = None
    found = False
    best_dist = 0

    while queue:
        cur = queue.popleft()
        if found and dist[cur] > best_dist:
            break
        if cur in in_range and cur != start:
            if not found:
                found = True
                best_dist = dist[cur]
                chosen = cur
        cx, cy = cur
        for dx, dy in READING_DIRS:
            np = (cx + dx, cy + dy)
            if not is_open(grid, occ, np[0], np[1]):
                continue
            if np in dist:
                continue
            dist[np] = dist[cur] + 1
            first_step[np] = np if cur == start else first_step[cur]
            queue.append(np)

    if not found:
        return None
    return first_step[chosen]


def combat(instr: str, elf_ap: int, stop_on_elf_death: bool) -> int:
    grid, units = parse_cave(instr)

    rounds = 0
    while True:
        units.sort(key=lambda u: (u.y, u.x))

        for u in units:
            if not u.alive:
                continue

            # Any enemies left? If not, combat ends mid-round.
            if not any(e.alive and e.kind != u.kind for e in units):
                total = sum(e.hp for e in units if e.alive)
                return rounds * total

            occ = {(e.x, e.y): e for e in units if e.alive}

            # Move unless already adjacent to an enemy.
            if adjacent_enemy(units, u.x, u.y, u.kind) is None:
                step = step_toward(grid, occ, units, u)
                if step is not None:
                    u.x, u.y = step

            # Attack the weakest adjacent enemy (reading order breaks HP ties).
            target = adjacent_enemy(units, u.x, u.y, u.kind)
            if target is not None:
                ap = elf_ap if u.kind == "E" else 3
                target.hp -= ap
                if target.hp <= 0:
                    target.alive = False
                    if stop_on_elf_death and target.kind == "E":
                        return -1

        rounds += 1


# Exercise for Advent of Code 2018 day 15.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return combat(instr, 3, False)

    @staticmethod
    def two(instr: str) -> int:
        ap = 4
        while True:
            outcome = combat(instr, ap, True)
            if outcome != -1:
                return outcome
            ap += 1
