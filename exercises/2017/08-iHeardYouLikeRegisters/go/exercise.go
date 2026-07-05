package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 8.
type Exercise struct {
	common.BaseExercise
}

func holds(a int, op string, b int) bool {
	switch op {
	case ">":
		return a > b
	case "<":
		return a < b
	case ">=":
		return a >= b
	case "<=":
		return a <= b
	case "==":
		return a == b
	case "!=":
		return a != b
	}
	return false
}

// run executes the instructions and returns the largest final register value
// and the highest value held in any register at any point.
func run(instr string) (int, int) {
	var final, peak int
	regs := map[string]int{}

	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		f := strings.Fields(line)
		if len(f) < 7 {
			continue
		}
		// f: reg inc|dec amount if condReg op value
		amount, _ := strconv.Atoi(f[2])
		condVal, _ := strconv.Atoi(f[6])
		if !holds(regs[f[4]], f[5], condVal) {
			continue
		}
		if f[1] == "dec" {
			amount = -amount
		}
		regs[f[0]] += amount
		if regs[f[0]] > peak {
			peak = regs[f[0]]
		}
	}

	for _, v := range regs {
		if v > final {
			final = v
		}
	}
	return final, peak
}

// One returns the largest register value after all instructions run.
func (e Exercise) One(instr string) (any, error) {
	final, _ := run(instr)
	return final, nil
}

// Two returns the highest value held in any register during execution.
func (e Exercise) Two(instr string) (any, error) {
	_, peak := run(instr)
	return peak, nil
}
