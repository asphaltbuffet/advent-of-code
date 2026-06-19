#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 8.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# Note: input lines contain literal backslashes (\x, \\, \"). Read them with
# `read -r` and never pass them through echo/printf %b, or bash will interpret
# the escapes and corrupt the counts.

part_one() {
  local line total=0 lines
  read_lines lines
  for line in "${lines[@]}"; do
    local code=${#line}
    # Walk chars between the surrounding quotes, counting decoded length.
    local mem=0 i=1 last=$((code - 1)) c
    while ((i < last)); do
      c=${line:i:1}
      if [[ $c == '\' ]]; then
        if [[ ${line:i+1:1} == 'x' ]]; then
          ((i += 4)) # \xNN -> 1 char
        else
          ((i += 2)) # \\ or \" -> 1 char
        fi
      else
        ((i++))
      fi
      ((mem++))
    done
    ((total += code - mem))
  done
  echo "$total"
}

part_two() {
  local line total=0 lines
  read_lines lines
  for line in "${lines[@]}"; do
    local code=${#line}
    # Count chars needing escaping (" and \); each adds 1, plus 2 new quotes.
    local stripped_q=${line//\"/}
    local stripped_b=${line//\\/}
    local quotes=$((code - ${#stripped_q}))
    local backslashes=$((code - ${#stripped_b}))
    ((total += 2 + quotes + backslashes))
  done
  echo "$total"
}
