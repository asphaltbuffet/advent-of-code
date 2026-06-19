#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 13.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# Globals built by load(): h["$a,$b"]=happiness change, people[] = distinct names.
declare -A h
declare -a people
declare -A _seen
best=0
have_best=0

load() {
  h=()
  people=()
  _seen=()
  local lines line f a b n
  read_lines lines
  for line in "${lines[@]}"; do
    line=${line%.}
    # shellcheck disable=SC2206
    f=($line)
    # <A> would <gain|lose> <N> happiness units by sitting next to <B>
    a=${f[0]} b=${f[10]} n=${f[3]}
    [[ ${f[2]} == lose ]] && n=$((-n))
    h["$a,$b"]=$n
    if [[ -z ${_seen[$a]+x} ]]; then people+=("$a"); _seen[$a]=1; fi
  done
}

# Recursively permute people[k..] in place; the first element is fixed by
# starting recursion at index 1 (factors out circular rotations). At each full
# permutation, sum adjacent happiness (both directions) and track the max.
permute() {
  local k=$1 n=${#people[@]}
  if ((k == n)); then
    local total=0 i a b
    for ((i = 0; i < n; i++)); do
      a=${people[i]}
      b=${people[(i + 1) % n]}
      total=$((total + ${h["$a,$b"]:-0} + ${h["$b,$a"]:-0}))
    done
    if ((have_best == 0 || total > best)); then
      best=$total
      have_best=1
    fi
    return
  fi
  local i tmp
  for ((i = k; i < n; i++)); do
    tmp=${people[k]}; people[k]=${people[i]}; people[i]=$tmp
    permute $((k + 1))
    tmp=${people[k]}; people[k]=${people[i]}; people[i]=$tmp
  done
}

part_one() {
  load
  best=0
  have_best=0
  permute 1 # fix people[0] to factor out rotations
  echo "$best"
}

part_two() {
  load
  # Add "me" with zero happiness in both directions to everyone.
  local p
  for p in "${people[@]}"; do
    h["me,$p"]=0
    h["$p,me"]=0
  done
  people+=("me")

  best=0
  have_best=0
  permute 1
  echo "$best"
}
