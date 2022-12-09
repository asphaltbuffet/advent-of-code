package aoc21_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/cmd/aoc21"
)

func Test_Day2Part1(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 1 Example", []string{"forward 5", "down 5", "forward 8", "up 3", "down 8", "forward 2"}, "150"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc21.D2P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day2Part2(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 2 Example", []string{"forward 5", "down 5", "forward 8", "up 3", "down 8", "forward 2"}, "900"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc21.D2P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
