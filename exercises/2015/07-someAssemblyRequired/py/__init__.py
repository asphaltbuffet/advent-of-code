from typing import *
from aocpy import BaseExercise


def solve(circuit: Dict[str, str], solved: Dict[str, str], wire: str) -> str:
    # print(f"solving for {wire}: {circuit[wire]}")

    if wire in solved:
        return solved[wire]

    raw = []
    for token in circuit[wire].split():
        if token in solved:
            raw.append(solved[token])

        elif token not in circuit:
            # skip operators and numbers
            raw.append(token)

        else:
            # we need to solve for this token
            raw.append(solve(circuit, solved, token))

    # print(f"solved {wire} as: {' '.join(raw)}")
    # print("=====================================")

    # wire is now filled in with actual values
    solved[wire] = str(eval(" ".join(raw)))
    return solved[wire]


# Exercise for Advent of Code 2015 day 7.
class Exercise(BaseExercise):
    @staticmethod
    def one(instr: str) -> int:
        circuit = {}
        solved = {}

        for line in instr.splitlines():
            signal, wire = line.split(" -> ")
            signal = signal.replace("AND", "&")
            signal = signal.replace("OR", "|")
            signal = signal.replace("LSHIFT", "<<")
            signal = signal.replace("RSHIFT", ">>")
            signal = signal.replace("NOT", "~")

            if signal.isdigit():
                solved[wire] = signal

            circuit[wire] = signal

        answer = int(solve(circuit, solved, "a"))

        # need to get the answer; use two's complement if it's negative
        if answer < 0:
            answer = answer + (1 << 16)

        return answer

    @staticmethod
    def two(instr: str) -> int:
        circuit = {}
        solved = {}

        for line in instr.splitlines():
            signal, wire = line.split(" -> ")
            signal = signal.replace("AND", "&")
            signal = signal.replace("OR", "|")
            signal = signal.replace("LSHIFT", "<<")
            signal = signal.replace("RSHIFT", ">>")
            signal = signal.replace("NOT", "~")

            if signal.isdigit():
                solved[wire] = signal

            circuit[wire] = signal

        # override b
        solved["b"] = "46065"
        del circuit["b"]
        answer = int(solve(circuit, solved, "a"))

        # need to get the answer; use two's complement if it's negative
        if answer < 0:
            answer = answer + (1 << 16)

        return answer
