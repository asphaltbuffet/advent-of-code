#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 11.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# The password is held as the global array `pw` of 0..25 letter offsets.
declare -a pw

# str_to_pw STRING — fill pw[] with letter offsets (a=0 .. z=25).
str_to_pw() {
  local s="$1" i c code
  pw=()
  for ((i = 0; i < ${#s}; i++)); do
    c="${s:i:1}"
    printf -v code '%d' "'$c"
    pw+=("$((code - 97))")
  done
}

# pw_to_str — echo the current pw[] as a letter string.
pw_to_str() {
  local out="" off
  for off in "${pw[@]}"; do
    printf -v out '%s%b' "$out" "\\$(printf '%03o' $((off + 97)))"
  done
  printf '%s' "$out"
}

# increment — advance pw[] by one as a base-26 odometer with carry.
increment() {
  local i
  for ((i = ${#pw[@]} - 1; i >= 0; i--)); do
    if ((pw[i] == 25)); then
      pw[i]=0
    else
      ((pw[i]++))
      break
    fi
  done
}

# valid — return 0 if pw[] satisfies all three rules.
valid() {
  local n=${#pw[@]} i

  # Rule 2: no i (8), o (14), l (11).
  for ((i = 0; i < n; i++)); do
    case ${pw[i]} in 8 | 14 | 11) return 1 ;; esac
  done

  # Rule 1: an increasing straight of at least three letters.
  local straight=1
  for ((i = 0; i + 2 < n; i++)); do
    if ((pw[i + 1] == pw[i] + 1 && pw[i + 2] == pw[i] + 2)); then
      straight=0
      break
    fi
  done
  ((straight == 0)) || return 1

  # Rule 3: at least two different, non-overlapping pairs.
  local -A seen=()
  i=0
  while ((i + 1 < n)); do
    if ((pw[i] == pw[i + 1])); then
      seen[${pw[i]}]=1
      ((i += 2))
    else
      ((i++))
    fi
  done
  ((${#seen[@]} >= 2))
}

# skip_forbidden — if pw[] contains i (8), o (14), or l (11), bump the first
# such position to the next letter and zero everything after it. This leaps over
# the entire range of passwords that would keep that forbidden letter, which is
# the dominant cost of the search. Returns 0 if a skip was performed.
skip_forbidden() {
  local n=${#pw[@]} i j
  for ((i = 0; i < n; i++)); do
    case ${pw[i]} in
      8 | 14 | 11)
        ((pw[i]++))
        for ((j = i + 1; j < n; j++)); do pw[j]=0; done
        return 0
        ;;
    esac
  done
  return 1
}

# next_password — advance pw[] to the next valid password.
next_password() {
  increment
  skip_forbidden
  while ! valid; do
    increment
    skip_forbidden
  done
}

part_one() {
  local input
  read_input input
  str_to_pw "$input"
  next_password
  pw_to_str
}

part_two() {
  local input
  read_input input
  str_to_pw "$input"
  next_password
  next_password
  pw_to_str
}
