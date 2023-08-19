package exercises

import "github.com/asphaltbuffet/advent-of-code/internal/common"

// Exercise for Advent of Code 2015 day 3.
type Exercise struct {
	common.BaseExercise
}

type coord struct {
	x, y int
}

// One returns the answer to the first part of the exercise.
// answer:
func (e Exercise) One(instr string) (any, error) {
	santa := coord{0, 0}
	visited := make(map[coord]bool)

	visited[santa] = true

	for _, r := range instr {
		santa = move(santa, r)
		visited[santa] = true
	}

	return len(visited), nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (e Exercise) Two(instr string) (any, error) {
	santa := coord{0, 0}
	robo := coord{0, 0}
	visited := make(map[coord]bool)

	visited[santa] = true

	for i, r := range instr {
		if i%2 == 0 {
			santa = move(santa, r)
			visited[santa] = true
		} else {
			robo = move(robo, r)
			visited[robo] = true
		}
	}

	return len(visited), nil
}

func move(c coord, r rune) coord {
	switch r {
	case '^':
		c.y--
	case 'v':
		c.y++
	case '>':
		c.x++
	case '<':
		c.x--
	}

	return c
}
