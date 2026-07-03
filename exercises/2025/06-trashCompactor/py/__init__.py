from math import prod
from typing import *
from aocpy import BaseExercise


def _apply(op: str, nums: list[int]) -> int:
    return sum(nums) if op == "+" else prod(nums)


# Exercise for Advent of Code 2025 day 6.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        # Whitespace-delimited tokens form a column table: the last row is the
        # operators, and each column's numbers are combined by its operator.
        *number_rows, op_row = instr.rstrip("\n").splitlines()
        ops = op_row.split()
        columns: list[list[int]] = [[] for _ in ops]
        for row in number_rows:
            for i, tok in enumerate(row.split()):
                columns[i].append(int(tok))
        return sum(_apply(op, col) for op, col in zip(ops, columns))

    @staticmethod
    def two(instr: str) -> int:
        # Read digits stacked vertically in each character column, right to left.
        # Operator columns delimit problems; a blank column separates them.
        lines = instr.rstrip("\n").splitlines()
        *number_rows, op_row = lines
        width = max(map(len, lines))
        number_rows = [row.ljust(width) for row in number_rows]
        op_row = op_row.ljust(width)

        total = 0
        nums: list[int] = []
        for col in range(width - 1, -1, -1):
            digits = "".join(r[col] for r in number_rows if r[col] != " ")
            if digits:
                nums.append(int(digits))
            if op_row[col] != " ":
                total += _apply(op_row[col], nums)
                nums = []
        return total
