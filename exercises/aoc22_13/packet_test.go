package aoc22_13_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_13"
)

func TestIsOrdered_Slice(t *testing.T) {
	tests := []struct {
		name  string
		left  any
		right any
		want  bool
	}{
		{"Equal", []any{1.}, []any{1.}, true},
		{"Greater", []any{2.}, []any{1.}, false},
		{"Less", []any{1.}, []any{2.}, true},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			got := aoc22_13.IsOrdered(tt.left.([]interface{}), tt.right.([]interface{}))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParsePacket(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      interface{}
		assertion assert.ErrorAssertionFunc
	}{
		{name: "List of Integers", input: "[1,2,3]", want: []any{1., 2., 3.}, assertion: assert.NoError},
		{name: "Nested Integers", input: "[1,[2,[3]]]", want: []any{1., []any{2., []any{3.}}}, assertion: assert.NoError},
		{name: "Bookend Integers", input: "[1,[2,3], 4]", want: []any{1., []any{2., 3.}, 4.}, assertion: assert.NoError},
		{name: "Empty Nested List", input: "[[[]]]", want: []any{[]any{[]any{}}}, assertion: assert.NoError},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			got, err := aoc22_13.ParsePacket(tt.input)

			tt.assertion(t, err)

			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
