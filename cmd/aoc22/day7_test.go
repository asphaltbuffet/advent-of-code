package aoc22_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/cmd/aoc22"
)

func Test_Day7Part1(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		// root: 48381165
		// ↪ a: 94853
		//   ↪ e: 584
		// ↪ d: 24933642
		{"Part 1 Example", []string{
			"$ cd /",
			"$ ls",
			"dir a",
			"14848514 b.txt",
			"8504156 c.dat",
			"dir d",
			"$ cd a",
			"$ ls",
			"dir e",
			"29116 f",
			"2557 g",
			"62596 h.lst",
			"$ cd e",
			"$ ls",
			"584 i",
			"$ cd ..",
			"$ cd ..",
			"$ cd d",
			"$ ls",
			"4060174 j",
			"8033020 d.log",
			"5626152 d.ext",
			"7214296 k",
		}, "95437"},

		// total: 95600 + 93600 + 600 -> 189800
		// root: 2000 + a(93600) ->95600
		// ↪ a: 93000 + e(600) -> 93600
		//   ↪ e: 600
		{"Small files in root", []string{
			"$ cd /", // root is 98600: 2000 + (a: 93600) + (e: 600))
			"$ ls",
			"dir a",
			"500 b.txt",
			"1500 c.dat",
			"$ cd a", // a is 93600: 93000 + (e: 600)
			"$ ls",
			"dir e",
			"30000 f",
			"3000 g",
			"60000 h.lst",
			"$ cd e", // e is 600
			"$ ls",
			"600 i",
		}, "189800"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22.D7P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day7Part2(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 1 Example", []string{
			"$ cd /",
			"$ ls",
			"dir a",
			"14848514 b.txt",
			"8504156 c.dat",
			"dir d",
			"$ cd a",
			"$ ls",
			"dir e",
			"29116 f",
			"2557 g",
			"62596 h.lst",
			"$ cd e",
			"$ ls",
			"584 i",
			"$ cd ..",
			"$ cd ..",
			"$ cd d",
			"$ ls",
			"4060174 j",
			"8033020 d.log",
			"5626152 d.ext",
			"7214296 k",
		}, "24933642"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22.D7P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

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
			got := aoc22.SumSizes(tt.input)
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
