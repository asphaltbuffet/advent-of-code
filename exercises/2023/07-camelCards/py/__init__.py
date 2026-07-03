from collections import Counter
from typing import *
from aocpy import BaseExercise


def _hand_type(hand: str, joker: bool) -> tuple[int, ...]:
    counts = Counter(hand)
    if joker and "J" in counts and len(counts) > 1:
        # Jokers join whichever card is most common, maximizing the type.
        wild = counts.pop("J")
        best = max(counts, key=counts.get)
        counts[best] += wild
    # The sorted multiset of counts uniquely ranks the hand type: five of a kind
    # (5,) beats four (4,1) beats full house (3,2), and so on.
    return tuple(sorted(counts.values(), reverse=True))


def _sort_key(hand: str, joker: bool) -> tuple:
    order = "J23456789TQKA" if joker else "23456789TJQKA"
    return (_hand_type(hand, joker), tuple(order.index(c) for c in hand))


def _winnings(instr: str, joker: bool) -> int:
    hands = []
    for line in instr.splitlines():
        cards, bid = line.split()
        hands.append((_sort_key(cards, joker), int(bid)))
    hands.sort()
    return sum(rank * bid for rank, (_, bid) in enumerate(hands, start=1))


# Exercise for Advent of Code 2023 day 7.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return _winnings(instr, joker=False)

    @staticmethod
    def two(instr: str) -> int:
        return _winnings(instr, joker=True)
