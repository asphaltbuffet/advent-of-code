package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 21.
type Exercise struct {
	common.BaseExercise
}

// grid is a square pattern held as rows of runes ('#'/'.').
type grid []string

// parseRules builds a lookup from every rotation/flip of each rule's input to
// its output pattern.
func parseRules(instr string) map[string]grid {
	rules := map[string]grid{}
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " => ", 2)
		in := grid(strings.Split(parts[0], "/"))
		out := grid(strings.Split(parts[1], "/"))
		for _, o := range orientations(in) {
			rules[o.key()] = out
		}
	}
	return rules
}

// key joins the rows so a grid can index a map.
func (g grid) key() string { return strings.Join(g, "/") }

// rotate returns g rotated 90° clockwise.
func (g grid) rotate() grid {
	n := len(g)
	out := make(grid, n)
	for r := 0; r < n; r++ {
		var sb strings.Builder
		for c := 0; c < n; c++ {
			sb.WriteByte(g[n-1-c][r])
		}
		out[r] = sb.String()
	}
	return out
}

// flip returns g mirrored horizontally.
func (g grid) flip() grid {
	out := make(grid, len(g))
	for r, row := range g {
		b := []byte(row)
		for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
			b[i], b[j] = b[j], b[i]
		}
		out[r] = string(b)
	}
	return out
}

// orientations returns the 8 rotations/flips of g.
func orientations(g grid) []grid {
	var out []grid
	cur := g
	for i := 0; i < 4; i++ {
		out = append(out, cur, cur.flip())
		cur = cur.rotate()
	}
	return out
}

// enhance applies one iteration: split into 2x2 or 3x3 blocks and replace each.
func enhance(g grid, rules map[string]grid) grid {
	n := len(g)
	var block int
	if n%2 == 0 {
		block = 2
	} else {
		block = 3
	}
	newBlock := block + 1
	count := n / block
	out := make(grid, count*newBlock)

	for br := 0; br < count; br++ {
		rows := make([]strings.Builder, newBlock)
		for bc := 0; bc < count; bc++ {
			// Extract the sub-block.
			sub := make(grid, block)
			for r := 0; r < block; r++ {
				sub[r] = g[br*block+r][bc*block : bc*block+block]
			}
			rep := rules[sub.key()]
			for r := 0; r < newBlock; r++ {
				rows[r].WriteString(rep[r])
			}
		}
		for r := 0; r < newBlock; r++ {
			out[br*newBlock+r] = rows[r].String()
		}
	}
	return out
}

// run applies iters enhancement steps and returns the final grid.
func run(instr string, iters int) grid {
	rules := parseRules(instr)
	g := grid{".#.", "..#", "###"}
	for i := 0; i < iters; i++ {
		g = enhance(g, rules)
	}
	return g
}

// countOn returns the number of '#' pixels in a grid.
func countOn(g grid) int {
	n := 0
	for _, row := range g {
		n += strings.Count(row, "#")
	}
	return n
}

// iterations picks the count: the small example (few rules) runs 2 steps, the
// real puzzle 5 for Part One.
func onePasses(instr string) int {
	if isExample(instr) {
		return 2
	}
	return 5
}

// isExample detects the tiny two-rule example set.
func isExample(instr string) bool {
	return strings.Count(strings.TrimSpace(instr), "\n") <= 1
}

// One returns the number of pixels on after the Part One iteration count.
func (e Exercise) One(instr string) (any, error) {
	return countOn(run(instr, onePasses(instr))), nil
}

// Two returns the number of pixels on after 18 iterations.
func (e Exercise) Two(instr string) (any, error) {
	return countOn(run(instr, 18)), nil
}

// Vis animates all 18 enhancement passes (GIF) on a canvas sized to the final
// grid, so no frame is ever downsampled. Every cell maps to a whole number of
// pixels: early grids are integer-scaled up and centred, and the final 2187×2187
// grid lands at exactly one pixel per cell.
func (e Exercise) Vis(instr, outdir string) error {
	const passes = 18

	rules := parseRules(instr)
	g := grid{".#.", "..#", "###"}

	frames := []grid{g}
	for i := 0; i < passes; i++ {
		g = enhance(g, rules)
		frames = append(frames, g)
	}
	canvas := len(frames[len(frames)-1]) // 2187: fits the final grid at 1px/cell

	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff}, // 0 off
		color.RGBA{0x2f, 0xd0, 0x9a, 0xff}, // 1 on
	}

	anim := &gif.GIF{}
	for fi, fr := range frames {
		n := len(fr)
		cell := canvas / n // whole pixels per cell; no downsampling
		if cell < 1 {
			cell = 1
		}
		span := n * cell
		off := (canvas - span) / 2 // centre the scaled grid

		img := image.NewPaletted(image.Rect(0, 0, canvas, canvas), pal)
		for r := 0; r < n; r++ {
			row := fr[r]
			for c := 0; c < n; c++ {
				if row[c] != '#' {
					continue
				}
				x0, y0 := off+c*cell, off+r*cell
				for dy := 0; dy < cell; dy++ {
					for dx := 0; dx < cell; dx++ {
						img.SetColorIndex(x0+dx, y0+dy, 1)
					}
				}
			}
		}
		anim.Image = append(anim.Image, img)
		delay := 40
		if fi == 0 || fi == len(frames)-1 {
			delay = 200
		}
		anim.Delay = append(anim.Delay, delay)
	}

	f, err := os.Create(filepath.Join(outdir, "fractal-art.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}
