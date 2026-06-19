#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 17.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# The 2^20 subset sweep is far too much for pure-bash arithmetic, so delegate
# the whole enumeration to one awk process (the established perf tactic). awk
# reads the container sizes, sweeps every bitmask, and tallies a histogram of
# valid-subset counts by popcount. It prints two lines: part 1 (total) and
# part 2 (count at the minimum container count). `solve` selects the line.
solve() {
  local input
  read_input input
  awk -v WANT="$1" '
    { sz[n++] = $1 }
    END {
      target = (n <= 5) ? 25 : 150
      total = 0
      for (mask = 0; mask < 2 ^ n; mask++) {
        sum = 0; bits = 0
        for (i = 0; i < n; i++) {
          if (int(mask / 2 ^ i) % 2 == 1) { sum += sz[i]; bits++ }
        }
        if (sum == target) { cnt[bits]++; total++ }
      }
      if (WANT == 1) { print total; exit }
      for (k = 1; k <= n; k++) if (cnt[k] > 0) { print cnt[k]; exit }
    }
  ' <<<"$input"
}

part_one() { solve 1; }
part_two() { solve 2; }
