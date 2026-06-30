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

// Exercise for Advent of Code 2024 day 10.
type Exercise struct {
	common.BaseExercise
}

type cell struct{ r, c int }

var step4 = [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

// topo holds the parsed height map. Cells off-grid or non-digit read as -1.
type topo struct {
	h          [][]int
	rows, cols int
}

func parseTopo(instr string) topo {
	lines := strings.Fields(instr)
	t := topo{rows: len(lines)}
	for _, line := range lines {
		row := make([]int, len(line))
		for c := 0; c < len(line); c++ {
			if line[c] >= '0' && line[c] <= '9' {
				row[c] = int(line[c] - '0')
			} else {
				row[c] = -1
			}
		}
		if len(row) > t.cols {
			t.cols = len(row)
		}
		t.h = append(t.h, row)
	}
	return t
}

func (t topo) at(r, c int) int {
	if r < 0 || r >= t.rows || c < 0 || c >= len(t.h[r]) {
		return -1
	}
	return t.h[r][c]
}

// reachable9s adds every height-9 cell reachable from (r,c) on an increasing
// trail to the set ends.
func (t topo) reachable9s(r, c int, ends map[cell]bool) {
	if t.at(r, c) == 9 {
		ends[cell{r, c}] = true
		return
	}
	for _, d := range step4 {
		nr, nc := r+d[0], c+d[1]
		if t.at(nr, nc) == t.at(r, c)+1 {
			t.reachable9s(nr, nc, ends)
		}
	}
}

// ratings returns the number of distinct trails from (r,c) to any 9, memoized.
func (t topo) ratings(r, c int, memo map[cell]int) int {
	if t.at(r, c) == 9 {
		return 1
	}
	if v, ok := memo[cell{r, c}]; ok {
		return v
	}
	total := 0
	for _, d := range step4 {
		nr, nc := r+d[0], c+d[1]
		if t.at(nr, nc) == t.at(r, c)+1 {
			total += t.ratings(nr, nc, memo)
		}
	}
	memo[cell{r, c}] = total
	return total
}

func (t topo) trailheads() []cell {
	var heads []cell
	for r := 0; r < t.rows; r++ {
		for c := 0; c < len(t.h[r]); c++ {
			if t.h[r][c] == 0 {
				heads = append(heads, cell{r, c})
			}
		}
	}
	return heads
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	t := parseTopo(instr)
	score := 0
	for _, h := range t.trailheads() {
		ends := make(map[cell]bool)
		t.reachable9s(h.r, h.c, ends)
		score += len(ends)
	}
	return score, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	t := parseTopo(instr)
	memo := make(map[cell]int)
	rating := 0
	for _, h := range t.trailheads() {
		rating += t.ratings(h.r, h.c, memo)
	}
	return rating, nil
}

// --- Visualization ---

const (
	visCell    = 18
	visPadding = 24
)

// elevationStops anchor a topographic ramp from lowland to peak. heightColor
// interpolates between the two stops bracketing a normalized height.
var elevationStops = []color.RGBA{
	{0x10, 0x2a, 0x55, 0xff}, // 0.00 deep blue
	{0x1f, 0x77, 0x8f, 0xff}, // 0.25 teal
	{0x3f, 0x9d, 0x57, 0xff}, // 0.50 green
	{0xc8, 0xb0, 0x6b, 0xff}, // 0.75 tan
	{0xf4, 0xf0, 0xe6, 0xff}, // 1.00 near-white peak
}

// heightColor maps a height 0..9 to the topographic ramp.
func heightColor(h int) color.RGBA {
	if h < 0 {
		return color.RGBA{0x0a, 0x0a, 0x14, 0xff}
	}
	t := float64(h) / 9.0
	seg := t * float64(len(elevationStops)-1)
	i := int(seg)
	if i >= len(elevationStops)-1 {
		return elevationStops[len(elevationStops)-1]
	}
	f := seg - float64(i)
	lerp := func(a, b uint8) uint8 { return uint8(float64(a) + (float64(b)-float64(a))*f) }
	a, b := elevationStops[i], elevationStops[i+1]
	return color.RGBA{lerp(a.R, b.R), lerp(a.G, b.G), lerp(a.B, b.B), 0xff}
}

var (
	visTrail = color.RGBA{R: 0xff, G: 0x55, B: 0x33, A: 0xff} // trail lines
	visHead  = color.RGBA{R: 0x33, G: 0xff, B: 0xcc, A: 0xff} // trailhead marker
	visPeak  = color.RGBA{R: 0xff, G: 0x2e, B: 0x9e, A: 0xff} // height-9 marker (magenta)
)

const visTrailAlpha = 0.30

// enumTrails walks every distinct trail from (r,c) and calls emit with each
// full path (slice of cells from trailhead to a 9).
func (t topo) enumTrails(r, c int, path []cell, emit func([]cell)) {
	path = append(path, cell{r, c})
	if t.at(r, c) == 9 {
		emit(path)
		return
	}
	for _, d := range step4 {
		nr, nc := r+d[0], c+d[1]
		if t.at(nr, nc) == t.at(r, c)+1 {
			t.enumTrails(nr, nc, path, emit)
		}
	}
}

// Vis renders the topographic map with all hiking trails overlaid.
func (e Exercise) Vis(instr string, outdir string) error {
	t := parseTopo(instr)
	w := t.cols*visCell + 2*visPadding
	h := t.rows*visCell + 2*visPadding
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Background.
	bg := color.RGBA{0x06, 0x06, 0x10, 0xff}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	// Height cells as filled squares.
	for r := 0; r < t.rows; r++ {
		for c := 0; c < len(t.h[r]); c++ {
			fillCell(img, r, c, heightColor(t.h[r][c]))
		}
	}

	center := func(r, c int) (float64, float64) {
		return float64(visPadding) + (float64(c)+0.5)*visCell,
			float64(visPadding) + (float64(r)+0.5)*visCell
	}

	// Overlay every trail as a translucent line; overlaps glow.
	for _, head := range t.trailheads() {
		t.enumTrails(head.r, head.c, nil, func(path []cell) {
			for i := 0; i+1 < len(path); i++ {
				x0, y0 := center(path[i].r, path[i].c)
				x1, y1 := center(path[i+1].r, path[i+1].c)
				drawLine(img, x0, y0, x1, y1, visTrail, 1)
			}
		})
	}

	// Mark trailheads (cyan) and peaks (white).
	for r := 0; r < t.rows; r++ {
		for c := 0; c < len(t.h[r]); c++ {
			switch t.h[r][c] {
			case 0:
				markDot(img, r, c, visHead)
			case 9:
				markDot(img, r, c, visPeak)
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "hoof-it.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func fillCell(img *image.RGBA, r, c int, col color.RGBA) {
	x0 := visPadding + c*visCell
	y0 := visPadding + r*visCell
	for y := y0; y < y0+visCell; y++ {
		for x := x0; x < x0+visCell; x++ {
			img.SetRGBA(x, y, col)
		}
	}
}

func markDot(img *image.RGBA, r, c int, col color.RGBA) {
	cx := visPadding + c*visCell + visCell/2
	cy := visPadding + r*visCell + visCell/2
	const rad = 3
	for y := -rad; y <= rad; y++ {
		for x := -rad; x <= rad; x++ {
			if x*x+y*y <= rad*rad {
				img.SetRGBA(cx+x, cy+y, col)
			}
		}
	}
}

// drawLine alpha-blends a thick line; each covered pixel blended once per line.
func drawLine(img *image.RGBA, x0, y0, x1, y1 float64, col color.RGBA, width int) {
	dx, dy := x1-x0, y1-y0
	steps := int(math.Max(math.Abs(dx), math.Abs(dy))) * 2
	if steps == 0 {
		steps = 1
	}
	imgW := img.Bounds().Dx()
	covered := make(map[int]struct{})
	for i := 0; i <= steps; i++ {
		tt := float64(i) / float64(steps)
		px := int(math.Round(x0 + dx*tt))
		py := int(math.Round(y0 + dy*tt))
		for oy := -width; oy <= width; oy++ {
			for ox := -width; ox <= width; ox++ {
				x, y := px+ox, py+oy
				key := y*imgW + x
				if _, seen := covered[key]; seen {
					continue
				}
				covered[key] = struct{}{}
				blendPixel(img, x, y, col, visTrailAlpha)
			}
		}
	}
}

func blendPixel(img *image.RGBA, x, y int, col color.RGBA, a float64) {
	if !image.Pt(x, y).In(img.Bounds()) {
		return
	}
	dst := img.RGBAAt(x, y)
	blend := func(s, d uint8) uint8 { return uint8(float64(s)*a + float64(d)*(1-a)) }
	img.SetRGBA(x, y, color.RGBA{blend(col.R, dst.R), blend(col.G, dst.G), blend(col.B, dst.B), 0xff})
}
