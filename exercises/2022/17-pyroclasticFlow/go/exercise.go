package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

const (
	rock    = "#"
	space   = "."
	falling = "@"

	partOneRocks = 2022
	partTwoRocks = 1000000000000
)

// Exercise for Advent of Code 2022 day 17.
type Exercise struct {
	common.BaseExercise
}

type point struct {
	x, y int
}

type state struct {
	settled           [][]string
	highestSettledRow int
	fallingCoords     [][2]int
	nextRockIndex     int
	steps             []string
	stepIndex         int
}

// One returns the answer to the first part of the exercise.
// answer:
func (c Exercise) One(instr string) (any, error) {
	s := newState(instr)

	for i := 0; i < partOneRocks; i++ {
		s.dropRock()
	}

	// height after 2022 rocks fall
	return s.highestSettledRow + 1, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (c Exercise) Two(instr string) (any, error) {
	s := newState(instr)
	pastStates := map[string][3]int{}

	dupeRows := 0

	rocksDropped := 0

	for rocksDropped < partTwoRocks {
		s.dropRock()
		rocksDropped++

		h := s.hash(20)

		if past, ok := pastStates[h]; ok {
			pastSteps, pastRocksDropped, pastHighRow := past[0], past[1], past[2]

			stepsToSkip := s.stepIndex - pastSteps
			rocksToSkip := rocksDropped - pastRocksDropped
			rowsToAdd := s.highestSettledRow - pastHighRow

			iterationsToSkip := (partTwoRocks - rocksDropped) / rocksToSkip

			dupeRows += rowsToAdd * iterationsToSkip
			s.stepIndex += stepsToSkip * iterationsToSkip
			rocksDropped += rocksToSkip * iterationsToSkip
		} else {
			pastStates[h] = [3]int{
				s.stepIndex, rocksDropped, s.highestSettledRow,
			}
		}
	}

	// height after 2022 rocks fall
	return s.highestSettledRow + 1 + dupeRows, nil
}

func newState(input string) state {
	s := state{
		settled:           [][]string{},
		highestSettledRow: -1,
		fallingCoords:     nil,
		nextRockIndex:     0,
		steps:             strings.Split(input, ""),
		stepIndex:         0,
	}

	return s
}

// knew I'd need this for debugging...
func (s state) printState() {
	copySettled := [][]string{}

	for _, row := range s.settled {
		copyRow := make([]string, len(row))
		copy(copyRow, row)
		copySettled = append(copySettled, copyRow)
	}

	for _, coord := range s.fallingCoords {
		copySettled[coord[0]][coord[1]] = "@"
	}

	var sb strings.Builder

	for r := len(copySettled) - 1; r >= 0; r-- {
		sb.WriteString(strings.Join(copySettled[r], ""))
		sb.WriteString(fmt.Sprintf("%d\n", r))
	}

	fmt.Println(sb.String())
}

func (s *state) dropRock() {
	s.populateNextBaseCoords()

	highestRow := 0
	for _, c := range s.fallingCoords {
		if c[0] > highestRow {
			highestRow = c[0]
		}
	}

	for len(s.settled) <= highestRow {
		s.settled = append(s.settled, newEmptyRow())
	}

	// will be set back to nil when settled
	for s.fallingCoords != nil {
		switch s.steps[s.stepIndex%len(s.steps)] {
		case ">":
			// check if can move right
			canMoveRight := true

			for _, c := range s.fallingCoords {
				if c[1] == 6 || s.settled[c[0]][c[1]+1] != "." {
					canMoveRight = false
				}
			}

			if canMoveRight {
				for i := range s.fallingCoords {
					s.fallingCoords[i][1]++
				}
			}
		case "<":
			// check if can move left
			canMoveLeft := true

			for _, c := range s.fallingCoords {
				if c[1] == 0 || s.settled[c[0]][c[1]-1] != "." {
					canMoveLeft = false
				}
			}

			if canMoveLeft {
				for i := range s.fallingCoords {
					s.fallingCoords[i][1]--
				}
			}
		default:
			panic(s.steps[s.stepIndex])
		}

		s.stepIndex++

		// move down
		canMoveDown := true

		for _, c := range s.fallingCoords {
			if c[0] == 0 || s.settled[c[0]-1][c[1]] != "." {
				canMoveDown = false
			}
		}
		// is blocked, draw onto settled then make nil
		if !canMoveDown {
			for _, c := range s.fallingCoords {
				s.settled[c[0]][c[1]] = "#"
			}

			s.fallingCoords = nil

			for r := len(s.settled) - 1; r >= 0; r-- {
				if strings.Join(s.settled[r], "") != "......." {
					s.highestSettledRow = r
					break
				}
			}
		} else {
			for i := range s.fallingCoords {
				s.fallingCoords[i][0]--
			}
		}
	}
}

func newEmptyRow() []string {
	row := make([]string, 7)

	for i := range row {
		row[i] = "."
	}

	return row
}

var baseCoords = [][][2]int{
	{
		// line ####
		{0, 0},
		{0, 1},
		{0, 2},
		{0, 3},
	}, {
		// plus
		{0, 1},
		{1, 0},
		{1, 1},
		{1, 2},
		{2, 1},
	}, {
		// flipped L
		{0, 0},
		{0, 1},
		{0, 2},
		{1, 2},
		{2, 2},
	}, {
		// vert line
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
	}, {
		// square
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	},
}

func init() {
	// add 2 cols to all baseCoords because they fall 2 off of left wall
	for i := range baseCoords {
		for j := range baseCoords[i] {
			baseCoords[i][j][1] += 2
		}
	}
}

func (s *state) populateNextBaseCoords() {
	copyCoords := make([][2]int, len(baseCoords[s.nextRockIndex]))
	copy(copyCoords, baseCoords[s.nextRockIndex])
	s.nextRockIndex++
	s.nextRockIndex %= 5

	// lowest row of baseCoords...

	for i := range copyCoords {
		copyCoords[i][0] += s.highestSettledRow + 1 + 3
	}

	s.fallingCoords = copyCoords
}

// for part 2 to find return states
// NOTE: had to play with the number of rows to be hashed... 20 seems to
// work on the example input
func (s *state) hash(topRowsToHash int) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprint(s.nextRockIndex))

	for r := s.highestSettledRow; r >= 0 && r > s.highestSettledRow-topRowsToHash; r-- {
		sb.WriteString("\n" + strings.Join(s.settled[r], ""))
	}

	return sb.String()
}
