package aoc22_24_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_24"
)

func Test_Day24Part1(t *testing.T) {
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
			got := aoc22_24.D24P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day24Part2(t *testing.T) {
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
			got := aoc22_24.D24P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
