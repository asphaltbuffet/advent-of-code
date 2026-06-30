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

// Exercise for Advent of Code 2016 day 18.
type Exercise struct {
	common.BaseExercise
}

// countSafe generates `rows` rows from the first row and counts safe tiles. A
// tile is a trap iff its previous-row left and right neighbours differ
// (boundaries are safe).
func countSafe(first string, rows int) int {
	row := []byte(first)
	width := len(row)
	safe := 0
	trap := func(b []byte, i int) bool {
		var l, r byte = '.', '.'
		if i > 0 {
			l = b[i-1]
		}
		if i < len(b)-1 {
			r = b[i+1]
		}
		return l != r
	}

	for r := 0; r < rows; r++ {
		for _, t := range row {
			if t == '.' {
				safe++
			}
		}
		if r == rows-1 {
			break
		}
		next := make([]byte, width)
		for i := 0; i < width; i++ {
			if trap(row, i) {
				next[i] = '^'
			} else {
				next[i] = '.'
			}
		}
		row = next
	}
	return safe
}

// nextRow returns the row following row, per the trap rule.
func nextRow(row []byte) []byte {
	next := make([]byte, len(row))
	for i := range row {
		var l, r byte = '.', '.'
		if i > 0 {
			l = row[i-1]
		}
		if i < len(row)-1 {
			r = row[i+1]
		}
		if l != r {
			next[i] = '^'
		} else {
			next[i] = '.'
		}
	}
	return next
}

// Vis renders the trap map as a grid (traps dark, safe tiles bright). It draws
// a near-square window of rows to show the cellular-automaton pattern.
func (e Exercise) Vis(instr string, outdir string) error {
	first := strings.TrimSpace(instr)
	width := len(first)
	rows := width // square-ish window of the pattern

	const scale = 6
	const pad = 8
	img := image.NewRGBA(image.Rect(0, 0, width*scale+2*pad, rows*scale+2*pad))

	bg := color.RGBA{0x08, 0x0a, 0x10, 0xff}
	trapC := color.RGBA{0x1a, 0x20, 0x34, 0xff} // traps: dark
	safeC := color.RGBA{0x8a, 0xe6, 0x9a, 0xff} // safe: bright green
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	row := []byte(first)
	for r := 0; r < rows; r++ {
		for c := 0; c < width; c++ {
			col := trapC
			if row[c] == '.' {
				col = safeC
			}
			x0, y0 := pad+c*scale, pad+r*scale
			for y := y0; y < y0+scale; y++ {
				for x := x0; x < x0+scale; x++ {
					img.SetRGBA(x, y, col)
				}
			}
		}
		row = nextRow(row)
	}

	f, err := os.Create(filepath.Join(outdir, "like-a-rogue.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// rowsFor returns the row count: 10 for the example first row, else `real`.
func rowsFor(first string, real int) int {
	if first == ".^^.^.^^^^" {
		return 10
	}
	return real
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	first := strings.TrimSpace(instr)
	return countSafe(first, rowsFor(first, 40)), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	first := strings.TrimSpace(instr)
	return countSafe(first, rowsFor(first, 400000)), nil
}
