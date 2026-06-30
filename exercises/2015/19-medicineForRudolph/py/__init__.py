import random
from typing import *
from aocpy import BaseExercise


def parse(instr: str) -> Tuple[List[Tuple[str, str]], str]:
    """Split into replacement rules and the target molecule."""
    block, mol = instr.replace("\r\n", "\n").split("\n\n", 1)
    rules = [tuple(line.split(" => ")) for line in block.strip().splitlines()]
    return rules, mol.strip()


# Exercise for Advent of Code 2015 day 19.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        rules, mol = parse(instr)
        seen = set()
        for frm, to in rules:
            # Replace each occurrence of frm independently.
            start = 0
            while (i := mol.find(frm, start)) != -1:
                seen.add(mol[:i] + to + mol[i + len(frm):])
                start = i + 1
        return len(seen)

    @staticmethod
    def two(instr: str) -> int:
        rules, mol = parse(instr)
        rng = random.Random(1)

        # Work backwards greedily: collapse any production back to its source
        # until only "e" remains; reshuffle and retry on a dead end.
        while True:
            cur, steps, stuck = mol, 0, False
            while cur != "e":
                for frm, to in rules:
                    if (i := cur.find(to)) != -1:
                        cur = cur[:i] + frm + cur[i + len(to):]
                        steps += 1
                        break
                else:
                    stuck = True
                    break
            if not stuck:
                return steps
            rng.shuffle(rules)
