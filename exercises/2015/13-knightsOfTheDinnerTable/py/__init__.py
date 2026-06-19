from typing import *
from aocpy import BaseExercise

from itertools import permutations


def parse(instr: str) -> Tuple[Dict[Tuple[str, str], int], List[str]]:
    """Build happiness[(a, b)] = a's happiness change next to b, plus people."""
    happiness: Dict[Tuple[str, str], int] = {}
    people: Set[str] = set()
    for line in instr.strip().splitlines():
        f = line.rstrip(".").split()
        # <A> would <gain|lose> <N> happiness units by sitting next to <B>
        a, b = f[0], f[10]
        n = int(f[3]) * (1 if f[2] == "gain" else -1)
        happiness[(a, b)] = n
        people.add(a)
    return happiness, sorted(people)


def best(happiness: Dict[Tuple[str, str], int], people: List[str]) -> int:
    """Max total happiness over all circular seatings (first person fixed)."""
    head, rest = people[0], people[1:]
    best_total = None
    for order in permutations(rest):
        seating = [head, *order]
        n = len(seating)
        total = sum(
            happiness.get((seating[i], seating[(i + 1) % n]), 0)
            + happiness.get((seating[(i + 1) % n], seating[i]), 0)
            for i in range(n)
        )
        if best_total is None or total > best_total:
            best_total = total
    return best_total


# Exercise for Advent of Code 2015 day 13.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        happiness, people = parse(instr)
        return best(happiness, people)

    @staticmethod
    def two(instr: str) -> int:
        happiness, people = parse(instr)
        for p in people:
            happiness[("me", p)] = 0
            happiness[(p, "me")] = 0
        return best(happiness, people + ["me"])
