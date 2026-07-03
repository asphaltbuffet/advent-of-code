from aocpy import BaseExercise


def _parse(instr: str) -> tuple[set[int], set[str]]:
    """Return the set of planted pot indices and the neighborhoods that grow a plant."""
    lines = instr.strip().split("\n")
    initial = lines[0].split("initial state:", 1)[-1].strip()
    pots = {i for i, c in enumerate(initial) if c == "#"}
    rules = set()
    for line in lines[1:]:
        line = line.strip()
        if not line:
            continue
        lhs, rhs = line.split(" => ")
        if rhs == "#":
            rules.add(lhs)
    return pots, rules


def _step(pots: set[int], rules: set[str]) -> set[int]:
    """Advance one generation; only pots within two of a live pot can change."""
    lo, hi = min(pots), max(pots)
    nxt = set()
    for i in range(lo - 2, hi + 3):
        window = "".join("#" if j in pots else "." for j in range(i - 2, i + 3))
        if window in rules:
            nxt.add(i)
    return nxt


def _shape(pots: set[int]) -> tuple[str, int]:
    """Pattern normalized to start at the leftmost pot, plus that offset."""
    lo, hi = min(pots), max(pots)
    return "".join("#" if i in pots else "." for i in range(lo, hi + 1)), lo


# Exercise for Advent of Code 2018 day 12.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        pots, rules = _parse(instr)
        for _ in range(20):
            pots = _step(pots, rules)
        return sum(pots)

    @staticmethod
    def two(instr: str) -> int:
        pots, rules = _parse(instr)
        target = 50_000_000_000

        # Once a shape repeats, the planted count is fixed and the pattern only
        # drifts sideways by a constant offset per generation, so the index sum
        # grows linearly and we extrapolate to the target generation.
        seen: dict[str, tuple[int, int]] = {}
        for gen in range(target):
            sh, lo = _shape(pots)
            if sh in seen:
                prev_gen, prev_lo = seen[sh]
                period = gen - prev_gen
                drift = lo - prev_lo
                remaining = target - gen
                return sum(pots) + remaining * len(pots) * drift // period
            seen[sh] = (gen, lo)
            pots = _step(pots, rules)
        return sum(pots)
