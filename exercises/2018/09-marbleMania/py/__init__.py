import re
from collections import deque

from aocpy import BaseExercise


def _high_score(players: int, last: int) -> int:
    """Play the full game and return the winning player's score.

    A ``deque`` models the ring: rotating it keeps the current marble at the
    right end, so the normal two-clockwise insertion and the scoring
    seven-counter-clockwise removal are both O(1). ``deque.rotate`` is the
    idiomatic, fast tool here — a plain list would shift elements and turn the
    part-two game into an O(n^2) crawl.
    """
    scores = [0] * players
    ring = deque([0])

    for m in range(1, last + 1):
        if m % 23 == 0:
            ring.rotate(7)
            scores[m % players] += m + ring.pop()
            ring.rotate(-1)
        else:
            ring.rotate(-1)
            ring.append(m)

    return max(scores)


# Exercise for Advent of Code 2018 day 9.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        players, last = map(int, re.findall(r"\d+", instr))
        return _high_score(players, last)

    @staticmethod
    def two(instr: str) -> int:
        players, last = map(int, re.findall(r"\d+", instr))
        return _high_score(players, last * 100)
