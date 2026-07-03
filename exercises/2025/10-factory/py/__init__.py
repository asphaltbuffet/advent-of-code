from itertools import combinations
from typing import *
from aocpy import BaseExercise


def _parse(line: str) -> tuple[int, list[int], tuple[int, ...]]:
    tokens = line.split()
    lights = tokens[0].strip("[]")
    # Light i (left to right) is bit i; buttons XOR-toggle a set of these bits.
    target_lights = sum(1 << i for i, c in enumerate(lights) if c == "#")
    joltage = tuple(int(x) for x in tokens[-1].strip("{}").split(","))
    buttons = []
    for tok in tokens[1:-1]:
        mask = 0
        for out in tok.strip("()").split(","):
            idx = int(out)
            if idx < len(lights):
                mask |= 1 << idx
        buttons.append(mask)
    return target_lights, buttons, joltage


def _min_button_presses(target: int, buttons: list[int]) -> int:
    # Fewest buttons whose XOR reproduces the target light pattern; try sizes in
    # increasing order so the first hit is the minimum.
    for k in range(1, len(buttons) + 1):
        for combo in combinations(buttons, k):
            x = 0
            for b in combo:
                x ^= b
            if x == target:
                return k
    return -1


def _min_joltage_presses(buttons: list[int], joltage: tuple[int, ...]) -> int:
    # Reach the joltage vector treating presses in base 2: the buttons pressed at
    # a level set the low bit of every component (parities must match), then the
    # remainder is halved and solved recursively at double cost. Precompute each
    # button-subset's component-sum vector so the recursion just subtracts.
    rlen = len(joltage)
    limit = 1 << len(buttons)
    button_vecs = [[(b >> i) & 1 for i in range(rlen)] for b in buttons]

    mask_sum: list[list[int]] = [[0] * rlen for _ in range(limit)]
    for mask in range(1, limit):
        low = mask & -mask
        b = low.bit_length() - 1
        prev = mask_sum[mask ^ low]
        mask_sum[mask] = [prev[i] + button_vecs[b][i] for i in range(rlen)]
    popcount = [bin(m).count("1") for m in range(limit)]

    memo: dict[tuple[int, ...], int] = {}

    def solve(target: tuple[int, ...]) -> int:
        if not any(target):
            return 0
        if target in memo:
            return memo[target]
        best = -1
        for mask in range(limit):
            sums = mask_sum[mask]
            rem = [target[i] - sums[i] for i in range(rlen)]
            if any(r < 0 or r & 1 for r in rem):
                continue
            sub = solve(tuple(r >> 1 for r in rem))
            if sub != -1:
                total = popcount[mask] + 2 * sub
                if best == -1 or total < best:
                    best = total
        memo[target] = best
        return best

    return solve(joltage)


# Exercise for Advent of Code 2025 day 10.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return sum(
            _min_button_presses(target, buttons)
            for target, buttons, _ in map(_parse, instr.splitlines())
        )

    @staticmethod
    def two(instr: str) -> int:
        return sum(
            _min_joltage_presses(buttons, joltage)
            for _, buttons, joltage in map(_parse, instr.splitlines())
        )
