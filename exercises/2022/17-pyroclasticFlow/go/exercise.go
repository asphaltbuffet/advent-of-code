package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

const (
	rock    string = "#"
	space   string = "."
	falling string = "@"

	partOneRocks int = 2022
	partTwoRocks int = 1000000000000
)

// Exercise for Advent of Code 2022 day 17.
type Exercise struct {
	common.BaseExercise
}

type point struct {
	x, y int
}

type shape int

const (
	horizLine shape = iota
	plus
	ninety
	vertLine
	square
)

// One returns the answer to the first part of the exercise.
// answer:
func (c Exercise) One(instr string) (any, error) {
	s := parse(instr)

	for i := 0; i < partOneRocks; i++ {
		s.moveRock()
	}

	// height after 2022 rocks fall
	return s.settledHeight + 1, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (c Exercise) Two(instr string) (any, error) {
	s := parse(instr)
	pastStates := map[string][3]int{}

	dupeRows := 0

	dropCount := 0

	for dropCount < partTwoRocks {
		s.moveRock()
		dropCount++

		h := s.hash(20)

		if past, ok := pastStates[h]; ok {
			pastSteps, pastRocksDropped, pastHighRow := past[0], past[1], past[2]

			stepsToSkip := s.jetNum - pastSteps
			rocksToSkip := dropCount - pastRocksDropped
			rowsToAdd := s.settledHeight - pastHighRow

			iterationsToSkip := (partTwoRocks - dropCount) / rocksToSkip

			dupeRows += rowsToAdd * iterationsToSkip
			s.jetNum += stepsToSkip * iterationsToSkip
			dropCount += rocksToSkip * iterationsToSkip
		} else {
			pastStates[h] = [3]int{
				s.jetNum, dropCount, s.settledHeight,
			}
		}
	}

	// height after 2022 rocks fall
	return s.settledHeight + 1 + dupeRows, nil
}

func parse(input string) state {
	s := state{
		chamber:       [][]string{},
		settledHeight: -1,
		curShape:      nil,
		nextShape:     0,
		jets:          strings.Split(input, ""),
		jetNum:        0,
	}

	return s
}
