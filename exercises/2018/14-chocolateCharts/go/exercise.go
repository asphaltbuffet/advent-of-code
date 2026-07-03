package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 14.
type Exercise struct {
	common.BaseExercise
}

// scoreboard holds the recipe digits (0-9) and the two elves' positions. Each
// step appends the digits of the two current recipes' sum, then advances each
// elf forward by one plus its current recipe.
type scoreboard struct {
	scores []byte
	a, b   int
}

func newScoreboard() *scoreboard {
	return &scoreboard{scores: []byte{3, 7}, a: 0, b: 1}
}

// step appends the next one or two recipes and moves both elves.
func (s *scoreboard) step() {
	sum := s.scores[s.a] + s.scores[s.b]
	if sum >= 10 {
		s.scores = append(s.scores, sum/10)
	}
	s.scores = append(s.scores, sum%10)

	s.a = (s.a + 1 + int(s.scores[s.a])) % len(s.scores)
	s.b = (s.b + 1 + int(s.scores[s.b])) % len(s.scores)
}

// One returns the answer to the first part of the exercise.
// Answer: 5371393113
func (e Exercise) One(instr string) (any, error) {
	n, err := strconv.Atoi(strings.TrimSpace(instr))
	if err != nil {
		return nil, err
	}

	sb := newScoreboard()
	for len(sb.scores) < n+10 {
		sb.step()
	}

	var b strings.Builder
	for _, d := range sb.scores[n : n+10] {
		b.WriteByte('0' + d)
	}

	return b.String(), nil
}

// Two returns the answer to the second part of the exercise.
// Answer: 20286858
func (e Exercise) Two(instr string) (any, error) {
	input := strings.TrimSpace(instr)
	target := make([]byte, len(input))
	for i := range input {
		target[i] = input[i] - '0'
	}

	sb := newScoreboard()

	// A step appends one or two recipes, so check for the target after each append
	// by comparing the tail. `checked` tracks how far we have already compared so
	// no suffix is missed across a two-digit step.
	checked := 0
	for {
		sb.step()
		for ; checked+len(target) <= len(sb.scores); checked++ {
			if matchesAt(sb.scores, target, checked) {
				return checked, nil
			}
		}
	}
}

// matchesAt reports whether target equals scores starting at index i.
func matchesAt(scores, target []byte, i int) bool {
	for j := range target {
		if scores[i+j] != target[j] {
			return false
		}
	}
	return true
}
