package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 11.
type Exercise struct {
	common.BaseExercise
}

type grid struct {
	e    [][]int
	rows int
	cols int
}

func parse(instr string) (grid, error) {
	var g grid

	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		row := make([]int, len(line))
		for i, c := range line {
			if c < '0' || c > '9' {
				return g, fmt.Errorf("non-digit %q in grid", c)
			}
			row[i] = int(c - '0')
		}
		g.e = append(g.e, row)
	}

	g.rows = len(g.e)
	if g.rows > 0 {
		g.cols = len(g.e[0])
	}

	return g, nil
}

var neighbors8 = [8][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

// step advances the grid one iteration and returns the number of flashes. Every
// octopus gains energy; any over 9 flashes, bumping its neighbors via a work
// queue so cascades are handled, and all flashed octopuses reset to 0.
func (g *grid) step() int {
	var queue [][2]int
	flashed := make([][]bool, g.rows)
	for r := range flashed {
		flashed[r] = make([]bool, g.cols)
	}

	// Charge everyone; queue any that cross the flash threshold.
	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			g.e[r][c]++
			if g.e[r][c] > 9 {
				queue = append(queue, [2]int{r, c})
				flashed[r][c] = true
			}
		}
	}

	// Propagate the cascade.
	count := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		count++

		for _, d := range neighbors8 {
			nr, nc := cur[0]+d[0], cur[1]+d[1]
			if nr < 0 || nr >= g.rows || nc < 0 || nc >= g.cols {
				continue
			}
			g.e[nr][nc]++
			if g.e[nr][nc] > 9 && !flashed[nr][nc] {
				flashed[nr][nc] = true
				queue = append(queue, [2]int{nr, nc})
			}
		}
	}

	// Reset everyone who flashed.
	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			if flashed[r][c] {
				g.e[r][c] = 0
			}
		}
	}

	return count
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	g, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	total := 0
	for i := 0; i < 100; i++ {
		total += g.step()
	}

	return fmt.Sprintf("%d", total), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	g, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	all := g.rows * g.cols
	for step := 1; ; step++ {
		if g.step() == all {
			return fmt.Sprintf("%d", step), nil
		}
	}
}
