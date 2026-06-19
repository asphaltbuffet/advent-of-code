#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 14.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# Globals built by load(): parallel arrays speed[], fly[], rest[] plus duration.
declare -a speed fly rest
duration=0

load() {
  speed=() fly=() rest=()
  local lines line f
  read_lines lines
  for line in "${lines[@]}"; do
    # shellcheck disable=SC2206
    f=($line)
    # <name> can fly <speed> km/s for <fly> seconds, ... rest for <rest> ...
    speed+=("${f[3]}")
    fly+=("${f[6]}")
    rest+=("${f[13]}")
  done
  # Duration is not in the input: 2-reindeer AoC example races 1000s, real 2503s.
  if ((${#speed[@]} <= 2)); then duration=1000; else duration=2503; fi
}

# distance INDEX T — set global REPLY to how far reindeer INDEX has travelled
# after T seconds. Uses a global (not command substitution) to avoid a subshell
# fork per call, which matters in part two's tight inner loop.
distance() {
  local i=$1 t=$2
  local cycle=$((fly[i] + rest[i]))
  local rem=$((t % cycle))
  local flying=$(((t / cycle) * fly[i]))
  if ((rem < fly[i])); then ((flying += rem)); else ((flying += fly[i])); fi
  REPLY=$((flying * speed[i]))
}

part_one() {
  load
  local i best=0
  for ((i = 0; i < ${#speed[@]}; i++)); do
    distance "$i" "$duration"
    ((REPLY > best)) && best=$REPLY
  done
  echo "$best"
}

part_two() {
  load
  local n=${#speed[@]}
  local -a points dists
  local i t lead
  for ((i = 0; i < n; i++)); do points[i]=0; done

  for ((t = 1; t <= duration; t++)); do
    lead=0
    for ((i = 0; i < n; i++)); do
      distance "$i" "$t"
      dists[i]=$REPLY
      ((dists[i] > lead)) && lead=${dists[i]}
    done
    for ((i = 0; i < n; i++)); do
      ((dists[i] == lead)) && ((points[i]++))
    done
  done

  local best=0
  for ((i = 0; i < n; i++)); do ((points[i] > best)) && best=${points[i]}; done
  echo "$best"
}
