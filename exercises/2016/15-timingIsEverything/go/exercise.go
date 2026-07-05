package exercises

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 15.
type Exercise struct {
	common.BaseExercise
}

type disc struct {
	positions, start int
}

var numRe = regexp.MustCompile(`\d+`)

// parseDiscs reads each line's (disc#, positions, time, position); only the
// position count and starting position are needed.
func parseDiscs(instr string) []disc {
	var discs []disc
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		n := numRe.FindAllString(line, -1)
		if len(n) < 4 {
			continue
		}
		positions, _ := strconv.Atoi(n[1])
		start, _ := strconv.Atoi(n[3])
		discs = append(discs, disc{positions, start})
	}
	return discs
}

// firstTime sieves for the earliest button-press time t where, for every disc i
// (1-indexed), (start + t + i) mod positions == 0. Each satisfied disc lets us
// step by its (coprime) position count.
func firstTime(discs []disc) int {
	t, step := 0, 1
	for i, d := range discs {
		for (d.start+t+i+1)%d.positions != 0 {
			t += step
		}
		step *= d.positions
	}
	return t
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return firstTime(parseDiscs(instr)), nil
}

// Two returns the answer to the second part of the exercise: an extra disc with
// 11 positions starting at position 0 is added at the bottom.
func (e Exercise) Two(instr string) (any, error) {
	discs := append(parseDiscs(instr), disc{positions: 11, start: 0})
	return firstTime(discs), nil
}
