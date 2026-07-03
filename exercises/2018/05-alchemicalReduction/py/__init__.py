from string import ascii_lowercase

from aocpy import BaseExercise


def react(polymer: str, skip: str = "") -> int:
    """Collapse the polymer with a list as a stack: push each unit, but pop when
    it is the same letter as the top in the opposite case. Units whose lowercase
    form equals ``skip`` are dropped; an empty ``skip`` keeps every unit."""
    stack: list[str] = []

    for unit in polymer:
        if skip and unit.lower() == skip:
            continue
        if stack and stack[-1] != unit and stack[-1] == unit.swapcase():
            stack.pop()
        else:
            stack.append(unit)

    return len(stack)


# Exercise for Advent of Code 2018 day 5.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return react(instr.strip())

    @staticmethod
    def two(instr: str) -> int:
        polymer = instr.strip()
        return min(react(polymer, unit) for unit in ascii_lowercase)
