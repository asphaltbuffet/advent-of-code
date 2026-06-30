#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 10.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# Look-and-say grows the string ~30% per round, reaching millions of chars by
# round 50. A per-char bash scan would be hopelessly slow, so the hot loop is
# delegated to a single awk process that iterates the transform $rounds times
# and prints the final length. $1 = number of rounds.
solve() {
  local rounds="$1" input
  read_input input
  awk -v s="$input" -v rounds="$rounds" '
    BEGIN {
      for (r = 0; r < rounds; r++) {
        out = ""
        n = length(s)
        i = 1
        while (i <= n) {
          c = substr(s, i, 1)
          j = i + 1
          while (j <= n && substr(s, j, 1) == c) j++
          out = out (j - i) c
          i = j
        }
        s = out
      }
      print length(s)
    }
  '
}

part_one() {
  solve 40
}

part_two() {
  solve 50
}
