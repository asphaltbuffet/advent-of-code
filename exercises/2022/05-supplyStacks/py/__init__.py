from dataclasses import dataclass
from typing import *
from aocpy import BaseExercise

from queue import LifoQueue as Stack
import re


@dataclass
class Action:
    qty: int
    src: int
    tgt: int


action_regex = re.compile(r"move (\d+) from (\d+) to (\d+)")


def parse(instr: str) -> Tuple[List[Stack], List[Action]]:
    stackstr, actionstr = instr.rstrip().split("\n\n")

    stackInput = stackstr.splitlines()[:-1]
    # stackInput = stackInput[:-1]

    stackInput = [
        [line[i + 1 : i + 2] for i in range(0, len(line), 4)] for line in stackInput
    ]

    stacks: List[Stack] = []
    for _ in range(len(stackInput[0])):
        stacks.append(Stack())

    for line in reversed(stackInput):
        if not line:
            raise ValueError("unexpected blank line")

        for i, crate in enumerate(line):
            if crate != " ":
                stacks[i].put(crate)

    actions: List[Action] = []

    for line in actionstr.strip().splitlines():
        match = action_regex.match(line)

        if not match:
            raise ValueError(f"parsing action line: {line}")

        qty, src, tgt = map(int, match.groups())

        # convert to 0-indexed
        actions.append(Action(qty, src - 1, tgt - 1))

    return stacks, actions


class Exercise(BaseExercise):
    """Exercise for Advent of Code 2022 day 5."""

    @staticmethod
    def one(instr: str) -> int:
        stacks, actions = parse(instr)

        for action in actions:
            for _ in range(action.qty):
                if not stacks[action.src]:
                    return ValueError(f"stack {action.src} is empty")

                c = stacks[action.src].get()
                stacks[action.tgt].put(c)

        return "".join(stack.get() for stack in stacks)

    @staticmethod
    def two(instr: str) -> int:
        stacks, actions = parse(instr)

        for action in actions:
            crates = []

            for _ in range(action.qty):
                crates.append(stacks[action.src].get())

            for c in reversed(crates):
                stacks[action.tgt].put(c)

        return "".join(stack.get() for stack in stacks)
