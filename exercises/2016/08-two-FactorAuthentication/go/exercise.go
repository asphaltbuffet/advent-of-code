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

// Exercise for Advent of Code 2016 day 8.
type Exercise struct {
	common.BaseExercise
}

// screen is a lit/unlit pixel grid.
type screen struct {
	w, h int
	px   [][]bool
}

// run builds the screen, applies every instruction, and returns it. The real
// display is 50x6; the puzzle example is 7x3, detected by its small indices.
func run(instr string) *screen {
	w, h := screenSize(instr)
	s := &screen{w: w, h: h, px: make([][]bool, h)}
	for r := range s.px {
		s.px[r] = make([]bool, w)
	}
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		f := strings.Fields(line)
		switch {
		case f[0] == "rect":
			var a, b int
			parseDims(f[1], &a, &b)
			s.rect(a, b)
		case f[1] == "row":
			y := atoiAfter(f[2], "y=")
			by := atoi(f[4])
			s.rotateRow(y, by)
		case f[1] == "column":
			x := atoiAfter(f[2], "x=")
			by := atoi(f[4])
			s.rotateCol(x, by)
		}
	}
	return s
}

// screenSize returns 7x3 for the example, 50x6 for the real display, based on
// the largest indices the instructions reference.
func screenSize(instr string) (int, int) {
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		f := strings.Fields(line)
		if f[0] == "rect" {
			var a, b int
			parseDims(f[1], &a, &b)
			if a > 7 || b > 3 {
				return 50, 6
			}
		}
		if len(f) > 2 && f[1] == "row" && atoiAfter(f[2], "y=") > 2 {
			return 50, 6
		}
		if len(f) > 2 && f[1] == "column" && atoiAfter(f[2], "x=") > 6 {
			return 50, 6
		}
	}
	return 7, 3
}

func (s *screen) rect(a, b int) {
	for r := 0; r < b && r < s.h; r++ {
		for c := 0; c < a && c < s.w; c++ {
			s.px[r][c] = true
		}
	}
}

func (s *screen) rotateRow(y, by int) {
	by %= s.w
	row := make([]bool, s.w)
	for c := range s.w {
		row[(c+by)%s.w] = s.px[y][c]
	}
	s.px[y] = row
}

func (s *screen) rotateCol(x, by int) {
	by %= s.h
	col := make([]bool, s.h)
	for r := range s.h {
		col[(r+by)%s.h] = s.px[r][x]
	}
	for r := range s.h {
		s.px[r][x] = col[r]
	}
}

func (s *screen) lit() int {
	n := 0
	for _, row := range s.px {
		for _, p := range row {
			if p {
				n++
			}
		}
	}
	return n
}

// render returns the screen as text with '#' lit and '.' unlit.
func (s *screen) render() string {
	var b strings.Builder
	for r := range s.h {
		for c := range s.w {
			if s.px[r][c] {
				b.WriteByte('#')
			} else {
				b.WriteByte('.')
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return run(instr).lit(), nil
}

// Two returns the answer to the second part of the exercise: the letters spelled
// on the screen.
func (e Exercise) Two(instr string) (any, error) {
	return ocr(run(instr)), nil
}

// Vis renders the final screen, with lit pixels forming the decoded letters.
func (e Exercise) Vis(instr string, outdir string) error {
	s := run(instr)

	const cell = 18
	const pad = 16
	img := image.NewRGBA(image.Rect(0, 0, s.w*cell+2*pad, s.h*cell+2*pad))

	bg := color.RGBA{0x0a, 0x12, 0x0a, 0xff}
	off := color.RGBA{0x14, 0x22, 0x14, 0xff}
	on := color.RGBA{0x5c, 0xff, 0x7a, 0xff}
	h, w := img.Bounds().Dy(), img.Bounds().Dx()
	for y := range h {
		for x := range w {
			img.SetRGBA(x, y, bg)
		}
	}

	for r := range s.h {
		for c := range s.w {
			col := off
			if s.px[r][c] {
				col = on
			}
			x0, y0 := pad+c*cell, pad+r*cell
			// Leave a 1px gutter so pixels read as a grid.
			for y := y0; y < y0+cell-1; y++ {
				for x := x0; x < x0+cell-1; x++ {
					img.SetRGBA(x, y, col)
				}
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "two-factor-auth.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// --- small parse helpers ---

func parseDims(s string, a, b *int) {
	parts := strings.SplitN(s, "x", 2)
	*a, *b = atoi(parts[0]), atoi(parts[1])
}

func atoiAfter(s, prefix string) int {
	return atoi(strings.TrimPrefix(s, prefix))
}

func atoi(s string) int {
	n := 0
	for i := range len(s) {
		if s[i] >= '0' && s[i] <= '9' {
			n = n*10 + int(s[i]-'0')
		}
	}
	return n
}

// glyphs maps the standard AoC 5x6 letter shapes to their characters. Only the
// letters that appear in this puzzle's output are needed.
var glyphs = map[string]byte{
	".##..#..#.#....#....#..#..##..": 'C',
	"####.#....###..#....#....#....": 'F',
	"#....#....#....#....#....####.": 'L',
	"####.#....###..#....#....####.": 'E',
	".##..#..#.#..#.#..#.#..#..##..": 'O',
	"#...##...#.#.#...#....#....#..": 'Y',
	".###.#....#.....##.....#.###..": 'S',
}

// ocr decodes the 5x6 letter glyphs (10 letters across a 50x6 screen).
func ocr(s *screen) string {
	if s.w%5 != 0 {
		return s.render()
	}
	var msg strings.Builder
	for k := 0; k*5 < s.w; k++ {
		var g strings.Builder
		for r := range s.h {
			for c := k * 5; c < k*5+5; c++ {
				if s.px[r][c] {
					g.WriteByte('#')
				} else {
					g.WriteByte('.')
				}
			}
		}
		if ch, ok := glyphs[g.String()]; ok {
			msg.WriteByte(ch)
		} else {
			msg.WriteByte('?')
		}
	}
	return msg.String()
}
