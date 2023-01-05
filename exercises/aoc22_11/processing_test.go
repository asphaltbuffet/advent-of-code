package aoc22_11_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_11"
)

func Test_ProcessRound(t *testing.T) {
	tests := []struct {
		name  string
		input aoc22_11.Day11
		want  aoc22_11.Day11
	}{
		{
			"Example",
			aoc22_11.Day11{
				Monkeys: []*aoc22_11.Monkey{
					{
						Items:     []int{79, 98},
						Divisor:   23,
						Operator:  "*",
						Scalar:    19,
						TargetOne: 2,
						TargetTwo: 3,
						Count:     0,
					},
					{
						Items:     []int{54, 65, 75, 74},
						Divisor:   19,
						Operator:  "+",
						Scalar:    6,
						TargetOne: 2,
						TargetTwo: 0,
						Count:     0,
					},
					{
						Items:     []int{79, 60, 97},
						Divisor:   13,
						Operator:  "^",
						Scalar:    2,
						TargetOne: 1,
						TargetTwo: 3,
						Count:     0,
					},
					{
						Items:     []int{74},
						Divisor:   17,
						Operator:  "+",
						Scalar:    3,
						TargetOne: 0,
						TargetTwo: 1,
						Count:     0,
					},
				},
			},
			aoc22_11.Day11{
				Monkeys: []*aoc22_11.Monkey{
					{
						Items:     []int{20, 23, 27, 26},
						Divisor:   23,
						Operator:  "*",
						Scalar:    19,
						TargetOne: 2,
						TargetTwo: 3,
						Count:     2,
					},
					{
						Items:     []int{2080, 25, 167, 207, 401, 1046},
						Divisor:   19,
						Operator:  "+",
						Scalar:    6,
						TargetOne: 2,
						TargetTwo: 0,
						Count:     4,
					},
					{
						Items:     nil,
						Divisor:   13,
						Operator:  "^",
						Scalar:    2,
						TargetOne: 1,
						TargetTwo: 3,
						Count:     3,
					},
					{
						Items:     nil,
						Divisor:   17,
						Operator:  "+",
						Scalar:    3,
						TargetOne: 0,
						TargetTwo: 1,
						Count:     5,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.input

			err := d.ProcessRound()
			assert.NoError(t, err)

			assert.Equal(t, tt.want, d)
		})
	}
}
