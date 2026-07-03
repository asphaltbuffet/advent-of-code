from typing import *

import numpy as np
from aocpy import BaseExercise

GRID_SIZE = 300


def _summed_area_table(serial: int) -> np.ndarray:
    """Build a summed-area table so any square's total power is an O(1) lookup.

    ``sat[y, x]`` holds the sum of all cells from (1,1) to (x,y) inclusive, with
    a zero-filled first row and column so edge squares need no special-casing.
    """
    coords = np.arange(1, GRID_SIZE + 1)
    x, y = np.meshgrid(coords, coords)  # 1-indexed grid, shape (GRID_SIZE, GRID_SIZE)
    rack_id = x + 10
    power = (rack_id * y + serial) * rack_id
    cells = (power // 100) % 10 - 5

    sat = np.zeros((GRID_SIZE + 1, GRID_SIZE + 1), dtype=np.int64)
    sat[1:, 1:] = cells.cumsum(axis=0).cumsum(axis=1)
    return sat


def _best(sat: np.ndarray, size: int) -> Tuple[int, int, int]:
    """Top-left corner and total of the highest-power square of ``size``.

    Vectorized square-sum over every top-left corner via the summed-area table.
    """
    n = GRID_SIZE - size + 1
    sums = (
        sat[size:, size:]
        - sat[:n, size:]
        - sat[size:, :n]
        + sat[:n, :n]
    )
    idx = int(sums.argmax())
    y, x = divmod(idx, n)
    return x + 1, y + 1, int(sums.flat[idx])


class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> str:
        serial = int(instr.strip())
        sat = _summed_area_table(serial)
        x, y, _ = _best(sat, 3)
        return f"{x},{y}"

    @staticmethod
    def two(instr: str) -> str:
        serial = int(instr.strip())
        sat = _summed_area_table(serial)

        # Scan every square size; each search is O(1) per corner thanks to the
        # summed-area table, so the whole thing is O(GRID_SIZE^3) but vectorized.
        best = (0, 0, 0, -(1 << 62))
        for size in range(1, GRID_SIZE + 1):
            x, y, total = _best(sat, size)
            if total > best[3]:
                best = (x, y, size, total)
        return f"{best[0]},{best[1]},{best[2]}"
