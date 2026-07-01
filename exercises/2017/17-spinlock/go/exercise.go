package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 17.
type Exercise struct {
	common.BaseExercise
}

func parseStep(instr string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(instr))
	return n
}

// One builds the buffer through 2017 insertions and returns the value that ends
// up immediately after the last inserted value, 2017.
func (e Exercise) One(instr string) (any, error) {
	step := parseStep(instr)

	buf := []int{0}
	pos := 0
	for i := 1; i <= 2017; i++ {
		pos = (pos+step)%len(buf) + 1
		// Insert i at index pos.
		buf = append(buf, 0)
		copy(buf[pos+1:], buf[pos:])
		buf[pos] = i
	}
	return buf[(pos+1)%len(buf)], nil
}

// Two runs 50 million insertions without a real buffer: value 0 stays at index
// 0, so it tracks only whatever value most recently landed at index 1.
func (e Exercise) Two(instr string) (any, error) {
	step := parseStep(instr)

	pos, afterZero := 0, 0
	for i := 1; i <= 50_000_000; i++ {
		pos = (pos+step)%i + 1
		if pos == 1 {
			afterZero = i
		}
	}
	return afterZero, nil
}
