from typing import *
from aocpy import BaseExercise
from collections import deque


def parse(instr):
    maxX, maxY, maxZ = 0, 0, 0
    cubes = []

    for line in instr.splitlines():
        c = tuple(map(int, line.split(",")))
        cubes.append(c)

        # print(f"parsed cube={c}")

        maxX = max(maxX, c[0])
        maxY = max(maxY, c[1])
        maxZ = max(maxZ, c[2])

    return (cubes, (maxX, maxY, maxZ))


def faces_that_can_reach_edge(cur, cube_set, bounds):
    count = 0

    adjacent = [
        (0, 0, -1),
        (0, -1, 0),
        (-1, 0, 0),
        (0, 0, 1),
        (0, 1, 0),
        (1, 0, 0),
    ]

    for c in adjacent:
        next_cube = add_cubes(cur, c)

        if can_reach_edge(next_cube, cube_set, bounds):
            count += 1

    return count


def can_reach_edge(start, cube_set, bounds):
    queue = deque([start])
    visited = {}

    adjacent = [
        (0, 0, -1),
        (0, -1, 0),
        (-1, 0, 0),
        (0, 0, 1),
        (0, 1, 0),
        (1, 0, 0),
    ]

    while queue:
        cur_cube = queue.popleft()

        if cur_cube in visited or cur_cube in cube_set:
            continue

        visited[cur_cube] = True

        maxX, maxY, maxZ = bounds

        if (
            cur_cube[0] <= 0
            or cur_cube[0] >= maxX
            or cur_cube[1] <= 0
            or cur_cube[1] >= maxY
            or cur_cube[2] <= 0
            or cur_cube[2] >= maxZ
        ):
            return True

        for a in adjacent:
            next_cube = add_cubes(cur_cube, a)
            queue.append(next_cube)

    return False


def add_cubes(c1, c2):
    return (c1[0] + c2[0], c1[1] + c2[1], c1[2] + c2[2])


# Exercise for Advent of Code 2022 day 18.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        # print("starting part 1")
        cubes, bounds = parse(instr)

        cube_set = set(cubes)

        # start with all cube sides exposed
        total_exposure = len(cubes) * 6

        adjacent = [
            (0, 0, -1),
            (0, -1, 0),
            (-1, 0, 0),
            (0, 0, 1),
            (0, 1, 0),
            (1, 0, 0),
        ]

        for c in cubes:
            for a in adjacent:
                tmp = add_cubes(c, a)

                if tmp in cube_set:
                    total_exposure -= 1

        return total_exposure

    @staticmethod
    def two(instr: str) -> int:
        # print("starting part 2")
        cubes, bounds = parse(instr)

        cube_set = set(cubes)

        # start with all cube sides exposed
        total_exposure = 0

        adjacent = [
            (0, 0, -1),
            (0, -1, 0),
            (-1, 0, 0),
            (0, 0, 1),
            (0, 1, 0),
            (1, 0, 0),
        ]

        for c in cube_set:
            total_exposure += faces_that_can_reach_edge(c, cube_set, bounds)

        return total_exposure
