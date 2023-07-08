from typing import *
from aocpy import BaseExercise


# Exercise for Advent of Code 2022 day 1
class Exercise(BaseExercise):
    @staticmethod
    def parse(data):
        sum_value = 0
        calories = []
        lines = data.split("\n")

        for line in lines:
            line = line.strip()

            if line:
                try:
                    n = int(line)
                    sum_value += n
                except ValueError:
                    raise ValueError("Invalid input")

            else:
                calories.append(sum_value)
                sum_value = 0

        calories.sort(reverse=True)

        return calories

    @staticmethod
    def one(instr: str) -> int:
        cal = Exercise.parse(instr)
        if cal is None:
            raise ValueError("Failed to parse input")

        return cal[0]

    @staticmethod
    def two(instr: str) -> int:
        cal = Exercise.parse(instr)
        if cal is None:
            raise ValueError("Failed to parse input")

        return sum(cal[:3])
