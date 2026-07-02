package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 20.
type Exercise struct {
	common.BaseExercise
}

type enhImage struct {
	pixels [][]bool // known window
	bg     bool     // value of the infinite background outside the window
}

func parse(instr string) (string, enhImage, error) {
	algoPart, imgPart, ok := strings.Cut(strings.TrimSpace(instr), "\n\n")
	if !ok {
		return "", enhImage{}, fmt.Errorf("input missing blank-line separator")
	}
	algo := strings.ReplaceAll(strings.TrimSpace(algoPart), "\n", "")
	if len(algo) != 512 {
		return "", enhImage{}, fmt.Errorf("algorithm must be 512 chars, got %d", len(algo))
	}

	var pixels [][]bool
	for _, line := range strings.Split(strings.TrimSpace(imgPart), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		row := make([]bool, len(line))
		for i, c := range line {
			row[i] = c == '#'
		}
		pixels = append(pixels, row)
	}

	return algo, enhImage{pixels: pixels, bg: false}, nil
}

// at returns the pixel at (r,c), falling back to the infinite background when
// the coordinate lies outside the known window.
func (im enhImage) at(r, c int) bool {
	if r < 0 || r >= len(im.pixels) || c < 0 || c >= len(im.pixels[0]) {
		return im.bg
	}
	return im.pixels[r][c]
}

// enhance applies the algorithm once, growing the window by one cell on every
// side and toggling the background if algo[0] flips empty space on.
func enhance(algo string, im enhImage) enhImage {
	rows := len(im.pixels) + 2
	cols := len(im.pixels[0]) + 2
	out := make([][]bool, rows)

	for r := 0; r < rows; r++ {
		out[r] = make([]bool, cols)
		for c := 0; c < cols; c++ {
			// Map output (r,c) back to input coords (offset by -1).
			idx := 0
			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					idx <<= 1
					if im.at(r-1+dr, c-1+dc) {
						idx |= 1
					}
				}
			}
			out[r][c] = algo[idx] == '#'
		}
	}

	nbg := im.bg
	if im.bg {
		nbg = algo[511] == '#' // all-lit neighborhood
	} else {
		nbg = algo[0] == '#' // all-dark neighborhood
	}

	return enhImage{pixels: out, bg: nbg}
}

func litCount(algo string, im enhImage, steps int) (int, error) {
	for i := 0; i < steps; i++ {
		im = enhance(algo, im)
	}

	// After a finite number of steps the background is finite too; if it is lit
	// the count would be infinite, but the puzzle inputs are chosen so it is dark
	// at the sampled steps.
	if im.bg {
		return 0, fmt.Errorf("infinite background is lit after %d steps", steps)
	}

	count := 0
	for _, row := range im.pixels {
		for _, p := range row {
			if p {
				count++
			}
		}
	}
	return count, nil
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	algo, im, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	n, err := litCount(algo, im, 2)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%d", n), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	algo, im, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	n, err := litCount(algo, im, 50)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%d", n), nil
}
