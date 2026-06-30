from typing import *
from aocpy import BaseExercise


def _expand(s: str) -> List[int]:
    blocks: List[int] = []
    for i, ch in enumerate(s):
        n = int(ch)
        val = i // 2 if i % 2 == 0 else -1
        blocks.extend([val] * n)
    return blocks


# Exercise for Advent of Code 2024 day 9.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        blocks = _expand(instr.strip())
        l, r = 0, len(blocks) - 1
        while l < r:
            if blocks[l] != -1:
                l += 1
            elif blocks[r] == -1:
                r -= 1
            else:
                blocks[l], blocks[r] = blocks[r], -1
                l += 1
                r -= 1
        return sum(i * v for i, v in enumerate(blocks) if v != -1)

    @staticmethod
    def two(instr: str) -> int:
        s = instr.strip()
        files: List[List[int]] = []  # [id, start, length]
        frees: List[List[int]] = []  # [start, length]
        pos = 0
        for i, ch in enumerate(s):
            n = int(ch)
            if i % 2 == 0:
                files.append([i // 2, pos, n])
            elif n > 0:
                frees.append([pos, n])
            pos += n

        for f in reversed(files):
            _, start, length = f
            for g in frees:
                if g[0] >= start:
                    break
                if g[1] >= length:
                    f[1] = g[0]
                    g[0] += length
                    g[1] -= length
                    break

        checksum = 0
        for fid, start, length in files:
            for k in range(length):
                checksum += (start + k) * fid
        return checksum
