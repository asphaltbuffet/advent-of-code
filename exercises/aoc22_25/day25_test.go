package aoc22_25_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_25"
)

func Test_Day25Part1(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 1 Example", []string{
			"input",
			"strings",
		}, "answer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_25.D25P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day25Part2(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 2 Example", []string{
			"input",
			"strings",
		}, "answer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_25.D25P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
