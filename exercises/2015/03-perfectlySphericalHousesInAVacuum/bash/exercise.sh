#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 3.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

part_one() {
  local input
  read_input input

  local x=0 y=0
  declare -A visited
  visited["0,0"]=1

  local i c
  for ((i = 0; i < ${#input}; i++)); do
    c="${input:i:1}"
    case "$c" in
      '^') ((y--)) ;;
      'v') ((y++)) ;;
      '>') ((x++)) ;;
      '<') ((x--)) ;;
      *) continue ;;
    esac
    visited["$x,$y"]=1
  done

  echo "${#visited[@]}"
}

part_two() {
  local input
  read_input input

  local sx=0 sy=0 rx=0 ry=0 turn=0
  declare -A visited
  visited["0,0"]=1

  local i c
  for ((i = 0; i < ${#input}; i++)); do
    c="${input:i:1}"
    case "$c" in
      '^' | 'v' | '>' | '<') ;;
      *) continue ;;
    esac

    if ((turn % 2 == 0)); then
      case "$c" in
        '^') ((sy--)) ;; 'v') ((sy++)) ;; '>') ((sx++)) ;; '<') ((sx--)) ;;
      esac
      visited["$sx,$sy"]=1
    else
      case "$c" in
        '^') ((ry--)) ;; 'v') ((ry++)) ;; '>') ((rx++)) ;; '<') ((rx--)) ;;
      esac
      visited["$rx,$ry"]=1
    fi
    ((turn++))
  done

  echo "${#visited[@]}"
}
