#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 12.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.
#
# jq is guaranteed present (elf's bash runner requires it), so we use it to walk
# the JSON structurally rather than hand-parsing in bash.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

part_one() {
  local input
  read_input input
  # Sum every number anywhere in the document.
  printf '%s' "$input" | jq '([.. | numbers] | add) // 0'
}

part_two() {
  local input
  read_input input
  # Replace any object that has a "red" value with {}, then sum all numbers.
  printf '%s' "$input" |
    jq '(walk(if type == "object" and any(.[]; . == "red") then {} else . end) | [.. | numbers] | add) // 0'
}
