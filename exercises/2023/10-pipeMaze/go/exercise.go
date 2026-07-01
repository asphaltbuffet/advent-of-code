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

// Exercise for Advent of Code 2023 day 10.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	m, start, err := parseInput(instr)
	if err != nil {
		return nil, err
	}

	_, path, err := findPath(m, start)
	if err != nil {
		return nil, err
	}

	return len(path) / 2, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	count := countInside(instr)

	return count, nil
}

// Vis renders the maze: the main loop drawn as a bright pipe, the tiles enclosed
// by it (the Part Two answer) filled in, and everything outside left dark.
func (e Exercise) Vis(instr, outdir string) error {
	m, start, err := parseInput(instr)
	if err != nil {
		return err
	}
	path, onpath, err := findPath(m, start)
	if err != nil {
		return err
	}

	lines := strings.Split(instr, "\n")
	h := len(lines)
	w := len(lines[0])

	const cell = 6
	img := image.NewRGBA(image.Rect(0, 0, w*cell, h*cell))

	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	loop := color.RGBA{0x2f, 0xd0, 0x9a, 0xff}
	inside := color.RGBA{0xff, 0x8a, 0x3a, 0xff}
	startC := color.RGBA{0xff, 0x44, 0x55, 0xff}

	fill := func(x, y int, col color.RGBA) {
		for dy := 0; dy < cell; dy++ {
			for dx := 0; dx < cell; dx++ {
				img.SetRGBA(x*cell+dx, y*cell+dy, col)
			}
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := Point{x, y}
			switch {
			case p == start:
				fill(x, y, startC)
			case onpath[p]:
				fill(x, y, loop)
			case windingNumber(p, path):
				fill(x, y, inside)
			default:
				fill(x, y, bg)
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "pipe-maze.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
