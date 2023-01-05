package aoc22_11_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_11"
)

func Test_Day11Part1(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 1 Example", []string{
			"Monkey 0:",
			"  Starting items: 79, 98",
			"  Operation: new = old * 19",
			"  Test: divisible by 23",
			"    If true: throw to monkey 2",
			"    If false: throw to monkey 3",
			"",
			"Monkey 1:",
			"  Starting items: 54, 65, 75, 74",
			"  Operation: new = old + 6",
			"  Test: divisible by 19",
			"    If true: throw to monkey 2",
			"    If false: throw to monkey 0",
			"",
			"Monkey 2:",
			"  Starting items: 79, 60, 97",
			"  Operation: new = old * old",
			"  Test: divisible by 13",
			"    If true: throw to monkey 1",
			"    If false: throw to monkey 3",
			"",
			"Monkey 3:",
			"  Starting items: 74",
			"  Operation: new = old + 3",
			"  Test: divisible by 17",
			"    If true: throw to monkey 0",
			"    If false: throw to monkey 1",
		}, "10605"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_11.D11P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day11Part2(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 2 Example", []string{
			"Monkey 0:",
			"  Starting items: 79, 98",
			"  Operation: new = old * 19",
			"  Test: divisible by 23",
			"    If true: throw to monkey 2",
			"    If false: throw to monkey 3",
			"",
			"Monkey 1:",
			"  Starting items: 54, 65, 75, 74",
			"  Operation: new = old + 6",
			"  Test: divisible by 19",
			"    If true: throw to monkey 2",
			"    If false: throw to monkey 0",
			"",
			"Monkey 2:",
			"  Starting items: 79, 60, 97",
			"  Operation: new = old * old",
			"  Test: divisible by 13",
			"    If true: throw to monkey 1",
			"    If false: throw to monkey 3",
			"",
			"Monkey 3:",
			"  Starting items: 74",
			"  Operation: new = old + 3",
			"  Test: divisible by 17",
			"    If true: throw to monkey 0",
			"    If false: throw to monkey 1",
		}, "2713310158"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_11.D11P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
