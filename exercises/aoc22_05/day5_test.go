package aoc22_05_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_05"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func Test_Day5(t *testing.T) {
	tests := []struct {
		name  string
		code  common.ExerciseFunc
		input []string
		want  string
	}{
		{"No Change - Single", aoc22_05.D5P1, []string{
			"    [D]    ",
			"[N] [C]    ",
			"[Z] [M] [P]",
			" 1   2   3 ",
			"",
			"move 1 from 2 to 1",
			"move 1 from 1 to 2",
		}, "NDP"},
		{"No Change - Bulk", aoc22_05.D5P2, []string{
			"    [D]    ",
			"[N] [C]    ",
			"[Z] [M] [P]",
			" 1   2   3 ",
			"",
			"move 2 from 2 to 1",
			"move 2 from 1 to 2",
		}, "NDP"},
		{"One Change - Single", aoc22_05.D5P1, []string{
			"    [D]    ",
			"[N] [C]    ",
			"[Z] [M] [P]",
			" 1   2   3 ",
			"",
			"move 1 from 2 to 1",
		}, "DCP"},
		{"One Change - Bulk", aoc22_05.D5P2, []string{
			"    [D]    ",
			"[N] [C]    ",
			"[Z] [M] [P]",
			" 1   2   3 ",
			"",
			"move 2 from 2 to 1",
		}, "DMP"},
		{"Part 1 Example", aoc22_05.D5P1, []string{
			//                                    Z |     Z |     Z
			"    [D]    ", //   D   | D     |     N |     N |     N
			"[N] [C]    ", // N C   | N C   |   C D | M   D |     D
			"[Z] [M] [P]", // Z M P | Z M P |   M P | C   P | C M P
			" 1   2   3 ",
			"",
			"move 1 from 2 to 1",
			"move 3 from 1 to 3",
			"move 2 from 2 to 1",
			"move 1 from 1 to 2",
		}, "CMZ"},
		{"Part 2 Example", aoc22_05.D5P2, []string{
			//                                    D |     D |     D
			"    [D]    ", //   D   | D     |     N |     N |     N
			"[N] [C]    ", // N C   | N C   |   C Z | C   Z |     Z
			"[Z] [M] [P]", // Z M P | Z M P |   M P | M   P | M C P
			" 1   2   3 ",
			"",
			"move 1 from 2 to 1",
			"move 3 from 1 to 3",
			"move 2 from 2 to 1",
			"move 1 from 1 to 2",
		}, "MCD"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.code(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GetMovementSectionLine(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  int
	}{
		{"Part 1 Example", []string{
			"    [D]    ",
			"[N] [C]    ",
			"[Z] [M] [P]",
			" 1   2   3 ",
			"",
			"move 1 from 2 to 1",
			"move 1 from 1 to 2",
		}, 5},
		{"All One", []string{
			"[Z] [M] [P]",
			" 1   2   3 ",
			"",
			"move 1 from 2 to 1",
			"move 3 from 1 to 3",
			"move 2 from 2 to 1",
			"move 1 from 1 to 2",
		}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_05.GetMovementSectionLine(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GetNumberOfStacks(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"Example", "    [D]    ", 3},
		{"Day Input", "[G]                 [D] [R]        ", 9},
		{"One", "[G]", 1},
		{"Two", "[G] [D]", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_05.GetNumberOfStacks(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_ParseStack(t *testing.T) {
	type args struct {
		input  []string
		stacks int
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"Part 1 Example",
			args{
				[]string{
					"    [D]    ",
					"[N] [C]    ",
					"[Z] [M] [P]",
					" 1   2   3 ",
					"",
					"move 1 from 2 to 1",
					"move 1 from 1 to 2",
				},
				1,
			},
			[]string{"M", "C", "D"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_05.ParseStack(tt.args.input, tt.args.stacks)
			assert.Equal(t, tt.want, got)
		})
	}
}
