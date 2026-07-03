import re
from copy import deepcopy
from dataclasses import dataclass, field
from itertools import count

from aocpy import BaseExercise

LINE = re.compile(
    r"(\d+) units each with (\d+) hit points (?:\((.*?)\) )?"
    r"with an attack that does (\d+) (\w+) damage at initiative (\d+)"
)


@dataclass
class Group:
    army: int  # 0 = immune system, 1 = infection
    units: int
    hp: int
    damage: int
    attack: str
    initiative: int
    weak: frozenset[str] = field(default_factory=frozenset)
    immune: frozenset[str] = field(default_factory=frozenset)

    @property
    def power(self) -> int:
        return self.units * self.damage

    def damage_to(self, other: "Group") -> int:
        if self.attack in other.immune:
            return 0
        return self.power * (2 if self.attack in other.weak else 1)


def _parse(instr: str) -> list[Group]:
    groups: list[Group] = []
    army = 0
    for line in instr.strip().splitlines():
        line = line.strip()
        if line.startswith("Immune System"):
            army = 0
        elif line.startswith("Infection"):
            army = 1
        elif (m := LINE.match(line)):
            units, hp, mods, dmg, atk, init = m.groups()
            weak: set[str] = set()
            immune: set[str] = set()
            for clause in (mods or "").split("; "):
                if clause.startswith("weak to "):
                    weak = set(clause[len("weak to "):].split(", "))
                elif clause.startswith("immune to "):
                    immune = set(clause[len("immune to "):].split(", "))
            groups.append(
                Group(army, int(units), int(hp), int(dmg), atk, int(init),
                      frozenset(weak), frozenset(immune))
            )
    return groups


def _fight(groups: list[Group]) -> tuple[int, int]:
    """Run to completion. Returns (winning army, its surviving units); a round that
    kills nobody is a stalemate, reported as an infection win."""
    while True:
        living = [g for g in groups if g.units > 0]
        by_army = {0: 0, 1: 0}
        for g in living:
            by_army[g.army] += g.units
        if by_army[0] == 0 or by_army[1] == 0:
            return (0, by_army[0]) if by_army[0] else (1, by_army[1])

        # Target selection.
        targets: dict[int, Group] = {}
        taken: set[int] = set()
        for atk in sorted(living, key=lambda g: (-g.power, -g.initiative)):
            candidates = [
                d for d in living
                if d.army != atk.army and id(d) not in taken and atk.damage_to(d) > 0
            ]
            if candidates:
                best = max(candidates, key=lambda d: (atk.damage_to(d), d.power, d.initiative))
                targets[id(atk)] = best
                taken.add(id(best))

        # Attack in decreasing initiative order.
        killed = 0
        for atk in sorted(living, key=lambda g: -g.initiative):
            if atk.units <= 0 or id(atk) not in targets:
                continue
            defender = targets[id(atk)]
            dead = min(defender.units, atk.damage_to(defender) // defender.hp)
            defender.units -= dead
            killed += dead
        if killed == 0:
            return 1, by_army[1]  # stalemate


# Exercise for Advent of Code 2018 day 24.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _fight(_parse(instr))[1]

    @staticmethod
    def two(instr: str) -> int:
        base = _parse(instr)
        for boost in count(0):
            groups = deepcopy(base)
            for g in groups:
                if g.army == 0:
                    g.damage += boost
            winner, units = _fight(groups)
            if winner == 0:
                return units
