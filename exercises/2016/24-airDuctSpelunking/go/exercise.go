package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 24.
type Exercise struct {
	common.BaseExercise
}

type point struct{ r, c int }

// parseMaze returns the grid rows and the coordinates of each numbered target,
// indexed by digit.
func parseMaze(instr string) ([]string, map[int]point) {
	grid := strings.Split(strings.TrimRight(instr, "\n"), "\n")
	targets := map[int]point{}
	for r, row := range grid {
		for c, ch := range row {
			if ch >= '0' && ch <= '9' {
				targets[int(ch-'0')] = point{r, c}
			}
		}
	}
	return grid, targets
}

// bfsFrom returns the shortest-path distance from src to every other numbered
// target, keyed by digit.
func bfsFrom(grid []string, src point, targets map[int]point) map[int]int {
	dist := map[point]int{src: 0}
	queue := []point{src}
	out := map[int]int{}
	// reverse lookup: coordinate -> digit
	at := map[point]int{}
	for d, p := range targets {
		at[p] = d
	}
	if d, ok := at[src]; ok {
		out[d] = 0
	}

	dirs := [4]point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, d := range dirs {
			np := point{cur.r + d.r, cur.c + d.c}
			if np.r < 0 || np.r >= len(grid) || np.c < 0 || np.c >= len(grid[np.r]) {
				continue
			}
			if grid[np.r][np.c] == '#' {
				continue
			}
			if _, seen := dist[np]; seen {
				continue
			}
			dist[np] = dist[cur] + 1
			if d, ok := at[np]; ok {
				out[d] = dist[np]
			}
			queue = append(queue, np)
		}
	}
	return out
}

// distMatrix builds the pairwise shortest-path distances between all numbered
// targets.
func distMatrix(grid []string, targets map[int]point) [][]int {
	n := len(targets)
	m := make([][]int, n)
	for i := range n {
		m[i] = make([]int, n)
		row := bfsFrom(grid, targets[i], targets)
		for j := range n {
			m[i][j] = row[j]
		}
	}
	return m
}

// tsp returns the fewest steps to start at 0 and visit every target, plus the
// visiting order (digits after the initial 0). If returnToStart is set, the
// tour must also return to 0.
func tsp(m [][]int, returnToStart bool) (int, []int) {
	n := len(m)
	others := make([]int, 0, n-1)
	for i := 1; i < n; i++ {
		others = append(others, i)
	}

	best := -1
	var bestOrder []int
	permute(others, 0, func(order []int) {
		cost, prev := 0, 0
		for _, next := range order {
			cost += m[prev][next]
			prev = next
		}
		if returnToStart {
			cost += m[prev][0]
		}
		if best < 0 || cost < best {
			best = cost
			bestOrder = append([]int(nil), order...)
		}
	})
	return best, bestOrder
}

// permute yields every ordering of s (in-place swap generator).
func permute(s []int, k int, yield func([]int)) {
	if k == len(s) {
		yield(s)
		return
	}
	for i := k; i < len(s); i++ {
		s[k], s[i] = s[i], s[k]
		permute(s, k+1, yield)
		s[k], s[i] = s[i], s[k]
	}
}

// One returns the fewest steps to visit every numbered location from 0.
func (e Exercise) One(instr string) (any, error) {
	grid, targets := parseMaze(instr)
	steps, _ := tsp(distMatrix(grid, targets), false)
	return steps, nil
}

// Two returns the fewest steps to visit every location and return to 0.
func (e Exercise) Two(instr string) (any, error) {
	grid, targets := parseMaze(instr)
	steps, _ := tsp(distMatrix(grid, targets), true)
	return steps, nil
}

// bfsPath reconstructs the shortest path (inclusive of both ends) between two
// grid points, avoiding walls.
func reconstructPath(prev map[point]point, src, dst point) []point {
	var path []point
	for p := dst; ; p = prev[p] {
		path = append([]point{p}, path...)
		if p == src {
			break
		}
	}
	return path
}

func bfsPath(grid []string, src, dst point) []point {
	prev := map[point]point{}
	dist := map[point]int{src: 0}
	queue := []point{src}
	dirs := [4]point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == dst {
			return reconstructPath(prev, src, dst)
		}
		for _, d := range dirs {
			np := point{cur.r + d.r, cur.c + d.c}
			if np.r < 0 || np.r >= len(grid) || np.c < 0 || np.c >= len(grid[np.r]) {
				continue
			}
			if grid[np.r][np.c] == '#' {
				continue
			}
			if _, seen := dist[np]; seen {
				continue
			}
			dist[np] = dist[cur] + 1
			prev[np] = cur
			queue = append(queue, np)
		}
	}
	return nil
}

// Vis renders the maze, the numbered targets, and the optimal collection route
// (returning to 0, as in Part Two).
func cellColor(grid []string, route map[point]bool, r, c int, wallC, routeC, open color.RGBA) color.RGBA {
	p := point{r, c}
	switch {
	case grid[r][c] == '#':
		return wallC
	case route[p]:
		return routeC
	default:
		return open
	}
}

func buildRoute(grid []string, targets map[int]point) map[point]bool {
	_, order := tsp(distMatrix(grid, targets), true)
	stops := append([]int{0}, order...)
	stops = append(stops, 0)
	route := map[point]bool{}
	for i := 0; i+1 < len(stops); i++ {
		for _, p := range bfsPath(grid, targets[stops[i]], targets[stops[i+1]]) {
			route[p] = true
		}
	}
	return route
}

// Vis renders the maze with the optimal route and target positions highlighted.
// Vis renders the maze with the optimal route and target positions highlighted.
func (e Exercise) Vis(instr, outdir string) error {
	grid, targets := parseMaze(instr)
	route := buildRoute(grid, targets)

	rows := len(grid)
	cols := 0
	for _, r := range grid {
		if len(r) > cols {
			cols = len(r)
		}
	}

	const cell = 6
	const pad = 8
	img := image.NewRGBA(image.Rect(0, 0, cols*cell+2*pad, rows*cell+2*pad))

	bg := color.RGBA{0x12, 0x14, 0x20, 0xff}
	wallC := color.RGBA{0x2a, 0x30, 0x44, 0xff}
	open := color.RGBA{0x1a, 0x1e, 0x2e, 0xff}
	routeC := color.RGBA{0x2f, 0x8a, 0x86, 0xff}
	startC := color.RGBA{0xff, 0x44, 0x55, 0xff}
	targetC := color.RGBA{0xff, 0xc8, 0x4a, 0xff}

	fill := func(c, r int, col color.RGBA) {
		x0, y0 := pad+c*cell, pad+r*cell
		for yy := y0; yy < y0+cell; yy++ {
			for xx := x0; xx < x0+cell; xx++ {
				img.SetRGBA(xx, yy, col)
			}
		}
	}

	ih, iw := img.Bounds().Dy(), img.Bounds().Dx()
	for y := range ih {
		for x := range iw {
			img.SetRGBA(x, y, bg)
		}
	}

	for r := range rows {
		for c := range len(grid[r]) {
			fill(c, r, cellColor(grid, route, r, c, wallC, routeC, open))
		}
	}
	// Draw targets on top: 0 red, the rest gold.
	for d, p := range targets {
		col := targetC
		if d == 0 {
			col = startC
		}
		fill(p.c, p.r, col)
	}

	f, err := os.Create(filepath.Join(outdir, "air-duct-spelunking.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
