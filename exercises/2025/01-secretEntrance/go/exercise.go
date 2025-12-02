package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 1.
type Exercise struct {
	common.BaseExercise
}

const (
	TotalPositions   int = 100
	StartingPosition int = 50

	Left  rune = 'L'
	Right rune = 'R'
)

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	var zero int
	pos := StartingPosition

	for _, l := range strings.Split(instr, "\n") {
		m, err := ToMove(l)
		if err != nil {
			return nil, err
		}

		pos, _ = NewPosition(pos, m)

		if pos == 0 {
			zero++
		}
	}

	return fmt.Sprintf("%d", zero), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	var zero int
	pos := StartingPosition

	for _, l := range strings.Split(instr, "\n") {
		var z int

		m, err := ToMove(l)
		if err != nil {
			return nil, err
		}

		pos, z = NewPosition(pos, m)

		zero += z
	}

	return fmt.Sprintf("%d", zero), nil
}

func ToMove(s string) (int, error) {
	n, err := strconv.Atoi(s[1:])
	if err != nil {
		return 0, err
	}

	switch rune(s[0]) {
	case Left:
		return -n, nil
	case Right:
		return n, nil
	default:
		return 0, fmt.Errorf("invalid direction %q", s[0])
	}
}

// NewPosition calculates the final position after a move and the number of
// times the dial passed '0'.
func NewPosition(p, move int) (int, int) {
	pos := p + move
	z := Abs(pos / TotalPositions)

	if pos > 0 && pos%TotalPositions == 0 {
		// if we land on mult of totalpositions, need to decrease count
		z--
	}

	pos %= TotalPositions

	switch {
	case pos > 0:
		// do nothing
		break
	case pos == 0:
		z++
	case pos < 0:
		pos += TotalPositions
		if p != 0 {
			z++
		}
	default:
		panic("this isn't possible")
	}

	// fmt.Printf("p=%d,move=%d,pos=%d,crossed=%d\n", p, move, pos, z)

	return pos, z
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
