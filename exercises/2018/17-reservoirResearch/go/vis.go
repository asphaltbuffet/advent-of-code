package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Grayscale-safe fills (Okabe-Ito derived). Meaning is carried by brightness:
// settled water brightest, then flowing, then clay, over near-black sand — so the
// reservoir reads without relying on hue.
const (
	visSand    = "#141418" // near-black background
	visClay    = "#6e645a" // mid-gray rock
	visFlowing = "#56b4e9" // sky blue
	visSettled = "#f0e442" // bright yellow, the retained water
)

// Vis renders the final reservoir as an SVG. Clay veins are drawn as thin stroked
// lines (one per input segment) and water as greedily-merged solid rectangles, so
// there are no per-tile seams and the file stays small.
func (e Exercise) Vis(instr, outdir string) error {
	g := simulate(instr)

	// Crop horizontally to the wet+clay span to avoid a mostly-empty margin.
	minCol, maxCol := len(g.tiles[g.minY]), 0
	for y := g.minY; y <= g.maxY; y++ {
		for x, t := range g.tiles[y] {
			if t != sand {
				if x < minCol {
					minCol = x
				}
				if x > maxCol {
					maxCol = x
				}
			}
		}
	}
	if minCol > maxCol {
		minCol, maxCol = 0, len(g.tiles[g.minY])-1
	}

	w := maxCol - minCol + 1
	h := g.maxY - g.minY + 1
	ox, oy := minCol+g.minX, g.minY // absolute coords of the top-left corner

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" `+
		`shape-rendering="crispEdges">`, w, h)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="%s"/>`, w, h, visSand)

	visWater(&sb, g, minCol, maxCol, oy)
	visClayLines(&sb, instr, ox, oy)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "reservoir-research.svg"), []byte(sb.String()), 0o600)
}

// visWater emits greedily-merged maximal rectangles covering the settled and
// flowing water, so each basin is a few large seam-free fills.
func visWater(sb *strings.Builder, g *grid, minCol, maxCol, oy int) {
	h := g.maxY - g.minY + 1
	w := maxCol - minCol + 1
	// water[y'][x'] holds the tile type (or sand) in cropped coordinates.
	done := make([][]bool, h)
	for i := range done {
		done[i] = make([]bool, w)
	}
	isWater := func(t byte) bool { return t == flowing || t == settled }

	for yy := 0; yy < h; yy++ {
		row := g.tiles[oy+yy]
		for xx := 0; xx < w; xx++ {
			if done[yy][xx] {
				continue
			}
			t := row[minCol+xx]
			if !isWater(t) {
				continue
			}
			// Extend right while the same type and not yet covered.
			rw := 1
			for xx+rw < w && !done[yy][xx+rw] && row[minCol+xx+rw] == t {
				rw++
			}
			// Extend down while every cell of that width matches.
			rh := 1
			for yy+rh < h {
				ry := g.tiles[oy+yy+rh]
				ok := true
				for k := 0; k < rw; k++ {
					if done[yy+rh][xx+k] || ry[minCol+xx+k] != t {
						ok = false
						break
					}
				}
				if !ok {
					break
				}
				rh++
			}
			for dy := 0; dy < rh; dy++ {
				for dx := 0; dx < rw; dx++ {
					done[yy+dy][xx+dx] = true
				}
			}
			fill := visFlowing
			if t == settled {
				fill = visSettled
			}
			fmt.Fprintf(sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`,
				xx, yy, rw, rh, fill)
		}
	}
}

// visClayLines draws each clay vein from the input as one thin, integer-aligned
// rect — a 1-wide column for vertical veins, a 1-tall bar for horizontal ones.
// Keeping every coordinate an integer (no fractional stroke centerlines) lets
// crispEdges snap all shapes to the pixel grid, so nothing anti-aliases.
func visClayLines(sb *strings.Builder, instr string, ox, oy int) {
	re := regexp.MustCompile(`-?\d+`)
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		if line == "" {
			continue
		}
		nums := re.FindAllString(line, -1)
		a, _ := strconv.Atoi(nums[0])
		b, _ := strconv.Atoi(nums[1])
		c, _ := strconv.Atoi(nums[2])

		var x, y, w, h int
		if line[0] == 'x' { // x=a, y=b..c (vertical)
			x, y, w, h = a-ox, b-oy, 1, c-b+1
		} else { // y=a, x=b..c (horizontal)
			x, y, w, h = b-ox, a-oy, c-b+1, 1
		}
		fmt.Fprintf(sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`,
			x, y, w, h, visClay)
	}
}
