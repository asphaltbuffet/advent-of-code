package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 25.
type Exercise struct {
	common.BaseExercise
}

// parse splits the input into lock and key height profiles. A block whose top
// row is all '#' is a lock; otherwise a key. Heights are '#' count per column
// minus the filled base row (so each height is in [0,5]).
func parse(instr string) (locks, keys [][5]int) {
	for _, block := range strings.Split(strings.TrimRight(instr, "\n"), "\n\n") {
		rows := strings.Split(strings.TrimSpace(block), "\n")
		if len(rows) < 7 {
			continue
		}
		var h [5]int
		for _, row := range rows {
			for c := 0; c < 5 && c < len(row); c++ {
				if row[c] == '#' {
					h[c]++
				}
			}
		}
		for c := range h {
			h[c]-- // discount the solid base row
		}
		if rows[0] == "#####" {
			locks = append(locks, h)
		} else {
			keys = append(keys, h)
		}
	}
	return locks, keys
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	locks, keys := parse(instr)
	count := 0
	for _, l := range locks {
		for _, k := range keys {
			fits := true
			for c := 0; c < 5; c++ {
				if l[c]+k[c] > 5 {
					fits = false
					break
				}
			}
			if fits {
				count++
			}
		}
	}
	return count, nil
}

// Two returns the answer to the second part of the exercise. Day 25 has no
// second puzzle — the final star is claimed once every other day is complete.
func (e Exercise) Two(instr string) (any, error) {
	return "Merry Christmas!", nil
}
