package aoc22_21_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_21"
)

func Test_Day21Part1(t *testing.T) {
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
			got := aoc22_21.D21P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day21Part2(t *testing.T) {
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
			got := aoc22_21.D21P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
