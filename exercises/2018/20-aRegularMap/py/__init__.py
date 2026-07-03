from aocpy import BaseExercise

# Compass moves as complex-number deltas: real = east/west, imag = south/north.
MOVES = {"N": -1j, "S": 1j, "E": 1, "W": -1}


def _walk(instr: str) -> dict[complex, int]:
    """Follow the route regex, returning the shortest door-distance to each room.

    Every door adds exactly one step, so tracking the running distance and keeping
    the minimum on revisits gives the true shortest distances directly.
    """
    route = instr.strip()
    dist: dict[complex, int] = {0: 0}
    stack: list[complex] = []
    pos = 0 + 0j

    for c in route:
        if c == "(":
            stack.append(pos)
        elif c == "|":
            pos = stack[-1]
        elif c == ")":
            pos = stack.pop()
        elif c in MOVES:
            prev = pos
            pos += MOVES[c]
            nd = dist[prev] + 1
            if pos not in dist or nd < dist[pos]:
                dist[pos] = nd
    return dist


# Exercise for Advent of Code 2018 day 20.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return max(_walk(instr).values())

    @staticmethod
    def two(instr: str) -> int:
        return sum(1 for d in _walk(instr).values() if d >= 1000)
