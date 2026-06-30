#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 9.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# Globals built by load(): dist["$a,$b"]=n, and cities[] = list of city names.
declare -A dist
declare -a cities
declare -A _seen_city
best_min=0
best_max=0

load() {
  dist=()
  cities=()
  _seen_city=()
  local lines line c1 c2 d
  read_lines lines
  for line in "${lines[@]}"; do
    read -r c1 _ c2 _ d <<<"$line"
    dist["$c1,$c2"]=$d
    dist["$c2,$c1"]=$d
    if [[ -z ${_seen_city[$c1]+x} ]]; then cities+=("$c1"); _seen_city[$c1]=1; fi
    if [[ -z ${_seen_city[$c2]+x} ]]; then cities+=("$c2"); _seen_city[$c2]=1; fi
  done
}

# Recursively permute cities[k..] in place; at each full permutation, sum the
# route and update best_min/best_max. Routes using a missing edge are skipped.
permute() {
  local k=$1 n=${#cities[@]}
  if ((k == n)); then
    local total=0 i a b leg
    for ((i = 0; i < n - 1; i++)); do
      a=${cities[i]}; b=${cities[i + 1]}
      leg=${dist["$a,$b"]+x}
      [[ -z $leg ]] && return # missing edge -> invalid route
      ((total += dist["$a,$b"]))
    done
    ((best_min == 0 || total < best_min)) && best_min=$total
    ((total > best_max)) && best_max=$total
    return
  fi
  local i tmp
  for ((i = k; i < n; i++)); do
    tmp=${cities[k]}; cities[k]=${cities[i]}; cities[i]=$tmp   # swap k,i
    permute $((k + 1))
    tmp=${cities[k]}; cities[k]=${cities[i]}; cities[i]=$tmp   # swap back
  done
}

solve() {
  load
  best_min=0
  best_max=0
  permute 0
}

part_one() {
  solve
  echo "$best_min"
}

part_two() {
  solve
  echo "$best_max"
}
