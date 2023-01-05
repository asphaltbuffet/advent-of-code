package aoc22_11_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_11"
)

func Test_ParseInput(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  aoc22_11.Day11
	}{
		{"Single Monkey", []string{
			"Monkey 0:",
			"  Starting items: 79, 98",
			"  Operation: new = old * 19",
			"  Test: divisible by 23",
			"    If true: throw to monkey 2",
			"    If false: throw to monkey 3",
		}, aoc22_11.Day11{
			Monkeys: []*aoc22_11.Monkey{
				{
					ID:        0,
					Items:     []int{79, 98},
					Divisor:   23,
					Operator:  "*",
					Scalar:    19,
					TargetOne: 2,
					TargetTwo: 3,
					Count:     0,
				},
			},
			Product: 23,
		}},
		{"Example", []string{
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
		}, aoc22_11.Day11{
			Monkeys: []*aoc22_11.Monkey{
				{
					ID:        0,
					Items:     []int{79, 98},
					Divisor:   23,
					Operator:  "*",
					Scalar:    19,
					TargetOne: 2,
					TargetTwo: 3,
					Count:     0,
				},
				{
					ID:        1,
					Items:     []int{54, 65, 75, 74},
					Divisor:   19,
					Operator:  "+",
					Scalar:    6,
					TargetOne: 2,
					TargetTwo: 0,
					Count:     0,
				},
				{
					ID:        2,
					Items:     []int{79, 60, 97},
					Divisor:   13,
					Operator:  "^",
					Scalar:    2,
					TargetOne: 1,
					TargetTwo: 3,
					Count:     0,
				},
				{
					ID:        3,
					Items:     []int{74},
					Divisor:   17,
					Operator:  "+",
					Scalar:    3,
					TargetOne: 0,
					TargetTwo: 1,
					Count:     0,
				},
			},
			Product: 96577,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := aoc22_11.Day11{Product: 1}

			err := d.ParseInput(tt.input)
			assert.NoError(t, err)

			assert.Equal(t, tt.want, d)
		})
	}
}
