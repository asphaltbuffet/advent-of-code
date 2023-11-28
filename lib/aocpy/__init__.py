from __future__ import annotations
from typing import *
from collections.abc import Sequence


class BaseExercise:
    @staticmethod
    def one(instr: str) -> Any:
        raise NotImplementedError

    @staticmethod
    def two(instr: str) -> Any:
        raise NotImplementedError

    @staticmethod
    def vis(instr: str, outputDir: str) -> Any:
        raise NotImplementedError
