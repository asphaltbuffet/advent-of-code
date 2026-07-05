package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 2.
type Exercise struct {
	common.BaseExercise
}

// keypad layouts: ' ' marks an invalid (off-pad) cell.
var pad1 = []string{
	"123",
	"456",
	"789",
}

var pad2 = []string{
	"  1  ",
	" 234 ",
	"56789",
	" ABC ",
	"  D  ",
}

// solve walks the move lines over the given keypad, starting on '5', and
// returns the code formed by the key pressed at the end of each line.
func solve(instr string, pad []string) string {
	// Locate the starting '5'.
	var r, c int
	for i := range pad {
		if j := strings.IndexByte(pad[i], '5'); j >= 0 {
			r, c = i, j
		}
	}

	valid := func(nr, nc int) bool {
		return nr >= 0 && nr < len(pad) && nc >= 0 && nc < len(pad[nr]) && pad[nr][nc] != ' '
	}

	var code strings.Builder
	for line := range strings.FieldsSeq(instr) {
		for i := range len(line) {
			nr, nc := r, c
			switch line[i] {
			case 'U':
				nr--
			case 'D':
				nr++
			case 'L':
				nc--
			case 'R':
				nc++
			}
			if valid(nr, nc) {
				r, c = nr, nc
			}
		}
		code.WriteByte(pad[r][c])
	}
	return code.String()
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return solve(instr, pad1), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return solve(instr, pad2), nil
}
