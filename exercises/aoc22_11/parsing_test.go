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
			Monkeys: []aoc22_11.Monkey{
				{
					Items:     []int{79, 98},
					Operation: func(old int) int { return old * 19 },
					Test: func(val int) int {
						if val%23 == 0 {
							return 2
						}

						return 3
					},
					Count: 0,
				},
			},
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
		}, aoc22_11.Day11{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := aoc22_11.Day11{}

			err := d.ParseInput(tt.input)
			assert.NoError(t, err)

			assert.Equal(t, tt.want, d)
		})
	}
}
