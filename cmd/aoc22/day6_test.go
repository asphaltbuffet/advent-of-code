package aoc22_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/cmd/aoc22"
)

func Test_Day6Part1(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Only a marker", []string{"bvwc"}, "4"},
		{"No marker", []string{"babababababababababababababaab"}, "processing datastream: not found"},

		{"Part 1 Example 1", []string{"bvwbjplbgvbhsrlpgdmjqwftvncz"}, "5"},
		{"Part 1 Example 2", []string{"nppdvjthqldpwncqszvftbrmjlhg"}, "6"},
		{"Part 1 Example 3", []string{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg"}, "10"},
		{"Part 1 Example 4", []string{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw"}, "11"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22.D6P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day6Part2(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Only a marker", []string{"abcdefghijklmn"}, "14"},
		{"No marker", []string{"babababababababababababababaab"}, "processing datastream: not found"},

		{"Part 2 Example 1", []string{"mjqjpqmgbljsphdztnvjfqwrcgsmlb"}, "19"},
		{"Part 2 Example 2", []string{"bvwbjplbgvbhsrlpgdmjqwftvncz"}, "23"},
		{"Part 2 Example 3", []string{"nppdvjthqldpwncqszvftbrmjlhg"}, "23"},
		{"Part 2 Example 4", []string{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg"}, "29"},
		{"Part 2 Example 5", []string{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw"}, "26"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22.D6P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
