from typing import *
from aocpy import BaseExercise


def calculate_fuel(mass: int, include_fuel: bool) -> int:
    fuel = mass // 3 - 2

    if not include_fuel:
        return fuel
    elif fuel <= 0:
        return 0
    else:
        return fuel + calculate_fuel(fuel, include_fuel)


# Exercise for Advent of Code 2019 day 1.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        sum = 0

        for line in instr.splitlines():
            sum += calculate_fuel(int(line), False)

        return sum

    @staticmethod
    def two(instr: str) -> int:
        sum = 0

        for line in instr.splitlines():
            sum += calculate_fuel(int(line), True)

        return sum
