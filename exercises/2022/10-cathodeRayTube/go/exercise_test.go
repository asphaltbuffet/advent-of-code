//go:build test
// +build test

package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	singleAddTestData  = []string{"addx 5"}
	singleNoopTestData = []string{"noop"}
)

func Test_Parse(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		want      Day10
		assertion assert.ErrorAssertionFunc
	}{
		{
			"Single Add",
			singleAddTestData,
			Day10{Cycle: 0, X: map[int]int{0: 0, 1: 1}, Commands: []Command{{Instruction: Add, Value: 5}}},
			assert.NoError,
		},
		{
			"Single Noop",
			singleNoopTestData,
			Day10{Cycle: 0, X: map[int]int{0: 0, 1: 1}, Commands: []Command{{Instruction: Noop, Value: 0}}},
			assert.NoError,
		},
		{
			"Mixed",
			[]string{"noop", "addx 5", "addx -3", "noop", "addx 1"},
			Day10{Cycle: 0, X: map[int]int{0: 0, 1: 1}, Commands: []Command{
				{Instruction: Noop, Value: 0},
				{Instruction: Add, Value: 5},
				{Instruction: Add, Value: -3},
				{Instruction: Noop, Value: 0},
				{Instruction: Add, Value: 1},
			}},
			assert.NoError,
		},
		{
			"Unknown Instruction",
			[]string{"foo"},
			Day10{},
			assert.Error,
		},
		{
			"Unknown Instruction with value",
			[]string{"foo 5"},
			Day10{},
			assert.Error,
		},
		{
			"Add with no value",
			[]string{"addx"},
			Day10{},
			assert.Error,
		},
		{
			"Add with non-numeric value",
			[]string{"addx one"},
			Day10{},
			assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Day10{
				Cycle:    0,
				X:        map[int]int{0: 0, 1: 1},
				Commands: []Command{},
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
		d := Day10{
			Cycle:    0,
			X:        map[int]int{0: 0, 1: 1},
			Commands: []Command{},
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
		input     Day10
		want      Day10
		assertion assert.ErrorAssertionFunc
	}{
		{
			"Single Add",
			Day10{Cycle: 0, X: map[int]int{0: 0, 1: 1}, Commands: []Command{{Instruction: Add, Value: 5}}},
			Day10{Cycle: 0, X: map[int]int{0: 0, 1: 1}, Commands: []Command{}, Cycle20: 0, Cycle60: 0, Cycle100: 0, Cycle140: 0, Cycle180: 0, Cycle220: 0},
			assert.NoError,
		},
		{
			"Single Noop",
			Day10{Cycle: 0, X: map[int]int{0: 0, 1: 1}, Commands: []Command{{Instruction: Noop, Value: 0}}},
			Day10{Cycle: 0, X: map[int]int{0: 0, 1: 1}, Commands: []Command{}, Cycle20: 0, Cycle60: 0, Cycle100: 0, Cycle140: 0, Cycle180: 0, Cycle220: 0},
			assert.NoError,
		},
		{
			"Mixed",
			Day10{Cycle: 0, X: map[int]int{0: 0, 1: 1}, Commands: []Command{
				{Instruction: Noop, Value: 0},
				{Instruction: Add, Value: 3},
				{Instruction: Add, Value: -5},
			}},
			Day10{Cycle: 0, X: map[int]int{0: 0, 1: 1}, Commands: []Command{}, Cycle20: 0, Cycle60: 0, Cycle100: 0, Cycle140: 0, Cycle180: 0, Cycle220: 0},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Day10{
				Cycle:    0,
				X:        map[int]int{0: 0, 1: 1},
				Commands: []Command{},
			}

			err := got.Process()
			tt.assertion(t, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
