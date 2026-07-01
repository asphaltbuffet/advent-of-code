package exercises

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 3.
type Exercise struct {
	common.BaseExercise
}

type coord struct{ x, y int }

// spiralWalk yields the grid coordinates of the Ulam-style spiral in order,
// starting at the origin (square 1) and turning counterclockwise. It calls
// yield for each square; yield returns false to stop.
func spiralWalk(yield func(coord) bool) {
	// Direction order: right, up, left, down (counterclockwise).
	dirs := []coord{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	pos := coord{0, 0}
	if !yield(pos) {
		return
	}
	dir, runLen, sinceTurn, legsAtLen := 0, 1, 0, 0
	for {
		d := dirs[dir]
		pos.x += d.x
		pos.y += d.y
		if !yield(pos) {
			return
		}
		sinceTurn++
		if sinceTurn == runLen {
			sinceTurn = 0
			dir = (dir + 1) % 4
			// Run length grows every two legs: 1,1,2,2,3,3,...
			legsAtLen++
			if legsAtLen == 2 {
				legsAtLen = 0
				runLen++
			}
		}
	}
}

func parseTarget(instr string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(instr))
	return n
}

// One returns the Manhattan distance from square N back to square 1.
func (e Exercise) One(instr string) (any, error) {
	target := parseTarget(instr)
	dist := 0
	i := 1
	spiralWalk(func(c coord) bool {
		if i == target {
			dist = abs(c.x) + abs(c.y)
			return false
		}
		i++
		return true
	})
	return dist, nil
}

// Two returns the first value written that exceeds the puzzle input, where each
// square stores the sum of its already-filled (including diagonal) neighbours.
func (e Exercise) Two(instr string) (any, error) {
	target := parseTarget(instr)
	vals := map[coord]int{}
	result := 0
	first := true
	spiralWalk(func(c coord) bool {
		if first {
			vals[c] = 1
			first = false
			if 1 > target {
				result = 1
				return false
			}
			return true
		}
		sum := 0
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}
				sum += vals[coord{c.x + dx, c.y + dy}]
			}
		}
		vals[c] = sum
		if sum > target {
			result = sum
			return false
		}
		return true
	})
	return result, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Vis renders the Part Two summed-neighbour spiral as a heat map. Each square
// stores the sum of its filled neighbours; cells are coloured by log-scaled
// magnitude so the exponential growth radiating from the centre is visible.
//
// The spiral is filled well past the point where a value first exceeds the
// puzzle input (out to a fixed window) so the pattern is large and detailed;
// the answer square — the first value over the target — is marked distinctly.
func (e Exercise) Vis(instr, outdir string) error {
	target := parseTarget(instr)

	// Render a fixed window so the image is large regardless of how early the
	// target is exceeded. 27x27 keeps the whole spiral square and centred.
	const half = 13 // window spans -half..+half on each axis
	limit := (2*half + 1) * (2*half + 1)

	vals := map[coord]int{}
	first := true
	answer := coord{} // first square whose value exceeds the target
	foundAnswer := false
	maxV := 1
	count := 0
	spiralWalk(func(c coord) bool {
		if c.x < -half || c.x > half || c.y < -half || c.y > half {
			return true // outside the window: skip drawing but keep walking
		}
		count++
		if first {
			vals[c] = 1
				first = false
			return count < limit
		}
		sum := 0
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}
				sum += vals[coord{c.x + dx, c.y + dy}]
			}
		}
		vals[c] = sum
		if sum > maxV {
			maxV = sum
		}
		if !foundAnswer && sum > target {
			answer = c
			foundAnswer = true
		}
		return count < limit
	})

	minX, maxX, minY, maxY := -half, half, -half, half
	cols, rows := maxX-minX+1, maxY-minY+1

	const cell = 30
	const pad = 16
	img := image.NewRGBA(image.Rect(0, 0, cols*cell+2*pad, rows*cell+2*pad))

	bg := color.RGBA{0x10, 0x12, 0x1c, 0xff}
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	// heat maps a normalized 0..1 value onto a dark-blue -> cyan -> gold ramp.
	heat := func(t float64) color.RGBA {
		switch {
		case t < 0.5:
			u := t / 0.5
			return color.RGBA{uint8(0x1a + u*0x15), uint8(0x2c + u*0x5e), uint8(0x50 + u*0x36), 0xff}
		default:
			u := (t - 0.5) / 0.5
			return color.RGBA{uint8(0x2f + u*0xd0), uint8(0x8a + u*0x3e), uint8(0x86 - u*0x3c), 0xff}
		}
	}

	logMax := math.Log(float64(maxV) + 1)
	fill := func(c coord, col color.RGBA) {
		gx, gy := c.x-minX, maxY-c.y // flip y so +y is up
		x0, y0 := pad+gx*cell, pad+gy*cell
		for yy := y0; yy < y0+cell; yy++ {
			for xx := x0; xx < x0+cell; xx++ {
				img.SetRGBA(xx, yy, col)
			}
		}
	}

	for c, v := range vals {
		t := math.Log(float64(v)+1) / logMax
		fill(c, heat(t))
	}
	// Mark the centre (square 1) in red.
	fill(coord{0, 0}, color.RGBA{0xff, 0x44, 0x55, 0xff})
	// Mark the answer square (first value > target) in white, if it lies within
	// the rendered window.
	if foundAnswer {
		fill(answer, color.RGBA{0xff, 0xff, 0xff, 0xff})
	}

	f, err := os.Create(filepath.Join(outdir, "spiral-memory.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
