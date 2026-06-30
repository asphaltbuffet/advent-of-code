#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 20.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# The divisor sieve runs over ~target/10 houses with a harmonic inner loop
# (tens of millions of additions) — far beyond pure bash, so delegate to one
# awk process holding the houses[] array. PART selects the gift rule: part 1
# gives 10*n to every multiple; part 2 gives 11*n to only the first 50.
# Values stay well under 2^53, so awk's doubles are exact here.
solve() {
  local input
  read_input input
  awk -v PART="$1" '
    {
      target = $1
      limit = int(target / 10) + 1
      for (h = 1; h <= limit; h++) houses[h] = 0
      if (PART == 1) {
        for (n = 1; n <= limit; n++)
          for (h = n; h <= limit; h += n) houses[h] += 10 * n
      } else {
        for (n = 1; n <= limit; n++) {
          c = 0
          for (h = n; h <= limit && c < 50; h += n) { houses[h] += 11 * n; c++ }
        }
      }
      for (h = 1; h <= limit; h++)
        if (houses[h] >= target) { print h; exit }
    }
  ' <<<"$input"
}

part_one() { solve 1; }
part_two() { solve 2; }
