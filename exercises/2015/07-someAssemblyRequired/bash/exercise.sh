#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 7.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# Globals populated by load(): circuit[wire]=signal-expression, solved[wire]=value.
declare -A circuit
declare -A solved

# Read "<signal> -> <wire>" lines from stdin into the circuit map.
load() {
  circuit=()
  solved=()
  local line signal wire lines
  read_lines lines
  for line in "${lines[@]}"; do
    signal="${line%% -> *}"
    wire="${line##* -> }"
    circuit["$wire"]="$signal"
  done
}

# Resolve an operand into the global REPLY: a literal number passes through,
# a wire name is evaluated (recursively, populating the memo).
value() {
  local tok="$1"
  if [[ $tok == [0-9]* ]]; then
    REPLY=$((tok & 0xFFFF))
  else
    eval_wire "$tok"
  fi
}

# Memoized evaluation of the signal feeding $1; result is left in global REPLY.
# Communicating via a global (not command substitution) keeps the memo alive
# across the recursion — no subshell forks, so each wire is evaluated once.
eval_wire() {
  local wire="$1"
  if [[ -n ${solved[$wire]+x} ]]; then
    REPLY="${solved[$wire]}"
    return
  fi

  # shellcheck disable=SC2206
  local t=(${circuit[$wire]})
  local r a b

  case ${#t[@]} in
    1)
      value "${t[0]}"; r=$REPLY
      ;;
    2) # NOT x
      value "${t[1]}"; r=$(( ~REPLY & 0xFFFF ))
      ;;
    3)
      value "${t[0]}"; a=$REPLY
      value "${t[2]}"; b=$REPLY
      case "${t[1]}" in
        AND)    r=$(( (a & b) & 0xFFFF )) ;;
        OR)     r=$(( (a | b) & 0xFFFF )) ;;
        LSHIFT) r=$(( (a << b) & 0xFFFF )) ;;
        RSHIFT) r=$(( (a >> b) & 0xFFFF )) ;;
      esac
      ;;
  esac

  solved["$wire"]="$r"
  REPLY="$r"
}

part_one() {
  load
  [[ -n ${circuit[a]+x} ]] || { echo ""; return; }
  eval_wire a
  echo "$REPLY"
}

part_two() {
  load
  [[ -n ${circuit[a]+x} ]] || { echo ""; return; }
  eval_wire a
  local a=$REPLY
  # Re-solve with wire b overridden by part one's answer.
  solved=()
  solved["b"]="$a"
  eval_wire a
  echo "$REPLY"
}
