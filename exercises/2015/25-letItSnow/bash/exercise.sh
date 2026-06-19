#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 25.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

readonly FIRST=20151125 MULT=252533 MOD=33554393

# Pure bash: extract row/col, find the diagonal index n, then modular
# exponentiation. Bash arithmetic is 64-bit; the largest intermediate
# (base*base ~ 1.1e15) stays well under 2^63, so no awk needed.
part_one() {
  local input
  read_input input

  # Pull the two integers from the prose.
  local nums
  read -r -a nums <<<"${input//[!0-9]/ }"
  local row=${nums[0]} col=${nums[1]}

  local diag=$((row + col - 2))
  local n=$((diag * (diag + 1) / 2 + col))

  # code = FIRST * MULT^(n-1) mod MOD, via fast exponentiation.
  local result=1 base=$((MULT % MOD)) exp=$((n - 1))
  while ((exp > 0)); do
    ((exp & 1)) && result=$((result * base % MOD))
    base=$((base * base % MOD))
    exp=$((exp >> 1))
  done

  echo $((FIRST * result % MOD))
}

# Day 25 has no part 2 — the final star comes from finishing the rest.
part_two() {
  echo "Merry Christmas!"
}
