import re

import numpy as np

from aocpy import BaseExercise

_INT = re.compile(r"-?\d+")


def _parse(instr: str) -> tuple[np.ndarray, np.ndarray]:
    # Scan each line's two integers, tolerant of spacing.
    pts = [
        (int(nums[0]), int(nums[1]))
        for line in instr.strip().splitlines()
        if (nums := _INT.findall(line))
    ]
    xy = np.array(pts, dtype=np.int64)
    return xy[:, 0], xy[:, 1]


def _distances(cx: np.ndarray, cy: np.ndarray, pad: int = 0) -> np.ndarray:
    # (H, W, N) Manhattan distances from every grid cell to every coordinate.
    ys = np.arange(cy.min() - pad, cy.max() + pad + 1)
    xs = np.arange(cx.min() - pad, cx.max() + pad + 1)
    gy, gx = np.meshgrid(ys, xs, indexing="ij")
    return np.abs(gx[..., None] - cx) + np.abs(gy[..., None] - cy)


# Exercise for Advent of Code 2018 day 6.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        cx, cy = _parse(instr)
        dist = _distances(cx, cy)

        owner = dist.argmin(axis=-1)
        # Cells whose minimum is shared leave no owner.
        tied = (dist == dist.min(axis=-1, keepdims=True)).sum(axis=-1) > 1
        owner[tied] = -1

        # Any coordinate reaching the border owns an unbounded region.
        border = np.concatenate(
            [owner[0], owner[-1], owner[:, 0], owner[:, -1]]
        )
        infinite = np.unique(border[border >= 0])

        counts = np.bincount(owner[owner >= 0], minlength=cx.size)
        counts[infinite] = 0
        return int(counts.max())

    @staticmethod
    def two(instr: str) -> int:
        cx, cy = _parse(instr)
        threshold = 32 if cx.size <= 10 else 10000
        pad = threshold // cx.size + 1

        total = _distances(cx, cy, pad).sum(axis=-1)
        return int((total < threshold).sum())
