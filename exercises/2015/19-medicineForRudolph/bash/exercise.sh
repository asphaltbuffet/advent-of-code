#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 19.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# Both parts are substring search/replace heavy — bash's weakest area — so the
# work goes to awk, which has index()/substr(), associative-array sets, and
# rand() for the shuffle. WANT selects the part. Part 1 counts the distinct
# molecules from one replacement; part 2 greedily reduces the molecule back to
# "e", reshuffling rules on a dead end.
solve() {
  local input
  read_input input
  awk -v WANT="$1" '
    BEGIN { nr = 0; blank = 0 }
    {
      if ($0 == "") { blank = 1; next }
      if (!blank) {
        # "from => to"
        split($0, p, " => ")
        nr++
        FR[nr] = p[1]
        TO[nr] = p[2]
        ord[nr] = nr
      } else {
        mol = $0
      }
    }
    END {
      if (WANT == 1) { print part1(); exit }
      print part2()
    }
    function part1(    i, j, frm, to, fl, pos, cur, cnt) {
      cnt = 0
      for (i = 1; i <= nr; i++) {
        frm = FR[i]; to = TO[i]; fl = length(frm)
        pos = 1
        while ((j = index(substr(mol, pos), frm)) > 0) {
          pos = pos + j - 1
          cur = substr(mol, 1, pos - 1) to substr(mol, pos + fl)
          if (!(cur in seen)) { seen[cur] = 1; cnt++ }
          pos = pos + 1
        }
      }
      return cnt
    }
    function part2(    cur, steps, stuck, i, k, j, tmp, idx, applied) {
      srand(1)
      while (1) {
        cur = mol; steps = 0; stuck = 0
        while (cur != "e") {
          applied = 0
          for (k = 1; k <= nr; k++) {
            i = ord[k]
            j = index(cur, TO[i])
            if (j > 0) {
              cur = substr(cur, 1, j - 1) FR[i] substr(cur, j + length(TO[i]))
              steps++
              applied = 1
              break
            }
          }
          if (!applied) { stuck = 1; break }
        }
        if (!stuck) return steps
        # Fisher-Yates shuffle of the rule order, then retry.
        for (k = nr; k > 1; k--) {
          idx = int(rand() * k) + 1
          tmp = ord[k]; ord[k] = ord[idx]; ord[idx] = tmp
        }
      }
    }
  ' <<<"$input"
}

part_one() { solve 1; }
part_two() { solve 2; }
