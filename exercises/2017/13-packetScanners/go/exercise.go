package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 13.
type Exercise struct {
	common.BaseExercise
}

type layer struct{ depth, rng int }

// parseLayers reads "depth: range" lines into layers.
func parseLayers(instr string) []layer {
	var layers []layer
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		d, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		r, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		layers = append(layers, layer{d, r})
	}
	return layers
}

// caught reports whether a scanner of the given range sits at the top when a
// packet enters at the given time. The scanner's period is 2*(range-1).
func caught(t, rng int) bool {
	return t%(2*(rng-1)) == 0
}

// One returns the total severity of crossing with no delay.
func (e Exercise) One(instr string) (any, error) {
	severity := 0
	for _, l := range parseLayers(instr) {
		if caught(l.depth, l.rng) {
			severity += l.depth * l.rng
		}
	}
	return severity, nil
}

// Two returns the fewest picoseconds of delay that avoid every scanner.
func (e Exercise) Two(instr string) (any, error) {
	layers := parseLayers(instr)
	for delay := 0; ; delay++ {
		safe := true
		for _, l := range layers {
			if caught(delay+l.depth, l.rng) {
				safe = false
				break
			}
		}
		if safe {
			return delay, nil
		}
	}
}
