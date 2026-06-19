package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 10.
type Exercise struct {
	common.BaseExercise
}

// lookAndSay applies one round of the look-and-say transform: each run of
// identical digits becomes <count><digit>.
func lookAndSay(s []byte) []byte {
	var b strings.Builder
	b.Grow(len(s) * 2)

	for i := 0; i < len(s); {
		j := i
		for j < len(s) && s[j] == s[i] {
			j++
		}
		b.WriteString(strconv.Itoa(j - i))
		b.WriteByte(s[i])
		i = j
	}

	return []byte(b.String())
}

// iterate applies lookAndSay n times to the trimmed input and returns the
// resulting length.
func iterate(instr string, n int) int {
	s := []byte(strings.TrimSpace(instr))

	for range n {
		s = lookAndSay(s)
	}

	return len(s)
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return iterate(instr, 40), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return iterate(instr, 50), nil
}
