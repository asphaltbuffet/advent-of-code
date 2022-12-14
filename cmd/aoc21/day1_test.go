package aoc21_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/cmd/aoc21"
)

func Test_Day1Part1(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 1 Example", []string{"199", "200", "208", "210", "200", "207", "240", "269", "260", "263"}, "7"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc21.D1P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day1Part2(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 2 Example", []string{"199", "200", "208", "210", "200", "207", "240", "269", "260", "263"}, "5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc21.D1P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
