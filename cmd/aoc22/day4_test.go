package aoc22_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/cmd/aoc22"
	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

func Test_Day4(t *testing.T) {
	tests := []struct {
		name  string
		code  aoc.ExerciseFunc
		input []string
		want  string
	}{
		{"Part 1 Example", aoc22.D4P1, []string{
			"2-4,6-8",
			"2-3,4-5",
			"5-7,7-9",
			"2-8,3-7",
			"6-6,4-6",
			"2-6,4-8",
		}, "2"},
		{"Part 1 - Larger IDs", aoc22.D4P1, []string{
			"20-40,60-80",
			"2-3,4-5",
			"5-7,70-90",
			"2-8,3-7",     // overlap
			"60-60,40-60", // overlap
			"2-6,4-8",
		}, "2"},
		{"Part 2 Example", aoc22.D4P2, []string{
			"2-4,6-8",
			"2-3,4-5",
			"5-7,7-9",
			"2-8,3-7",
			"6-6,4-6",
			"2-6,4-8",
		}, "4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code(tt.input); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
