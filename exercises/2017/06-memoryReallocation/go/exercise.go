package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 6.
type Exercise struct {
	common.BaseExercise
}

func parseBanks(instr string) []int {
	fields := strings.Fields(instr)
	banks := make([]int, len(fields))
	for i, f := range fields {
		banks[i], _ = strconv.Atoi(f)
	}
	return banks
}

// reallocate runs the redistribution routine until a configuration repeats. It
// returns the number of cycles to reach that repeat and the loop size (cycles
// since the repeated configuration was first seen).
func reallocate(banks []int) (int, int) {
	var cycles int
	state := make([]int, len(banks))
	copy(state, banks)

	seen := map[string]int{}
	for {
		key := key(state)
		if first, ok := seen[key]; ok {
			return cycles, cycles - first
		}
		seen[key] = cycles

		// Pick the fullest bank (lowest index wins ties), then spread it.
		m, idx := state[0], 0
		for i, v := range state {
			if v > m {
				m, idx = v, i
			}
		}
		state[idx] = 0
		for i := 1; i <= m; i++ {
			state[(idx+i)%len(state)]++
		}
		cycles++
	}
}

// key renders a bank configuration as a stable map key.
func key(state []int) string {
	var b strings.Builder
	for i, v := range state {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.Itoa(v))
	}
	return b.String()
}

// One returns the cycles completed before a configuration repeats.
func (e Exercise) One(instr string) (any, error) {
	cycles, _ := reallocate(parseBanks(instr))
	return cycles, nil
}

// Two returns the size of the loop the configurations settle into.
func (e Exercise) Two(instr string) (any, error) {
	_, loop := reallocate(parseBanks(instr))
	return loop, nil
}
