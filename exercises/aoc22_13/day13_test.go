package aoc22_13_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_13"
)

func Test_Day13Part1(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 1 Example", []string{
			"[1,1,3,1,1]",
			"[1,1,5,1,1]",
			"",
			"[[1],[2,3,4]]",
			"[[1],4]",
			"",
			"[9]",
			"[[8,7,6]]",
			"",
			"[[4,4],4,4]",
			"[[4,4],4,4,4]",
			"",
			"[7,7,7,7]",
			"[7,7,7]",
			"",
			"[]",
			"[3]",
			"",
			"[[[]]]",
			"[[]]",
			"",
			"[1,[2,[3,[4,[5,6,7]]]],8,9]",
			"[1,[2,[3,[4,[5,6,0]]]],8,9]",
		}, "13"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_13.D13P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day13Part2(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 2 Example", []string{
			"[1,1,3,1,1]",
			"[1,1,5,1,1]",
			"",
			"[[1],[2,3,4]]",
			"[[1],4]",
			"",
			"[9]",
			"[[8,7,6]]",
			"",
			"[[4,4],4,4]",
			"[[4,4],4,4,4]",
			"",
			"[7,7,7,7]",
			"[7,7,7]",
			"",
			"[]",
			"[3]",
			"",
			"[[[]]]",
			"[[]]",
			"",
			"[1,[2,[3,[4,[5,6,7]]]],8,9]",
			"[1,[2,[3,[4,[5,6,0]]]],8,9]",
		}, "140"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_13.D13P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
