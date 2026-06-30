#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 6.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

# The grid is 1000x1000 with ~300 rectangle ops (up to ~300M cell updates).
# A pure-bash associative-array grid would run for hours, so the hot loop is
# delegated to a single awk process that holds the grid as a flat array keyed by
# x*1000+y. `mode` selects part 1 (boolean) or part 2 (brightness).
solve() {
  local mode="$1"
  awk -v mode="$mode" '
    {
      # Determine op (0=off, 1=on, 2=toggle) and shift coords to fields $2,$3.
      if ($1 == "turn") {
        op = ($2 == "on") ? 1 : 0
        c1 = $3; c2 = $5
      } else {            # toggle
        op = 2
        c1 = $2; c2 = $4
      }
      split(c1, a, ",")
      split(c2, b, ",")
      x1 = a[1]; y1 = a[2]
      x2 = b[1]; y2 = b[2]

      for (x = x1; x <= x2; x++) {
        base = x * 1000
        for (y = y1; y <= y2; y++) {
          k = base + y
          if (mode == 1) {                      # part 1: boolean
            if (op == 1)      g[k] = 1
            else if (op == 0) g[k] = 0
            else              g[k] = !g[k]
          } else {                              # part 2: brightness
            if (op == 1)      g[k]++
            else if (op == 0) { if (g[k] > 0) g[k]-- }
            else              g[k] += 2
          }
        }
      }
    }
    END {
      total = 0
      if (mode == 1) {
        for (k in g) if (g[k]) total++
      } else {
        for (k in g) total += g[k]
      }
      print total
    }
  '
}

part_one() {
  solve 1
}

part_two() {
  solve 2
}
