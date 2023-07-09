//go:build test
// +build test

package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SumSizes(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]int
		want  int
	}{
		{"Empty", map[string]int{}, 0},
		{"Single", map[string]int{"a": 100}, 100},
		{"Single + Large", map[string]int{"a": 200000}, 0},
		{"Multiple", map[string]int{"a": 100, "b": 200, "c": 300}, 600},
		{"Multiple + Large", map[string]int{"a": 100, "b": 200, "c": 300000}, 300},

		// a: 100
		// ↪ b: 200000
		// c: 300
		{"Nested", map[string]int{"a": 200100, "a/b": 200000, "c": 300}, 300},

		// a: 100000
		// ↪ b: 200
		// c: 300
		{"Nested - Parent too large", map[string]int{"a": 100200, "a/b": 200, "c": 300}, 500},

		// a: 100500
		// ↪ b: 200
		// ↪ c: 300
		{"Nested with root - root too large", map[string]int{"a": 100500, "a/b": 200, "a/c": 300}, 500},

		{"files in root - sum all", map[string]int{"root": 95600, "root/a": 93600, "root/a/e": 600}, 189800},
		{"example - part 1", map[string]int{"root": 48381165, "root/a": 94853, "root/a/e": 584, "root/d": 24933642}, 95437},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sumSizes(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

// // Property 1: Sum of all sizes is equal to the sum of all files
// func Test_SumSizesProperty1(t *testing.T) {
// 	properties := gopter.NewProperties(nil)

// 	properties.Property("Sum of all sizes is equal to the sum of all files", prop.ForAll(

// }

// Property based testing
// Property 1: Sum of all sizes is less than or equal to the sum of all files
// Property 2: Key values in the map do not affect the sum of all sizes
// Property 3: Value order does not change the sum
