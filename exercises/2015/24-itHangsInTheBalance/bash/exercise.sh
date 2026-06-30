#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 24.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# The k-subset search (hundreds of thousands of combinations) is far beyond
# pure-bash arithmetic, so it runs in awk with a recursive pick(). GROUPS is the
# number of equal piles. awk's doubles hold the quantum entanglement (~1e10)
# exactly. Search by increasing first-group size; the smallest size that sums to
# the target gives the fewest packages, and its minimum product is the answer.
solve() {
  local input
  read_input input
  awk -v GROUPS="$1" '
    { w[n++] = $1; total += $1 }
    END {
      target = total / GROUPS
      for (size = 1; size <= n; size++) {
        best = -1
        pick(0, target, size, 1)
        if (best >= 0) { printf "%d\n", best; exit }
      }
    }
    # Choose `count` more weights from index `start`; track the minimum product
    # (qe) of selections that exactly reach `remaining`.
    function pick(start, remaining, count, qe,    i) {
      if (count == 0) {
        if (remaining == 0 && (best < 0 || qe < best)) best = qe
        return
      }
      for (i = start; i <= n - count; i++) {
        if (w[i] <= remaining) pick(i + 1, remaining - w[i], count - 1, qe * w[i])
      }
    }
  ' <<<"$input"
}

part_one() { solve 3; }
part_two() { solve 4; }
