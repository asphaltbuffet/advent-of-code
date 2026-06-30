#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 2.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

part_one() {
  local total=0 line lines l w h
  read_lines lines
  for line in "${lines[@]}"; do
    IFS=x read -r l w h <<<"$line"
    local a=$((l * w)) b=$((w * h)) c=$((h * l))
    local min=$a
    ((b < min)) && min=$b
    ((c < min)) && min=$c
    total=$((total + 2 * (a + b + c) + min))
  done
  echo "$total"
}

part_two() {
  local total=0 line lines l w h
  read_lines lines
  for line in "${lines[@]}"; do
    IFS=x read -r l w h <<<"$line"
    local max=$l
    ((w > max)) && max=$w
    ((h > max)) && max=$h
    local perimeter=$((2 * (l + w + h - max)))
    local volume=$((l * w * h))
    total=$((total + perimeter + volume))
  done
  echo "$total"
}
