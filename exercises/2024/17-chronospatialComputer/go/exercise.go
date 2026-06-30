package exercises

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 17.
type Exercise struct {
	common.BaseExercise
}

var numRe = regexp.MustCompile(`\d+`)

// parse reads registers A, B, C and the program from the input.
func parse(instr string) (a, b, c int, program []int) {
	nums := numRe.FindAllString(instr, -1)
	vals := make([]int, len(nums))
	for i, s := range nums {
		vals[i], _ = strconv.Atoi(s)
	}
	return vals[0], vals[1], vals[2], vals[3:]
}

// run executes the program with the given registers and returns its output.
func run(a, b, c int, program []int) []int {
	var out []int
	combo := func(op int) int {
		switch op {
		case 4:
			return a
		case 5:
			return b
		case 6:
			return c
		default:
			return op
		}
	}

	for ip := 0; ip+1 < len(program); {
		opcode, operand := program[ip], program[ip+1]
		switch opcode {
		case 0: // adv
			a >>= uint(combo(operand))
		case 1: // bxl
			b ^= operand
		case 2: // bst
			b = combo(operand) & 7
		case 3: // jnz
			if a != 0 {
				ip = operand
				continue
			}
		case 4: // bxc
			b ^= c
		case 5: // out
			out = append(out, combo(operand)&7)
		case 6: // bdv
			b = a >> uint(combo(operand))
		case 7: // cdv
			c = a >> uint(combo(operand))
		}
		ip += 2
	}
	return out
}

func joinInts(nums []int) string {
	parts := make([]string, len(nums))
	for i, n := range nums {
		parts[i] = strconv.Itoa(n)
	}
	return strings.Join(parts, ",")
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	a, b, c, program := parse(instr)
	return joinInts(run(a, b, c, program)), nil
}

// findQuine searches for the lowest A whose output equals the program, building
// A in base-8 from the most-significant digit (each loop emits one digit, then
// A /= 8). pos indexes which program value the next digit must reproduce.
func findQuine(program []int, b, c, soFar, pos int) (int, bool) {
	if pos < 0 {
		return soFar, true
	}
	best := -1
	for d := 0; d < 8; d++ {
		a := soFar<<3 | d
		out := run(a, b, c, program)
		// We need the output to match program[pos:].
		if len(out) == len(program)-pos && out[0] == program[pos] {
			if got, ok := findQuine(program, b, c, a, pos-1); ok {
				if best == -1 || got < best {
					best = got
				}
			}
		}
	}
	if best == -1 {
		return 0, false
	}
	return best, true
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	_, b, c, program := parse(instr)
	ans, _ := findQuine(program, b, c, 0, len(program)-1)
	return ans, nil
}
