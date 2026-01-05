package exercises

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 10.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	mm := ParseMachines(instr)

	sum := 0
	for _, m := range mm {
		sum += m.GetButtonPresses()
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	sum := 0

	for _, line := range strings.Split(instr, "\n") {
		// parse each line
		tokens := strings.Fields(line)
		n := len(tokens)

		jj := strings.Split(strings.Trim(tokens[n-1], "{}"), ",")

		req := make([]int, len(jj))

		for i, j := range jj {
			req[i], _ = strconv.Atoi(j)
		}

		buttonsRaw := tokens[1 : n-1]
		buttons := make([][]int, len(buttonsRaw))

		rlen := len(req)

		for i, b := range buttonsRaw {
			outTokens := strings.Split(strings.Trim(b, "()"), ",")
			outputs := make([]int, rlen)

			for _, o := range outTokens {
				outIdx, _ := strconv.Atoi(o)
				if outIdx < rlen {
					outputs[outIdx] = 1
				}
			}

			buttons[i] = outputs
		}

		memo := map[string]int{fmt.Sprint(slices.Repeat([]int{0}, rlen)): 0}
		sum += minPresses(buttons, req, memo)
	}

	return sum, nil
}
