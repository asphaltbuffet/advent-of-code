import re
from math import prod
from typing import *
from aocpy import BaseExercise

TEASPOONS = 100


def parse(instr: str) -> List[List[int]]:
    """Each ingredient becomes [capacity, durability, flavor, texture, calories]."""
    return [
        [int(n) for n in re.findall(r"-?\d+", line)]
        for line in instr.strip().splitlines()
    ]


def best(ingredients: List[List[int]], calorie_goal: Optional[int]) -> int:
    """Maximum score over every distribution of TEASPOONS across ingredients.
    If calorie_goal is set, only recipes hitting that calorie count count."""
    n = len(ingredients)
    amounts = [0] * n
    highest = 0

    def rec(i: int, remaining: int) -> None:
        nonlocal highest
        if i == n - 1:
            amounts[i] = remaining
            s = score(ingredients, amounts, calorie_goal)
            if s > highest:
                highest = s
            return
        for a in range(remaining + 1):
            amounts[i] = a
            rec(i + 1, remaining - a)

    rec(0, TEASPOONS)
    return highest


def score(ingredients: List[List[int]], amounts: List[int], calorie_goal: Optional[int]) -> int:
    # Sum each property column weighted by amount; calories is the 5th column.
    totals = [
        sum(ing[p] * a for ing, a in zip(ingredients, amounts))
        for p in range(5)
    ]
    if calorie_goal is not None and totals[4] != calorie_goal:
        return 0
    # Negative property totals clamp to 0; calories excluded from the product.
    return prod(max(0, t) for t in totals[:4])


# Exercise for Advent of Code 2015 day 15.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return best(parse(instr), None)

    @staticmethod
    def two(instr: str) -> int:
        return best(parse(instr), 500)
