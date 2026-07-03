from typing import *
from aocpy import BaseExercise


def _deps(instr: str) -> dict[str, set[str]]:
    """Map each step to the set of steps that must precede it."""
    d: dict[str, set[str]] = {}
    for line in instr.strip().splitlines():
        parts = line.split()
        before, after = parts[1], parts[7]
        d.setdefault(before, set())
        d.setdefault(after, set()).add(before)
    return d


# Exercise for Advent of Code 2018 day 7.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> str:
        d = _deps(instr)
        done: set[str] = set()
        order: list[str] = []

        while len(done) < len(d):
            ready = sorted(s for s, pre in d.items() if s not in done and pre <= done)
            nxt = ready[0]
            done.add(nxt)
            order.append(nxt)

        return "".join(order)

    @staticmethod
    def two(instr: str) -> int:
        d = _deps(instr)

        # The small example runs 2 workers with no base cost; the real puzzle
        # runs 5 workers with a 60-second base per step.
        workers, base = (2, 0) if len(d) <= 6 else (5, 60)

        done: set[str] = set()
        in_progress: dict[str, int] = {}  # step -> finish second

        t = 0
        while True:
            # Retire any steps that have finished by the start of this second.
            for step in [s for s, finish in in_progress.items() if finish <= t]:
                done.add(step)
                del in_progress[step]

            if len(done) == len(d):
                return t

            # Assign idle workers to ready steps, alphabetically first.
            ready = sorted(
                s
                for s, pre in d.items()
                if s not in done and s not in in_progress and pre <= done
            )
            for step in ready:
                if len(in_progress) >= workers:
                    break
                in_progress[step] = t + base + (ord(step) - ord("A") + 1)

            t += 1
