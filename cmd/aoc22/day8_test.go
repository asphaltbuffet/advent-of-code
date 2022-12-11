package aoc22_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/cmd/aoc22"
)

func Test_Day8Part1(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 1 Example", []string{
			"30373",
			"25512",
			"65332",
			"33549",
			"35390",
		}, "21"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22.D8P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day8Part2(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 1 Example", []string{
			"30373",
			"25512",
			"65332",
			"33549",
			"35390",
		}, "8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22.D8P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
