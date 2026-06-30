package exercises

import (
	"image"
	"image/color"
	"image/png"
	"math/bits"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 13.
type Exercise struct {
	common.BaseExercise
}

type point struct{ x, y int }

// isOpen reports whether (x,y) is an open space for the given favorite number.
func isOpen(x, y, fav int) bool {
	if x < 0 || y < 0 {
		return false
	}
	v := x*x + 3*x + 2*x*y + y + y*y + fav
	return bits.OnesCount(uint(v))%2 == 0
}

var steps4 = [4]point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

// bfs explores from (1,1), returning steps-to-target (or -1) and the count of
// distinct cells reachable within maxSteps.
func bfs(fav, targetX, targetY, maxSteps int) (toTarget, reachable int) {
	start := point{1, 1}
	dist := map[point]int{start: 0}
	queue := []point{start}
	toTarget = -1
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		d := dist[cur]
		if cur.x == targetX && cur.y == targetY {
			toTarget = d
		}
		if d >= maxSteps {
			continue
		}
		for _, s := range steps4 {
			n := point{cur.x + s.x, cur.y + s.y}
			if !isOpen(n.x, n.y, fav) {
				continue
			}
			if _, ok := dist[n]; !ok {
				dist[n] = d + 1
				queue = append(queue, n)
			}
		}
	}
	return toTarget, len(dist)
}

// target returns the goal coordinate: (7,4) for the favorite-number-10 example,
// (31,39) for the real input.
func target(fav int) (int, int) {
	if fav == 10 {
		return 7, 4
	}
	return 31, 39
}

func parseFav(instr string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(instr))
	return n
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	fav := parseFav(instr)
	tx, ty := target(fav)
	toTarget, _ := bfs(fav, tx, ty, 1<<30)
	return toTarget, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	fav := parseFav(instr)
	tx, ty := target(fav)
	_, reachable := bfs(fav, tx, ty, 50)
	return reachable, nil
}

// --- Visualization ---

// shortestPath returns the cells of a shortest path from (1,1) to (tx,ty).
func shortestPath(fav, tx, ty int) []point {
	start, goal := point{1, 1}, point{tx, ty}
	prev := map[point]point{start: start}
	queue := []point{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == goal {
			var path []point
			for c := goal; ; c = prev[c] {
				path = append(path, c)
				if c == start {
					break
				}
			}
			return path
		}
		for _, s := range steps4 {
			n := point{cur.x + s.x, cur.y + s.y}
			if !isOpen(n.x, n.y, fav) {
				continue
			}
			if _, ok := prev[n]; !ok {
				prev[n] = cur
				queue = append(queue, n)
			}
		}
	}
	return nil
}

// reachableSet returns the cells reachable within maxSteps of (1,1).
func reachableSet(fav, maxSteps int) map[point]bool {
	start := point{1, 1}
	dist := map[point]int{start: 0}
	out := map[point]bool{start: true}
	queue := []point{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if dist[cur] >= maxSteps {
			continue
		}
		for _, s := range steps4 {
			n := point{cur.x + s.x, cur.y + s.y}
			if !isOpen(n.x, n.y, fav) {
				continue
			}
			if _, ok := dist[n]; !ok {
				dist[n] = dist[cur] + 1
				out[n] = true
				queue = append(queue, n)
			}
		}
	}
	return out
}

// Vis renders the maze with the ≤50-step reachable region shaded and the
// shortest path to the target highlighted.
func (e Exercise) Vis(instr string, outdir string) error {
	fav := parseFav(instr)
	tx, ty := target(fav)
	path := shortestPath(fav, tx, ty)
	reach := reachableSet(fav, 50)

	// Window: enough to show the target and the reach region.
	const w, h = 50, 50
	pathSet := map[point]bool{}
	for _, p := range path {
		pathSet[p] = true
	}

	const cell = 16
	const pad = 12
	img := image.NewRGBA(image.Rect(0, 0, w*cell+2*pad, h*cell+2*pad))

	wall := color.RGBA{0x12, 0x14, 0x20, 0xff}
	open := color.RGBA{0x26, 0x2c, 0x3e, 0xff}
	reachC := color.RGBA{0x2f, 0x5a, 0x86, 0xff}
	pathC := color.RGBA{0xff, 0xc8, 0x4a, 0xff}
	startC := color.RGBA{0x44, 0xff, 0x88, 0xff}
	goalC := color.RGBA{0xff, 0x44, 0x55, 0xff}

	fill := func(x, y int, col color.RGBA) {
		x0, y0 := pad+x*cell, pad+y*cell
		for yy := y0; yy < y0+cell; yy++ {
			for xx := x0; xx < x0+cell; xx++ {
				img.SetRGBA(xx, yy, col)
			}
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := point{x, y}
			switch {
			case !isOpen(x, y, fav):
				fill(x, y, wall)
			case pathSet[p]:
				fill(x, y, pathC)
			case reach[p]:
				fill(x, y, reachC)
			default:
				fill(x, y, open)
			}
		}
	}
	fill(1, 1, startC)
	fill(tx, ty, goalC)

	f, err := os.Create(filepath.Join(outdir, "twisty-cubicles.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
