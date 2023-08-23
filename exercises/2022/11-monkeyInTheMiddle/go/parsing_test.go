//go:build test
// +build test

package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Day11
	}{
		{
			"Single Monkey",
			"Monkey 0:\n" +
				"  Starting items: 79, 98\n" +
				"  Operation: new = old * 19\n" +
				"  Test: divisible by 23\n" +
				"    If true: throw to monkey 2\n" +
				"    If false: throw to monkey 3",
			Day11{
				Monkeys: []*Monkey{
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
			},
		},
		{
			"Example",
			"Monkey 0:\n" +
				"  Starting items: 79, 98\n" +
				"  Operation: new = old * 19\n" +
				"  Test: divisible by 23\n" +
				"    If true: throw to monkey 2\n" +
				"    If false: throw to monkey 3\n\n" +
				"Monkey 1:\n" +
				"  Starting items: 54, 65, 75, 74\n" +
				"  Operation: new = old + 6\n" +
				"  Test: divisible by 19\n" +
				"    If true: throw to monkey 2\n" +
				"    If false: throw to monkey 0\n\n" +
				"Monkey 2:\n" +
				"  Starting items: 79, 60, 97\n" +
				"  Operation: new = old * old\n" +
				"  Test: divisible by 13\n" +
				"    If true: throw to monkey 1\n" +
				"    If false: throw to monkey 3\n\n" +
				"Monkey 3:\n" +
				"  Starting items: 74\n" +
				"  Operation: new = old + 3\n" +
				"  Test: divisible by 17\n" +
				"    If true: throw to monkey 0\n" +
				"    If false: throw to monkey 1",
			Day11{
				Monkeys: []*Monkey{
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Day11{Product: 1}

			err := d.ParseInput(tt.input)
			assert.NoError(t, err)

			assert.Equal(t, tt.want, d)
		})
	}
}
