package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 5.
type Exercise struct {
	common.BaseExercise
}

// seatID decodes a boarding pass as a 10-bit binary number: B/R are 1 bits, F/L
// are 0. The row*8+col seat ID is exactly this value.
func seatID(pass string) int {
	id := 0
	for i := 0; i < len(pass); i++ {
		id <<= 1
		if pass[i] == 'B' || pass[i] == 'R' {
			id |= 1
		}
	}
	return id
}

func parse(instr string) []int {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	ids := make([]int, 0, len(lines))
	for _, line := range lines {
		ids = append(ids, seatID(strings.TrimSpace(line)))
	}
	return ids
}

// One returns the highest seat ID on any boarding pass.
func (e Exercise) One(instr string) (any, error) {
	max := 0
	for _, id := range parse(instr) {
		if id > max {
			max = id
		}
	}
	return fmt.Sprintf("%d", max), nil
}

// Two finds the single empty seat whose neighbors (ID-1 and ID+1) are both
// occupied — the missing ID inside the otherwise-contiguous range.
func (e Exercise) Two(instr string) (any, error) {
	present := map[int]bool{}
	lo, hi := 1<<10, 0
	for _, id := range parse(instr) {
		present[id] = true
		if id < lo {
			lo = id
		}
		if id > hi {
			hi = id
		}
	}

	for id := lo + 1; id < hi; id++ {
		if !present[id] && present[id-1] && present[id+1] {
			return fmt.Sprintf("%d", id), nil
		}
	}

	return nil, fmt.Errorf("no missing seat found")
}
