package aoc22_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/cmd/aoc22"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func Test_Day6(t *testing.T) {
	tests := []struct {
		name  string
		code  common.ExerciseFunc
		input []string
		want  string
	}{
		{"Part 1 Example 1", aoc22.D6P1, []string{"bvwbjplbgvbhsrlpgdmjqwftvncz"}, "5"},
		{"Part 1 Example 1", aoc22.D6P1, []string{"nppdvjthqldpwncqszvftbrmjlhg"}, "6"},
		{"Part 1 Example 1", aoc22.D6P1, []string{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg"}, "10"},
		{"Part 1 Example 1", aoc22.D6P1, []string{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw"}, "11"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.code(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
