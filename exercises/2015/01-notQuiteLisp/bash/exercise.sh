#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 1.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

part_one() {
  local input ups downs
  read_input input

  # Strip everything that isn't the target paren, then the leftover length is
  # the count. ${var//pattern/} deletes all matches; "()" here means a literal
  # close paren, so ups keeps only "(" and downs keeps only ")".
  ups="${input//[^(]/}"
  downs="${input//[^)]/}"

  echo $(( ${#ups} - ${#downs} ))
}

part_two() {
  local input floor i char
  read_input input

  floor=0
  for (( i = 0; i < ${#input}; i++ )); do
    char="${input:i:1}"
    case "$char" in
      '(') (( floor++ )) ;;
      ')') (( floor-- )) ;;
    esac

    # Position is 1-based: the character that pushed Santa to -1.
    if (( floor < 0 )); then
      echo $(( i + 1 ))
      return 0
    fi
  done

  echo "santa never enters the basement" >&2
  return 1
}
