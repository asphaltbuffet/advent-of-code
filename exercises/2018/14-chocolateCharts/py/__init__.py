from typing import *
from aocpy import BaseExercise


# Exercise for Advent of Code 2018 day 14.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> str:
        n = int(instr.strip())

        scores = bytearray((3, 7))
        a, b = 0, 1
        while len(scores) < n + 10:
            total = scores[a] + scores[b]
            if total >= 10:
                scores.append(total // 10)
            scores.append(total % 10)
            a = (a + 1 + scores[a]) % len(scores)
            b = (b + 1 + scores[b]) % len(scores)

        return "".join(str(d) for d in scores[n:n + 10])

    @staticmethod
    def two(instr: str) -> int:
        instr = instr.strip()
        target = bytes(int(c) for c in instr)
        tlen = len(target)

        scores = bytearray((3, 7))
        a, b = 0, 1

        # A step appends one or two recipes. Since only the last one or two
        # digits are new, compare just the two tail windows each step instead of
        # slicing every prefix. Grow the board in bulk, then scan the batch with
        # ``bytes.find`` (a C-level search) so the hot loop stays cheap.
        while True:
            prev = len(scores)
            # Extend the board in a tight loop with locals bound for speed.
            append = scores.append
            for _ in range(1 << 16):
                total = scores[a] + scores[b]
                if total >= 10:
                    append(total // 10)
                append(total % 10)
                n = len(scores)
                a = (a + 1 + scores[a]) % n
                b = (b + 1 + scores[b]) % n

            # Search the newly grown region, overlapping the previous tail so a
            # match straddling the batch boundary is not missed.
            start = max(0, prev - tlen + 1)
            idx = scores.find(target, start)
            if idx != -1:
                return idx
