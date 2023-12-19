package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 13.
type Exercise struct {
	common.BaseExercise
}

const RowMultiplier int = 100

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	var total int

	p := getPatterns(instr)

	for i, pattern := range p {
		if hr, hRes := findMirror(pattern.Row, false); hRes == 0 {
			total += hr * RowMultiplier
			continue
		}

		if vr, vRes := findMirror(pattern.Col, false); vRes == 0 {
			total += vr
			continue
		}

		return nil, fmt.Errorf("no mirrors found: pattern=%d", i)
	}

	return total, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	var total int

	p := getPatterns(instr)

	for _, pattern := range p {
		hr, hRes := findMirror(pattern.Row, true)
		vr, vRes := findMirror(pattern.Col, true)

		switch {
		case hRes == 0 && vRes == 0:
			return nil, fmt.Errorf("multiple mirrors found: v=%d, h=%d", vr, hr)
		case hRes >= 0 && vRes <= 0:
			total += hr * RowMultiplier
		case vRes >= 0 && hRes <= 0:
			total += vr
		default:
			return nil, fmt.Errorf("unknown state: v[%d]=%d, h[%d]=%d", vr, vRes, hr, hRes)
		}
	}

	return total, nil
}
