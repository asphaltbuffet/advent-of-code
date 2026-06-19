#!/usr/bin/env bash
# Solution for Advent of Code 2015 day 23.
#
# Each part reads the puzzle input on stdin and echoes the answer on stdout.
# Define part_one and part_two; the elf bash runner sources this file and
# dispatches to them.

source "$(dirname "${BASH_SOURCE[0]}")/../../../../lib/bash/aoc.bash"

# The program is tiny (~48 instructions), so a pure-bash interpreter is fine.
# Instructions are decoded into parallel arrays; the loop dispatches on opcode.
# Note: jio jumps when the register equals one, not when it is odd.
run() {
  local start_a=$1
  local -a lines
  read_lines lines

  local -a op a1 a2
  local n=0 line f
  for line in "${lines[@]}"; do
    line=${line//,/}
    # shellcheck disable=SC2206
    f=($line)
    # Operands are optional (e.g. jmp has no register), so default-expand to
    # avoid set -u errors on missing fields.
    op[n]=${f[0]}
    a1[n]=${f[1]:-}
    a2[n]=${f[2]:-}
    ((n++))
  done

  local -A reg=([a]="$start_a" [b]=0)
  local pc=0 r
  while ((pc >= 0 && pc < n)); do
    r=${a1[pc]}
    case ${op[pc]} in
      hlf) reg[$r]=$((reg[$r] / 2)); ((pc++)) ;;
      tpl) reg[$r]=$((reg[$r] * 3)); ((pc++)) ;;
      inc) reg[$r]=$((reg[$r] + 1)); ((pc++)) ;;
      jmp) ((pc += a1[pc])) ;;
      jie) if ((reg[$r] % 2 == 0)); then ((pc += a2[pc])); else ((pc++)); fi ;;
      jio) if ((reg[$r] == 1)); then ((pc += a2[pc])); else ((pc++)); fi ;;
    esac
  done

  echo "${reg[b]}"
}

part_one() { run 0; }
part_two() { run 1; }
