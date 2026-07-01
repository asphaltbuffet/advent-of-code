package exercises

import (
	"fmt"
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 9.
type Exercise struct {
	common.BaseExercise
}

type grid struct {
	h    [][]int
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
				return g, fmt.Errorf("non-digit %q in heightmap", c)
			}
			row[i] = int(c - '0')
		}
		g.h = append(g.h, row)
	}

	g.rows = len(g.h)
	if g.rows > 0 {
		g.cols = len(g.h[0])
	}

	return g, nil
}

var dirs = [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

// lowPoints returns the coordinates of cells lower than all orthogonal neighbors.
func (g grid) lowPoints() [][2]int {
	var lows [][2]int

	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			low := true
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1]
				if nr < 0 || nr >= g.rows || nc < 0 || nc >= g.cols {
					continue
				}
				if g.h[nr][nc] <= g.h[r][c] {
					low = false
					break
				}
			}
			if low {
				lows = append(lows, [2]int{r, c})
			}
		}
	}

	return lows
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	g, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	risk := 0
	for _, p := range g.lowPoints() {
		risk += 1 + g.h[p[0]][p[1]]
	}

	return fmt.Sprintf("%d", risk), nil
}

// basinSizes flood-fills every non-9 connected component and returns their sizes.
// Height-9 cells are walls that separate basins, so each component is one basin.
func (g grid) basinSizes() []int {
	seen := make([][]bool, g.rows)
	for i := range seen {
		seen[i] = make([]bool, g.cols)
	}

	var sizes []int
	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			if seen[r][c] || g.h[r][c] == 9 {
				continue
			}

			// BFS the basin containing (r,c).
			size := 0
			queue := [][2]int{{r, c}}
			seen[r][c] = true
			for len(queue) > 0 {
				cur := queue[0]
				queue = queue[1:]
				size++
				for _, d := range dirs {
					nr, nc := cur[0]+d[0], cur[1]+d[1]
					if nr < 0 || nr >= g.rows || nc < 0 || nc >= g.cols {
						continue
					}
					if seen[nr][nc] || g.h[nr][nc] == 9 {
						continue
					}
					seen[nr][nc] = true
					queue = append(queue, [2]int{nr, nc})
				}
			}
			sizes = append(sizes, size)
		}
	}

	return sizes
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	g, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	sizes := g.basinSizes()
	if len(sizes) < 3 {
		return nil, fmt.Errorf("expected at least 3 basins, found %d", len(sizes))
	}

	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	return fmt.Sprintf("%d", sizes[0]*sizes[1]*sizes[2]), nil
}
