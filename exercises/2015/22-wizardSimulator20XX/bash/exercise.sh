#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 22.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# A recursive branch-and-bound DFS over ~8 state variables is far beyond bash's
# comfort zone, so the whole search runs in awk, which has recursive functions
# and a global best. HARD adds the part-two rule (player loses 1 HP at the start
# of each of their turns). Effects tick at the start of every turn.
solve() {
  local input
  read_input input
  awk -v HARD="$1" '
    /Hit Points/ { bhp0 = $3 }
    /Damage/     { bdmg = $2 }
    END {
      best = 1000000000
      # spell table columns: cost dmg heal shield poison recharge
      ns = split("53 73 113 173 229", scost)
      split("4 2 0 0 0", sdmg); split("0 2 0 0 0", sheal)
      split("0 0 6 0 0", ssh); split("0 0 0 6 0", spo); split("0 0 0 0 5", sre)
      search(50, 500, bhp0, 0, 0, 0, 0)
      print best
    }
    # tick: globals via the passed copies. awk has no by-ref scalars, so the
    # tick is inlined at each call site below instead of a helper.
    function search(php, mana, bhp, shield, poison, recharge, spent,
                    i, np, nm, nb, nsh, npo, nre, nsp, hit, armor) {
      if (spent >= best) return

      # --- Player turn ---
      if (HARD == 1) { php--; if (php <= 0) return }
      # tick effects
      if (poison > 0)   { bhp -= 3; poison-- }
      if (recharge > 0) { mana += 101; recharge-- }
      if (shield > 0)   shield--
      if (bhp <= 0) { if (spent < best) best = spent; return }

      for (i = 1; i <= ns; i++) {
        if (scost[i] > mana) continue
        if (ssh[i] > 0 && shield > 0) continue
        if (spo[i] > 0 && poison > 0) continue
        if (sre[i] > 0 && recharge > 0) continue

        np = php + sheal[i]
        nm = mana - scost[i]
        nb = bhp - sdmg[i]
        nsh = (ssh[i] > 0) ? ssh[i] : shield
        npo = (spo[i] > 0) ? spo[i] : poison
        nre = (sre[i] > 0) ? sre[i] : recharge
        nsp = spent + scost[i]

        if (nb <= 0) { if (nsp < best) best = nsp; continue }

        # --- Boss turn: tick, check, attack ---
        if (npo > 0) { nb -= 3; npo-- }
        if (nre > 0) { nm += 101; nre-- }
        if (nsh > 0) nsh--
        if (nb <= 0) { if (nsp < best) best = nsp; continue }
        armor = (nsh > 0) ? 7 : 0
        hit = bdmg - armor; if (hit < 1) hit = 1
        np -= hit
        if (np <= 0) continue

        search(np, nm, nb, nsh, npo, nre, nsp)
      }
    }
  ' <<<"$input"
}

part_one() { solve 0; }
part_two() { solve 1; }
