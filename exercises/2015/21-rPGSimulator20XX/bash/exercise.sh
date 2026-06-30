#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 21.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# Only 660 loadouts with a closed-form fight, but the hardcoded shop tables and
# the weapon/armor/ring enumeration read more cleanly in awk than in bash
# arithmetic. WANT=1 -> least gold to win; WANT=2 -> most gold and still lose.
solve() {
  local input
  read_input input
  awk -v WANT="$1" '
    function ceildiv(a, b) { return int((a + b - 1) / b) }
    function wins(pd, pa,    bt, pt) {
      bt = ceildiv(bhp, (pd - barm) > 1 ? (pd - barm) : 1)
      pt = ceildiv(100, (bdmg - pa) > 1 ? (bdmg - pa) : 1)
      return bt <= pt
    }
    /Hit Points/ { bhp = $3 }
    /Damage/     { bdmg = $2 }
    /Armor/      { barm = $2 }
    END {
      # Shop tables: cost / damage / armor.
      nw = split("8 10 25 40 74", wc); split("4 5 6 7 8", wd)
      # Armor includes the "none" option (cost 0).
      na = split("0 13 31 53 75 102", ac); split("0 1 2 3 4 5", aar)
      nr = split("25 50 100 20 40 80", rc)
      split("1 2 3 0 0 0", rd); split("0 0 0 1 2 3", rar)

      best = -1; worst = 0
      for (w = 1; w <= nw; w++)
        for (a = 1; a <= na; a++)
          # Ring sets: i==0 means no rings; i>0 picks ring i, and j>i adds a
          # distinct second ring (j==i means only the single ring i).
          for (i = 0; i <= nr; i++)
            for (j = i; j <= (i == 0 ? 0 : nr); j++) {
              cost = wc[w] + ac[a]
              dmg  = wd[w]
              arm  = aar[a]
              if (i > 0) { cost += rc[i]; dmg += rd[i]; arm += rar[i] }
              if (j > i) { cost += rc[j]; dmg += rd[j]; arm += rar[j] }
              if (wins(dmg, arm)) {
                if (best < 0 || cost < best) best = cost
              } else {
                if (cost > worst) worst = cost
              }
            }
      print (WANT == 1 ? best : worst)
    }
  ' <<<"$input"
}

part_one() { solve 1; }
part_two() { solve 2; }
