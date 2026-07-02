package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 2.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer: 5727
func (e Exercise) One(instr string) (any, error) {
	twos, threes := 0, 0

	for _, id := range strings.Fields(instr) {
		var counts [26]int
		for _, ch := range id {
			if ch >= 'a' && ch <= 'z' {
				counts[ch-'a']++
			}
		}

		has2, has3 := false, false
		for _, n := range counts {
			switch n {
			case 2:
				has2 = true
			case 3:
				has3 = true
			}
		}

		if has2 {
			twos++
		}
		if has3 {
			threes++
		}
	}

	return twos * threes, nil
}

// Two returns the answer to the second part of the exercise.
// answer: uwfmdjxyxlbgnrotcfpvswaqh
func (e Exercise) Two(instr string) (any, error) {
	ids := strings.Fields(instr)

	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			a, b := ids[i], ids[j]
			if len(a) != len(b) {
				continue
			}

			diff, at := 0, -1
			for k := 0; k < len(a) && diff <= 1; k++ {
				if a[k] != b[k] {
					diff++
					at = k
				}
			}

			if diff == 1 {
				return a[:at] + a[at+1:], nil
			}
		}
	}

	return nil, fmt.Errorf("no pair differing by a single character")
}
