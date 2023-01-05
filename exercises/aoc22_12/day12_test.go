package aoc22_12_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_12"
)

func Test_Day12Part1(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 1 Example", []string{
			"Sabqponm",
			"abcryxxl",
			"accszExk",
			"acctuvwj",
			"abdefghi",
		}, "31"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_12.D12P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day12Part2(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 2 Example", []string{
			"Sabqponm",
			"abcryxxl",
			"accszExk",
			"acctuvwj",
			"abdefghi",
		}, "29"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_12.D12P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
