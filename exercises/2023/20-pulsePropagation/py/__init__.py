from collections import deque
from math import lcm
from typing import *
from aocpy import BaseExercise


def _parse(instr: str):
    kinds: dict[str, str] = {}       # name -> '%', '&', or '' (broadcaster)
    outputs: dict[str, list[str]] = {}
    for line in instr.splitlines():
        left, right = line.split(" -> ")
        name = left.lstrip("%&")
        kinds[name] = left[0] if left[0] in "%&" else ""
        outputs[name] = right.split(", ")

    # Flip-flop on/off state and conjunction memory of last pulse per input.
    flip: dict[str, bool] = {n: False for n, k in kinds.items() if k == "%"}
    conj: dict[str, dict[str, int]] = {n: {} for n, k in kinds.items() if k == "&"}
    for src, dsts in outputs.items():
        for dst in dsts:
            if dst in conj:
                conj[dst][src] = 0

    return kinds, outputs, flip, conj


def _press(kinds, outputs, flip, conj, watch: str = "") -> tuple[int, int, set[str]]:
    # Run one button press. Return (#low, #high) and the set of sources that sent
    # a high pulse into the `watch` module during this press.
    low = high = 0
    high_into_watch: set[str] = set()
    queue = deque([("button", "broadcaster", 0)])

    while queue:
        src, name, pulse = queue.popleft()
        if pulse:
            high += 1
            if name == watch:
                high_into_watch.add(src)
        else:
            low += 1

        kind = kinds.get(name)
        if kind is None:
            continue  # untyped sink (rx, output, ...)

        if kind == "%":
            if pulse:
                continue
            flip[name] = not flip[name]
            out = 1 if flip[name] else 0
        elif kind == "&":
            conj[name][src] = pulse
            out = 0 if all(conj[name].values()) else 1
        else:  # broadcaster
            out = pulse

        for dst in outputs[name]:
            queue.append((name, dst, out))

    return low, high, high_into_watch


# Exercise for Advent of Code 2023 day 20.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        kinds, outputs, flip, conj = _parse(instr)
        low_total = high_total = 0
        for _ in range(1000):
            low, high, _ = _press(kinds, outputs, flip, conj)
            low_total += low
            high_total += high
        return low_total * high_total

    @staticmethod
    def two(instr: str) -> int:
        kinds, outputs, flip, conj = _parse(instr)
        # rx receives from one conjunction; it emits low only when all of that
        # conjunction's inputs are high. Each input cycles, so the first press at
        # which each sends a high pulse gives its period, and the answer is the LCM.
        feeder = next(src for src, dsts in outputs.items() if "rx" in dsts)
        inputs = set(conj[feeder])
        periods: dict[str, int] = {}

        press = 0
        while len(periods) < len(inputs):
            press += 1
            _, _, highs = _press(kinds, outputs, flip, conj, watch=feeder)
            for src in highs:
                periods.setdefault(src, press)

        return lcm(*periods.values())
