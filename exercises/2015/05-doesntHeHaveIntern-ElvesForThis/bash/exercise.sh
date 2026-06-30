#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 5.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

is_nice_one() {
  local s="$1" i n=${#1}

  # At least three vowels.
  local rest="${s//[^aeiou]/}"
  ((${#rest} >= 3)) || return 1

  # A letter appearing twice in a row.
  local doubled=1
  for ((i = 0; i < n - 1; i++)); do
    [[ "${s:i:1}" == "${s:i+1:1}" ]] && { doubled=0; break; }
  done
  ((doubled == 0)) || return 1

  # None of the forbidden pairs.
  case "$s" in
    *ab* | *cd* | *pq* | *xy*) return 1 ;;
  esac

  return 0
}

is_nice_two() {
  local s="$1" i n=${#1}

  # A pair of two letters that appears at least twice without overlapping.
  local pair=1
  for ((i = 0; i < n - 1; i++)); do
    local p="${s:i:2}"
    [[ "${s:i+2}" == *"$p"* ]] && { pair=0; break; }
  done
  ((pair == 0)) || return 1

  # A letter that repeats with exactly one letter between.
  local repeat=1
  for ((i = 0; i < n - 2; i++)); do
    [[ "${s:i:1}" == "${s:i+2:1}" ]] && { repeat=0; break; }
  done
  ((repeat == 0)) || return 1

  return 0
}

part_one() {
  local count=0 line lines
  read_lines lines
  for line in "${lines[@]}"; do
    is_nice_one "$line" && ((count++))
  done
  echo "$count"
}

part_two() {
  local count=0 line lines
  read_lines lines
  for line in "${lines[@]}"; do
    is_nice_two "$line" && ((count++))
  done
  echo "$count"
}
