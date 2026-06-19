#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 15.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# The composition search over 100 teaspoons does ~176k scorings; bash's
# per-operation overhead makes that ~20s/part. Delegate the whole enumeration to
# a single awk process (one long-lived interpreter, no per-leaf forks) — the
# established tactic for compute-heavy days. CALORIE_GOAL=-1 disables the filter.
#
# awk parses each ingredient's five signed ints positionally (it sees the same
# space-delimited tokens bash would), recurses over compositions of 100, scores
# each complete recipe (negative property totals clamped to 0), and prints the
# max.
solve() {
  local input
  read_input input
  awk -v GOAL="$1" '
    function score(amounts,    p, cap, dur, fla, tex, cal, i) {
      cap = dur = fla = tex = cal = 0
      for (i = 1; i <= n; i++) {
        cap += C[i] * amounts[i]
        dur += D[i] * amounts[i]
        fla += F[i] * amounts[i]
        tex += T[i] * amounts[i]
        cal += K[i] * amounts[i]
      }
      if (GOAL >= 0 && cal != GOAL) return 0
      if (cap < 0) cap = 0
      if (dur < 0) dur = 0
      if (fla < 0) fla = 0
      if (tex < 0) tex = 0
      return cap * dur * fla * tex
    }
    function rec(i, remaining, amounts,    a, s) {
      if (i == n) {
        amounts[i] = remaining
        s = score(amounts)
        if (s > max) max = s
        return
      }
      for (a = 0; a <= remaining; a++) {
        amounts[i] = a
        rec(i + 1, remaining - a, amounts)
      }
    }
    {
      # capacity N, durability N, flavor N, texture N, calories N (values follow
      # each keyword); pull the five signed ints in order.
      c = 0
      for (i = 1; i <= NF; i++) {
        t = $i; gsub(/,/, "", t)
        if (t ~ /^-?[0-9]+$/) { c++; v[c] = t + 0 }
      }
      n++
      C[n] = v[1]; D[n] = v[2]; F[n] = v[3]; T[n] = v[4]; K[n] = v[5]
    }
    END {
      max = 0
      rec(1, 100, amounts)
      print max
    }
  ' <<<"$input"
}

part_one() { solve -1; }
part_two() { solve 500; }
