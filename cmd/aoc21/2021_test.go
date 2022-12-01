package aoc21_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/cmd/aoc21"
	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

func Test_2021(t *testing.T) {
	tests := []struct {
		name  string
		code  aoc.ExerciseFunc
		input []string
		want  string
	}{
		{"2021-1a Example", aoc21.Day1part1, []string{"199", "200", "208", "210", "200", "207", "240", "269", "260", "263"}, "7"},
		{"2021-1b Example", aoc21.Day1part2, []string{"199", "200", "208", "210", "200", "207", "240", "269", "260", "263"}, "5"},
		// {"2021-1a", Day2part1, []string{"forward 5", "down 5", "forward 8", "up 3", "down 8", "forward 2"}, 150},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code(tt.input); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
