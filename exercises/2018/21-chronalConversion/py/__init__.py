from collections.abc import Iterator

from aocpy import BaseExercise

MASK = 0xFFFFFF


def _constants(instr: str) -> tuple[int, int]:
    """Read the two per-input magic numbers out of the program: the seed the
    compared register is reset to, and the multiplier applied to it."""
    prog = [line.split() for line in instr.strip().splitlines() if line[0] != "#"]
    # The value register is the operand of `eqrr … 0` that isn't register 0.
    val_reg = next(
        (b if a == "0" else a) for op, a, b, _ in prog if op == "eqrr" and "0" in (a, b)
    )
    seed = max(int(a) for op, a, _, c in prog if op == "seti" and c == val_reg)
    mult = next(int(b) for op, _, b, c in prog if op == "muli" and c == val_reg)
    return seed, mult


def _halt_values(instr: str) -> Iterator[int]:
    """Yield each distinct value register 0 is compared against, in order, until
    the sequence repeats. The inner loop is evaluated directly rather than
    interpreted opcode by opcode."""
    seed, mult = _constants(instr)
    acc = 0
    seen: set[int] = set()
    while True:
        hi = acc | 0x10000
        acc = seed
        while True:
            acc = ((acc + (hi & 0xFF)) & MASK) * mult & MASK
            if hi < 256:
                break
            hi //= 256
        if acc in seen:
            return
        seen.add(acc)
        yield acc


# Exercise for Advent of Code 2018 day 21.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        return next(_halt_values(instr))

    @staticmethod
    def two(instr: str) -> int:
        last = 0
        for last in _halt_values(instr):
            pass
        return last
