from typing import *
from aocpy import BaseExercise


# Exercise for Advent of Code 2015 day 6.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        lights = {}
        for i in range(1000):
            for j in range(1000):
                lights[(i, j)] = False

        for line in instr.splitlines():
            if line.startswith("turn on"):
                _, _, start, _, end = line.split()
                # print(f"on: {start} -> {end}")
                start = tuple(map(int, start.split(",")))
                end = tuple(map(int, end.split(",")))
                for i in range(start[0], end[0] + 1):
                    for j in range(start[1], end[1] + 1):
                        lights[(i, j)] = True
            elif line.startswith("turn off"):
                _, _, start, _, end = line.split()
                # print(f"off: {start} -> {end}")
                start = tuple(map(int, start.split(",")))
                end = tuple(map(int, end.split(",")))
                for i in range(start[0], end[0] + 1):
                    for j in range(start[1], end[1] + 1):
                        lights[(i, j)] = False
            elif line.startswith("toggle"):
                _, start, _, end = line.split()
                # print(f"toggle: {start} -> {end}")
                start = tuple(map(int, start.split(",")))
                end = tuple(map(int, end.split(",")))
                for i in range(start[0], end[0] + 1):
                    for j in range(start[1], end[1] + 1):
                        lights[(i, j)] = not lights[(i, j)]

        return sum([1 for light in lights.values() if light])

    @staticmethod
    def two(instr: str) -> int:
        lights = {}
        for i in range(1000):
            for j in range(1000):
                lights[(i, j)] = 0

        for line in instr.splitlines():
            if line.startswith("turn on"):
                _, _, start, _, end = line.split()
                # print(f"on: {start} -> {end}")
                start = tuple(map(int, start.split(",")))
                end = tuple(map(int, end.split(",")))
                for i in range(start[0], end[0] + 1):
                    for j in range(start[1], end[1] + 1):
                        lights[(i, j)] += 1

            elif line.startswith("turn off"):
                _, _, start, _, end = line.split()
                # print(f"off: {start} -> {end}")
                start = tuple(map(int, start.split(",")))
                end = tuple(map(int, end.split(",")))
                for i in range(start[0], end[0] + 1):
                    for j in range(start[1], end[1] + 1):
                        lights[(i, j)] -= 1
                        if lights[(i, j)] < 0:
                            lights[(i, j)] = 0

            elif line.startswith("toggle"):
                _, start, _, end = line.split()
                # print(f"toggle: {start} -> {end}")
                start = tuple(map(int, start.split(",")))
                end = tuple(map(int, end.split(",")))
                for i in range(start[0], end[0] + 1):
                    for j in range(start[1], end[1] + 1):
                        lights[(i, j)] += 2

        # 952172 (too low)
        return sum(lights.values())
