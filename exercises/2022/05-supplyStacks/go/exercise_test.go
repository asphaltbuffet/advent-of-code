//go:build test
// +build test

package exercises_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetNumberOfStacks(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"Example", "    [D]    \n", 3},
		{"Day Input", "[G]                 [D] [R]        \n", 9},
		{"One", "[G]\n", 1},
		{"Two", "[G] [D]\n", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getNumberOfStacks(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_ParseStack(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Stacks
	}{
		{
			"Part 1 Example",
			"    [D]    \n" +
				"[N] [C]    \n" +
				"[Z] [M] [P]",
			Stacks{[]string{"Z", "N"}, []string{"M", "C", "D"}, []string{"P"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStacks(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
