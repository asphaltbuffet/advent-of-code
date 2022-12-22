package aoc22_14_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_14"
)

func Test_Day14Part1(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 1 Example", []string{
			"498,4 -> 498,6 -> 496,6",
			"503,4 -> 502,4 -> 502,9 -> 494,9",
		}, "24"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_14.D14P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day14Part2(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 2 Example", []string{
			"498,4 -> 498,6 -> 496,6",
			"503,4 -> 502,4 -> 502,9 -> 494,9",
		}, "93"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_14.D14P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
