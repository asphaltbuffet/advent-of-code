package exercises

import (
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 3.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	// logger.Level = log.DebugLevel
	s := New(instr)

	// for each part, check if the bounds are a symbol
	var sum int
	for _, p := range s.parts {
		p := p
		for _, b := range p.bounds {
			if s.isSymbol[b] {
				// fmt.Printf("part[%d] %q has a symbol in its bounds at %s\n", p.id, p.value, b)

				sum += p.number
				break
			}
		}
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	// logger.Level = log.DebugLevel
	s := New(instr)

	counts := make(map[point][]int)

	// for each part, check if the bounds are a symbol
	for gp := range s.isSymbol {
		for _, p := range s.parts {
			for _, b := range p.bounds {
				if b.x != gp.x || b.y != gp.y {
					continue
				}

				// fmt.Printf("symbol %s has part %d\n", gp, p.number)
				counts[gp] = append(counts[gp], p.number)
			}
		}
	}

	var sum int
	for _, c := range counts {
		if len(c) == 2 {
			sum += c[0] * c[1]
		}
	}

	return sum, nil
}
