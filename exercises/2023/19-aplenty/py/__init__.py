import re
from math import prod
from typing import *
from aocpy import BaseExercise

# A rule is (category, operator, threshold, target); a bare fallthrough is
# (None, None, None, target).
Rule = tuple[Optional[str], Optional[str], Optional[int], str]


def _parse(instr: str) -> tuple[dict[str, list[Rule]], str]:
    wf_block, part_block = instr.split("\n\n")
    workflows: dict[str, list[Rule]] = {}
    for line in wf_block.splitlines():
        name, body = re.match(r"(\w+)\{(.+)\}", line).groups()
        rules: list[Rule] = []
        for token in body.split(","):
            if ":" in token:
                cond, target = token.split(":")
                rules.append((cond[0], cond[1], int(cond[2:]), target))
            else:
                rules.append((None, None, None, token))
        workflows[name] = rules
    return workflows, part_block


def _accepts(part: dict[str, int], workflows: dict[str, list[Rule]]) -> bool:
    name = "in"
    while name not in ("A", "R"):
        for cat, op, threshold, target in workflows[name]:
            if cat is None or (
                part[cat] < threshold if op == "<" else part[cat] > threshold
            ):
                name = target
                break
    return name == "A"


def _count(ranges: dict[str, tuple[int, int]], name: str,
           workflows: dict[str, list[Rule]]) -> int:
    # Recursively split the 4D range through a workflow, counting accepted volume.
    if name == "R":
        return 0
    if name == "A":
        return prod(hi - lo + 1 for lo, hi in ranges.values())

    total = 0
    ranges = dict(ranges)
    for cat, op, threshold, target in workflows[name]:
        if cat is None:
            total += _count(ranges, target, workflows)
            break
        lo, hi = ranges[cat]
        # Split [lo, hi] into the part matching this rule and the remainder.
        if op == "<":
            matched, rest = (lo, min(hi, threshold - 1)), (max(lo, threshold), hi)
        else:
            matched, rest = (max(lo, threshold + 1), hi), (lo, min(hi, threshold))
        if matched[0] <= matched[1]:
            total += _count({**ranges, cat: matched}, target, workflows)
        if rest[0] > rest[1]:
            break  # nothing falls through to later rules
        ranges[cat] = rest
    return total


# Exercise for Advent of Code 2023 day 19.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        workflows, part_block = _parse(instr)
        total = 0
        for line in part_block.splitlines():
            part = {m[0]: int(m[1]) for m in re.findall(r"(\w)=(\d+)", line)}
            if _accepts(part, workflows):
                total += sum(part.values())
        return total

    @staticmethod
    def two(instr: str) -> int:
        workflows, _ = _parse(instr)
        full = {cat: (1, 4000) for cat in "xmas"}
        return _count(full, "in", workflows)
