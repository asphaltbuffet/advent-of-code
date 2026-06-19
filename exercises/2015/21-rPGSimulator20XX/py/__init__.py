from itertools import combinations
from math import ceil
from typing import *
from aocpy import BaseExercise

PLAYER_HP = 100

# Shop tables (cost, damage, armor) from the problem statement, not the input.
WEAPONS = [(8, 4, 0), (10, 5, 0), (25, 6, 0), (40, 7, 0), (74, 8, 0)]
ARMORS = [(0, 0, 0), (13, 0, 1), (31, 0, 2), (53, 0, 3), (75, 0, 4), (102, 0, 5)]
RINGS = [(25, 1, 0), (50, 2, 0), (100, 3, 0), (20, 0, 1), (40, 0, 2), (80, 0, 3)]


def parse_boss(instr: str) -> Tuple[int, int, int]:
    vals = [int(line.split(": ")[1]) for line in instr.strip().splitlines()]
    return vals[0], vals[1], vals[2]  # hp, damage, armor


def loadouts() -> Iterator[Tuple[int, int, int]]:
    """Yield (cost, damage, armor) for every legal equipment combination:
    one weapon, zero/one armor, and zero/one/two distinct rings."""
    ring_sets = [()]
    for k in (1, 2):
        ring_sets.extend(combinations(RINGS, k))
    for w in WEAPONS:
        for a in ARMORS:
            for rs in ring_sets:
                items = (w, a, *rs)
                yield (sum(i[0] for i in items),
                       sum(i[1] for i in items),
                       sum(i[2] for i in items))


def player_wins(p_dmg: int, p_arm: int, b_hp: int, b_dmg: int, b_arm: int) -> bool:
    # Each side needs ceil(targetHP / effectiveDamage) hits; player strikes
    # first, so ties go to the player.
    boss_turns = ceil(b_hp / max(1, p_dmg - b_arm))
    player_turns = ceil(PLAYER_HP / max(1, b_dmg - p_arm))
    return boss_turns <= player_turns


# Exercise for Advent of Code 2015 day 21.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        b = parse_boss(instr)
        return min(cost for cost, dmg, arm in loadouts() if player_wins(dmg, arm, *b))

    @staticmethod
    def two(instr: str) -> int:
        b = parse_boss(instr)
        return max(cost for cost, dmg, arm in loadouts() if not player_wins(dmg, arm, *b))
