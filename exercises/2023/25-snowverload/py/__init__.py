from collections import defaultdict, deque
from typing import *
from aocpy import BaseExercise


def _parse(instr: str) -> tuple[list[list[int]], int]:
    ids: dict[str, int] = {}
    adj: list[list[int]] = []

    def intern(name: str) -> int:
        if name not in ids:
            ids[name] = len(adj)
            adj.append([])
        return ids[name]

    for line in instr.splitlines():
        left, rights = line.split(": ")
        a = intern(left)
        for right in rights.split():
            b = intern(right)
            adj[a].append(b)
            adj[b].append(a)
    return adj, len(adj)


def _max_flow(adj: list[list[int]], s: int, t: int) -> tuple[int, list[bool]]:
    # Unit-capacity max flow from s to t via BFS augmenting paths. When no more
    # paths exist, the residual-reachable set from s is the source side of the
    # min cut.
    cap: dict[tuple[int, int], int] = defaultdict(int)
    for u, nbrs in enumerate(adj):
        for v in nbrs:
            cap[(u, v)] += 1

    flow = 0
    while True:
        prev = {s: s}
        queue = deque([s])
        while queue:
            u = queue.popleft()
            if u == t:
                break
            for v in adj[u]:
                if v not in prev and cap[(u, v)] > 0:
                    prev[v] = u
                    queue.append(v)
        if t not in prev:
            reachable = [i in prev for i in range(len(adj))]
            return flow, reachable
        # Push one unit along the augmenting path.
        v = t
        while v != s:
            u = prev[v]
            cap[(u, v)] -= 1
            cap[(v, u)] += 1
            v = u
        flow += 1


# Exercise for Advent of Code 2023 day 25.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        adj, n = _parse(instr)
        # The clusters are joined by exactly three wires, so the min cut between
        # two nodes on opposite sides is 3. Fix source 0 and find such a node.
        for t in range(1, n):
            flow, reachable = _max_flow(adj, 0, t)
            if flow == 3:
                size = sum(reachable)
                return size * (n - size)
        return 0

    @staticmethod
    def two(instr: str) -> str:
        # Day 25's final star is granted for completing every other puzzle.
        return ""
