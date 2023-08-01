from typing import *
from aocpy import BaseExercise


def parse(instr: str) -> Dict[str, str]:
    out = {}

    for line in instr.split("\n"):
        name, action = line.split(": ")
        out[name] = action

    return out


def calc(name: str, raw: Dict[str, str], done: Dict[str, int]):
    if name in done:
        if done[name] == "die die die die die die":
            return 0, False
        return done[name], True

    action = raw[name].split()

    if len(action) == 1:
        done[name] = int(action[0])
    elif len(action) == 3:
        left, ok = calc(action[0], raw, done)
        if not ok:
            return 0, False

        right, ok = calc(action[2], raw, done)
        if not ok:
            return 0, False

        if action[1] == "+":
            done[name] = left + right
        elif action[1] == "-":
            done[name] = left - right
        elif action[1] == "*":
            done[name] = left * right
        elif action[1] == "/":
            done[name] = left / right
        else:
            raise ValueError(f"Unknown operator: {action[1]}")
    else:
        return 0, False

    return int(done[name]), True


# Exercise for Advent of Code 2022 day 21.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        raw = parse(instr)

        v, ok = calc("root", raw, {})
        return v

    @staticmethod
    def two(instr: str) -> int:
        raw = parse(instr)
        results = {}

        # force error when calculating "humn"
        raw["humn"] = "die die die die die die"

        # change root equation to be left / right = 1
        reversed = {"root": "1"}

        orig = raw["root"].split()
        orig[1] = "/"
        raw["root"] = " ".join(orig)
        currentName = "root"

        while currentName != "humn":
            left, operator, right = raw[currentName].split()

            leftVal, left_ok = calc(left, raw, results)
            if left_ok:
                reversed[left] = str(leftVal)

            rightVal, right_ok = calc(right, raw, results)
            if right_ok:
                reversed[right] = str(rightVal)

            if operator == "+":
                if not left_ok:
                    reversed[left] = f"{currentName} - {right}"
                    currentName = left
                elif not right_ok:
                    reversed[right] = f"{currentName} - {left}"
                    currentName = right
                else:
                    raise ValueError("no error path found")
            elif operator == "-":
                if not left_ok:
                    reversed[left] = f"{currentName} + {right}"
                    currentName = left
                elif not right_ok:
                    reversed[right] = f"{left} - {currentName}"
                    currentName = right
                else:
                    raise ValueError("no error path found")
            elif operator == "*":
                if not left_ok:
                    reversed[left] = f"{currentName} / {right}"
                    currentName = left
                elif not right_ok:
                    reversed[right] = f"{currentName} / {left}"
                    currentName = right
                else:
                    raise ValueError("no error path found")
            elif operator == "/":
                if not left_ok:
                    reversed[left] = f"{currentName} * {right}"
                    currentName = left
                elif not right_ok:
                    reversed[right] = f"{left} / {currentName}"
                    currentName = right
                else:
                    raise ValueError("no error path found")
            else:
                raise ValueError("invalid operator")

        v, _ = calc("humn", reversed, {})

        return v
