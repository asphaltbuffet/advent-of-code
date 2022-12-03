package aoc22_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/cmd/aoc22"
	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

func Test_Day2(t *testing.T) {
	tests := []struct {
		name  string
		code  aoc.ExerciseFunc
		input []string
		want  string
	}{
		{"2022-2a Example", aoc22.D2P1, []string{"A Y", "B X", "C Z"}, "15"},
		{"2022-2a Example", aoc22.D2P2, []string{"A Y", "B X", "C Z"}, "12"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code(tt.input); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
