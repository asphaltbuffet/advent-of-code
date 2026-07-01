package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 9.
type Exercise struct {
	common.BaseExercise
}

// process scans the stream once and returns the total group score and the count
// of non-cancelled characters inside garbage.
func process(stream string) (score, garbage int) {
	depth := 0
	inGarbage := false

	for i := 0; i < len(stream); i++ {
		c := stream[i]
		switch {
		case c == '!':
			i++ // skip the next character (cancellation)
		case inGarbage:
			if c == '>' {
				inGarbage = false
			} else {
				garbage++
			}
		case c == '<':
			inGarbage = true
		case c == '{':
			depth++
			score += depth
		case c == '}':
			depth--
		}
	}
	return score, garbage
}

// One returns the total score for all groups in the stream.
func (e Exercise) One(instr string) (any, error) {
	score, _ := process(strings.TrimSpace(instr))
	return score, nil
}

// Two returns the number of non-cancelled characters within garbage.
func (e Exercise) Two(instr string) (any, error) {
	_, garbage := process(strings.TrimSpace(instr))
	return garbage, nil
}
