from typing import *
from aocpy import BaseExercise


# Exercise for Advent of Code 2015 day 20.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        target = int(instr.strip())
        limit = target // 10 + 1
        # House h receives 10 * (sum of divisors of h). Sieve each elf's gift
        # onto its multiples; the bound target/10 suffices since elf h alone
        # delivers 10*h to house h.
        houses = [0] * (limit + 1)
        for n in range(1, limit + 1):
            for h in range(n, limit + 1, n):
                houses[h] += 10 * n
        for h in range(1, limit + 1):
            if houses[h] >= target:
                return h
        return -1

    @staticmethod
    def two(instr: str) -> int:
        target = int(instr.strip())
        limit = target // 10 + 1
        # Each elf gives 11*n but visits only its first 50 multiples.
        houses = [0] * (limit + 1)
        for n in range(1, limit + 1):
            for h in range(n, min(n * 50, limit) + 1, n):
                houses[h] += 11 * n
        for h in range(1, limit + 1):
            if houses[h] >= target:
                return h
        return -1
