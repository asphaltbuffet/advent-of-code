package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 12.
type Exercise struct {
	common.BaseExercise
}

// move is one navigation instruction: an action letter and its value.
type move struct {
	action byte
	value  int
}

func parse(instr string) ([]move, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	moves := make([]move, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		v, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %w", line, err)
		}
		moves = append(moves, move{line[0], v})
	}
	return moves, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// One steers the ship directly: N/S/E/W shift it, L/R rotate its heading, F moves
// it forward along that heading. Returns the Manhattan distance from the start.
func (e Exercise) One(instr string) (any, error) {
	moves, err := parse(instr)
	if err != nil {
		return nil, err
	}

	x, y := 0, 0
	// Heading as a unit vector; start facing east.
	dx, dy := 1, 0

	for _, m := range moves {
		switch m.action {
		case 'N':
			y += m.value
		case 'S':
			y -= m.value
		case 'E':
			x += m.value
		case 'W':
			x -= m.value
		case 'F':
			x += dx * m.value
			y += dy * m.value
		case 'L':
			dx, dy = rotate(dx, dy, m.value)
		case 'R':
			dx, dy = rotate(dx, dy, 360-m.value)
		}
	}

	return fmt.Sprintf("%d", abs(x)+abs(y)), nil
}

// rotate turns the vector (dx, dy) counterclockwise by deg (a multiple of 90).
func rotate(dx, dy, deg int) (int, int) {
	for deg = ((deg % 360) + 360) % 360; deg > 0; deg -= 90 {
		dx, dy = -dy, dx
	}
	return dx, dy
}

// Two steers by a waypoint: N/S/E/W move the waypoint (relative to the ship),
// L/R rotate it around the ship, and F moves the ship toward the waypoint that
// many times. Returns the Manhattan distance from the start.
func (e Exercise) Two(instr string) (any, error) {
	moves, err := parse(instr)
	if err != nil {
		return nil, err
	}

	x, y := 0, 0
	// Waypoint offset from the ship; starts 10 east and 1 north.
	wx, wy := 10, 1

	for _, m := range moves {
		switch m.action {
		case 'N':
			wy += m.value
		case 'S':
			wy -= m.value
		case 'E':
			wx += m.value
		case 'W':
			wx -= m.value
		case 'F':
			x += wx * m.value
			y += wy * m.value
		case 'L':
			wx, wy = rotate(wx, wy, m.value)
		case 'R':
			wx, wy = rotate(wx, wy, 360-m.value)
		}
	}

	return fmt.Sprintf("%d", abs(x)+abs(y)), nil
}
