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
	count := 0
	lights := map[coord]bool{}

	for _, a := range actions {
		switch a.action {
		case "on":
			for i := a.start.x; i <= a.end.x; i++ {
				for j := a.start.y; j <= a.end.y; j++ {
					lights[coord{i, j}] = true
				}
			}
		case "off":
			for i := a.start.x; i <= a.end.x; i++ {
				for j := a.start.y; j <= a.end.y; j++ {
					lights[coord{i, j}] = false
				}
			}
		case "toggle":
			for i := a.start.x; i <= a.end.x; i++ {
				for j := a.start.y; j <= a.end.y; j++ {
					lights[coord{i, j}] = !lights[coord{i, j}]
				}
			}
		}
	}

	for _, v := range lights {
		if v {
			count++
		}
	}

	return count, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	actions := parse(instr)
	count := 0
	lights := map[coord]int{}

	for _, a := range actions {
		switch a.action {
		case "on":
			for i := a.start.x; i <= a.end.x; i++ {
				for j := a.start.y; j <= a.end.y; j++ {
					lights[coord{i, j}]++
				}
			}
		case "off":
			for i := a.start.x; i <= a.end.x; i++ {
				for j := a.start.y; j <= a.end.y; j++ {
					lights[coord{i, j}]--

					if lights[coord{i, j}] < 0 {
						lights[coord{i, j}] = 0
					}
				}
			}
		case "toggle":
			for i := a.start.x; i <= a.end.x; i++ {
				for j := a.start.y; j <= a.end.y; j++ {
					lights[coord{i, j}] += 2
				}
			}
		}
	}

	for _, v := range lights {
		count += v
	}

	return count, nil
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

		if strings.HasPrefix(line, "turn") {
			_, _ = fmt.Sscanf(line, "turn %s %d,%d through %d,%d", &a, &x1, &y1, &x2, &y2)
		} else {
			_, _ = fmt.Sscanf(line, "%s %d,%d through %d,%d", &a, &x1, &y1, &x2, &y2)
		}

		out = append(out, action{
			action: a,
			start:  coord{x: x1, y: y1},
			end:    coord{x: x2, y: y2},
		})
	}

	return out
}
