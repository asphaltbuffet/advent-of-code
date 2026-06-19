from typing import *
from aocpy import BaseExercise


def parse_boss(instr: str) -> Tuple[int, int]:
    vals = [int(line.split(": ")[1]) for line in instr.strip().splitlines()]
    return vals[0], vals[1]  # hp, damage


# Spells: (cost, instant boss damage, instant self heal, shield, poison, recharge)
# where the last three set an effect timer when nonzero.
SPELLS = [
    (53, 4, 0, 0, 0, 0),    # Magic Missile
    (73, 2, 2, 0, 0, 0),    # Drain
    (113, 0, 0, 6, 0, 0),   # Shield
    (173, 0, 0, 0, 6, 0),   # Poison
    (229, 0, 0, 0, 0, 5),   # Recharge
]


def solve(instr: str, hard: bool) -> int:
    bhp0, bdmg = parse_boss(instr)
    best = [1 << 30]

    def tick(mana, bhp, shield, poison, recharge):
        """Apply active effects once and decrement their timers."""
        if poison > 0:
            bhp -= 3
            poison -= 1
        if recharge > 0:
            mana += 101
            recharge -= 1
        if shield > 0:
            shield -= 1
        return mana, bhp, shield, poison, recharge

    def search(php, mana, bhp, shield, poison, recharge, spent):
        if spent >= best[0]:
            return

        # Player turn.
        if hard:
            php -= 1
            if php <= 0:
                return
        mana, bhp, shield, poison, recharge = tick(mana, bhp, shield, poison, recharge)
        if bhp <= 0:
            best[0] = min(best[0], spent)
            return

        for cost, dmg, heal, sh, po, re in SPELLS:
            if cost > mana:
                continue
            if (sh and shield > 0) or (po and poison > 0) or (re and recharge > 0):
                continue

            nphp = php + heal
            nmana = mana - cost
            nbhp = bhp - dmg
            nshield = sh if sh else shield
            npoison = po if po else poison
            nrecharge = re if re else recharge
            nspent = spent + cost

            if nbhp <= 0:
                best[0] = min(best[0], nspent)
                continue

            # Boss turn.
            nmana, nbhp, nshield, npoison, nrecharge = tick(
                nmana, nbhp, nshield, npoison, nrecharge)
            if nbhp <= 0:
                best[0] = min(best[0], nspent)
                continue
            hit = max(1, bdmg - (7 if nshield > 0 else 0))
            nphp -= hit
            if nphp <= 0:
                continue

            search(nphp, nmana, nbhp, nshield, npoison, nrecharge, nspent)

    search(50, 500, bhp0, 0, 0, 0, 0)
    return best[0]


# Exercise for Advent of Code 2015 day 22.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return solve(instr, False)

    @staticmethod
    def two(instr: str) -> int:
        return solve(instr, True)
