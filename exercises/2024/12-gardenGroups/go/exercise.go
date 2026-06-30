package exercises

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 12.
type Exercise struct {
	common.BaseExercise
}

// region holds the accumulated metrics for one connected same-letter region.
type region struct {
	area, perimeter, sides int
}

// regions flood-fills the grid and returns one region per connected component.
func regions(instr string) []region {
	grid := strings.Fields(instr)
	rows := len(grid)
	at := func(r, c int) byte {
		if r < 0 || r >= rows || c < 0 || c >= len(grid[r]) {
			return 0
		}
		return grid[r][c]
	}

	seen := make([][]bool, rows)
	for r := range seen {
		seen[r] = make([]bool, len(grid[r]))
	}

	var out []region
	for sr := 0; sr < rows; sr++ {
		for sc := 0; sc < len(grid[sr]); sc++ {
			if seen[sr][sc] {
				continue
			}
			letter := grid[sr][sc]
			reg := region{}
			stack := [][2]int{{sr, sc}}
			seen[sr][sc] = true
			for len(stack) > 0 {
				cur := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				r, c := cur[0], cur[1]
				reg.area++

				same := func(nr, nc int) bool { return at(nr, nc) == letter }
				// Perimeter: edges to a non-matching neighbour.
				for _, d := range [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
					nr, nc := r+d[0], c+d[1]
					if !same(nr, nc) {
						reg.perimeter++
					} else if !seen[nr][nc] {
						seen[nr][nc] = true
						stack = append(stack, [2]int{nr, nc})
					}
				}

				// Corners (== number of sides). For each of the 4 diagonal
				// quadrants, count a convex or concave corner.
				for _, q := range [4][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}} {
					vert := same(r+q[0], c)   // vertical neighbour
					horiz := same(r, c+q[1])  // horizontal neighbour
					diag := same(r+q[0], c+q[1])
					if !vert && !horiz {
						reg.sides++ // convex corner
					} else if vert && horiz && !diag {
						reg.sides++ // concave corner
					}
				}
			}
			out = append(out, reg)
		}
	}
	return out
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	total := 0
	for _, reg := range regions(instr) {
		total += reg.area * reg.perimeter
	}
	return total, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	total := 0
	for _, reg := range regions(instr) {
		total += reg.area * reg.sides
	}
	return total, nil
}

// --- Visualization ---

const (
	visCell  = 10 // pixels per garden plot
	visPad   = 16 // outer margin
	visFence = 2  // fence line thickness in pixels
)

var (
	visBG    = color.RGBA{0x10, 0x10, 0x14, 0xff}
	visFenceC = color.RGBA{0x10, 0x10, 0x14, 0xff} // fence matches bg for crisp borders
)

// letterColor maps a plant letter to a stable, well-separated color via its
// hue. Same letter -> same color.
func letterColor(ch byte) color.RGBA {
	// Spread 26+ letters around the hue circle using a large coprime step.
	hue := math.Mod(float64(ch)*47.0, 360.0)
	return hsv(hue, 0.55, 0.85)
}

// hsv converts H (0-360), S, V (0-1) to an RGBA color.
func hsv(h, s, v float64) color.RGBA {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c
	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	return color.RGBA{
		uint8((r + m) * 255),
		uint8((g + m) * 255),
		uint8((b + m) * 255),
		0xff,
	}
}

// Vis renders the garden coloured by plant letter with fences on region edges.
func (e Exercise) Vis(instr string, outdir string) error {
	grid := strings.Fields(instr)
	rows := len(grid)
	cols := 0
	for _, row := range grid {
		if len(row) > cols {
			cols = len(row)
		}
	}
	at := func(r, c int) byte {
		if r < 0 || r >= rows || c < 0 || c >= len(grid[r]) {
			return 0
		}
		return grid[r][c]
	}

	w := cols*visCell + 2*visPad
	h := rows*visCell + 2*visPad
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, visBG)
		}
	}

	// Fill each plot with its letter colour.
	for r := 0; r < rows; r++ {
		for c := 0; c < len(grid[r]); c++ {
			fillRect(img, visPad+c*visCell, visPad+r*visCell, visCell, visCell, letterColor(grid[r][c]))
		}
	}

	// Draw fences: every cell edge where the neighbour differs.
	for r := 0; r < rows; r++ {
		for c := 0; c < len(grid[r]); c++ {
			letter := grid[r][c]
			x0 := visPad + c*visCell
			y0 := visPad + r*visCell
			if at(r-1, c) != letter { // top
				fillRect(img, x0, y0, visCell, visFence, visFenceC)
			}
			if at(r+1, c) != letter { // bottom
				fillRect(img, x0, y0+visCell-visFence, visCell, visFence, visFenceC)
			}
			if at(r, c-1) != letter { // left
				fillRect(img, x0, y0, visFence, visCell, visFenceC)
			}
			if at(r, c+1) != letter { // right
				fillRect(img, x0+visCell-visFence, y0, visFence, visCell, visFenceC)
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "garden-groups.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func fillRect(img *image.RGBA, x0, y0, w, h int, col color.RGBA) {
	for y := y0; y < y0+h; y++ {
		for x := x0; x < x0+w; x++ {
			img.SetRGBA(x, y, col)
		}
	}
}
