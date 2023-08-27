package exercises

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 8.
type Exercise struct {
	common.BaseExercise
}

var re = regexp.MustCompile(`\\x[0-9a-f]{2}`)

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	sum := 0

	for _, line := range strings.Split(instr, "\n") {
		ol := len(line)
		s := line

		s = strings.ReplaceAll(s, `\"`, `"`)
		s = strings.ReplaceAll(s, `\\`, `\`)
		s = re.ReplaceAllString(s, "x")

		sum += ol - (len(s) - 2)
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	sum := 0

	for _, line := range strings.Split(instr, "\n") {
		ol := len(line)

		line = fmt.Sprintf("%q", line)

		sum += len(line) - ol
	}

	return sum, nil
}
