from typing import *
from aocpy import BaseExercise

Rng = tuple[int, int]  # half-open interval [start, end)


def _parse(instr: str) -> tuple[list[int], list[list[tuple[int, int, int]]]]:
    seed_block, *map_blocks = instr.split("\n\n")
    seeds = [int(n) for n in seed_block.split(":")[1].split()]
    maps = []
    for block in map_blocks:
        rules = []
        for line in block.splitlines()[1:]:
            dst, src, length = (int(n) for n in line.split())
            rules.append((dst, src, length))
        maps.append(rules)
    return seeds, maps


def _map_value(value: int, rules: list[tuple[int, int, int]]) -> int:
    for dst, src, length in rules:
        if src <= value < src + length:
            return dst + (value - src)
    return value


def _map_ranges(ranges: list[Rng], rules: list[tuple[int, int, int]]) -> list[Rng]:
    # Split each interval against every rule; unmatched fragments pass through.
    out: list[Rng] = []
    work = list(ranges)
    while work:
        start, end = work.pop()
        for dst, src, length in rules:
            lo = max(start, src)
            hi = min(end, src + length)
            if lo < hi:  # overlap maps through this rule
                out.append((lo - src + dst, hi - src + dst))
                if start < lo:
                    work.append((start, lo))
                if hi < end:
                    work.append((hi, end))
                break
        else:
            out.append((start, end))  # no rule matched: identity
    return out


# Exercise for Advent of Code 2023 day 5.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        seeds, maps = _parse(instr)
        best = None
        for seed in seeds:
            value = seed
            for rules in maps:
                value = _map_value(value, rules)
            best = value if best is None else min(best, value)
        return best

    @staticmethod
    def two(instr: str) -> int:
        seeds, maps = _parse(instr)
        ranges = [
            (start, start + length)
            for start, length in zip(seeds[::2], seeds[1::2])
        ]
        for rules in maps:
            ranges = _map_ranges(ranges, rules)
        return min(start for start, _ in ranges)
