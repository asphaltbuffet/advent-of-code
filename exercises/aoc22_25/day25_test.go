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
			"1=-0-2",
			"12111",
			"2=0=",
			"21",
			"2=01",
			"111",
			"20012",
			"112",
			"1=-1=",
			"1-12",
			"12",
			"1=",
			"122",
		}, "2=-1=0"},
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
		// {"Part 2 Example", []string{
		// 	"1=-0-2",
		// 	"12111",
		// 	"2=0=",
		// 	"21",
		// 	"2=01",
		// 	"111",
		// 	"20012",
		// 	"112",
		// 	"1=-1=",
		// 	"1-12",
		// 	"12",
		// 	"1=",
		// 	"122",
		// }, "answer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_25.D25P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
