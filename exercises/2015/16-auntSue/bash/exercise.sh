#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 16.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# MFCSAM ticker-tape reading from the problem statement (not the input).
declare -A target=(
  [children]=3 [cats]=7 [samoyeds]=2 [pomeranians]=3 [akitas]=0
  [vizslas]=0 [goldfish]=5 [trees]=3 [cars]=2 [perfumes]=1
)

# find PART — echo the 1-based number of the first Sue consistent with target.
# PART 2 treats cats/trees as lower bounds and pomeranians/goldfish as upper
# bounds; otherwise compounds must match exactly.
find() {
  local part=$1
  local lines line rest pair k v want n=0 match
  read_lines lines
  for line in "${lines[@]}"; do
    ((n++))
    # Strip the "Sue N: " prefix, then split the "k: v" pairs on ", ".
    rest=${line#*: }
    match=1
    local IFS=,
    for pair in $rest; do
      pair=${pair# }                 # drop leading space after the comma
      k=${pair%%:*}
      v=${pair##*: }
      want=${target[$k]}
      if ((part == 2)) && [[ $k == cats || $k == trees ]]; then
        ((v > want)) || { match=0; break; }
      elif ((part == 2)) && [[ $k == pomeranians || $k == goldfish ]]; then
        ((v < want)) || { match=0; break; }
      else
        ((v == want)) || { match=0; break; }
      fi
    done
    if ((match)); then echo "$n"; return; fi
  done
}

part_one() { find 1; }
part_two() { find 2; }
