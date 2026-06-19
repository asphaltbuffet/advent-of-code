#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 18.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# 100x100 grid over 100 steps (~8M cell updates) is far beyond pure-bash
# arithmetic, so delegate the whole Game of Life to one awk process (the day-6
# grid tactic). STUCK=1 forces the four corners on, initially and after each
# step. Step count: the 6x6 example runs 4 (part 1) or 5 (part 2), real 100.
solve() {
  local input
  read_input input
  awk -v STUCK="$1" '
    function stick() {
      g[0,0] = 1; g[0,C-1] = 1; g[R-1,0] = 1; g[R-1,C-1] = 1
    }
    { for (c = 0; c < length($0); c++) g[NR-1, c] = (substr($0,c+1,1) == "#"); C = length($0) }
    END {
      R = NR
      steps = (R > 6) ? 100 : (STUCK ? 5 : 4)
      if (STUCK) stick()
      for (s = 0; s < steps; s++) {
        for (r = 0; r < R; r++) {
          for (c = 0; c < C; c++) {
            on = 0
            for (dr = -1; dr <= 1; dr++)
              for (dc = -1; dc <= 1; dc++) {
                if (dr == 0 && dc == 0) continue
                nr = r + dr; nc = c + dc
                if (nr >= 0 && nr < R && nc >= 0 && nc < C && g[nr,nc]) on++
              }
            if (g[r,c]) nx[r,c] = (on == 2 || on == 3)
            else        nx[r,c] = (on == 3)
          }
        }
        for (r = 0; r < R; r++) for (c = 0; c < C; c++) g[r,c] = nx[r,c]
        if (STUCK) stick()
      }
      total = 0
      for (r = 0; r < R; r++) for (c = 0; c < C; c++) total += g[r,c]
      print total
    }
  ' <<<"$input"
}

part_one() { solve 0; }
part_two() { solve 1; }
