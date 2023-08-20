package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 6.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	actions := parse(instr)
	fmt.Printf("%v\n", actions)

	lights := map[coord]bool{}

	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			lights[coord{i, j}] = false
		}
	}

	count := 0
	for _, v := range lights {
		if v {
			count++
		}
	}

	return count, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return nil, fmt.Errorf("part 2 not implemented")
}

type action struct {
	action string
	start  coord
	end    coord
}

type coord struct {
	x, y int
}

func parse(s string) []action {
	out := []action{}
	for _, line := range strings.Split(s, "\n") {
		var x1, y1, x2, y2 int
		var a string

		_, _ = fmt.Sscanf(line, "%s %d,%d through %d,%d", &a, &x1, &y1, &x2, &y2)
		out = append(out, action{
			action: a,
			start:  coord{x: x1, y: y1},
			end:    coord{x: x2, y: y2},
		})
	}

	return out
}
