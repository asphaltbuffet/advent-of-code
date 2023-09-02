from typing import *
from aocpy import BaseExercise

from itertools import permutations


def parse(instr: str) -> Dict[str, Dict[str, int]]:
    distances = {}
    for line in instr.splitlines():
        c1, _, c2, _, dist = line.split()
        dist = int(dist)
        if c1 not in distances:
            distances[c1] = {}
        if c2 not in distances:
            distances[c2] = {}
        distances[c1][c2] = dist
        distances[c2][c1] = dist

    # print(f"distances: {distances}")
    return distances


def unique_perms(ss: List[str]) -> List[List[str]]:
    """Get all unique permutations of a list of strings."""
    seen = set()
    result = []

    for perm in list(permutations(ss)):
        str_perm = "-".join(perm)

        # print(f"checking {str_perm}")

        if str_perm not in seen:
            result.append(perm)

            rev_str_perm = "-".join(perm[::-1])
            seen.add(str_perm)
            seen.add(rev_str_perm)

    return result


def get_paths(distances: List[str]):
    """Get all valid paths that visit all cities."""
    routes = {}
    paths = unique_perms(distances.keys())
    # print(f"{len(paths)} unique paths found")

    # validate the path is possible
    for path in paths:
        valid = True
        sum = 0
        for i in range(0, len(path) - 1):
            if path[i + 1] not in distances[path[i]]:
                valid = False
                sum = 0
                break

            sum += distances[path[i]][path[i + 1]]

        # print(f"validating: {path} {'✔' if valid else '❌'}")

        if valid:
            # print(f"valid path: {path} = {sum}")
            routes[path] = sum

    # print(f"{len(routes)} valid paths found")

    return routes


# Exercise for Advent of Code 2015 day 9.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        distances = parse(instr)
        routes = get_paths(distances)

        return min(routes.values())

    @staticmethod
    def two(instr: str) -> int:
        distances = parse(instr)
        routes = get_paths(distances)

        return max(routes.values())
