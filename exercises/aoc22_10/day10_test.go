package aoc22_10_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_10"
)

var (
	singleAddTestData  = []string{"addx 5"}
	singleNoopTestData = []string{"noop"}
)

func Test_Day10Part1(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 1 Example", []string{
			"addx 15",
			"addx -11",
			"addx 6",
			"addx -3",
			"addx 5",
			"addx -1",
			"addx -8",
			"addx 13",
			"addx 4",
			"noop",
			"addx -1",
			"addx 5",
			"addx -1",
			"addx 5",
			"addx -1",
			"addx 5",
			"addx -1",
			"addx 5",
			"addx -1",
			"addx -35",
			"addx 1",
			"addx 24",
			"addx -19",
			"addx 1",
			"addx 16",
			"addx -11",
			"noop",
			"noop",
			"addx 21",
			"addx -15",
			"noop",
			"noop",
			"addx -3",
			"addx 9",
			"addx 1",
			"addx -3",
			"addx 8",
			"addx 1",
			"addx 5",
			"noop",
			"noop",
			"noop",
			"noop",
			"noop",
			"addx -36",
			"noop",
			"addx 1",
			"addx 7",
			"noop",
			"noop",
			"noop",
			"addx 2",
			"addx 6",
			"noop",
			"noop",
			"noop",
			"noop",
			"noop",
			"addx 1",
			"noop",
			"noop",
			"addx 7",
			"addx 1",
			"noop",
			"addx -13",
			"addx 13",
			"addx 7",
			"noop",
			"addx 1",
			"addx -33",
			"noop",
			"noop",
			"noop",
			"addx 2",
			"noop",
			"noop",
			"noop",
			"addx 8",
			"noop",
			"addx -1",
			"addx 2",
			"addx 1",
			"noop",
			"addx 17",
			"addx -9",
			"addx 1",
			"addx 1",
			"addx -3",
			"addx 11",
			"noop",
			"noop",
			"addx 1",
			"noop",
			"addx 1",
			"noop",
			"noop",
			"addx -13",
			"addx -19",
			"addx 1",
			"addx 3",
			"addx 26",
			"addx -30",
			"addx 12",
			"addx -1",
			"addx 3",
			"addx 1",
			"noop",
			"noop",
			"noop",
			"addx -9",
			"addx 18",
			"addx 1",
			"addx 2",
			"noop",
			"noop",
			"addx 9",
			"noop",
			"noop",
			"noop",
			"addx -1",
			"addx 2",
			"addx -37",
			"addx 1",
			"addx 3",
			"noop",
			"addx 15",
			"addx -21",
			"addx 22",
			"addx -6",
			"addx 1",
			"noop",
			"addx 2",
			"addx 1",
			"noop",
			"addx -10",
			"noop",
			"noop",
			"addx 20",
			"addx 1",
			"addx 2",
			"addx 2",
			"addx -6",
			"addx -11",
			"noop",
			"noop",
			"noop",
		}, "13140"},
		{"Part 1 Mini", []string{
			"noop",
			"addx 3",
			"addx -5",
		}, "-720"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_10.D10P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day10Part2(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 2 Example", []string{
			"addx 15",
			"addx -11",
			"addx 6",
			"addx -3",
			"addx 5",
			"addx -1",
			"addx -8",
			"addx 13",
			"addx 4",
			"noop",
			"addx -1",
			"addx 5",
			"addx -1",
			"addx 5",
			"addx -1",
			"addx 5",
			"addx -1",
			"addx 5",
			"addx -1",
			"addx -35",
			"addx 1",
			"addx 24",
			"addx -19",
			"addx 1",
			"addx 16",
			"addx -11",
			"noop",
			"noop",
			"addx 21",
			"addx -15",
			"noop",
			"noop",
			"addx -3",
			"addx 9",
			"addx 1",
			"addx -3",
			"addx 8",
			"addx 1",
			"addx 5",
			"noop",
			"noop",
			"noop",
			"noop",
			"noop",
			"addx -36",
			"noop",
			"addx 1",
			"addx 7",
			"noop",
			"noop",
			"noop",
			"addx 2",
			"addx 6",
			"noop",
			"noop",
			"noop",
			"noop",
			"noop",
			"addx 1",
			"noop",
			"noop",
			"addx 7",
			"addx 1",
			"noop",
			"addx -13",
			"addx 13",
			"addx 7",
			"noop",
			"addx 1",
			"addx -33",
			"noop",
			"noop",
			"noop",
			"addx 2",
			"noop",
			"noop",
			"noop",
			"addx 8",
			"noop",
			"addx -1",
			"addx 2",
			"addx 1",
			"noop",
			"addx 17",
			"addx -9",
			"addx 1",
			"addx 1",
			"addx -3",
			"addx 11",
			"noop",
			"noop",
			"addx 1",
			"noop",
			"addx 1",
			"noop",
			"noop",
			"addx -13",
			"addx -19",
			"addx 1",
			"addx 3",
			"addx 26",
			"addx -30",
			"addx 12",
			"addx -1",
			"addx 3",
			"addx 1",
			"noop",
			"noop",
			"noop",
			"addx -9",
			"addx 18",
			"addx 1",
			"addx 2",
			"noop",
			"noop",
			"addx 9",
			"noop",
			"noop",
			"noop",
			"addx -1",
			"addx 2",
			"addx -37",
			"addx 1",
			"addx 3",
			"noop",
			"addx 15",
			"addx -21",
			"addx 22",
			"addx -6",
			"addx 1",
			"noop",
			"addx 2",
			"addx 1",
			"noop",
			"addx -10",
			"noop",
			"noop",
			"addx 20",
			"addx 1",
			"addx 2",
			"addx 2",
			"addx -6",
			"addx -11",
			"noop",
			"noop",
			"noop",
		}, "##..##..##..##..##..##..##..##..##..##..\n###...###...###...###...###...###...###.\n####....####....####....####....####....\n#####.....#####.....#####.....#####.....\n######......######......######......####\n#######.......#######.......#######.....\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_10.D10P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Parse(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		want      aoc22_10.Day10
		assertion assert.ErrorAssertionFunc
	}{
		{
			"Single Add",
			singleAddTestData,
			aoc22_10.Day10{Cycle: 0, X: map[int]int{-1: 0}, Commands: []aoc22_10.Command{{Instruction: aoc22_10.Add, Value: 5}}},
			assert.NoError,
		},
		{
			"Single Noop",
			singleNoopTestData,
			aoc22_10.Day10{Cycle: 0, X: map[int]int{-1: 0}, Commands: []aoc22_10.Command{{Instruction: aoc22_10.Noop, Value: 0}}},
			assert.NoError,
		},
		{
			"Mixed",
			[]string{"noop", "addx 5", "addx -3", "noop", "addx 1"},
			aoc22_10.Day10{Cycle: 0, X: map[int]int{-1: 0}, Commands: []aoc22_10.Command{
				{Instruction: aoc22_10.Noop, Value: 0},
				{Instruction: aoc22_10.Add, Value: 5},
				{Instruction: aoc22_10.Add, Value: -3},
				{Instruction: aoc22_10.Noop, Value: 0},
				{Instruction: aoc22_10.Add, Value: 1},
			}},
			assert.NoError,
		},
		{
			"Unknown Instruction",
			[]string{"foo"},
			aoc22_10.Day10{},
			assert.Error,
		},
		{
			"Unknown Instruction with value",
			[]string{"foo 5"},
			aoc22_10.Day10{},
			assert.Error,
		},
		{
			"Add with no value",
			[]string{"addx"},
			aoc22_10.Day10{},
			assert.Error,
		},
		{
			"Add with non-numeric value",
			[]string{"addx one"},
			aoc22_10.Day10{},
			assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_10.Day10{
				Cycle:    0,
				X:        map[int]int{0: 0},
				Commands: []aoc22_10.Command{},
			}

			err := got.Parse(tt.input)
			tt.assertion(t, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

var errResult error

func BenchmarkParse(b *testing.B) {
	var err error

	for i := 0; i < b.N; i++ {
		d := aoc22_10.Day10{
			Cycle:    0,
			X:        map[int]int{0: 0},
			Commands: []aoc22_10.Command{},
		}

		// record the result to prevent the compiler eliminating the function call.
		err = d.Parse(singleAddTestData)
	}

	// store result to a package level variable to prevent the compiler eliminating the benchmark itself.
	errResult = err
}

func Test_Process(t *testing.T) {
	tests := []struct {
		name      string
		input     aoc22_10.Day10
		want      aoc22_10.Day10
		assertion assert.ErrorAssertionFunc
	}{
		{
			"Single Add",
			aoc22_10.Day10{Cycle: 0, X: map[int]int{-1: 0}, Commands: []aoc22_10.Command{{Instruction: aoc22_10.Add, Value: 5}}},
			aoc22_10.Day10{Cycle: 2, X: map[int]int{-1: 0, 1: 1, 2: 1, 3: 5}, Commands: []aoc22_10.Command{{Instruction: aoc22_10.Add, Value: 5}}},
			assert.NoError,
		},
		{
			"Single Noop",
			aoc22_10.Day10{Cycle: 0, X: map[int]int{-1: 0}, Commands: []aoc22_10.Command{{Instruction: aoc22_10.Noop, Value: 0}}},
			aoc22_10.Day10{Cycle: 1, X: map[int]int{-1: 0, 1: 0}, Commands: []aoc22_10.Command{{Instruction: aoc22_10.Noop, Value: 0}}},
			assert.NoError,
		},
		{
			"Mixed",
			aoc22_10.Day10{Cycle: 0, X: map[int]int{-1: 0}, Commands: []aoc22_10.Command{
				{Instruction: aoc22_10.Noop, Value: 0},
				{Instruction: aoc22_10.Add, Value: 3},
				{Instruction: aoc22_10.Add, Value: -5},
			}},
			aoc22_10.Day10{Cycle: 0, X: map[int]int{-1: 0, 1: 1, 2: 1, 3: 1, 4: 4, 5: 4, 6: -1}, Commands: []aoc22_10.Command{
				{Instruction: aoc22_10.Noop, Value: 0},
				{Instruction: aoc22_10.Add, Value: 3},
				{Instruction: aoc22_10.Add, Value: -5},
			}},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_10.Day10{
				Cycle:    0,
				X:        map[int]int{0: 0},
				Commands: []aoc22_10.Command{},
			}

			err := got.Process()
			tt.assertion(t, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
