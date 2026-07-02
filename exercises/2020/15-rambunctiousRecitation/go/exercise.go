package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 15.
type Exercise struct {
	common.BaseExercise
}

func parse(instr string) ([]int, error) {
	fields := strings.Split(strings.TrimSpace(instr), ",")
	nums := make([]int, 0, len(fields))
	for _, f := range fields {
		n, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %w", f, err)
		}
		nums = append(nums, n)
	}
	return nums, nil
}

// play runs the memory game to turn `target` and returns the number spoken then.
// lastSeen[v] holds the (1-based) turn on which v was most recently spoken; a
// flat slice sized to target is used because values never exceed the turn count.
func play(start []int, target int) int {
	lastSeen := make([]int32, target)
	for i := 0; i < len(start)-1; i++ {
		lastSeen[start[i]] = int32(i + 1)
	}

	last := start[len(start)-1]
	for turn := len(start); turn < target; turn++ {
		var next int
		if prev := lastSeen[last]; prev != 0 {
			next = turn - int(prev)
		}
		lastSeen[last] = int32(turn)
		last = next
	}
	return last
}

// One returns the 2020th number spoken.
func (e Exercise) One(instr string) (any, error) {
	start, err := parse(instr)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%d", play(start, 2020)), nil
}

// Two returns the 30,000,000th number spoken.
func (e Exercise) Two(instr string) (any, error) {
	start, err := parse(instr)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%d", play(start, 30_000_000)), nil
}
