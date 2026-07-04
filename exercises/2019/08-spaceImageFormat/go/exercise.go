package exercises

import (
	"bytes"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 8.
type Exercise struct {
	common.BaseExercise
}

const (
	width  = 25
	height = 6
	pixels = width * height
)

// One returns the answer to the first part of the exercise.
// Find the layer with fewest '0' digits; return count('1') * count('2') on that layer.
func (e Exercise) One(instr string) (any, error) {
	data := bytes.TrimSpace([]byte(instr))

	minZeros := pixels + 1
	result := 0

	for i := 0; i+pixels <= len(data); i += pixels {
		layer := data[i : i+pixels]
		zeros, ones, twos := 0, 0, 0
		for _, b := range layer {
			switch b {
			case '0':
				zeros++
			case '1':
				ones++
			case '2':
				twos++
			}
		}
		if zeros < minZeros {
			minZeros = zeros
			result = ones * twos
		}
	}

	return result, nil
}

// Two returns the answer to the second part of the exercise.
// Composite all layers front-to-back: first non-transparent pixel wins.
// '0'=black (░), '1'=white (█), '2'=transparent.
func (e Exercise) Two(instr string) (any, error) {
	data := bytes.TrimSpace([]byte(instr))

	// composite: start with all transparent
	image := make([]byte, pixels)
	for i := range image {
		image[i] = '2'
	}

	for i := 0; i+pixels <= len(data); i += pixels {
		layer := data[i : i+pixels]
		for j, b := range layer {
			if image[j] == '2' {
				image[j] = b
			}
		}
	}

	// render
	var sb strings.Builder
	for row := range height {
		if row > 0 {
			sb.WriteByte('\n')
		}
		for col := range width {
			if image[row*width+col] == '1' {
				sb.WriteString("█")
			} else {
				sb.WriteString("░")
			}
		}
	}

	return sb.String(), nil
}
