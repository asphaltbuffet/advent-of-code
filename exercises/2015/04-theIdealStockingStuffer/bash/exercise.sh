#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 4.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# Scan candidates [start, end] for the lowest i whose MD5(key+i) hex digest
# begins with $prefix; print that i and exit 0, or exit 1 if none in range.
# Exported so GNU parallel's worker shells can call it.
mine_chunk() {
  local key="$1" prefix="$2" start="$3" end="$4"
  local n=${#prefix} i hash
  for ((i = start; i <= end; i++)); do
    hash="$(printf '%s%d' "$key" "$i" | md5sum)"
    if [[ "${hash:0:$n}" == "$prefix" ]]; then
      echo "$i"
      return 0
    fi
  done
  return 1
}
export -f mine_chunk

# Serial fallback: forks md5sum once per candidate, one process total.
mine_serial() {
  local key="$1" prefix="$2" i=0
  while true; do
    ((i++))
    mine_chunk "$key" "$prefix" "$i" "$i" && return 0
  done
}

# Find the lowest positive integer that, appended to $key, yields an MD5 hash
# (hex) starting with $prefix. Uses GNU parallel when available: each round
# hands $jobs contiguous chunks to parallel with --halt now,success=1, so the
# first worker to hit stops the round; we take the min hit across that round to
# guarantee the lowest index. Falls back to a serial loop without GNU parallel.
mine() {
  local key="$1" prefix="$2"

  # moreutils ships a different `parallel`; require the GNU one (has --version).
  if ! parallel --version >/dev/null 2>&1; then
    mine_serial "$key" "$prefix"
    return
  fi

  local jobs chunk round start ends hits
  jobs=$(nproc 2>/dev/null || echo 4)
  chunk=25000 # candidates per worker per round

  for ((round = 0; ; round++)); do
    start=$((round * jobs * chunk + 1))
    # Build this round's per-job [start,end] ranges.
    local starts=() endlist=()
    local j s
    for ((j = 0; j < jobs; j++)); do
      s=$((start + j * chunk))
      starts+=("$s")
      endlist+=("$((s + chunk - 1))")
    done

    hits=$(parallel --halt now,success=1 -j "$jobs" \
      mine_chunk "$key" "$prefix" {1} {2} \
      ::: "${starts[@]}" :::+ "${endlist[@]}" 2>/dev/null)

    if [[ -n $hits ]]; then
      # Multiple workers in the round may have hit; the answer is the smallest.
      echo "$hits" | sort -n | head -1
      return 0
    fi
  done
}

part_one() {
  local input
  read_input input
  mine "$input" "00000"
}

part_two() {
  local input
  read_input input
  mine "$input" "000000"
}
